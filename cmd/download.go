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
	"golang.org/x/sync/semaphore"
)

const (
	defaultRetries = 10
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
		Retry:      cfg.DLOpt.Retry,
	})

	if d.Option.Info {
		log.Infof("Start fetching anime info")
	} else {
		log.Info("Start downloading ...")
	}

	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	//sem := semaphore.NewWeighted(int64(2))
	wg := sync.WaitGroup{}
	errs := make([]error, len(anis))
	for i, ani := range anis {
		wg.Add(1)

		go func(idx int, a *resolvers.HAnime) {
			defer wg.Done()

			if err := sem.Acquire(ctx, 1); err != nil {
				log.Errorf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			if err := d.Download(a, m); err != nil {
				errs[idx] = err
			}
		}(i, ani)
	}

	go func() {
		wg.Wait()
		p.Send(progressbar.ProgressCompleteMsg{})
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	for _, e := range errs {
		if e == nil {
			continue
		}
		log.Errorf("dl: %v", e)
	}

	return nil
}

func init() {
	// DL Opts
	dlCmd.Flags().StringP("output-dir", "o", "", "output directory")
	dlCmd.Flags().StringP("quality", "q", "", "specify video quality. e.g. 1080p, 720p, 480p ...")
	dlCmd.Flags().BoolP("info", "i", false, "get anime info only")
	dlCmd.Flags().Bool("low-quality", false, "download the lowest quality video")
	dlCmd.Flags().Uint8("retry", defaultRetries, "number of retries, max 255")

	_ = viper.BindPFlag("DLOpt.OutputDir", dlCmd.Flags().Lookup("output-dir"))
	_ = viper.BindPFlag("DLOpt.Quality", dlCmd.Flags().Lookup("quality"))
	_ = viper.BindPFlag("DLOpt.Info", dlCmd.Flags().Lookup("info"))
	_ = viper.BindPFlag("DLOpt.LowQuality", dlCmd.Flags().Lookup("low-quality"))
	_ = viper.BindPFlag("DLOpt.Retry", dlCmd.Flags().Lookup("retry"))

	// Resolver Opts
	dlCmd.Flags().BoolP("series", "s", false, "download full series")

	_ = viper.BindPFlag("ResolverOpt.Series", dlCmd.Flags().Lookup("series"))
}
