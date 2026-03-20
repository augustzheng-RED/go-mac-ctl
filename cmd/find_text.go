package cmd

import (
	"encoding/json"
	"fmt"

	"go-mac-ctl/internal/tools"

	"github.com/spf13/cobra"
)

var findTextCmd = &cobra.Command{
	Use:   "find-text [query]",
	Short: "Capture the screen and prepare a text-finding result",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := tools.FindText(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(findTextCmd)
}
