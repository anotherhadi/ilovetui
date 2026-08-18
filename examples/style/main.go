package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type model struct{ focused bool }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "tab":
			m.focused = !m.focused
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	box := style.S.Panel
	if m.focused {
		box = style.S.PanelFocused
	}
	return tea.NewView(style.RenderWithTitle(box, "Panel", "tab: toggle focus  q: quit", 30, 5))
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
