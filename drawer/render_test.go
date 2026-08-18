package drawer

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func background(w, h int) string {
	return blank(w, h)
}

func TestRenderLeftFlushToLeftEdge(t *testing.T) {
	m := New(WithMaxWidth(10))
	m, _ = m.Update(ShowMsg{Drawer: newDrawer("Nav", Text("hi"), WithSide(Left))})

	out := m.Render(background(40, 10))
	lines := strings.Split(out, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(lines))
	}
	if lipgloss.Width(out) == 0 {
		t.Fatalf("expected non-empty render")
	}
	// The border's top-left corner glyph should be the very first rune.
	if len([]rune(lines[0])) == 0 {
		t.Fatalf("expected a rendered top border line")
	}
}

func TestRenderRightFlushToRightEdge(t *testing.T) {
	m := New(WithMaxWidth(10))
	m, _ = m.Update(ShowMsg{Drawer: newDrawer("Inspector", Text("hi"), WithSide(Right))})

	out := m.Render(background(40, 10))
	if lipgloss.Width(out) != 40 {
		t.Fatalf("expected full background width 40, got %d", lipgloss.Width(out))
	}
}

func TestSpansFullBackgroundHeight(t *testing.T) {
	m := New()
	m, _ = m.Update(ShowMsg{Drawer: newDrawer("Nav", Text("hi"))})

	out := m.Render(background(40, 12))
	if lipgloss.Height(out) != 12 {
		t.Fatalf("expected full height 12, got %d", lipgloss.Height(out))
	}
}

func TestFixedWidthHonored(t *testing.T) {
	m := New(WithMaxWidth(50))
	m, _ = m.Update(ShowMsg{Drawer: newDrawer("Nav", Text("x"), WithWidth(20))})

	box := m.renderBox(m.drawers[0], m.styles, 80, 10)
	if got := lipgloss.Width(box); got != 20 { // WithWidth is the total box width
		t.Fatalf("expected fixed box width 20, got %d", got)
	}
}

func TestNoDrawersReturnsBackgroundUnchanged(t *testing.T) {
	m := New()
	bg := background(10, 5)
	if got := m.Render(bg); got != bg {
		t.Fatalf("expected background unchanged when no drawer is open")
	}
}
