package cmd

import (
	"os"
	"fmt"
	"strings"

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

	rootCmd.SilenceErrors = true
    rootCmd.SilenceUsage = true

    rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
        return err
    })
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
        if strings.Contains(err.Error(), "unknown command") {
            fmt.Printf("  %s  Unknown command\n", warnStyle.Render("✗"))
        } else {
            fmt.Fprintln(os.Stderr, err)
        }

        os.Exit(1)
	}
}
