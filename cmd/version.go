package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "unknown"

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
