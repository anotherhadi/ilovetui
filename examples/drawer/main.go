package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/drawer"
)

type model struct {
	d             drawer.Model
	width, height int
}

func newModel() model { return model{d: drawer.New()} }

func (m model) Init() tea.Cmd { return m.d.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			return m, drawer.Show("Nav", drawer.Text("Home\nSettings"), drawer.WithSide(drawer.Left))
		case "esc":
			if m.d.Open() {
				return m, drawer.Close()
			}
		}
	}
	var cmd tea.Cmd
	m.d, cmd = m.d.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "l: open  q: quit")
	view := tea.NewView(m.d.Render(background))
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
