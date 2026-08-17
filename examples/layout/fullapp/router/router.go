// Package router owns every page's own layout.Model for the lifetime of
// the whole app, so switching pages never discards their state - only
// which one is currently active/rendered changes. It slots into the root
// tree's "router" leaf exactly like a single nested layout.Model would
// (see examples/layout/nested), because Model implements layout.Navigable
// itself on top of layout.Pane, delegating everything to whichever page is
// active: ctrl+hjkl, the help bar, and SendMsg/RequestFocusMsg addressed to
// ids inside that page all keep working transparently.
package router

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/first"
	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/second"
	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/third"
	"github.com/anotherhadi/ilovetui/layout"
)

// SelectMsg asks Model to switch its active page, by id ("first", "second"
// or "third"). Sent by the sidebar via layout.SendMsg{Target: "router"} -
// it never holds a reference to Model either.
type SelectMsg struct{ Page string }

type Model struct {
	active string
	pages  map[string]layout.Model
}

func New() *Model {
	return &Model{
		active: "first",
		pages: map[string]layout.Model{
			"first":  first.NewWorkspace(),
			"second": second.NewWorkspace(),
			"third":  third.NewWorkspace(),
		},
	}
}

func (m *Model) current() layout.Model { return m.pages[m.active] }

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, p := range m.pages {
		if cmd := p.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	if sel, ok := msg.(SelectMsg); ok {
		if _, exists := m.pages[sel.Page]; exists {
			m.active = sel.Page
		}
		return m, nil
	}

	if size, ok := msg.(layout.SizeMsg); ok {
		// Every page needs to stay correctly sized even while hidden - the
		// outer tree only re-sends SizeMsg when the "router" leaf's own
		// Rect changes, never on a plain page switch, so a page that was
		// inactive during a resize must still learn about it here, or it
		// renders at a stale size whenever it becomes active again.
		var cmds []tea.Cmd
		for id, p := range m.pages {
			updated, cmd := p.Update(size)
			m.pages[id] = updated.(layout.Model)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	}

	updated, cmd := m.current().Update(msg)
	m.pages[m.active] = updated.(layout.Model)
	return m, cmd
}

func (m *Model) View() string {
	return m.current().View()
}

// Leaves, MoveFocus, Route, Focus and FocusedHelp implement
// layout.Navigable, delegating to whichever page is currently active.

func (m *Model) Leaves() []layout.LeafRect {
	return m.current().Leaves()
}

func (m *Model) MoveFocus(dir layout.FocusDirection) bool {
	return m.current().MoveFocus(dir)
}

func (m *Model) Route(target string, msg tea.Msg) (bool, tea.Cmd) {
	return m.current().Route(target, msg)
}

func (m *Model) Focus(id string) (bool, tea.Cmd) {
	return m.current().Focus(id)
}

func (m *Model) FocusedHelp() []key.Binding {
	return m.current().FocusedHelp()
}
