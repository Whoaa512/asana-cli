package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/whoaa512/asana-cli/internal/errors"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	gidStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type pickerModel struct {
	matches  []taskMatch
	cursor   int
	selected *taskMatch
	quitting bool
}

func newPickerModel(matches []taskMatch) pickerModel {
	return pickerModel{
		matches: matches,
		cursor:  0,
	}
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.matches)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = &m.matches[m.cursor]
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder

	b.WriteString("Select a task (j/k to move, enter to select, q to quit):\n\n")

	for i, match := range m.matches {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}

		name := match.task.Name
		if len(name) > 60 {
			name = name[:57] + "..."
		}

		line := fmt.Sprintf("%s%s %s\n", cursor, style.Render(name), gidStyle.Render("("+match.task.GID+")"))
		b.WriteString(line)
	}

	return b.String()
}

func pickTask(matches []taskMatch) (*taskMatch, error) {
	if len(matches) == 0 {
		return nil, errors.NewGeneralError("no tasks to pick from", nil)
	}

	p := tea.NewProgram(newPickerModel(matches))
	result, err := p.Run()
	if err != nil {
		return nil, errors.NewGeneralError("picker failed", err)
	}

	m := result.(pickerModel)
	if m.selected == nil {
		return nil, errors.NewGeneralError("no task selected", nil)
	}

	return m.selected, nil
}
