package bubbles

import "strings"

// SplitH splits totalHeight into top and bottom sections, accounting for the
// height of statusBar (measured by newline count).
func SplitH(totalHeight int, statusBar string, ratio float64) (top, bottom int) {
	statusH := strings.Count(statusBar, "\n") + 1
	available := totalHeight - statusH
	top = int(float64(available) * ratio)
	bottom = available - top
	return
}
