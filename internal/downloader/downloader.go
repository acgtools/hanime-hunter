package downloader

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/request"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(ani *resolvers.HAnime) error {
	videos := util.SortAniVideos(ani.Videos)

	err := d.save(videos[0])
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	return nil
}

func (d *Downloader) save(v *resolvers.Video) error {
	name := fmt.Sprintf("%s_%s.%s", v.Title, v.Quality, v.Ext)

	const basePath = "./test/dl/"
	err := os.MkdirAll(basePath, os.ModePerm)
	filePath := filepath.Join(basePath, name)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file %q: %w", name, err)
	}

	resp, err := request.Request(http.MethodGet, v.URL)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("wirte file: %w", err)
	}

	return nil
}
