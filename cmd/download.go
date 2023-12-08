package cmd

import (
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/downloader"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := NewCfg()
		if err != nil {
			return err
		}

		logLevel, err := log.ParseLevel(cfg.Log.Level)
		if err != nil {
			return fmt.Errorf("parse log level: %w", err)
		}
		log.SetLevel(logLevel)
		log.SetReportTimestamp(false)

		return download(args[0], cfg)
	},
}

func download(aniURL string, cfg *Config) error {
	anis, err := resolvers.Resolve(aniURL, &resolvers.Option{
		Series:   cfg.ResolverOpt.Series,
		PlayList: cfg.ResolverOpt.PlayList,
	})
	if err != nil {
		return err
	}

	d := downloader.NewDownloader(&downloader.Option{
		OutputDir: cfg.DLOpt.OutputDir,
	})

	for _, ani := range anis {
		err = d.Download(ani)
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
	}

	return nil
}

func init() {
	dlCmd.Flags().String("output-dir", "", "output directory")
	dlCmd.Flags().Bool("series", false, "download full series")

	_ = viper.BindPFlag("dlopt.outputdir", dlCmd.Flags().Lookup("output-dir"))
	_ = viper.BindPFlag("resolveropt.series", dlCmd.Flags().Lookup("series"))
}
