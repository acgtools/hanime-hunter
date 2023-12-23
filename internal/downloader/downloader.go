package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/acgtools/hanime-hunter/internal/request"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/resolvers/hanimetv"
	"github.com/acgtools/hanime-hunter/internal/tui/color"
	"github.com/acgtools/hanime-hunter/internal/tui/progressbar"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/grafov/m3u8"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const defaultGoRoutineNum = 20

type Downloader struct {
	p      *tea.Program
	Option *Option
}

type Option struct {
	OutputDir  string
	Quality    string
	Info       bool
	LowQuality bool
	Retry      uint8
}

func NewDownloader(p *tea.Program, opt *Option) *Downloader {
	return &Downloader{
		p:      p,
		Option: opt,
	}
}

func (d *Downloader) Download(ani *resolvers.HAnime, m *progressbar.Model) error {
	videos := resolvers.SortAniVideos(ani.Videos, d.Option.LowQuality)

	if d.Option.Info {
		log.Infof("Videos available:\n%s", sPrintVideosInfo(videos))
		return nil
	}

	video := videos[0] // by default, download the highest quality
	if d.Option.Quality != "" {
		if v, ok := ani.Videos[strings.ToLower(d.Option.Quality)]; ok {
			video = v
		}
	}

	err := d.save(video, ani.Title, m)
	if err != nil {
		return fmt.Errorf("download file %q: %w", video.Title, err)
	}

	return nil
}

func (d *Downloader) SendPbStatus(fileName, status string) {
	d.p.Send(progressbar.ProgressStatusMsg{
		FileName: fileName,
		Status:   status,
	})
}

func sPrintVideosInfo(vs []*resolvers.Video) string {
	var sb strings.Builder
	for _, v := range vs {
		sb.WriteString(fmt.Sprintf(" Title: %s, Quality: %s, Ext: %s\n", v.Title, v.Quality, v.Ext))
	}

	return sb.String()
}

func (d *Downloader) save(v *resolvers.Video, aniTitle string, m *progressbar.Model) error {
	outputDir := filepath.Join(d.Option.OutputDir, aniTitle)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return fmt.Errorf("create download dir %q: %w", outputDir, err)
	}

	fName := fmt.Sprintf("%s %s.%s", v.Title, v.Quality, v.Ext)
	fPath := filepath.Join(outputDir, fName)
	if f, err := os.Lstat(fPath); err == nil {
		if f.Size() == v.Size {
			log.Infof("File %q exists, Skip ...", fPath)
			return nil
		}
	}

	if !v.IsM3U8 {
		return d.saveSingleVideo(v, fPath, fName, m)
	}

	return d.saveM3U8(v, outputDir, fPath, fName, m)
}

func (d *Downloader) saveSingleVideo(v *resolvers.Video, fPath, fName string, m *progressbar.Model) error {
	pb := progressBar(d.p, v.Size, fName)
	m.AddPb(pb)

	file, err := os.Create(fPath)
	if err != nil {
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return fmt.Errorf("create file %q: %w", fPath, err)
	}
	defer file.Close()

	var curSize int64
	headers := map[string]string{}
	for i := 1; ; i++ {
		written, err := writeFile(d.p, pb.Pw, file, v.URL, headers)
		if err == nil {
			break
		} else if i-1 == int(d.Option.Retry) {
			d.SendPbStatus(fName, progressbar.ErrStatus)
			return err
		}

		curSize += written
		headers["Range"] = fmt.Sprintf("bytes=%d-", curSize)
		d.SendPbStatus(fName, progressbar.RetryStatus)

		time.Sleep(time.Duration(util.RandomInt63n(900, 3000)) * time.Millisecond) //nolint:gomnd
	}

	d.SendPbStatus(fName, progressbar.CompleteStatus)

	return nil
}

func writeFile(p *tea.Program, pw *progressbar.ProgressWriter, file *os.File, u string, headers map[string]string) (int64, error) {
	resp, err := request.Request(http.MethodGet, u, headers)
	if err != nil {
		return 0, fmt.Errorf("send request to %q: %w", u, err)
	}
	defer resp.Body.Close()

	pw.File = file
	pw.Reader = resp.Body
	return pw.Start(p) //nolint:wrapcheck
}

func (d *Downloader) saveM3U8(v *resolvers.Video, outputDir, fPath, fName string, m *progressbar.Model) error {
	segURIs, mediaPL, err := getSegURIs(v.URL)
	if err != nil {
		return err
	}

	tmpDir := filepath.Join(outputDir, "tmp-"+fName)
	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("make dir %q: %w", tmpDir, err)
	}
	defer os.RemoveAll(tmpDir)

	fileListPath := filepath.Join(tmpDir, "fileList.txt")
	err = createTmpFileList(fileListPath, len(segURIs))
	if err != nil {
		return err
	}

	key, iv, err := getKeyIV(mediaPL)
	if err != nil {
		return err
	}

	pb := countProgressBar(d.p, int64(len(segURIs)), fName)
	m.AddPb(pb)

	sem := semaphore.NewWeighted(defaultGoRoutineNum)
	group, ctx := errgroup.WithContext(context.Background())
	dlTS := func(idx, u string) func() error {
		return func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("download TS file: %w", err)
			}
			defer sem.Release(1)

			tsPath := filepath.Join(tmpDir, idx+".ts")

			for i := 1; ; i++ {
				err := saveTS(tsPath, u, key, iv)
				if err == nil {
					break
				} else if i-1 == int(d.Option.Retry) {
					return err
				}
				log.Debugf("err: %s", err)
				log.Debugf("retry download %s", tsPath)
			}

			time.Sleep(time.Duration(util.RandomInt63n(900, 3000)) * time.Millisecond) //nolint:gomnd

			pb.Pc.Increase()
			return nil
		}
	}

	for i, u := range segURIs {
		group.Go(dlTS(strconv.Itoa(i), u))
	}
	if err := group.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return d.mergeFiles(fileListPath, fName, fPath)
}

