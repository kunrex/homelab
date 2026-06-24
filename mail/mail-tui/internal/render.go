package internal

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(md string) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		return fmt.Errorf("creating renderer: %w", err)
	}
	out, err := renderer.Render(md)
	if err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}
	fmt.Print(out)
	return nil
}
