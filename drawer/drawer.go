// Package drawer renders a full-height panel flush against the left or
// right edge of an already-rendered background, on top of it dimmed -
// the sidebar/drawer equivalent of github.com/anotherhadi/ilovetui/modal,
// which this package otherwise mirrors closely: same stack of panels
// triggered from anywhere via an exported tea.Msg (see ShowMsg/Show)
// rather than a direct reference to the Model that ends up rendering it,
// same composite-over-an-already-rendered-string Render, same absence of
// any assumption about how the host builds that string.
//
// A drawer's content is a model, not a string: it is updated while it's on
// top of the stack, so it can hold anything a pane can - a file list, a
// filter form, a picker reporting its choice back with its own tea.Msg. See
// Show and Text.
package drawer

import tea "charm.land/bubbletea/v2"

// Model holds the currently open drawers (a stack: the most recently shown
// is drawn on top, everything beneath it - the background and any earlier
// drawer - dimmed, see Render) and the rendering config (max width, styles)
// they share. Build one with New.
type Model struct {
	drawers  []Drawer
	nextID   int
	maxWidth int
	styles   Styles
}

// Option configures a Model at construction. See WithMaxWidth, WithStyles.
type Option func(*Model)

// WithMaxWidth caps a drawer's total width (border and padding included). A
// drawer shown without WithWidth shrinks to fit its content instead of
// padding out to the cap; a drawer shown with WithWidth uses that width
// instead, still capped by this. 0 means only the background's own width
// caps it.
func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

// WithStyles overrides the default styles (see DefaultStyles).
func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

// New builds a Model. Defaults: a 30-column max width, DefaultStyles.
func New(opts ...Option) Model {
	m := Model{
		maxWidth: 30,
		styles:   DefaultStyles(),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowMsg:
		return m.show(msg.Drawer)
	case DismissMsg:
		return m.pop(), nil
	}
	return m.updateTop(msg)
}

// updateTop forwards msg to the topmost drawer's content - the only one the
// user can interact with, everything beneath it being dimmed (see Render). A
// drawer deeper in the stack is frozen until the ones above it close.
//
// This is what lets drawer content be a real model: it gets the key presses,
// the ticks and the results of its own commands, and can report back to the
// rest of the program with a tea.Msg of its own.
func (m Model) updateTop(msg tea.Msg) (Model, tea.Cmd) {
	i := len(m.drawers) - 1
	if i < 0 || m.drawers[i].Content == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.drawers[i].Content, cmd = m.drawers[i].Content.Update(msg)
	return m, cmd
}

// Open reports whether at least one drawer is currently shown - handy for a
// host that wants to route key presses to the drawer (e.g. esc to dismiss)
// instead of its normal UI while one is open.
func (m Model) Open() bool { return len(m.drawers) > 0 }

// show pushes a drawer on top of the stack and returns its content's Init - a
// drawer's body starts the same way any other model does.
func (m Model) show(d Drawer) (Model, tea.Cmd) {
	m.drawers = append(m.drawers, d)
	return m, initContent(d)
}

// initContent is d's content's Init, or nil for a drawer without content.
func initContent(d Drawer) tea.Cmd {
	if d.Content == nil {
		return nil
	}
	return d.Content.Init()
}

// pop closes the topmost drawer (see Close).
func (m Model) pop() Model {
	if len(m.drawers) == 0 {
		return m
	}
	m.drawers = m.drawers[:len(m.drawers)-1]
	return m
}
