package cmd

import (
	"fmt"
	"strings"

	"go-mac-ctl/internal/act"

	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key [name]",
	Short: "Press a supported key such as enter, tab, or escape",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])

		if err := act.PressKey(key); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Pressed:", key)
	},
}

func init() {
	rootCmd.AddCommand(keyCmd)
}
