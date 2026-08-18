package notification

type Position int

const (
	Top Position = iota
	TopLeft
	TopRight
	Bottom
	BottomLeft
	BottomRight
)

const margin = 1

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

func (pos Position) anchoredTop() bool {
	return pos == Top || pos == TopLeft || pos == TopRight
}