func (d *Downloader) mergeFiles(fileListPath, fName, fPath string) error {
	d.SendPbStatus(fName, progressbar.MergingStatus)

	err := util.MergeToMP4(fileListPath, fPath)
	if err != nil {
		return fmt.Errorf("file merge: %w", err)
	}

	d.SendPbStatus(fName, progressbar.CompleteStatus)

	return nil
}

func createTmpFileList(path string, num int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer file.Close()

	for i := 0; i < num; i++ {
		_, err := file.WriteString(fmt.Sprintf("file '%s.ts'\n", strconv.Itoa(i)))
		if err != nil {
			return fmt.Errorf("write file list: %w", err)
		}
	}

	return nil
}

func getSegURIs(u string) ([]string, *m3u8.MediaPlaylist, error) {
	m3u8Data, err := getM3U8Data(u)
	if err != nil {
		return nil, nil, err
	}

	list, listType, err := m3u8.DecodeFrom(bytes.NewReader(m3u8Data), true)
	if err != nil {
		return nil, nil, fmt.Errorf("parse m3u8 data: %w", err)
	}
	if listType != m3u8.MEDIA {
		return nil, nil, errors.New("no media data found")
	}
	mediaPL := list.(*m3u8.MediaPlaylist) //nolint:forcetypeassert

	segURIs := make([]string, 0)
	for _, s := range mediaPL.Segments {
		if s == nil {
			continue
		}
		segURIs = append(segURIs, s.URI)
	}

	return segURIs, mediaPL, nil
}

func saveTS(path, u string, key, iv []byte) error {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36",
	}

	resp, err := util.Get(hanimetv.NewClient(), u, headers)
	if err != nil {
		return fmt.Errorf("download TS file: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read data from %q: %w", u, err)
	}

	if len(data) == 0 { // if there is no data here, skip
		return nil
	}

	tsData, err := util.AESDecrypt(data, key, iv)
	if err != nil {
		return fmt.Errorf("decrypt data from %q: %w", u, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(tsData)
	if err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}

	return nil
}

func getKeyIV(mediaPL *m3u8.MediaPlaylist) ([]byte, []byte, error) {
	resp, err := request.Request(http.MethodGet, mediaPL.Key.URI, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("get m3u8 Key: %w", err)
	}
	defer resp.Body.Close()

	key, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read m3u8 Key: %w", err)
	}

	iv := key
	if mediaPL.Key.IV != "" {
		iv = []byte(mediaPL.Key.IV)
	}

	return key, iv, nil
}

func getM3U8Data(u string) ([]byte, error) {
	client := hanimetv.NewClient()
	headers := map[string]string{
		"User-Agent": resolvers.UA,
	}

	resp, err := util.Get(client, u, headers)
	if err != nil {
		return nil, fmt.Errorf("get m3u8 data: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read m3u8 data: %w", err)
	}

	return data, nil
}

func countProgressBar(p *tea.Program, total int64, fileName string) *progressbar.ProgressBar {
	pc := &progressbar.ProgressCounter{
		Total:      total,
		Downloaded: atomic.Int64{},
		FileName:   fileName,
		Onprogress: func(fileName string, ratio float64) {
			p.Send(progressbar.ProgressMsg{
				FileName: fileName,
				Ratio:    ratio,
			})
		},
	}

	colors := color.PbColors.Colors()

	pb := &progressbar.ProgressBar{
		Pc:       pc,
		Progress: progress.New(progress.WithGradient(colors[0], colors[1])),
		FileName: fileName,
		Status:   progressbar.DownloadingStatus,
	}

	return pb
}

func progressBar(p *tea.Program, total int64, fileName string) *progressbar.ProgressBar {
	pw := &progressbar.ProgressWriter{
		FileName: fileName,
		Total:    total,
		OnProgress: func(fileName string, ratio float64) {
			p.Send(progressbar.ProgressMsg{
				FileName: fileName,
				Ratio:    ratio,
			})
		},
	}

	colors := color.PbColors.Colors()

	pb := &progressbar.ProgressBar{
		Pw:       pw,
		Progress: progress.New(progress.WithGradient(colors[0], colors[1])),
		FileName: fileName,
		Status:   progressbar.DownloadingStatus,
	}

	return pb
}
