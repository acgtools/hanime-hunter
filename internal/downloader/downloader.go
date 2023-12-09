package downloader

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/request"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/tui/color"
	"github.com/acgtools/hanime-hunter/internal/tui/progressbar"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
	p      *tea.Program
	Option *Option
}

type Option struct {
	OutputDir string
}

func NewDownloader(p *tea.Program, opt *Option) *Downloader {
	return &Downloader{
		p:      p,
		Option: opt,
	}
}

func (d *Downloader) Download(ani *resolvers.HAnime, m *progressbar.Model) error {
	videos := resolvers.SortAniVideos(ani.Videos)

	err := d.save(videos[0], ani.Title, m)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	return nil
}

func (d *Downloader) save(v *resolvers.Video, aniTitle string, m *progressbar.Model) error {
	fPath := fmt.Sprintf("%s %s.%s", v.Title, v.Quality, v.Ext)

	if d.Option.OutputDir != "" {
		outputDir := filepath.Join(d.Option.OutputDir, aniTitle)
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create output dirs: %w", err)
		}
		fPath = filepath.Join(outputDir, fPath)
	}

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("create file %q: %w", fPath, err)
	}
	defer file.Close()

	resp, err := request.Request(http.MethodGet, v.URL)
	if err != nil {
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
