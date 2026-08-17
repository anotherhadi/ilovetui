package layout

import tea "charm.land/bubbletea/v2"

// FocusDirection is a screen-space direction for moving focus between
// panes: ctrl+h/j/k/l. Distinct from Direction (a Split's axis) because a
// split only ever has two sides, while focus needs to move one of four
// ways.
type FocusDirection int

const (
	FocusLeft FocusDirection = iota
	FocusRight
	FocusUp
	FocusDown
)

// FindNeighbor picks the leaf, among leaves, that's geometrically closest to
// from in dir - tmux's select-pane -L/-D/-U/-R, not a tree walk. A
// candidate must be strictly positioned beyond from's edge in dir and share
// some extent with it on the perpendicular axis. Ranked by, in order: edge
// gap (smaller wins), then shared perpendicular extent (larger wins - the
// real tie-breaker between two equally-close neighbors), then
// center-to-center distance on the perpendicular axis (last resort). Pure
// and independent of Model so it's testable on its own.
func FindNeighbor(leaves []LeafRect, from LeafRect, dir FocusDirection) (id string, ok bool) {
	type candidate struct {
		id                  string
		gap, overlap, cross int
	}
	var best *candidate

	for _, lr := range leaves {
		if lr.ID == from.ID {
			continue
		}
		gap, overlap, cross, ok := edgeScore(from.Rect, lr.Rect, dir)
		if !ok {
			continue
		}
		c := candidate{id: lr.ID, gap: gap, overlap: overlap, cross: cross}
		if best == nil ||
			c.gap < best.gap ||
			(c.gap == best.gap && c.overlap > best.overlap) ||
			(c.gap == best.gap && c.overlap == best.overlap && c.cross < best.cross) {
			best = &c
		}
	}

	if best == nil {
		return "", false
	}
	return best.id, true
}

// edgeScore returns the ranking tuple used by FindNeighbor for a single
// (from, other) pair, and ok=false if other isn't a valid candidate in dir
// at all (wrong side, or zero overlap on the perpendicular axis).
func edgeScore(from, other Rect, dir FocusDirection) (gap, overlap, cross int, ok bool) {
	switch dir {
	case FocusLeft:
		if other.X+other.W > from.X {
			return 0, 0, 0, false
		}
		gap = from.X - (other.X + other.W)
		overlap = spanOverlap(from.Y, from.Y+from.H, other.Y, other.Y+other.H)
		cross = abs((from.Y + from.H/2) - (other.Y + other.H/2))
	case FocusRight:
		if other.X < from.X+from.W {
			return 0, 0, 0, false
		}
		gap = other.X - (from.X + from.W)
		overlap = spanOverlap(from.Y, from.Y+from.H, other.Y, other.Y+other.H)
		cross = abs((from.Y + from.H/2) - (other.Y + other.H/2))
	case FocusUp:
		if other.Y+other.H > from.Y {
			return 0, 0, 0, false
		}
		gap = from.Y - (other.Y + other.H)
		overlap = spanOverlap(from.X, from.X+from.W, other.X, other.X+other.W)
		cross = abs((from.X + from.W/2) - (other.X + other.W/2))
	case FocusDown:
		if other.Y < from.Y+from.H {
			return 0, 0, 0, false
		}
		gap = other.Y - (from.Y + from.H)
		overlap = spanOverlap(from.X, from.X+from.W, other.X, other.X+other.W)
		cross = abs((from.X + from.W/2) - (other.X + other.W/2))
	}
	if overlap <= 0 {
		return 0, 0, 0, false
	}
	return gap, overlap, cross, true
}

// spanOverlap returns the length shared by [aStart,aEnd) and [bStart,bEnd).
func spanOverlap(aStart, aEnd, bStart, bEnd int) int {
	start := max(aStart, bStart)
	end := min(aEnd, bEnd)
	if end <= start {
		return 0
	}
	return end - start
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// leafRect looks up id's current Rect in this Model's own flat registry.
func (m Model) leafRect(id string) (LeafRect, bool) {
	for _, lr := range m.leaves {
		if lr.ID == id {
			return lr, true
		}
	}
	return LeafRect{}, false
}

// MoveFocus implements Navigable. If the currently focused leaf holds a
// nested Navigable (an embedded layout.Model), it's given first refusal -
// only once it reports being at its own edge (false) does this Model try
// moving focus among its own direct children instead. Cmds produced by the
// BlurMsg/FocusMsg this triggers can't be returned directly (Navigable's
// signature is bool-only); they're queued on m.state.pendingCmd instead
// (see focusState) and drained by the next Update that actually returns a
// tea.Cmd - a lag of at most one Update cycle, never user-visible.
func (m Model) MoveFocus(dir FocusDirection) bool {
	if focused, ok := findNode(m.root, m.state.id); ok {
		if nav, ok := focused.model.(Navigable); ok {
			if nav.MoveFocus(dir) {
				focused.model = nav
				return true
			}
		}
	}

	from, ok := m.leafRect(m.state.id)
	if !ok {
		return false
	}
	id, ok := FindNeighbor(m.leaves, from, dir)
	if !ok {
		return false
	}
	m.state.pendingCmd = tea.Batch(m.state.pendingCmd, m.setFocus(id))
	return true
}
