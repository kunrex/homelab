package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"mail-tui/cfg"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:          "mtui",
	Short:        "LLM-curated mail digest",
	SilenceUsage: true,
	SilenceErrors: true,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		return cfg.Init(cfgPath)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file (default: ~/.config/mtui/config.yaml)")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(llmCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
