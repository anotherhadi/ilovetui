package style

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func Overlay(base, top string, x, y int) string {
	baseLines := strings.Split(base, "\n")
	topLines := strings.Split(top, "\n")

	for i, tl := range topLines {
		row := y + i
		if row < 0 || row >= len(baseLines) {
			continue
		}
		bl := baseLines[row]
		tw := ansi.StringWidth(tl)
		left := ansi.Cut(bl, 0, x)
		right := ansi.Cut(bl, x+tw, ansi.StringWidth(bl))
		baseLines[row] = left + tl + right
	}

	return strings.Join(baseLines, "\n")
}
