package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"mail-tui/cfg"
	"mail-tui/internal"
	"mail-tui/llm"
	"mail-tui/models"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch, classify and digest unprocessed mails",
	RunE:  runDigest,
}

func runDigest(_ *cobra.Command, _ []string) error {
	mailServer := internal.NewServerClient()
	llmClient := llm.NewClient(cfg.Cfg.LLM)

	if up, _, err := llmClient.Status(); err != nil || !up {
		printLLMUnreachable(err)
		return fmt.Errorf("llm unreachable")
	}

	available, err := llmClient.Models()
	if err != nil {
		return fmt.Errorf("fetching models: %w", err)
	}
	if !containsModel(available, cfg.Cfg.LLM.Model) {
		printModelNotFound(cfg.Cfg.LLM.Model, available)
		return fmt.Errorf("model not available")
	}

	var mails []models.Mail
	if err := internal.WithSpinner("Fetching mails", func() error {
		var err error
		mails, err = mailServer.FetchMails()
		return err
	}); err != nil {
		return fmt.Errorf("fetch: %w", err)
	}
	internal.PrintStep("Fetching", fmt.Sprintf("%d unprocessed", len(mails)))

	if len(mails) == 0 {
		fmt.Println("  Nothing to do.")
		return nil
	}

	relevant, err := internal.ClassifyMails(llmClient, mails)
	if err != nil {
		return fmt.Errorf("classify: %w", err)
	}
	internal.PrintStep("Classifying", fmt.Sprintf("%d relevant", len(relevant)))

	if len(relevant) == 0 {
		fmt.Println("  No relevant mails.")
		internal.PromptMarkProcessed(mailServer)
		return nil
	}

	var digest string
	if err := internal.WithSpinner("Generating digest", func() error {
		var err error
		digest, err = internal.GenerateDigest(llmClient, relevant)
		return err
	}); err != nil {
		return fmt.Errorf("digest: %w", err)
	}
	internal.PrintStep("Digest", "ready")

	fmt.Println()
	if err := internal.RenderMarkdown(digest); err != nil {
		fmt.Println(digest)
	}

	internal.PromptMarkProcessed(mailServer)
	return nil
}
