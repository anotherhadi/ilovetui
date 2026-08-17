package layout

import "testing"

func TestComputeLayoutEvenSplit(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	leaves := computeLayout(root, Rect{W: 100, H: 40})

	want := map[string]Rect{
		"a": {X: 0, Y: 0, W: 50, H: 40},
		"b": {X: 50, Y: 0, W: 50, H: 40},
	}
	assertLeafRects(t, leaves, want)
}

func TestComputeLayoutVerticalRatio(t *testing.T) {
	root := VSplit(0.25, Leaf("top", newStub()), Leaf("bottom", newStub()))
	leaves := computeLayout(root, Rect{W: 80, H: 40})

	want := map[string]Rect{
		"top":    {X: 0, Y: 0, W: 80, H: 10},
		"bottom": {X: 0, Y: 10, W: 80, H: 30},
	}
	assertLeafRects(t, leaves, want)
}

func TestComputeLayoutMinimum(t *testing.T) {
	root := HSplit(0.1, Leaf("a", newStub()), Leaf("b", newStub())).WithMinimum(20)
	leaves := computeLayout(root, Rect{W: 100, H: 10})

	want := map[string]Rect{
		"a": {X: 0, Y: 0, W: 20, H: 10},
		"b": {X: 20, Y: 0, W: 80, H: 10},
	}
	assertLeafRects(t, leaves, want)
}

func TestComputeLayoutMaximum(t *testing.T) {
	root := HSplit(0.9, Leaf("a", newStub()), Leaf("b", newStub())).WithMaximum(20)
	leaves := computeLayout(root, Rect{W: 100, H: 10})

	want := map[string]Rect{
		"a": {X: 0, Y: 0, W: 20, H: 10},
		"b": {X: 20, Y: 0, W: 80, H: 10},
	}
	assertLeafRects(t, leaves, want)
}

func TestComputeLayoutFixedWhenMinEqualsMax(t *testing.T) {
	root := HSplit(0.9, Leaf("a", newStub()), Leaf("b", newStub())).WithMinimum(20).WithMaximum(20)
	leaves := computeLayout(root, Rect{W: 100, H: 10})

	want := map[string]Rect{
		"a": {X: 0, Y: 0, W: 20, H: 10},
		"b": {X: 20, Y: 0, W: 80, H: 10},
	}
	assertLeafRects(t, leaves, want)
}

func TestComputeLayoutNested(t *testing.T) {
	root := HSplit(0.3,
		Leaf("sidebar", newStub()),
		VSplit(0.5, Leaf("top", newStub()), Leaf("bottom", newStub())),
	)
	leaves := computeLayout(root, Rect{W: 100, H: 20})

	want := map[string]Rect{
		"sidebar": {X: 0, Y: 0, W: 30, H: 20},
		"top":     {X: 30, Y: 0, W: 70, H: 10},
		"bottom":  {X: 30, Y: 10, W: 70, H: 10},
	}
	assertLeafRects(t, leaves, want)
}

func assertLeafRects(t *testing.T, leaves []LeafRect, want map[string]Rect) {
	t.Helper()
	if len(leaves) != len(want) {
		t.Fatalf("got %d leaves, want %d (%v)", len(leaves), len(want), leaves)
	}
	for _, lr := range leaves {
		wr, ok := want[lr.ID]
		if !ok {
			t.Fatalf("unexpected leaf %q", lr.ID)
		}
		if lr.Rect != wr {
			t.Errorf("leaf %q: got %+v, want %+v", lr.ID, lr.Rect, wr)
		}
	}
}
