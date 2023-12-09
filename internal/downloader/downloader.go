package downloader

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/request"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/tui/color"
	"github.com/acgtools/hanime-hunter/internal/tui/progressbar"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Downloader struct {
	p      *tea.Program
	Option *Option
}

type Option struct {
	OutputDir  string
	Quality    string
	Info       bool
	LowQuality bool
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
		log.Infof("Videos available: %s", SPrintVideosInfo(videos))
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
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return fmt.Errorf("download file: %w", err)
	}

	return nil
}

func SPrintVideosInfo(vs []*resolvers.Video) string {
	var sb strings.Builder
	for _, v := range vs {
		sb.WriteString(fmt.Sprintf(" Title: %s, Quality: %s, Ext: %s\n", v.Title, v.Quality, v.Ext))
	}

	return sb.String()
}

func (d *Downloader) save(v *resolvers.Video, aniTitle string, m *progressbar.Model) error {
	fPath := fmt.Sprintf("%s %s.%s", v.Title, v.Quality, v.Ext)

	outputDir := filepath.Join(d.Option.OutputDir, aniTitle)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return err
	}
	fPath = filepath.Join(outputDir, fPath)

	if f, err := os.Lstat(fPath); err == nil {
		if f.Size() == v.Size {
			log.Infof("File %q exists, Skip ...", fPath)
			return nil
		}
	}

	file, err := os.Create(fPath)
	if err != nil {
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return fmt.Errorf("create file %q: %w", fPath, err)
	}
	defer file.Close()

	resp, err := request.Request(http.MethodGet, v.URL)
	if err != nil {
		d.p.Send(progressbar.ProgressErrMsg{Err: err})
		return fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	fileName := filepath.Base(file.Name())

	pw := &progressbar.ProgressWriter{
		Total:    resp.ContentLength,
		File:     file,
		Reader:   resp.Body,
		FileName: fileName,
		OnProgress: func(fileName string, ratio float64) {
			d.p.Send(progressbar.ProgressMsg{
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
	}

	m.Mux.Lock()
	m.Pbs[fileName] = pb
	m.Mux.Unlock()

	pw.Start(d.p)

	return nil
}
