package layout

import "testing"

func TestFindNeighborBasicDirections(t *testing.T) {
	leaves := []LeafRect{
		{ID: "left", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}},
		{ID: "right", Rect: Rect{X: 10, Y: 0, W: 10, H: 10}},
	}

	if id, ok := FindNeighbor(leaves, leaves[0], FocusRight); !ok || id != "right" {
		t.Fatalf("FocusRight from left: got (%q, %v), want (right, true)", id, ok)
	}
	if id, ok := FindNeighbor(leaves, leaves[1], FocusLeft); !ok || id != "left" {
		t.Fatalf("FocusLeft from right: got (%q, %v), want (left, true)", id, ok)
	}
	if _, ok := FindNeighbor(leaves, leaves[0], FocusLeft); ok {
		t.Fatalf("FocusLeft from left (edge): expected no neighbor")
	}
}

// Two candidates at the exact same gap: the one that actually shares more
// of from's edge should win, not an arbitrary pick. This is the overlap
// tie-break, the one that matters most for a real tmux/i3-style layout
// (picking the pane that's genuinely alongside you, not just "a" neighbor).
func TestFindNeighborPrefersLargerOverlapOnEqualGap(t *testing.T) {
	from := LeafRect{ID: "from", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}}
	leaves := []LeafRect{
		from,
		{ID: "full", Rect: Rect{X: 10, Y: 0, W: 10, H: 10}},
		{ID: "partial", Rect: Rect{X: 10, Y: 5, W: 10, H: 10}},
	}

	id, ok := FindNeighbor(leaves, from, FocusRight)
	if !ok || id != "full" {
		t.Fatalf("got (%q, %v), want (full, true): equal gap, full's overlap (10) beats partial's (5)", id, ok)
	}
}

// When gap AND overlap are tied, center-to-center distance on the
// perpendicular axis is the last resort.
func TestFindNeighborCrossIsLastResortTiebreak(t *testing.T) {
	from := LeafRect{ID: "from", Rect: Rect{X: 0, Y: 10, W: 10, H: 10}} // Y 10..20, center 15
	leaves := []LeafRect{
		from,
		{ID: "near", Rect: Rect{X: 10, Y: 8, W: 10, H: 20}}, // Y 8..28, overlap 10, center 18
		{ID: "far", Rect: Rect{X: 10, Y: 0, W: 10, H: 20}},  // Y 0..20, overlap 10, center 10
	}

	id, ok := FindNeighbor(leaves, from, FocusRight)
	if !ok || id != "near" {
		t.Fatalf("got (%q, %v), want (near, true): equal gap and overlap, near's center is closer (3 vs 5)", id, ok)
	}
}

func TestFindNeighborNoCandidateBeyondEdge(t *testing.T) {
	leaves := []LeafRect{
		{ID: "only", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}},
	}
	if _, ok := FindNeighbor(leaves, leaves[0], FocusDown); ok {
		t.Fatal("expected no neighbor with a single leaf")
	}
}

func TestFindNeighborRequiresPerpendicularOverlap(t *testing.T) {
	leaves := []LeafRect{
		{ID: "from", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}},
		// Directly to the right on the X axis, but no shared Y extent at
		// all: not a valid candidate even though it's the only one there.
		{ID: "diagonal", Rect: Rect{X: 10, Y: 10, W: 10, H: 10}},
	}
	if _, ok := FindNeighbor(leaves, leaves[0], FocusRight); ok {
		t.Fatal("expected no neighbor: candidate shares no Y extent with from")
	}
}
