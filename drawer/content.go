package drawer

import tea "charm.land/bubbletea/v2"

type Side int

const (
	Left Side = iota
	Right
)

type Drawer struct {
	Title   string
	Content tea.Model
	Side    Side
	Width   int
	Style   *Styles
}

type DrawerOption func(*Drawer)

func WithSide(s Side) DrawerOption {
	return func(d *Drawer) { d.Side = s }
}

func WithWidth(w int) DrawerOption {
	return func(d *Drawer) { d.Width = w }
}

func WithDrawerStyle(s Styles) DrawerOption {
	return func(d *Drawer) { d.Style = &s }
}

func newDrawer(title string, content tea.Model, opts ...DrawerOption) Drawer {
	d := Drawer{Title: title, Content: content}
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

type ShowMsg struct{ Drawer Drawer }

func Show(title string, content tea.Model, opts ...DrawerOption) tea.Cmd {
	d := newDrawer(title, content, opts...)
	return func() tea.Msg { return ShowMsg{Drawer: d} }
}

type DismissMsg struct{}

func Close() tea.Cmd {
	return func() tea.Msg { return DismissMsg{} }
}

type text string

func (t text) Init() tea.Cmd                       { return nil }
func (t text) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t text) View() tea.View                      { return tea.NewView(string(t)) }

func Text(s string) tea.Model { return text(s) }
