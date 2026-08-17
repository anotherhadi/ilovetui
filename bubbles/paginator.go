package bubbles

import (
	"charm.land/bubbles/v2/paginator"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// NewPaginator returns a dot-style paginator.Model styled with the active theme.
func NewPaginator() paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(style.S.Primary).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(style.S.Subtle).Render("•")
	return p
}
