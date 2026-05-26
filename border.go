package ilovetui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// ContentHeight returns the usable inner height for a bordered panel of totalH rows.
func ContentHeight(totalH int) int {
	h := totalH - 2
	if h < 0 {
		return 0
	}
	return h
}

// RenderWithTitle renders a bordered box with a title embedded in the top border.
// title may contain ANSI color codes. width and height are the total outer dimensions.
//
// Example:
//
//	box := ilovetui.RenderWithTitle(theme.Styles.PanelFocused, "Header", content, w, h)
func RenderWithTitle(border lipgloss.Style, title, content string, width, height int) string {
	boxH := height - 1
	if contentH := boxH - 1; contentH > 0 {
		lines := strings.Split(content, "\n")
		if len(lines) > contentH {
			content = strings.Join(lines[:contentH], "\n")
		}
	}
	box := border.BorderTop(false).Width(width).Height(boxH).Render(content)

	boxWidth := lipgloss.Width(strings.SplitN(box, "\n", 2)[0])
	titleW := lipgloss.Width(title)
	fillW := boxWidth - titleW - 4 // 4 = "╭ " + " " + "╮"
	if fillW < 0 {
		fillW = 0
	}
	bc := lipgloss.NewStyle().Foreground(border.GetBorderTopForeground())
	topLine := bc.Render("╭ ") + bc.Render(title) + bc.Render(" "+strings.Repeat("─", fillW)+"╮")

	return lipgloss.JoinVertical(lipgloss.Left, topLine, box)
}
