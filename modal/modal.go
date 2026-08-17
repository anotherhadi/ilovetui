// Package modal renders a centered popup box on top of a dimmed background,
// triggered from anywhere in a bubbletea program via an exported tea.Msg
// (see ShowMsg/Show) rather than a direct reference to the Model that ends
// up rendering it.
//
// It composites over an already-rendered string (see Model.Render), so it
// has no dependency on github.com/anotherhadi/ilovetui/layout: the same
// Model works whether the host uses layout for its main content or not (see
// Model.Render and Model.View).
package modal

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// Model holds the currently open modals (a stack: the most recently shown
// is drawn on top, everything beneath it - the background and any earlier
// modal - dimmed, see Render) and the rendering config (max size, styles)
// they share. Build one with New.
type Model struct {
	modals    []Modal
	nextID    int
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
		return m.show(msg.Modal), nil
	case DismissMsg:
		return m.remove(msg.ID), nil
	}
	return m, nil
}

// Open reports whether at least one modal is currently shown - handy for a
// host that wants to route key presses to the modal (e.g. esc to dismiss,
// enter to confirm) instead of its normal UI while one is open.
func (m Model) Open() bool { return len(m.modals) > 0 }

// TopID returns the ID of the topmost (currently interactive) modal, or ""
// if none is open.
func (m Model) TopID() string {
	if len(m.modals) == 0 {
		return ""
	}
	return m.modals[len(m.modals)-1].ID
}

// show pushes or replaces (see WithID) a modal on top of the stack.
func (m Model) show(mo Modal) Model {
	if mo.ID == "" {
		mo.ID = fmt.Sprintf("modal-%d", m.nextID)
		m.nextID++
	}
	for i, existing := range m.modals {
		if existing.ID == mo.ID {
			m.modals[i] = mo
			return m
		}
	}
	m.modals = append(m.modals, mo)
	return m
}

// remove closes the modal identified by id, or the topmost one if id is
// empty (see Close).
func (m Model) remove(id string) Model {
	if len(m.modals) == 0 {
		return m
	}
	if id == "" {
		m.modals = m.modals[:len(m.modals)-1]
		return m
	}
	for i, mo := range m.modals {
		if mo.ID == id {
			m.modals = append(m.modals[:i], m.modals[i+1:]...)
			break
		}
	}
	return m
}
