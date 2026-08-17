package bubbles

import (
	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// NewHelp returns a help.Model styled with the active theme.
func NewHelp() help.Model {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(style.S.Primary)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(style.S.Muted)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(style.S.Subtle)
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(style.S.Primary)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(style.S.Muted)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(style.S.Subtle)
	h.Styles.Ellipsis = lipgloss.NewStyle().Foreground(style.S.Subtle)
	return h
}
