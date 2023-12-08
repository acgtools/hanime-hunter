package downloader

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/request"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/tui/progress"
	"github.com/charmbracelet/log"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
	Option *Option
}

type Option struct {
	OutputDir string
}

func NewDownloader(opt *Option) *Downloader {
	return &Downloader{
		Option: opt,
	}
}

func (d *Downloader) Download(ani *resolvers.HAnime) error {
	videos := resolvers.SortAniVideos(ani.Videos)

	v := videos[0]
	log.Info("Start Downloading: ", "title", v.Title, "quality", v.Quality, "extension", v.Ext)

	err := d.save(videos[0])
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	return nil
}

func (d *Downloader) save(v *resolvers.Video) error {
	fPath := fmt.Sprintf("%s %s.%s", v.Title, v.Quality, v.Ext)

	if d.Option.OutputDir != "" {
		outputDir := d.Option.OutputDir
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

	err = progress.StartProgressBar(resp, file)
	if err != nil {
		return fmt.Errorf("start progress bar: %w", err)
	}

	return nil
}
