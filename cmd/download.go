package cmd

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/acgtools/hanime-hunter/internal/downloader"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/internal/tui/progressbar"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
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
		Series: cfg.ResolverOpt.Series,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	m := &progressbar.Model{
		Mux: sync.Mutex{},
		Pbs: make(map[string]*progressbar.ProgressBar),
	}
	p := tea.NewProgram(m)

	d := downloader.NewDownloader(p, &downloader.Option{
		OutputDir:  cfg.DLOpt.OutputDir,
		Quality:    cfg.DLOpt.Quality,
		Info:       cfg.DLOpt.Info,
		LowQuality: cfg.DLOpt.LowQuality,
	})

	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	group, ctx := errgroup.WithContext(context.Background())
	dl := func(ani *resolvers.HAnime, m *progressbar.Model) func() error {
		return func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("goroutine download %q acquire semaphore: %w", ani.Title, err)
			}
			defer sem.Release(1)

			err := d.Download(ani, m) //nolint:contextcheck
			if err != nil {
				return fmt.Errorf("download %q error: %w", ani.Title, err)
			}
			return nil
		}
	}

	if d.Option.Info {
		log.Infof("Start fetching anime info")
	} else {
		log.Info("Start downloading ...")
	}

	for _, ani := range anis {
		group.Go(dl(ani, m))
	}

	if _, err := p.Run(); err != nil {
		log.Errorf("Start progress bar %v", err)
	}

	return group.Wait() //nolint:wrapcheck
}

func init() {
	dlCmd.Flags().StringP("output-dir", "o", "", "output directory")
	dlCmd.Flags().StringP("quality", "q", "", "specify video quality. e.g. 1080p, 720p, 480p ...")

	dlCmd.Flags().BoolP("series", "s", false, "download full series")
	dlCmd.Flags().BoolP("info", "i", false, "get anime info only")
	dlCmd.Flags().Bool("low-quality", false, "download the lowest quality video")

	_ = viper.BindPFlag("DLOpt.OutputDir", dlCmd.Flags().Lookup("output-dir"))
	_ = viper.BindPFlag("DLOpt.Quality", dlCmd.Flags().Lookup("quality"))
	_ = viper.BindPFlag("DLOpt.Info", dlCmd.Flags().Lookup("info"))
	_ = viper.BindPFlag("DLOpt.LowQuality", dlCmd.Flags().Lookup("low-quality"))

	_ = viper.BindPFlag("ResolverOpt.Series", dlCmd.Flags().Lookup("series"))
}
