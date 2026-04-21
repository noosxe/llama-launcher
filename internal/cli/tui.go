package cli

import (
	"github.com/spf13/cobra"
	"github.com/noosxe/llama-launcher/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(cfgFile); err != nil {
			println("Error:", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
