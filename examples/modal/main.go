package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/modal"
)

type model struct {
	m             modal.Model
	width, height int
}

func newModel() model { return model{m: modal.New()} }

func (m model) Init() tea.Cmd { return m.m.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "o":
			return m, modal.Show("Hello", modal.Text("This is a modal.\n\nesc: close"))
		case "esc":
			if m.m.Open() {
				return m, modal.Close()
			}
		}
	}
	var cmd tea.Cmd
	m.m, cmd = m.m.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "o: open  q: quit")
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
