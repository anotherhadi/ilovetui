// Package notification renders toast-style notifications, triggered from
// anywhere in a bubbletea program via an exported tea.Msg (see ShowMsg/Show)
// rather than a direct reference to the Model that ends up rendering them.
//
// It composites over an already-rendered string (see Model.Render), so it
// has no dependency on github.com/anotherhadi/ilovetui/layout: the same
// Model works whether the host uses layout for its main content or not (see
// Model.Render and Model.View).
package notification

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Model holds the currently visible toasts and the rendering config
// (position, max width, per-Kind styles) they share. Build one with New.
type Model struct {
	toasts   []Toast
	nextID   int
	position Position
	maxWidth int
	styles   Styles
}

// Option configures a Model at construction. See WithPosition, WithMaxWidth,
// WithStyles.
type Option func(*Model)

// WithPosition sets which edge/corner the toast stack anchors to. TopRight
// by default.
func WithPosition(p Position) Option {
	return func(m *Model) { m.position = p }
}

// WithMaxWidth caps how wide a toast box can grow before its message
// wraps. A toast narrower than this shrinks to fit its content instead of
// padding out to the cap. 0 (also the zero-value Model's default without
// New) means unlimited.
func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

// WithStyles overrides the default per-Kind styles (see DefaultStyles).
func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

// New builds a Model. Defaults: TopRight, a 40-cell max width, DefaultStyles.
func New(opts ...Option) Model {
	m := Model{
		position: TopRight,
		maxWidth: 40,
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
		return m.show(msg.Toast)
	case DismissMsg:
		return m.remove(msg.ID), nil
	case expireMsg:
		return m.remove(msg.id), nil
	}
	return m, nil
}

// show adds or replaces (see WithID) a toast, and schedules its expiry via
// tea.Tick if it isn't sticky (Duration <= 0).
func (m Model) show(t Toast) (Model, tea.Cmd) {
	if t.ID == "" {
		t.ID = fmt.Sprintf("toast-%d", m.nextID)
		m.nextID++
	}

	replaced := false
	for i, existing := range m.toasts {
		if existing.ID == t.ID {
			m.toasts[i] = t
			replaced = true
			break
		}
	}
	if !replaced {
		m.toasts = append(m.toasts, t)
	}

	if t.Duration <= 0 {
		return m, nil
	}
	id := t.ID
	return m, tea.Tick(t.Duration, func(time.Time) tea.Msg {
		return expireMsg{id: id}
	})
}

func (m Model) remove(id string) Model {
	for i, t := range m.toasts {
		if t.ID == id {
			m.toasts = append(m.toasts[:i], m.toasts[i+1:]...)
			break
		}
	}
	return m
}
