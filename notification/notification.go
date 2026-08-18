package notification

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	toasts   []Toast
	nextID   int
	position Position
	maxWidth int
	styles   Styles
}

type Option func(*Model)

func WithPosition(p Position) Option {
	return func(m *Model) { m.position = p }
}

func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

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
