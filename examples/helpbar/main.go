package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/helpbar"
)

type keyMap struct {
	Inc, Help, Quit key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Inc:  key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "increment")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit: key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}

type model struct {
	help helpbar.Model
	keys keyMap
	n, w int
}

func newModel() model {
	keys := defaultKeyMap()
	return model{
		keys: keys,
		help: helpbar.New(helpbar.WithToggle(keys.Help), helpbar.WithGlobal(keys.Quit)),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.help.SetWidth(m.w)
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Inc):
			m.n++
		default:
			m.help, _ = m.help.Update(msg)
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	bar := m.help.View(m.keys.Inc)
	body := lipgloss.Place(m.w, 1, lipgloss.Center, lipgloss.Top, fmt.Sprintf("count: %d", m.n))
	view := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, body, bar))
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
