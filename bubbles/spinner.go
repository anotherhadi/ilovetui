package bubbles

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// NewSpinner returns a spinner.Model styled with the active theme. Any opts
// are forwarded to spinner.New before the theme is applied.
func NewSpinner(opts ...spinner.Option) spinner.Model {
	s := spinner.New(opts...)
	s.Style = lipgloss.NewStyle().Foreground(style.S.Primary)
	return s
}
