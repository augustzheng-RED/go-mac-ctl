package cmd

import (
	"fmt"
	"time"

	"go-mac-ctl/internal/tools"

	"github.com/spf13/cobra"
)

var chromeOpenCmd = &cobra.Command{
	Use:   "chrome-open [url]",
	Short: "Open a URL in Google Chrome",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := tools.ChromeOpenAndWait(args[0], 5*time.Second); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Opened in Google Chrome:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(chromeOpenCmd)
}
