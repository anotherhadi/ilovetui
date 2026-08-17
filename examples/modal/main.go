package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/style"
)

const confirmID = "confirm"

type model struct {
	m             modal.Model
	width, height int
}

func newModel() model {
	return model{m: modal.New()}
}

func (m model) Init() tea.Cmd { return m.m.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "m":
			return m, modal.Show("Delete file?", "This can't be undone.\n\ny: confirm  esc: cancel",
				modal.WithID(confirmID))

		case "n":
			if m.m.Open() {
				return m, modal.Show("Really sure?", "There's no undo for this one either.")
			}

		case "esc":
			if m.m.Open() {
				return m, modal.Close()
			}

		case "y":
			if m.m.TopID() == confirmID {
				return m, modal.Dismiss(confirmID)
			}
		}
	}

	var cmd tea.Cmd
	m.m, cmd = m.m.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	title := lipgloss.NewStyle().Bold(true).Foreground(style.S.Primary).Render("My App")
	body := lipgloss.NewStyle().Foreground(style.S.Text).Render(
		"Some regular content, styled with theme colors,\nso you can see it turn flat gray behind the modal.")
	help := lipgloss.NewStyle().Foreground(style.S.Subtle).Render(
		"m: open modal  n: open nested modal  y: confirm  esc: cancel  q: quit")

	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", help))

	view := tea.NewView(m.m.Render(background))
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
