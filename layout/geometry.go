package layout

import "math"

// Rect is an axis-aligned screen region in terminal cells, origin top-left.
type Rect struct {
	X, Y, W, H int
}

// LeafRect pairs a Leaf's id with the Rect it was allocated by the most
// recent layout pass.
type LeafRect struct {
	ID   string
	Rect Rect
}

// computeLayout descends the tree rooted at n, allocating r between its
// leaves according to each Split's ratio/min/max, and returns a flat
// registry of every leaf's resolved Rect. Order is deterministic (a
// depth-first walk, first child before second), which is what makes it safe
// to use directly as a stable iteration order elsewhere (Init, routing).
func computeLayout(n *Node, r Rect) []LeafRect {
	if n == nil {
		return nil
	}
	if n.leaf {
		return []LeafRect{{ID: n.id, Rect: r}}
	}

	var firstRect, secondRect Rect
	if n.dir == Horizontal {
		w1 := resolveSize(n, r.W)
		firstRect = Rect{X: r.X, Y: r.Y, W: w1, H: r.H}
		secondRect = Rect{X: r.X + w1, Y: r.Y, W: r.W - w1, H: r.H}
	} else {
		h1 := resolveSize(n, r.H)
		firstRect = Rect{X: r.X, Y: r.Y, W: r.W, H: h1}
		secondRect = Rect{X: r.X, Y: r.Y + h1, W: r.W, H: r.H - h1}
	}

	leaves := computeLayout(n.first, firstRect)
	return append(leaves, computeLayout(n.second, secondRect)...)
}

// resolveSize returns the cell size a Split's first child gets out of total,
// starting from n.ratio and then clamped to [n.min, n.max] (0 on either side
// means that bound is unset). n.min == n.max fixes the size outright,
// regardless of ratio.
func resolveSize(n *Node, total int) int {
	size := int(math.Round(n.ratio * float64(total)))
	if n.min > 0 && size < n.min {
		size = n.min
	}
	if n.max > 0 && size > n.max {
		size = n.max
	}
	if size < 0 {
		size = 0
	}
	if size > total {
		size = total
	}
	return size
}
