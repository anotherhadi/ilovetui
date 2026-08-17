package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/notification"
	"github.com/anotherhadi/ilovetui/style"
)

var positions = []struct {
	name string
	pos  notification.Position
}{
	{"top", notification.Top},
	{"top-left", notification.TopLeft},
	{"top-right", notification.TopRight},
	{"bottom", notification.Bottom},
	{"bottom-left", notification.BottomLeft},
	{"bottom-right", notification.BottomRight},
}

const stickyID = "sticky-demo"

type model struct {
	notif         notification.Model
	posIdx        int
	width, height int
}

func newModel() model {
	return model{notif: notification.New(notification.WithPosition(positions[0].pos))}
}

func (m model) Init() tea.Cmd {
	return m.notif.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1":
			return m, notification.Show("Info", "Just so you know.", notification.Info)
		case "2":
			return m, notification.Show("Success", "Config written to disk.", notification.Success)
		case "3":
			return m, notification.Show("Warning", "Disk space getting low on /dev/sda1.", notification.Warning)
		case "4":
			return m, notification.Show("Error", "Failed to reach the remote host.", notification.Error)

		case "s":
			return m, notification.Show("Sticky", "Stays until you press d.",
				notification.Info, notification.WithID(stickyID), notification.WithDuration(0))
		case "d":
			return m, notification.Dismiss(stickyID)

		case "p":
			m.posIdx = (m.posIdx + 1) % len(positions)
			m.notif = notification.New(notification.WithPosition(positions[m.posIdx].pos))
			return m, m.notif.Init()
		}
	}

	var cmd tea.Cmd
	m.notif, cmd = m.notif.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	help := lipgloss.NewStyle().Foreground(style.S.Subtle).Render(
		"1-4: info/success/warning/error  s: sticky  d: dismiss sticky  p: position (" +
			positions[m.posIdx].name + ")  q: quit")

	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, help)

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
