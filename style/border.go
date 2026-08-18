package style

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var borderTypes = map[string]lipgloss.Border{
	"rounded": lipgloss.RoundedBorder(),
	"normal":  lipgloss.NormalBorder(),
	"thick":   lipgloss.ThickBorder(),
	"double":  lipgloss.DoubleBorder(),
	"hidden":  lipgloss.HiddenBorder(),
	"ascii":   lipgloss.ASCIIBorder(),
}

func resolveBorderType(name string) lipgloss.Border {
	if b, ok := borderTypes[strings.ToLower(strings.TrimSpace(name))]; ok {
		return b
	}
	return lipgloss.RoundedBorder()
}

func ContentHeight(totalH int) int {
	h := totalH - 2
	if h < 0 {
		return 0
	}
	return h
}

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

	b, _, _, _, _ := border.GetBorder()
	topLeft, top, topRight := b.TopLeft, b.Top, b.TopRight

	fillW := boxWidth - titleW - lipgloss.Width(topLeft) - lipgloss.Width(topRight) - 2
	if fillW < 0 {
		fillW = 0
	}
	bc := lipgloss.NewStyle().Foreground(border.GetBorderTopForeground())
	topLine := bc.Render(topLeft+" ") + bc.Render(title) + bc.Render(" "+strings.Repeat(top, fillW)+topRight)

	return lipgloss.JoinVertical(lipgloss.Left, topLine, box)
}
