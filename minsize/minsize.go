package minsize

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type Model struct {
	MinWidth  int
	MinHeight int
	Style     lipgloss.Style
}

type Option func(*Model)

func WithStyle(s lipgloss.Style) Option {
	return func(m *Model) { m.Style = s }
}

func New(minWidth, minHeight int, opts ...Option) Model {
	m := Model{
		MinWidth:  minWidth,
		MinHeight: minHeight,
		Style:     lipgloss.NewStyle().Foreground(style.S.Muted),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m Model) Fits(width, height int) bool {
	return width >= m.MinWidth && height >= m.MinHeight
}

func (m Model) View(width, height int) string {
	msg := m.Style.Render(m.message(width, height))
	if width < 1 || height < 1 {
		return msg
	}
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, msg)
}

func (m Model) message(width, height int) string {
	tooNarrow := width < m.MinWidth
	tooShort := height < m.MinHeight
	switch {
	case tooNarrow && tooShort:
		return fmt.Sprintf("minimum is %dx%d", m.MinWidth, m.MinHeight)
	case tooNarrow:
		return fmt.Sprintf("minimum width: %d", m.MinWidth)
	case tooShort:
		return fmt.Sprintf("minimum height: %d", m.MinHeight)
	default:
		return ""
	}
}
