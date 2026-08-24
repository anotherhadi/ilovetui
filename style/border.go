package style

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
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

type LayoutBorderMode string

const (
	LayoutBorderFull      LayoutBorderMode = "full"
	LayoutBorderSidebar   LayoutBorderMode = "sidebar"
	LayoutBorderSeparator LayoutBorderMode = "separator"
)

func resolveLayoutBorder(name string) LayoutBorderMode {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case string(LayoutBorderSidebar):
		return LayoutBorderSidebar
	case string(LayoutBorderSeparator):
		return LayoutBorderSeparator
	default:
		return LayoutBorderFull
	}
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

	b, _, _, _, _ := border.GetBorder()
	topLeft, top, topRight := b.TopLeft, b.Top, b.TopRight
	overhead := lipgloss.Width(topLeft) + lipgloss.Width(topRight) + 2

	maxTitleW := max(boxWidth-overhead, 0)
	if titleW := lipgloss.Width(title); titleW > maxTitleW {
		title = ansi.Truncate(title, maxTitleW, "")
	}
	titleW := lipgloss.Width(title)

	fillW := boxWidth - titleW - overhead
	if fillW < 0 {
		fillW = 0
	}
	bc := lipgloss.NewStyle().Foreground(border.GetBorderTopForeground())
	topLine := bc.Render(topLeft+" ") + bc.Render(title) + bc.Render(" "+strings.Repeat(top, fillW)+topRight)

	return lipgloss.JoinVertical(lipgloss.Left, topLine, box)
}

func RenderPlain(border lipgloss.Style, title, content string, width, height int) string {
	if titleW := lipgloss.Width(title); titleW > width {
		title = ansi.Truncate(title, width, "")
	}
	titleLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(border.GetBorderTopForeground()).
		Width(width).
		Render(title)

	contentH := max(height-1, 0)
	lines := strings.Split(content, "\n")
	if len(lines) > contentH {
		content = strings.Join(lines[:contentH], "\n")
	}
	box := lipgloss.NewStyle().Width(width).Height(contentH).Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, titleLine, box)
}

func VerticalDivider(border lipgloss.Style, height int) string {
	b, _, _, _, _ := border.GetBorder()
	ch := b.Left
	if ch == "" {
		ch = "│"
	}
	line := lipgloss.NewStyle().Foreground(border.GetBorderLeftForeground()).Render(ch)
	rows := make([]string, max(height, 0))
	for i := range rows {
		rows[i] = line
	}
	return strings.Join(rows, "\n")
}
