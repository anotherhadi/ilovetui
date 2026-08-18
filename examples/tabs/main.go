package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/tabs"
)

type pane string

func (p pane) Init() tea.Cmd                      { return nil }
func (p pane) Update(tea.Msg) (tabs.Tab, tea.Cmd) { return p, nil }
func (p pane) View() string                       { return string(p) }

type model struct{ tabs tabs.Model }

func newModel() model {
	items := []tabs.Item{
		{Title: "First", Model: pane("first content")},
		{Title: "Second", Model: pane("second content")},
		{Title: "Third", Model: pane("third content")},
	}
	return model{tabs: tabs.New(items)}
}

func (m model) Init() tea.Cmd { return m.tabs.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "q" {
		return m, tea.Quit
	}
	if s, ok := msg.(tea.WindowSizeMsg); ok {
		m.tabs.SetSize(s.Width, s.Height)
		msg = tea.WindowSizeMsg{Width: m.tabs.ContentWidth(), Height: m.tabs.ContentHeight()}
	}
	var cmd tea.Cmd
	m.tabs, cmd = m.tabs.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	view := tea.NewView(m.tabs.View())
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
