package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/notification"
)

type model struct {
	notif         notification.Model
	width, height int
}

func newModel() model { return model{notif: notification.New()} }

func (m model) Init() tea.Cmd { return m.notif.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "s":
			return m, notification.Show("Saved", "Config written to disk.", notification.Success)
		}
	}
	var cmd tea.Cmd
	m.notif, cmd = m.notif.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "s: show  q: quit")
	view := tea.NewView(m.notif.Render(background))
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
