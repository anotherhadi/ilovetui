package notification

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type KindStyle struct {
	Border  lipgloss.Style
	Title   lipgloss.Style
	Message lipgloss.Style
}

type Styles struct {
	Info    KindStyle
	Success KindStyle
	Warning KindStyle
	Error   KindStyle
}

func DefaultStyles() Styles {
	return Styles{
		Info:    kindStyle(style.S.Primary),
		Success: kindStyle(style.S.Success),
		Warning: kindStyle(style.S.Warning),
		Error:   kindStyle(style.S.Error),
	}
}

func kindStyle(c color.Color) KindStyle {
	return KindStyle{
		Border: lipgloss.NewStyle().
			Border(style.S.BorderType).
			BorderForeground(c).
			Padding(0, 1),
		Title:   lipgloss.NewStyle().Bold(true).Foreground(c),
		Message: lipgloss.NewStyle().Foreground(style.S.Text),
	}
}

func (s Styles) forKind(t Toast) KindStyle {
	if t.Style != nil {
		return *t.Style
	}
	switch t.Kind {
	case Success:
		return s.Success
	case Warning:
		return s.Warning
	case Error:
		return s.Error
	default:
		return s.Info
	}
}
