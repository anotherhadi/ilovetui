package bubbles

import (
	"charm.land/bubbles/v2/paginator"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

func NewPaginator() paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(style.S.Primary).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(style.S.Subtle).Render("•")
	return p
}
