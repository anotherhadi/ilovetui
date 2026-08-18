// Package overview is the fullapp example's first page: the simplest pane
// there is - no keys, no commands, just text sized to whatever room the
// shell gives it.
package overview

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// Model is the page. It's an ordinary tea.Model: nothing about it knows it
// lives in a shell, and it never implements HelpBindings because it has no
// keys of its own to advertise.
type Model struct {
	width, height int
}

func New() Model { return Model{} }

func (m Model) Init() tea.Cmd { return nil }

// Update only tracks the size the shell hands down as a tea.WindowSizeMsg,
// same as if the page were the whole program.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
	}
	return m, nil
}

func (m Model) View() tea.View {
	body := lipgloss.JoinVertical(lipgloss.Center,
		style.S.Bold.Render("Overview"),
		"",
		style.S.Faint.Render("tab focuses this pane, but there's"),
		style.S.Faint.Render("nothing here to focus on."),
	)
	return tea.NewView(lipgloss.NewStyle().
		Width(m.width).Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(body))
}
