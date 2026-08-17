// Command messaging demonstrates the two ways panes talk without holding a
// reference to each other: "control" sends arbitrary commands to "editor"
// by id via SendMsg (1/2/3 keys), and asks layout to move focus there via
// RequestFocusMsg (enter key) - the same pattern a real sidebar would use
// to both drive and jump to a content pane it selected.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

// commandMsg is what "control" sends to "editor" - an app-defined message,
// entirely opaque to layout itself (see SendMsg).
type commandMsg struct{ text string }

type controlPane struct {
	id      string // learned from SizeMsg.ID, needed as RequestFocusMsg.Source
	w, h    int
	focused bool
}

func (p *controlPane) Init() tea.Cmd { return nil }

func (p *controlPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.id, p.w, p.h = msg.ID, msg.Width, msg.Height
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	case tea.KeyPressMsg:
		switch msg.String() {
		case "1", "2", "3":
			text := "command " + msg.String()
			return p, func() tea.Msg {
				return layout.SendMsg{Target: "editor", Msg: commandMsg{text: text}}
			}
		case "enter":
			return p, func() tea.Msg {
				return layout.RequestFocusMsg{Source: p.id, Target: "editor"}
			}
		}
	}
	return p, nil
}

func (p *controlPane) View() string {
	content := "control\n\n1/2/3: send a command\nenter: focus editor"
	inner := lipgloss.NewStyle().Width(p.w - 2).Height(p.h - 2).Render(content)
	return layout.Bordered(p.focused, p.w, p.h, inner)
}

type editorPane struct {
	w, h    int
	focused bool
	last    string
}

func (p *editorPane) Init() tea.Cmd { return nil }

func (p *editorPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.w, p.h = msg.Width, msg.Height
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	case commandMsg:
		p.last = msg.text
	}
	return p, nil
}

func (p *editorPane) View() string {
	last := p.last
	if last == "" {
		last = "(nothing yet)"
	}
	content := fmt.Sprintf("editor\n\nlast command received:\n%s", last)
	inner := lipgloss.NewStyle().Width(p.w - 2).Height(p.h - 2).Render(content)
	return layout.Bordered(p.focused, p.w, p.h, inner)
}

// model is the actual top-level tea.Model: layout itself reserves no quit
// key (that's an app policy, not layout's to make), so the host wraps it
// and handles ctrl+c/q itself, same as any other custom component in this
// repo (see examples/tabs).
type model struct {
	layout layout.Model
}

func (m model) Init() tea.Cmd { return m.layout.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	updated, cmd := m.layout.Update(msg)
	m.layout = updated.(layout.Model)
	return m, cmd
}

func (m model) View() tea.View {
	view := tea.NewView(m.layout.View())
	view.AltScreen = true
	return view
}

func main() {
	root := layout.HSplit(0.35,
		layout.Leaf("control", &controlPane{}),
		layout.Leaf("editor", &editorPane{}),
	)
	m := model{layout: layout.New(root, layout.AsRoot())}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
