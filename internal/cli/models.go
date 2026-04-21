package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/noosxe/llama-launcher/internal/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			fmt.Fprintf(stderr, "Failed to load config: %v\n", err)
			return
		}
		for _, m := range cfg.Models {
			fmt.Printf("%s - %s\n", m.Name, m.ModelPath)
		}
	},
}

var stderr = &writer{os.Stderr}

type writer struct {
	*os.File
}

func (w *writer) Write(p []byte) (n int, err error) {
	return fmt.Fprintln(w, string(p))
}

func init() {
	rootCmd.AddCommand(listCmd)
}
