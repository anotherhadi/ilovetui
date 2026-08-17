package sidebar

import (
	"charm.land/bubbles/v2/key"
	bubbleslist "charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/router"
	"github.com/anotherhadi/ilovetui/layout"
)

// page is a sidebar entry: just an id (matching one of router.Model's
// own page keys) and a label. Model has no idea what a page actually is or
// how it's built - that's router's business entirely, this only needs
// to name it.
type page struct {
	id    string
	title string
}

func (p page) Title() string       { return p.title }
func (p page) Description() string { return "" }
func (p page) FilterValue() string { return p.title }

type Model struct {
	id      string // learned from SizeMsg.ID, needed as RequestFocusMsg.Source
	list    bubbleslist.Model
	w, h    int
	focused bool
}

func New() *Model {
	items := []bubbleslist.Item{
		page{id: "first", title: "First"},
		page{id: "second", title: "Second"},
		page{id: "third", title: "Third"},
	}
	list := bubbles.NewList(items, 0, 0)
	list.Title = "Sidebar"
	// list's own built-in help footer is redundant - Model already feeds
	// layout's single centralized help bar via HelpBindings below - and at
	// a narrow sidebar width it word-wraps onto multiple lines, which
	// list.SetSize doesn't account for: list.View() ends up taller than
	// the height it was given, breaking layout's "render exactly h" rule
	// and pushing the border past the bottom of the pane.
	list.SetShowHelp(false)
	return &Model{list: list}
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.id, m.w, m.h = msg.ID, msg.Width, msg.Height
		m.list.SetSize(m.w-2, m.h-2)
	case layout.FocusMsg:
		m.focused = true
	case layout.BlurMsg:
		m.focused = false
	}

	if !m.focused {
		return m, nil
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
		if selected, ok := m.list.SelectedItem().(page); ok {
			// tea.Sequence, not tea.Batch: the page switch must land
			// before the focus jump, or RequestFocusMsg could cascade
			// into router while it's still showing the previous page
			// (Batch runs both concurrently with no ordering guarantee).
			return m, tea.Sequence(
				func() tea.Msg {
					return layout.SendMsg{Target: "router", Msg: router.SelectMsg{Page: selected.id}}
				},
				func() tea.Msg {
					return layout.RequestFocusMsg{Source: m.id, Target: "router"}
				},
			)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return layout.Bordered(m.focused, m.w, m.h, m.list.View())
}

func (m *Model) HelpBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open page")),
		key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "go down")),
		key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "go up")),
	}
}
