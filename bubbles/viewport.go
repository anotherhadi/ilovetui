package bubbles

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// NewViewport returns a viewport.Model with mouse wheel disabled.
func NewViewport() viewport.Model {
	vp := viewport.New()
	vp.MouseWheelEnabled = false
	return vp
}

// ViewportView renders the viewport and appends a subtle scroll indicator
// on the last visible line when the user has not reached the bottom.
func ViewportView(vp *viewport.Model) string {
	v := vp.View()
	if vp.AtBottom() {
		return v
	}
	lines := strings.Split(v, "\n")
	if len(lines) == 0 {
		return v
	}
	arrow := lipgloss.NewStyle().Foreground(style.S.Subtle).Render("↓")
	arrowW := lipgloss.Width(arrow)
	inner := vp.Width() - 2*arrowW
	if inner < 0 {
		inner = 0
	}
	lines[len(lines)-1] = arrow + strings.Repeat(" ", inner) + arrow
	return strings.Join(lines, "\n")
}
