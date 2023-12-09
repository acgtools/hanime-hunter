package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var version = "unknown"

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(cmd *cobra.Command, args []string) {
		goVersion := runtime.Version()
		platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

		log.SetTimeFunction(func() time.Time {
			return time.Time{}
		})
		log.Printf("Version: %s, Go: %s, Platform: %s", version, goVersion, platform)
	},
}
