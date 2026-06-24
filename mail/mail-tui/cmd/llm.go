package cmd

import (
	"fmt"
	"slices"
	"strings"

	"mail-tui/cfg"
	"mail-tui/llm"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	okStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	labelStyle   = lipgloss.NewStyle().Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	activeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
)

func containsModel(models []string, model string) bool {
	return slices.Contains(models, model)
}

func printModelNotFound(model string, available []string) {
	fmt.Printf("  %s  %s  %s\n", warnStyle.Render("✗"), labelStyle.Render("Model"), mutedStyle.Render(model+" — not loaded at "+cfg.Cfg.LLM.Endpoint))
	if len(available) > 0 {
		fmt.Printf("      %s  %s\n", mutedStyle.Render("available:"), mutedStyle.Render(strings.Join(available, ", ")))
	}
}

func printLLMUnreachable(err error) {
	fmt.Printf("  %s  %s  %s\n", warnStyle.Render("✗"), labelStyle.Render("LLM"), mutedStyle.Render(cfg.Cfg.LLM.Endpoint))
	if err != nil {
		fmt.Printf("      %s\n", mutedStyle.Render(err.Error()))
	}
}

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Inspect the LLM endpoint",
}

var llmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check whether the LLM endpoint is reachable",
	RunE: func(_ *cobra.Command, _ []string) error {
		client := llm.NewClient(cfg.Cfg.LLM)
		up, latency, err := client.Status()
		endpoint := mutedStyle.Render(cfg.Cfg.LLM.Endpoint)
		if err != nil {
			fmt.Printf("  %s  %s  %s\n", warnStyle.Render("✗"), labelStyle.Render("LLM"), endpoint)
			fmt.Printf("      %s\n", mutedStyle.Render(err.Error()))
			return nil
		}
		if up {
			fmt.Printf("  %s  %s  %s  %s\n",
				okStyle.Render("✓"),
				labelStyle.Render("LLM"),
				endpoint,
				mutedStyle.Render(fmt.Sprintf("%dms", latency.Milliseconds())),
			)
		} else {
			fmt.Printf("  %s  %s  %s\n", warnStyle.Render("✗"), labelStyle.Render("LLM"), endpoint)
		}
		return nil
	},
}

var llmModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List models available at the LLM endpoint",
	RunE: func(_ *cobra.Command, _ []string) error {
		client := llm.NewClient(cfg.Cfg.LLM)
		if up, _, err := client.Status(); err != nil || !up {
			printLLMUnreachable(err)
			return fmt.Errorf("llm unreachable")
		}
		available, err := client.Models()
		if err != nil {
			return fmt.Errorf("fetching models: %w", err)
		}
		if len(available) == 0 {
			fmt.Println(mutedStyle.Render("  No models found."))
			return nil
		}
		fmt.Printf("  %s  %s\n", labelStyle.Render("Endpoint"), mutedStyle.Render(cfg.Cfg.LLM.Endpoint))
		fmt.Println()
		for _, m := range available {
			if m == cfg.Cfg.LLM.Model {
				fmt.Printf("  %s  %s\n", activeStyle.Render("*"), activeStyle.Render(m))
			} else {
				fmt.Printf("     %s\n", mutedStyle.Render(m))
			}
		}
		return nil
	},
}

func init() {
	llmCmd.AddCommand(llmStatusCmd)
	llmCmd.AddCommand(llmModelsCmd)
}
