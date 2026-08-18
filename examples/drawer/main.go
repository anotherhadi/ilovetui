package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/drawer"
	"github.com/anotherhadi/ilovetui/style"
)

type model struct {
	d             drawer.Model
	width, height int
}

func newModel() model {
	return model{d: drawer.New()}
}

func (m model) Init() tea.Cmd { return m.d.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "l":
			return m, drawer.Show("Nav", drawer.Text("Home\nProjects\nSettings"),
				drawer.WithSide(drawer.Left), drawer.WithWidth(20))

		case "r":
			return m, drawer.Show("Inspector", drawer.Text("id: 42\nstatus: ok"),
				drawer.WithSide(drawer.Right), drawer.WithWidth(20))

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
	title := lipgloss.NewStyle().Bold(true).Foreground(style.S.Primary).Render("My App")
	body := lipgloss.NewStyle().Foreground(style.S.Text).Render(
		"Some regular content, styled with theme colors,\nso you can see it turn flat gray behind the drawer.")
	help := lipgloss.NewStyle().Foreground(style.S.Subtle).Render(
		"l: open left drawer  r: open right drawer  esc: close the top one  q: quit")

	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", help))

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
