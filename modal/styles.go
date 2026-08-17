package modal

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// Styles is the set of lipgloss styles, plus the dim color, used to render a
// modal and the background behind it. Border carries the box's border
// (shape + color, no size), Title and Content color the two pieces of text
// drawn inside it, DimColor is the single flat color every character of the
// background gets overwritten with while the modal is open.
type Styles struct {
	Border   lipgloss.Style
	Title    lipgloss.Style
	Content  lipgloss.Style
	DimColor color.Color
}

// DefaultStyles builds a Styles from style.S: the box borrows
// PanelFocused's border (the modal is what has focus while it's open).
// DimColor reuses Subtle - the base16 "comments/invisibles" role, already
// used across this repo for de-emphasized text (borders, placeholders,
// separators, see bubbles/*.go) - darker than Muted, which reads too bright
// once it's covering an entire screen instead of a single blurred field.
func DefaultStyles() Styles {
	return Styles{
		Border:   style.S.PanelFocused.Padding(0, 1),
		Title:    lipgloss.NewStyle().Bold(true).Foreground(style.S.Primary),
		Content:  lipgloss.NewStyle().Foreground(style.S.Text),
		DimColor: style.S.Subtle,
	}
}
