package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/style"
)

// confirmedMsg is what the confirmation modal reports back with. The app
// listens for it like any other message - it never holds a reference to the
// modal, and the modal never knows what confirming means.
type confirmedMsg struct{}

// confirm is the modal's content: a model, so it owns its own keys. The host
// no longer has to ask which modal is on top to know where "y" should go.
type confirm struct{}

func (c confirm) Init() tea.Cmd { return nil }

func (c confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "y" {
		// Report back and close itself: Close() is a package-level command,
		// so the content needs no reference to the modal.Model either.
		return c, tea.Batch(
			func() tea.Msg { return confirmedMsg{} },
			modal.Close(),
		)
	}
	return c, nil
}

func (c confirm) View() tea.View {
	return tea.NewView("This can't be undone.\n\ny: confirm  esc: cancel")
}

type model struct {
	m             modal.Model
	deleted       bool
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

	case confirmedMsg:
		m.deleted = true
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "m":
			return m, modal.Show("Delete file?", confirm{})

		case "n":
			if m.m.Open() {
				return m, modal.Show("Really sure?", modal.Text("There's no undo for this one either."))
			}

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
	title := lipgloss.NewStyle().Bold(true).Foreground(style.S.Primary).Render("My App")
	text := "Some regular content, styled with theme colors,\nso you can see it turn flat gray behind the modal."
	if m.deleted {
		text = "File deleted - the modal's content reported back\nwith its own message, and closed itself."
	}
	body := lipgloss.NewStyle().Foreground(style.S.Text).Render(text)
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
