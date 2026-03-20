package cmd

import (
	"fmt"

	"go-mac-ctl/internal/act"

	"github.com/spf13/cobra"
)

var typeCmd = &cobra.Command{
	Use:   "type [text]",
	Short: "Type text into the focused input",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := act.TypeText(args[0]); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Typed:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(typeCmd)
}
