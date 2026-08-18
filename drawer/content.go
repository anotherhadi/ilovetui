package drawer

import tea "charm.land/bubbletea/v2"

// Side is which edge of the background a Drawer is anchored to.
type Side int

const (
	// Left anchors the drawer to the left edge (the zero value, so a Drawer
	// built without WithSide opens on the left).
	Left Side = iota
	// Right anchors the drawer to the right edge.
	Right
)

// Drawer is one sidebar panel on the stack.
type Drawer struct {
	Title string
	// Content is the drawer's body: a full model, updated and rendered by
	// the drawer.Model while it's on top of the stack. Wrap a plain string
	// with Text for a drawer with nothing to interact with.
	Content tea.Model
	Side    Side
	// Width, if non-zero, fixes the drawer's total width (border and
	// padding included, same unit as the Model's WithMaxWidth) instead of
	// shrinking to fit what Content draws, and Title. Still capped by whatever actually
	// fits the background.
	Width int
	// Style, if non-nil, overrides the Model's default Styles for this
	// drawer alone.
	Style *Styles
}

// DrawerOption configures a Drawer built by Show.
type DrawerOption func(*Drawer)

// WithSide anchors the drawer to the given Side (default Left).
func WithSide(s Side) DrawerOption {
	return func(d *Drawer) { d.Side = s }
}

// WithWidth fixes the drawer's total width instead of shrinking to fit its
// content, still capped by whatever actually fits the background.
func WithWidth(w int) DrawerOption {
	return func(d *Drawer) { d.Width = w }
}

// WithDrawerStyle overrides the Model's default Styles for this drawer
// alone, for a one-off custom look instead of the shared theme.
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

// ShowMsg tells a drawer.Model to display Drawer, pushing it on top of the
// stack. Any component in the same bubbletea program can trigger one via
// Show, without holding a reference to the drawer.Model that will actually
// render it - that Model just needs to see every tea.Msg the program
// produces, same as any other child model.
type ShowMsg struct{ Drawer Drawer }

// Show returns a tea.Cmd that opens a new drawer on top of the stack,
// anchored to the left edge by default. content is a model, so a drawer can
// hold anything a pane can - a file list, a filter form, a picker that
// reports its choice with its own tea.Msg:
//
//	return m, drawer.Show("Files", newFileList(dir), drawer.WithSide(drawer.Right))
//	return m, drawer.Show("Nav", drawer.Text("Home\nProjects"))
//
// The content's Init runs when the drawer opens, and it receives every
// message while it's the topmost drawer (see Model.Update).
func Show(title string, content tea.Model, opts ...DrawerOption) tea.Cmd {
	d := newDrawer(title, content, opts...)
	return func() tea.Msg { return ShowMsg{Drawer: d} }
}

// DismissMsg closes the topmost drawer. The stack is a plain LIFO: a drawer
// is closed by being on top, never by being named.
type DismissMsg struct{}

// Close returns a tea.Cmd that closes the topmost drawer - which is also the
// only one that can act (see Model.Update), so a content model closes itself
// by returning it:
//
//	return c, drawer.Close()
func Close() tea.Cmd {
	return func() tea.Msg { return DismissMsg{} }
}

// text is a model wrapping a fixed string: a drawer body with nothing to
// update.
type text string

func (t text) Init() tea.Cmd                       { return nil }
func (t text) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t text) View() tea.View                      { return tea.NewView(string(t)) }

// Text wraps a plain string as drawer content, for the common drawer that has
// nothing to interact with:
//
//	drawer.Show("Nav", drawer.Text("Home\nProjects\nSettings"))
func Text(s string) tea.Model { return text(s) }
