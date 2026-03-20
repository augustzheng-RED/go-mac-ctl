package cmd

import (
	"fmt"

	"go-mac-ctl/internal/tools"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [url-or-path]",
	Short: "Open a URL, file, or folder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := tools.OpenURL(args[0]); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Opened: ", args[0])
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
