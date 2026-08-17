package notification

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

// Render composites the current toasts on top of background (already
// rendered, e.g. layout.Model.View() or any other component's View()) and
// returns the result. background is returned unchanged whenever there's
// nothing to draw (no toasts, or a background with no measurable size).
//
// This is what makes notification work identically with or without layout:
// the host just wraps whatever it would otherwise return from its own
// View() with this call.
func (m Model) Render(background string) string {
	if len(m.toasts) == 0 {
		return background
	}
	w, h := lipgloss.Width(background), lipgloss.Height(background)
	if w <= 0 || h <= 0 {
		return background
	}

	stack := clipToHeight(m.renderStack(effectiveMaxWidth(m.maxWidth, w)), h-2*margin, m.position.anchoredTop())
	if stack == "" {
		return background
	}
	sw, sh := lipgloss.Width(stack), lipgloss.Height(stack)
	x, y := placement(m.position, w, h, sw, sh)

	// Canvas.Compose(layer) alone ignores the layer's X/Y and draws it across
	// the canvas's whole bounds, not just its own footprint - that's what
	// made the toast layer blank out the entire background instead of
	// floating over it. Compositor is what actually resolves each layer's
	// absolute bounds (background at 0,0, the stack at x,y) before drawing
	// each one only within its own area.
	compositor := lipgloss.NewCompositor(
		lipgloss.NewLayer(background),
		lipgloss.NewLayer(stack).X(x).Y(y).Z(1),
	)
	return compositor.Render()
}

// View is a convenience for a pane whose sole purpose is showing toasts (e.g.
// a dedicated layout.Leaf): it draws the stack over a blank width x height
// area instead of an existing background.
func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

// renderStack stacks every visible toast into one block, newest closest to
// the anchored edge (see Position.anchoredTop), separated by a blank line,
// and aligned so the edge the stack anchors to stays flush across toasts of
// different widths.
func (m Model) renderStack(maxWidth int) string {
	ordered := m.orderedToasts()
	parts := make([]string, 0, len(ordered)*2-1)
	for i, t := range ordered {
		if i > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, m.renderToast(t, maxWidth))
	}
	return lipgloss.JoinVertical(stackAlign(m.position), parts...)
}

// orderedToasts returns the toasts in the order they should stack, newest
// nearest the anchored edge: reversed (newest first) for a top anchor,
// insertion order (oldest first, newest last) for a bottom anchor.
func (m Model) orderedToasts() []Toast {
	if !m.position.anchoredTop() {
		return m.toasts
	}
	ordered := make([]Toast, len(m.toasts))
	for i, t := range m.toasts {
		ordered[len(m.toasts)-1-i] = t
	}
	return ordered
}

func stackAlign(pos Position) lipgloss.Position {
	switch pos {
	case TopLeft, BottomLeft:
		return lipgloss.Left
	case TopRight, BottomRight:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}

// renderToast draws a single toast as a box with its title embedded in the
// top border (style.RenderWithTitle), shrunk to fit its content up to
// maxWidth.
func (m Model) renderToast(t Toast, maxWidth int) string {
	k := m.styles.forKind(t)

	inner := contentWidth(t, maxWidth)
	message := k.Message.Width(inner).Render(t.Message)

	boxWidth := inner + 4 // border (2) + Padding(0, 1) (2)
	boxHeight := lipgloss.Height(message) + 2

	return style.RenderWithTitle(k.Border, k.Title.Render(t.Title), message, boxWidth, boxHeight)
}

// contentWidth is the toast's inner (border/padding excluded) width: its
// natural size (long enough for the wider of title/message on one line),
// capped at maxWidth if positive.
func contentWidth(t Toast, maxWidth int) int {
	natural := max(lipgloss.Width(t.Title), lipgloss.Width(t.Message), 1)
	if maxWidth <= 0 {
		return natural
	}
	capped := max(maxWidth-4, 1)
	return min(natural, capped)
}

// effectiveMaxWidth resolves the cap actually used to render a toast:
// configured (Model.maxWidth, 0 = unlimited) narrowed down to whatever
// actually fits the background it's about to be drawn on, so a toast can
// never overflow past the edge of the background - or the terminal, when
// the background is a full-screen View() - regardless of how WithMaxWidth
// was set. bgWidth is background's own width, already measured by Render.
func effectiveMaxWidth(configured, bgWidth int) int {
	fits := max(bgWidth-2*margin, 1)
	if configured > 0 && configured < fits {
		return configured
	}
	return fits
}

// clipToHeight trims stack to at most maxHeight lines when it overflows,
// keeping the lines nearest the anchored edge (top rows for a top anchor,
// bottom rows for a bottom anchor) so the newest toasts - always nearest
// that edge, see orderedToasts - are the ones that stay visible.
func clipToHeight(stack string, maxHeight int, anchoredTop bool) string {
	lines := strings.Split(stack, "\n")
	if len(lines) <= maxHeight {
		return stack
	}
	if maxHeight <= 0 {
		return ""
	}
	if anchoredTop {
		lines = lines[:maxHeight]
	} else {
		lines = lines[len(lines)-maxHeight:]
	}
	return strings.Join(lines, "\n")
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
