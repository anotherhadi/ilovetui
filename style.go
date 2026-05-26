package ilovetui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// Styles holds both the raw Base16 palette and ready-to-use semantic colors
// and lipgloss styles. Access via the package-level variable S.
type Styles struct {
	// Raw Base16 palette
	Base00 color.Color
	Base01 color.Color
	Base02 color.Color
	Base03 color.Color
	Base04 color.Color
	Base05 color.Color
	Base06 color.Color
	Base07 color.Color
	Base08 color.Color
	Base09 color.Color
	Base0A color.Color
	Base0B color.Color
	Base0C color.Color
	Base0D color.Color
	Base0E color.Color
	Base0F color.Color

	// Semantic color aliases
	Background color.Color
	SubtleBg   color.Color
	Selection  color.Color
	Subtle     color.Color
	Muted      color.Color
	Text       color.Color
	Primary    color.Color
	Success    color.Color
	Warning    color.Color
	Error      color.Color

	// Pre-built text styles
	Bold  lipgloss.Style
	Faint lipgloss.Style

	// Pre-built panel styles (rounded border)
	Panel        lipgloss.Style
	PanelFocused lipgloss.Style

	// Pre-rendered pager dot strings
	PagerDotActive   string
	PagerDotInactive string
}

func newStyles(c colorsYAML) Styles {
	lc := func(s string) color.Color {
		s = strings.TrimSpace(s)
		if s != "" && s[0] != '#' {
			s = "#" + s
		}
		return lipgloss.Color(s)
	}

	b00 := lc(c.Base00)
	b01 := lc(c.Base01)
	b02 := lc(c.Base02)
	b03 := lc(c.Base03)
	b04 := lc(c.Base04)
	b05 := lc(c.Base05)
	b06 := lc(c.Base06)
	b07 := lc(c.Base07)
	b08 := lc(c.Base08)
	b09 := lc(c.Base09)
	b0A := lc(c.Base0A)
	b0B := lc(c.Base0B)
	b0C := lc(c.Base0C)
	b0D := lc(c.Base0D)
	b0E := lc(c.Base0E)
	b0F := lc(c.Base0F)

	return Styles{
		Base00: b00, Base01: b01, Base02: b02, Base03: b03,
		Base04: b04, Base05: b05, Base06: b06, Base07: b07,
		Base08: b08, Base09: b09, Base0A: b0A, Base0B: b0B,
		Base0C: b0C, Base0D: b0D, Base0E: b0E, Base0F: b0F,

		Background: b00,
		SubtleBg:   b01,
		Selection:  b02,
		Subtle:     b03,
		Muted:      b04,
		Text:       b05,
		Primary:    b0D,
		Success:    b0B,
		Warning:    b09,
		Error:      b08,

		Bold:  lipgloss.NewStyle().Bold(true),
		Faint: lipgloss.NewStyle().Foreground(b03).Faint(true),

		Panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(b03),

		PanelFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(b0D),

		PagerDotActive:   lipgloss.NewStyle().Foreground(b0D).SetString("•").String(),
		PagerDotInactive: lipgloss.NewStyle().Foreground(b03).SetString("•").String(),
	}
}
