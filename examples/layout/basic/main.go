// Command basic demonstrates the minimum layout needs: a 2x2 grid of plain
// panes, no custom border, ctrl+hjkl moving focus between them out of the
// box.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

// textPane is the simplest possible Pane: it just reports its own id, size
// and focus state. It still renders to exactly the width/height it was
// told (via lipgloss's own Width/Height), which is the one rule every Pane
// has to follow - layout composes View() output side by side and never
// pads it itself.
type textPane struct {
	id      string
	w, h    int
	focused bool
}

func newTextPane(id string) *textPane { return &textPane{id: id} }

func (p *textPane) Init() tea.Cmd { return nil }

func (p *textPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.w, p.h = msg.Width, msg.Height
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	}
	return p, nil
}

func (p *textPane) View() string {
	state := "blurred"
	if p.focused {
		state = "focused"
	}
	content := fmt.Sprintf("%s\n(%s, %dx%d)", p.id, state, p.w, p.h)
	return lipgloss.NewStyle().
		Width(p.w).
		Height(p.h).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
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
	root := layout.HSplit(0.5,
		layout.VSplit(0.5,
			layout.Leaf("top-left", newTextPane("top-left")),
			layout.Leaf("bottom-left", newTextPane("bottom-left")),
		),
		layout.VSplit(0.5,
			layout.Leaf("top-right", newTextPane("top-right")),
			layout.Leaf("bottom-right", newTextPane("bottom-right")),
		),
	)
	m := model{layout: layout.New(root, layout.AsRoot())}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
