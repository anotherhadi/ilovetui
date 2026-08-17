package notification

// Position anchors the toast stack to one of six spots on the rendered
// background. There's no center/middle variant: toasts always hug an edge or
// a corner, never the middle of the screen.
type Position int

const (
	Top Position = iota
	TopLeft
	TopRight
	Bottom
	BottomLeft
	BottomRight
)

// margin is the fixed gap, in cells, kept between the toast stack and the
// edge(s) of the background it's anchored to.
const margin = 1

// placement resolves the top-left (x, y) coordinate to draw a stack of size
// (sw, sh) at, given a background of size (w, h) and the anchor position.
func placement(pos Position, w, h, sw, sh int) (x, y int) {
	switch pos {
	case Top:
		x = (w - sw) / 2
		y = margin
	case TopLeft:
		x = margin
		y = margin
	case TopRight:
		x = w - sw - margin
		y = margin
	case Bottom:
		x = (w - sw) / 2
		y = h - sh - margin
	case BottomLeft:
		x = margin
		y = h - sh - margin
	case BottomRight:
		x = w - sw - margin
		y = h - sh - margin
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return x, y
}

// anchoredTop reports whether pos hugs the top edge, which decides both the
// stacking order (see Model.orderedToasts) and which side of an overflowing
// stack gets clipped (see clipToHeight).
func (pos Position) anchoredTop() bool {
	return pos == Top || pos == TopLeft || pos == TopRight
}
