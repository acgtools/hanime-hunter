package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var version = "unknown"

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		goVersion := runtime.Version()
		platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

		fmt.Printf("Version: %s, Go: %s, Platform: %s", version, goVersion, platform)
	},
}
