package cmd

import (
	"context"
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/downloader"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/tui/progressbar"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"runtime"
	"sync"
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

	m := &progressbar.Model{
		Mux: sync.Mutex{},
		Pbs: make(map[string]*progressbar.ProgressBar),
	}
	p := tea.NewProgram(m)

	d := downloader.NewDownloader(p, &downloader.Option{
		OutputDir: cfg.DLOpt.OutputDir,
		Info:      cfg.DLOpt.Info,
	})

	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	group, ctx := errgroup.WithContext(context.Background())
	dl := func(ani *resolvers.HAnime, m *progressbar.Model) func() error {
		return func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			err := d.Download(ani, m)
			if err != nil {
				return fmt.Errorf("download error: %w", err)
			}
			return nil
		}
	}

	log.Info("Start downloading ...")

	for _, ani := range anis {
		group.Go(dl(ani, m))
	}

	if _, err := p.Run(); err != nil {
		log.Errorf("Open progress bar %v", err)
	}

	return group.Wait()
}

func init() {
	dlCmd.Flags().StringP("output-dir", "o", "", "output directory")
	dlCmd.Flags().BoolP("series", "s", false, "download full series")
	dlCmd.Flags().BoolP("info", "i", false, "get anime info only")

	_ = viper.BindPFlag("dlopt.outputdir", dlCmd.Flags().Lookup("output-dir"))
	_ = viper.BindPFlag("dlopt.info", dlCmd.Flags().Lookup("info"))

	_ = viper.BindPFlag("resolveropt.series", dlCmd.Flags().Lookup("series"))
}
