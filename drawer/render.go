package drawer

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/style"
)

func (m Model) Render(background string) string {
	if len(m.drawers) == 0 {
		return background
	}

	result := background
	for _, d := range m.drawers {
		w, h := lipgloss.Width(result), lipgloss.Height(result)
		if w <= 0 || h <= 0 {
			return result
		}
		result = m.renderOne(d, result, w, h)
	}
	return result
}

func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

func (m Model) renderOne(d Drawer, background string, w, h int) string {
	s := m.styles
	if d.Style != nil {
		s = *d.Style
	}

	box := m.renderBox(d, s, w, h)
	bw := lipgloss.Width(box)
	x := 0
	if d.Side == Right {
		x = max(w-bw, 0)
	}

	compositor := lipgloss.NewCompositor(
		lipgloss.NewLayer(dim(background, s.DimColor)),
		lipgloss.NewLayer(box).X(x).Y(0).Z(1),
	)
	return compositor.Render()
}

func dim(s string, c color.Color) string {
	return lipgloss.NewStyle().Foreground(c).Render(ansi.Strip(s))
}

func (m Model) renderBox(d Drawer, s Styles, bgW, bgH int) string {
	widthCap := m.maxWidth
	if d.Width > 0 {
		widthCap = d.Width
	}
	maxW := effectiveMax(widthCap, bgW)

	body := contentView(d)
	inner := contentWidth(d, body, maxW)
	content := s.Content.Width(inner).Render(body)

	boxWidth := inner + 4

	return style.RenderWithTitle(s.Border, s.Title.Render(d.Title), content, boxWidth, bgH)
}

func contentView(d Drawer) string {
	if d.Content == nil {
		return ""
	}
	return d.Content.View().Content
}

func contentWidth(d Drawer, body string, maxWidth int) int {
	capped := max(maxWidth-4, 1)
	if d.Width > 0 {
		return capped
	}
	natural := max(naturalWidth(body), lipgloss.Width(d.Title), 1)
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
