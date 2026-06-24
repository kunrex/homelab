package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	stepStyle    = lipgloss.NewStyle().Bold(true)
	countStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	spinStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	checkmark = successStyle.Render("✓")
)

func PrintStep(label, detail string) {
	fmt.Printf("  %s  %s  %s\n", checkmark, stepStyle.Render(label), dimStyle.Render(detail))
}

func Confirm(question string) bool {
	p := tea.NewProgram(confirmModel{question: question})
	m, err := p.Run()
	if err != nil {
		return false
	}
	return m.(confirmModel).result
}

type confirmModel struct {
	question string
	yes      bool
	result   bool
}

var (
	confirmActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	confirmDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "h", "right", "l", "tab":
			m.yes = !m.yes
		case "y", "Y":
			m.result = true
			return m, tea.Quit
		case "n", "N":
			m.result = false
			return m, tea.Quit
		case "enter", " ":
			m.result = m.yes
			return m, tea.Quit
		case "ctrl+c", "esc", "q":
			m.result = false
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	yes := confirmDimStyle.Render("Yes")
	no := confirmDimStyle.Render("No")
	if m.yes {
		yes = confirmActiveStyle.Render("[ Yes ]")
	} else {
		no = confirmActiveStyle.Render("[ No ]")
	}
	return fmt.Sprintf("\n  %s   %s   %s\n", stepStyle.Render(m.question), yes, no)
}

func WithSpinner(label string, fn func() error) error {
	s := spinner.Dot
	done := make(chan error, 1)
	go func() { done <- fn() }()

	i := 0
	ticker := time.NewTicker(s.FPS)
	defer ticker.Stop()
	for {
		select {
		case err := <-done:
			fmt.Print("\r\033[K")
			return err
		case <-ticker.C:
			fmt.Printf("\r  %s  %s", spinStyle.Render(s.Frames[i%len(s.Frames)]), stepStyle.Render(label))
			i++
		}
	}
}
