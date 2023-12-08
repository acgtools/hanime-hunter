package cmd

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/downloader"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return download(args[0])
	},
}

func download(aniURL string) error {
	data, err := resolvers.Resolve(aniURL)
	if err != nil {
		return err
	}

	d := downloader.NewDownloader()

	err = d.Download(data[0])
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}

	return nil
}
