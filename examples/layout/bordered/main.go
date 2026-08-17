// Command bordered demonstrates that layout draws no chrome of its own:
// each pane here owns its border and picks its color from FocusMsg/BlurMsg,
// via the optional layout.Bordered helper.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

type borderedPane struct {
	id      string
	w, h    int
	focused bool
}

func newBorderedPane(id string) *borderedPane { return &borderedPane{id: id} }

func (p *borderedPane) Init() tea.Cmd { return nil }

func (p *borderedPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
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

func (p *borderedPane) View() string {
	inner := lipgloss.NewStyle().
		Width(p.w - 2).
		Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(p.id)
	// layout.Bordered already draws the border at exactly p.w x p.h, so the
	// content it wraps must already be sized to w-2 x h-2 (border eats one
	// cell on each side) - same rule as any other Pane, just one layer in.
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
	root := layout.HSplit(0.3,
		layout.Leaf("sidebar", newBorderedPane("sidebar")),
		layout.VSplit(0.6,
			layout.Leaf("main", newBorderedPane("main")),
			layout.Leaf("log", newBorderedPane("log")),
		),
	)
	m := model{layout: layout.New(root, layout.AsRoot())}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
