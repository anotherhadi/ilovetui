package sidebar

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
)

type NavItem struct {
	Icon string
	Name string
}

func (n NavItem) Title() string {
	if n.Icon == "" {
		return n.Name
	}
	return n.Icon + " " + n.Name
}

func (n NavItem) Description() string { return "" }
func (n NavItem) FilterValue() string { return n.Name }

type SelectMsg struct {
	Index int
	Item  NavItem
}

// BlurMsg is the sidebar asking to be given up: it goes out with the
// SelectMsg, on the grounds that picking an entry means you're done with the
// menu. Where focus lands instead is the host's call - the sidebar has no
// idea what else is on screen.
type BlurMsg struct{}

// blur is BlurMsg's command form.
func blur() tea.Msg { return BlurMsg{} }

type KeyMap struct {
	Select key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

type Model struct {
	KeyMap KeyMap

	list     list.Model
	items    []NavItem
	selected int
}

func New(items ...NavItem) Model {
	d := bubbles.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)

	l := bubbles.NewList(listItems(items), 0, 0)
	l.SetDelegate(d)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		KeyMap: DefaultKeyMap(),
		list:   l,
		items:  items,
	}
}

func listItems(items []NavItem) []list.Item {
	out := make([]list.Item, len(items))
	for i, item := range items {
		out[i] = item
	}
	return out
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if press, ok := msg.(tea.KeyPressMsg); ok && key.Matches(press, m.KeyMap.Select) {
		// Only the interactive path asks for focus to move on. Select on its
		// own doesn't, because the host also calls it to set the starting
		// entry - which must not steal focus from anything.
		return m, tea.Batch(m.Select(m.list.Index()), blur)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) Select(index int) tea.Cmd {
	if index < 0 || index >= len(m.items) {
		return nil
	}
	m.selected = index
	m.list.Select(index)

	item := m.items[index]
	return func() tea.Msg { return SelectMsg{Index: index, Item: item} }
}

func (m *Model) SetSize(width, height int) { m.list.SetSize(width, height) }

func (m Model) Selected() NavItem {
	if m.selected < 0 || m.selected >= len(m.items) {
		return NavItem{}
	}
	return m.items[m.selected]
}

func (m Model) SelectedIndex() int { return m.selected }

func (m Model) Cursor() int { return m.list.Index() }

func (m Model) View() string { return m.list.View() }

func (m Model) HelpBindings() []key.Binding {
	return []key.Binding{
		m.list.KeyMap.CursorUp,
		m.list.KeyMap.CursorDown,
		m.KeyMap.Select,
	}
}
