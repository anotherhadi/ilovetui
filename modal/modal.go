// Package modal renders a centered popup box on top of a dimmed background,
// triggered from anywhere in a bubbletea program via an exported tea.Msg
// (see ShowMsg/Show) rather than a direct reference to the Model that ends
// up rendering it.
//
// It composites over an already-rendered string (see Model.Render), so it
// makes no assumption at all about how the host builds that string: the same
// Model works whatever the host uses to lay out its main content (see
// Model.Render and Model.View).
//
// A modal's content is a model, not a string: it is updated while it's on
// top of the stack, so it can hold anything a pane can - a form, a list, a
// confirmation reporting its answer back with its own tea.Msg. See Show and
// Text.
package modal

import tea "charm.land/bubbletea/v2"

// Model holds the currently open modals (a stack: the most recently shown
// is drawn on top, everything beneath it - the background and any earlier
// modal - dimmed, see Render) and the rendering config (max size, styles)
// they share. Build one with New.
type Model struct {
	modals    []Modal
	maxWidth  int
	maxHeight int
	styles    Styles
}

// Option configures a Model at construction. See WithMaxWidth,
// WithMaxHeight, WithStyles.
type Option func(*Model)

// WithMaxWidth caps how wide a modal box can grow before its content wraps.
// A modal narrower than this shrinks to fit its content instead of padding
// out to the cap. 0 (also the zero-value Model's default without New) means
// only the background's own size caps it.
func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

// WithMaxHeight caps how tall a modal box can grow before its content is
// truncated. 0 means only the background's own size caps it.
func WithMaxHeight(h int) Option {
	return func(m *Model) { m.maxHeight = h }
}

// WithStyles overrides the default styles (see DefaultStyles).
func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

// New builds a Model. Defaults: a 60x20 max size, DefaultStyles.
func New(opts ...Option) Model {
	m := Model{
		maxWidth:  60,
		maxHeight: 20,
		styles:    DefaultStyles(),
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
		return m.show(msg.Modal)
	case DismissMsg:
		return m.pop(), nil
	}
	return m.updateTop(msg)
}

// updateTop forwards msg to the topmost modal's content - the only one the
// user can interact with, everything beneath it being dimmed (see Render). A
// modal deeper in the stack is frozen until the ones above it close.
//
// This is what lets modal content be a real model: it gets the key presses,
// the ticks and the results of its own commands, and can report back to the
// rest of the program with a tea.Msg of its own.
func (m Model) updateTop(msg tea.Msg) (Model, tea.Cmd) {
	i := len(m.modals) - 1
	if i < 0 || m.modals[i].Content == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.modals[i].Content, cmd = m.modals[i].Content.Update(msg)
	return m, cmd
}

// Open reports whether at least one modal is currently shown - handy for a
// host that wants to route key presses to the modal (e.g. esc to dismiss,
// enter to confirm) instead of its normal UI while one is open.
func (m Model) Open() bool { return len(m.modals) > 0 }

// show pushes a modal on top of the stack and returns its content's Init - a
// modal's body starts the same way any other model does.
func (m Model) show(mo Modal) (Model, tea.Cmd) {
	m.modals = append(m.modals, mo)
	return m, initContent(mo)
}

// initContent is mo's content's Init, or nil for a modal without content.
func initContent(mo Modal) tea.Cmd {
	if mo.Content == nil {
		return nil
	}
	return mo.Content.Init()
}

// pop closes the topmost modal (see Close).
func (m Model) pop() Model {
	if len(m.modals) == 0 {
		return m
	}
	m.modals = m.modals[:len(m.modals)-1]
	return m
}
