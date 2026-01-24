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

type Pickable interface {
	GetName() string
	GetGID() string
}

type pickerModel[T Pickable] struct {
	items    []T
	prompt   string
	cursor   int
	selected *T
	quitting bool
}

func newGenericPickerModel[T Pickable](items []T, prompt string) pickerModel[T] {
	return pickerModel[T]{
		items:  items,
		prompt: prompt,
		cursor: 0,
	}
}

func (m pickerModel[T]) Init() tea.Cmd {
	return nil
}

func (m pickerModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = &m.items[m.cursor]
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m pickerModel[T]) View() string {
	var b strings.Builder

	b.WriteString(m.prompt + " (j/k to move, enter to select, q to quit):\n\n")

	for i, item := range m.items {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}

		name := item.GetName()
		if len(name) > 60 {
			name = name[:57] + "..."
		}

		line := fmt.Sprintf("%s%s %s\n", cursor, style.Render(name), gidStyle.Render("("+item.GetGID()+")"))
		b.WriteString(line)
	}

	return b.String()
}

func pick[T Pickable](items []T, prompt string) (*T, error) {
	if len(items) == 0 {
		return nil, errors.NewGeneralError("no items to pick from", nil)
	}

	p := tea.NewProgram(newGenericPickerModel(items, prompt))
	result, err := p.Run()
	if err != nil {
		return nil, errors.NewGeneralError("picker failed", err)
	}

	m := result.(pickerModel[T])
	if m.selected == nil {
		return nil, errors.NewGeneralError("no item selected", nil)
	}

	return m.selected, nil
}

func pickTask(matches []taskMatch) (*taskMatch, error) {
	return pick(matches, "Select a task")
}
