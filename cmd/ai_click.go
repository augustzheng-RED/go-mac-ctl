package cmd

import (
	"fmt"

	"go-mac-ctl/internal/tools"

	"github.com/spf13/cobra"
)

var aiClickCmd = &cobra.Command{
	Use:   "ai-click [instruction]",
	Short: "Use OCR plus OpenAI to choose and click the best on-screen target for a natural-language instruction",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := tools.AIClick(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		centerX := result.Candidate.Bounds.X + result.Candidate.Bounds.Width/2
		centerY := result.Candidate.Bounds.Y + result.Candidate.Bounds.Height/2
		fmt.Printf("AI clicked %q at center (%d, %d)\n", result.Candidate.Text, centerX, centerY)
		if result.Reason != "" {
			fmt.Println("Reason:", result.Reason)
		}
	},
}

func init() {
	rootCmd.AddCommand(aiClickCmd)
}
