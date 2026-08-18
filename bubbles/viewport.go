package bubbles

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

const ViewportGutterWidth = 2

func NewViewport() viewport.Model {
	vp := viewport.New()
	vp.MouseWheelEnabled = false
	return vp
}

func ViewportView(vp *viewport.Model) string {
	height := vp.Height()
	total := vp.TotalLineCount()
	blank := lipgloss.NewStyle().Width(ViewportGutterWidth).Render("")

	if height <= 0 || total <= height {
		vp.LeftGutterFunc = func(viewport.GutterContext) string { return blank }
		return vp.View()
	}

	yOffset := vp.YOffset()
	thumbSize := max(1, height*height/total)
	thumbStart := 0
	if maxOffset := total - height; maxOffset > 0 {
		thumbStart = yOffset * (height - thumbSize) / maxOffset
	}

	trackStyle := lipgloss.NewStyle().Foreground(style.S.Subtle)
	thumbStyle := lipgloss.NewStyle().Foreground(style.S.Primary)

	vp.LeftGutterFunc = func(ctx viewport.GutterContext) string {
		pos := ctx.Index - yOffset
		if pos >= thumbStart && pos < thumbStart+thumbSize {
			return thumbStyle.Render("█") + " "
		}
		return trackStyle.Render("│") + " "
	}

	return vp.View()
}
