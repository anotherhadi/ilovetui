package modal

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type Styles struct {
	Border   lipgloss.Style
	Title    lipgloss.Style
	Content  lipgloss.Style
	DimColor color.Color
}

func DefaultStyles() Styles {
	return Styles{
		Border:   style.S.PanelFocused.Padding(0, 1),
		Title:    lipgloss.NewStyle().Bold(true).Foreground(style.S.Primary),
		Content:  lipgloss.NewStyle().Foreground(style.S.Text),
		DimColor: style.S.Subtle,
	}
}
