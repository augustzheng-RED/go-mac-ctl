package cmd

import (
	"fmt"
	"strconv"

	"go-mac-ctl/internal/tools"

	"github.com/spf13/cobra"
)

var clickTextCmd = &cobra.Command{
	Use:   "click-text [query] [index]",
	Short: "Find text on screen and click the indexed match",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: invalid target index")
			return
		}

		target, err := tools.ClickText(args[0], index)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		centerX := target.Bounds.X + target.Bounds.Width/2
		centerY := target.Bounds.Y + target.Bounds.Height/2
		fmt.Printf("Clicked text %q at center (%d, %d)\n", target.Text, centerX, centerY)
	},
}

func init() {
	rootCmd.AddCommand(clickTextCmd)
}
