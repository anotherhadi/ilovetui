package style

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// borderTypes maps the `border:` config value to a lipgloss.Border. Only the
// symmetric, general-purpose border kinds are exposed here; MarkdownBorder,
// BlockBorder and the half-block variants are content-specific rather than a
// theming choice.
var borderTypes = map[string]lipgloss.Border{
	"rounded": lipgloss.RoundedBorder(),
	"normal":  lipgloss.NormalBorder(),
	"thick":   lipgloss.ThickBorder(),
	"double":  lipgloss.DoubleBorder(),
	"hidden":  lipgloss.HiddenBorder(),
	"ascii":   lipgloss.ASCIIBorder(),
}

// resolveBorderType maps a `border:` config value to a lipgloss.Border,
// falling back to RoundedBorder for an empty or unrecognized name.
func resolveBorderType(name string) lipgloss.Border {
	if b, ok := borderTypes[strings.ToLower(strings.TrimSpace(name))]; ok {
		return b
	}
	return lipgloss.RoundedBorder()
}

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
//	box := style.RenderWithTitle(style.S.PanelFocused, "Header", content, w, h)
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

	// Pull the corner/fill glyphs from the style's own border spec instead of
	// hardcoding rounded-border characters, so this respects style.S.BorderType.
	b, _, _, _, _ := border.GetBorder()
	topLeft, top, topRight := b.TopLeft, b.Top, b.TopRight

	fillW := boxWidth - titleW - lipgloss.Width(topLeft) - lipgloss.Width(topRight) - 2 // 2 = the spaces around the title
	if fillW < 0 {
		fillW = 0
	}
	bc := lipgloss.NewStyle().Foreground(border.GetBorderTopForeground())
	topLine := bc.Render(topLeft+" ") + bc.Render(title) + bc.Render(" "+strings.Repeat(top, fillW)+topRight)

	return lipgloss.JoinVertical(lipgloss.Left, topLine, box)
}
