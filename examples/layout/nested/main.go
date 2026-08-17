// Command nested demonstrates composability: "workspace" is itself a full
// layout.Model (its own two-pane split) embedded as an ordinary Leaf inside
// the outer tree. ctrl+hjkl bubbles in and out of it transparently, and
// only the outer Model shows a help bar - the inner one is built without
// layout.AsRoot(), see newWorkspace.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

type pane struct {
	id      string
	w, h    int
	focused bool
}

func newPane(id string) *pane { return &pane{id: id} }

func (p *pane) Init() tea.Cmd { return nil }

func (p *pane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
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

func (p *pane) View() string {
	content := lipgloss.NewStyle().
		Width(p.w - 2).Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(p.id)
	return layout.Bordered(p.focused, p.w, p.h, content)
}

// HelpBindings implements layout.HelpProvider, so every pane - at the top
// level or nested three levels deep, doesn't matter - shows up correctly
// in the single, outer help bar.
func (p *pane) HelpBindings() []key.Binding {
	return []key.Binding{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "use "+p.id))}
}

// newWorkspace builds the nested layout.Model: note the lack of
// layout.AsRoot() here - only the outermost Model (see main) should render
// a help bar, or focused help would show up twice.
func newWorkspace() layout.Model {
	root := layout.VSplit(0.7,
		layout.Leaf("editor", newPane("editor")),
		layout.Leaf("terminal", newPane("terminal")),
	)
	return layout.New(root)
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
	root := layout.HSplit(0.25,
		layout.Leaf("sidebar", newPane("sidebar")).WithMinimum(20),
		layout.Leaf("workspace", newWorkspace()),
	)
	m := model{layout: layout.New(root, layout.AsRoot())}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
