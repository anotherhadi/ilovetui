// Package settings is the fullapp example's third page: a short list of
// toggles. It's the one page with keys of its own, so it's what proves the
// shell routes presses to the focused pane and lists that pane's bindings in
// the help bar.
package settings

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type toggle struct {
	name string
	on   bool
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
}

// Model is the page. The cursor is a plain int: three lines don't need a
// list.Model behind them.
type Model struct {
	keys    keyMap
	toggles []toggle
	cursor  int

	width, height int
}

func New() Model {
	return Model{
		keys: keyMap{
			Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Toggle: key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "toggle")),
		},
		toggles: []toggle{
			{name: "Nerd fonts", on: style.S.NerdFonts},
			{name: "Notifications", on: true},
			{name: "Telemetry", on: false},
		},
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		// The shell only sends these while this pane has focus, so there's
		// no focused check to make here.
		switch {
		case key.Matches(msg, m.keys.Up):
			m.cursor = max(m.cursor-1, 0)
		case key.Matches(msg, m.keys.Down):
			m.cursor = min(m.cursor+1, len(m.toggles)-1)
		case key.Matches(msg, m.keys.Toggle):
			m.toggles[m.cursor].on = !m.toggles[m.cursor].on
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	rows := make([]string, len(m.toggles))
	for i, t := range m.toggles {
		cursor, box := "  ", "[ ]"
		if i == m.cursor {
			cursor = "> "
		}
		if t.on {
			box = "[x]"
		}
		line := cursor + box + " " + t.name
		if i == m.cursor {
			line = lipgloss.NewStyle().Foreground(style.S.Primary).Render(line)
		}
		rows[i] = line
	}

	return tea.NewView(lipgloss.NewStyle().
		Width(m.width).Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, rows...)))
}

// HelpBindings makes the page's keys show up in the shell's help bar while
// it has focus. It's the optional half of the contract: overview and metrics
// have no keys and don't implement it.
func (m Model) HelpBindings() []key.Binding {
	return []key.Binding{m.keys.Up, m.keys.Down, m.keys.Toggle}
}
