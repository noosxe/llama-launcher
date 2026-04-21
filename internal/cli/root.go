package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var cfgFile string
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "llama-launcher",
	Short:   "Launch llama.cpp containers with preconfigured models",
	Long:    `A CLI tool to manage and launch llama.cpp Docker containers with predefined model configurations.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("llama-launcher: Use --help for usage information")
	},
}

func Execute() {
	// If version is not set by ldflags, try to get it from build info (for go install)
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				Version = info.Main.Version
				rootCmd.Version = Version
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("{{.Name}} version {{.Version}}\n")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/llama-launcher/config.toml or ./config.toml)")
}

func initConfig() {
	// Config loading will be implemented
}
