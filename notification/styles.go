package notification

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// KindStyle is the set of lipgloss styles used to render one toast: Border
// carries the box's border (shape + color, no size), Title and Message color
// the two pieces of text drawn inside it. Building a value directly (rather
// than through a constructor) is the intended way to hand WithToastStyle a
// custom, per-toast look.
type KindStyle struct {
	Border  lipgloss.Style
	Title   lipgloss.Style
	Message lipgloss.Style
}

// Styles maps each Kind to the KindStyle used to render it. Build one with
// DefaultStyles and tweak individual fields, or construct one from scratch
// for a fully custom palette across all kinds.
type Styles struct {
	Info    KindStyle
	Success KindStyle
	Warning KindStyle
	Error   KindStyle
}

// DefaultStyles builds a Styles from style.S: Info uses the theme's primary
// accent (style.S has no dedicated "info" color, Primary already fills that
// neutral-accent role elsewhere in this repo), Success/Warning/Error use
// their matching style.S alias.
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

// forKind resolves the KindStyle to render t with: its own Style override if
// set, otherwise s's preset for t.Kind (falling back to Info for an
// out-of-range Kind).
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
