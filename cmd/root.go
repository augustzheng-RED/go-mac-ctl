package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-mac-ctl",
	Short: "Control browser and desktop actions on macOS",
	Long:  "go-mac-ctl provides a small set of commands for opening pages, finding text on screen, clicking text, typing, and pressing keys.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
