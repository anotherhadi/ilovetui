package modal

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/style"
)

// margin is the fixed gap, in cells, kept between a modal box and the edges
// of the background it's centered on.
const margin = 2

// Render draws every open modal (see Model.Update/Show) on top of background
// (already rendered, e.g. layout.Model.View() or any other component's
// View()) and returns the result. Each modal in the stack first flattens
// whatever came before it - background plus any earlier modal - to a single
// flat DimColor (see dim), then draws its own box centered on top, so
// nesting a second modal on top of a first dims the first one too. background
// is returned unchanged whenever there's nothing to draw.
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

// View is a convenience for a pane whose sole purpose is showing modals
// (e.g. a dedicated layout.Leaf): it draws the stack over a blank
// width x height area instead of an existing background.
func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

// renderOne dims background flat and draws mo's box centered on top of it.
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

// dim flattens s to a single flat color: every existing style (colors,
// bold, underline...) is stripped, then every character - including
// whitespace, so highlighted/selected backgrounds vanish too - is
// repainted in c. Applying a Foreground style to a multi-line string styles
// each line independently (see lipgloss.Style.Render), so this keeps s's
// line structure intact.
func dim(s string, c color.Color) string {
	return lipgloss.NewStyle().Foreground(c).Render(ansi.Strip(s))
}

// renderBox draws mo as a bordered, title-embedded box (style.RenderWithTitle),
// shrunk to fit its content, capped by the Model's configured max size and by
// whatever actually fits inside a bgW x bgH background.
func (m Model) renderBox(mo Modal, s Styles, bgW, bgH int) string {
	maxW := effectiveMax(m.maxWidth, bgW-2*margin)
	maxH := effectiveMax(m.maxHeight, bgH-2*margin)

	inner := contentWidth(mo, maxW)
	content := s.Content.Width(inner).Render(mo.Content)

	boxWidth := inner + 4 // border (2) + Padding(0, 1) (2)
	boxHeight := min(lipgloss.Height(content)+2, maxH)

	return style.RenderWithTitle(s.Border, s.Title.Render(mo.Title), content, boxWidth, boxHeight)
}

// contentWidth is the modal's inner (border/padding excluded) width: its
// natural size (long enough for the widest line of title/content), capped
// at maxWidth.
func contentWidth(mo Modal, maxWidth int) int {
	natural := max(naturalWidth(mo.Content), lipgloss.Width(mo.Title), 1)
	capped := max(maxWidth-4, 1)
	return min(natural, capped)
}

// naturalWidth is the width of content's widest line.
func naturalWidth(content string) int {
	w := 0
	for _, line := range strings.Split(content, "\n") {
		if lw := lipgloss.Width(line); lw > w {
			w = lw
		}
	}
	return w
}

// effectiveMax resolves the cap actually used along one axis: configured
// (0 = unlimited) narrowed down to fits, whatever actually fits the
// background - a modal can never overflow past the edge of the background,
// or the terminal, when the background is a full-screen View(), regardless
// of how WithMaxWidth/WithMaxHeight was set.
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
