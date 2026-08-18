package notification

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

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

	compositor := lipgloss.NewCompositor(
		lipgloss.NewLayer(background),
		lipgloss.NewLayer(stack).X(x).Y(y).Z(1),
	)
	return compositor.Render()
}

func (m Model) View(width, height int) string {
	return m.Render(blank(width, height))
}

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

func (m Model) renderToast(t Toast, maxWidth int) string {
	k := m.styles.forKind(t)

	inner := contentWidth(t, maxWidth)
	message := k.Message.Width(inner).Render(t.Message)

	boxWidth := inner + 4
	boxHeight := lipgloss.Height(message) + 2

	return style.RenderWithTitle(k.Border, k.Title.Render(t.Title), message, boxWidth, boxHeight)
}

func contentWidth(t Toast, maxWidth int) int {
	natural := max(lipgloss.Width(t.Title), lipgloss.Width(t.Message), 1)
	if maxWidth <= 0 {
		return natural
	}
	capped := max(maxWidth-4, 1)
	return min(natural, capped)
}

func effectiveMaxWidth(configured, bgWidth int) int {
	fits := max(bgWidth-2*margin, 1)
	if configured > 0 && configured < fits {
		return configured
	}
	return fits
}

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
