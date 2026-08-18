package modal

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/style"
)

const margin = 2

func (m Model) Render(background string) string {
	if len(m.modals) == 0 {
		return background
	}

	result := background
	for _, mo := range m.modals {
		w, h := lipgloss.Width(result), lipgloss.Height(result)
		if w <= 0 || h <= 0 {
			return result
		}
		result = m.renderOne(mo, result, w, h)
	}
	return result
}

func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

func (m Model) renderOne(mo Modal, background string, w, h int) string {
	s := m.styles
	if mo.Style != nil {
		s = *mo.Style
	}

	box := m.renderBox(mo, s, w, h)
	bw, bh := lipgloss.Width(box), lipgloss.Height(box)
	x, y := max((w-bw)/2, 0), max((h-bh)/2, 0)

	compositor := lipgloss.NewCompositor(
		lipgloss.NewLayer(dim(background, s.DimColor)),
		lipgloss.NewLayer(box).X(x).Y(y).Z(1),
	)
	return compositor.Render()
}

func dim(s string, c color.Color) string {
	return lipgloss.NewStyle().Foreground(c).Render(ansi.Strip(s))
}

func (m Model) renderBox(mo Modal, s Styles, bgW, bgH int) string {
	maxW := effectiveMax(m.maxWidth, bgW-2*margin)
	maxH := effectiveMax(m.maxHeight, bgH-2*margin)

	body := contentView(mo)
	inner := contentWidth(body, mo.Title, maxW)
	content := s.Content.Width(inner).Render(body)

	boxWidth := inner + 4
	boxHeight := min(lipgloss.Height(content)+2, maxH)

	return style.RenderWithTitle(s.Border, s.Title.Render(mo.Title), content, boxWidth, boxHeight)
}

func contentView(mo Modal) string {
	if mo.Content == nil {
		return ""
	}
	return mo.Content.View().Content
}

func contentWidth(body, title string, maxWidth int) int {
	natural := max(naturalWidth(body), lipgloss.Width(title), 1)
	capped := max(maxWidth-4, 1)
	return min(natural, capped)
}

func naturalWidth(content string) int {
	w := 0
	for _, line := range strings.Split(content, "\n") {
		if lw := lipgloss.Width(line); lw > w {
			w = lw
		}
	}
	return w
}

func effectiveMax(configured, fits int) int {
	if fits < 1 {
		fits = 1
	}
	if configured > 0 && configured < fits {
		return configured
	}
	return fits
}

func blank(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	line := strings.Repeat(" ", width)
	lines := make([]string, height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}
