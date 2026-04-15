package cmd

import (
	"github.com/spf13/cobra"
	"llama-launcher/tui"
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
