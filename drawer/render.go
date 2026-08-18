package drawer

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/style"
)

// Render draws every open drawer (see Model.Update/Show) on top of
// background (already rendered, e.g. layout.Model.View() or any other
// component's View()) and returns the result. Each drawer in the stack
// first flattens whatever came before it - background plus any earlier
// drawer - to a single flat DimColor (see dim), then draws its own
// full-height box flush against its Side on top, so nesting a second
// drawer on top of a first dims the first one too. background is returned
// unchanged whenever there's nothing to draw.
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

// View is a convenience for a pane whose sole purpose is showing drawers
// (e.g. a dedicated layout.Leaf): it draws the stack over a blank
// width x height area instead of an existing background.
func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

// renderOne dims background flat and draws d's box, spanning its full
// height, flush against d.Side.
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

// dim flattens s to a single flat color: every existing style (colors,
// bold, underline...) is stripped, then every character - including
// whitespace, so highlighted/selected backgrounds vanish too - is
// repainted in c. Applying a Foreground style to a multi-line string styles
// each line independently (see lipgloss.Style.Render), so this keeps s's
// line structure intact.
func dim(s string, c color.Color) string {
	return lipgloss.NewStyle().Foreground(c).Render(ansi.Strip(s))
}

// renderBox draws d as a bordered, title-embedded box (style.RenderWithTitle)
// spanning the background's full height, capped by the Model's configured
// max width and by whatever actually fits inside a bgW-wide background.
func (m Model) renderBox(d Drawer, s Styles, bgW, bgH int) string {
	widthCap := m.maxWidth
	if d.Width > 0 {
		widthCap = d.Width
	}
	maxW := effectiveMax(widthCap, bgW)

	body := contentView(d)
	inner := contentWidth(d, body, maxW)
	content := s.Content.Width(inner).Render(body)

	boxWidth := inner + 4 // border (2) + Padding(0, 1) (2)

	return style.RenderWithTitle(s.Border, s.Title.Render(d.Title), content, boxWidth, bgH)
}

// contentView is the drawer body's rendered string, or "" for a drawer
// without content. The box shrinks to fit whatever the content model draws,
// so a content that wants a specific size sets it on itself - the drawer only
// ever sees the result.
func contentView(d Drawer) string {
	if d.Content == nil {
		return ""
	}
	return d.Content.View().Content
}

// contentWidth is the drawer's inner (border/padding excluded) width, given
// the resolved total-width cap maxWidth (see renderBox): maxWidth-4 if
// d.Width is set (a fixed total width), otherwise its natural size (long
// enough for the widest line of title/content), capped at maxWidth-4.
func contentWidth(d Drawer, body string, maxWidth int) int {
	capped := max(maxWidth-4, 1)
	if d.Width > 0 {
		return capped
	}
	natural := max(naturalWidth(body), lipgloss.Width(d.Title), 1)
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

// effectiveMax resolves the cap actually used along the width axis:
// configured (0 = unlimited) narrowed down to fits, whatever actually fits
// the background - a drawer can never overflow past the edge of the
// background, or the terminal, when the background is a full-screen View(),
// regardless of how WithMaxWidth/WithWidth was set.
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
