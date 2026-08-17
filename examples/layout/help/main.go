// Command help demonstrates the dynamic help bar: each pane implements
// layout.HelpProvider (or doesn't) and the bottom bar always reflects
// whichever one is currently focused, with no wiring beyond that.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

// helpfulPane implements layout.HelpProvider with a binding or two of its
// own, on top of the usual Size/Focus/Blur bookkeeping.
type helpfulPane struct {
	id       string
	w, h     int
	focused  bool
	bindings []key.Binding
}

func newHelpfulPane(id string, bindings []key.Binding) *helpfulPane {
	return &helpfulPane{id: id, bindings: bindings}
}

func (p *helpfulPane) Init() tea.Cmd { return nil }

func (p *helpfulPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
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

func (p *helpfulPane) View() string {
	content := lipgloss.NewStyle().
		Width(p.w - 2).Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(p.id)
	return layout.Bordered(p.focused, p.w, p.h, content)
}

// HelpBindings implements layout.HelpProvider.
func (p *helpfulPane) HelpBindings() []key.Binding { return p.bindings }

// silentPane deliberately does NOT implement layout.HelpProvider, to show
// that the help bar just falls back to layout's own controls (ctrl+hjkl,
// ?) instead of disappearing or erroring when it's focused.
type silentPane struct {
	w, h    int
	focused bool
}

func (p *silentPane) Init() tea.Cmd { return nil }

func (p *silentPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
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

func (p *silentPane) View() string {
	content := lipgloss.NewStyle().
		Width(p.w - 2).Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render("no help bindings\n(watch the bar below)")
	return layout.Bordered(p.focused, p.w, p.h, content)
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
	writer := newHelpfulPane("writer", []key.Binding{
		key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		key.NewBinding(key.WithKeys("ctrl+z"), key.WithHelp("ctrl+z", "undo")),
	})
	browser := newHelpfulPane("browser", []key.Binding{
		key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	})

	root := layout.HSplit(0.34,
		layout.Leaf("writer", writer),
		layout.HSplit(0.5,
			layout.Leaf("browser", browser),
			layout.Leaf("scratch", &silentPane{}),
		),
	)
	m := model{layout: layout.New(root, layout.AsRoot())}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
