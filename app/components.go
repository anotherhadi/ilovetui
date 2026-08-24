package app

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui/style"
)

func RenderToggle(on bool) string {
	if !style.S.NerdFonts {
		if on {
			return lipgloss.NewStyle().
				Foreground(style.S.Success).
				Bold(true).
				Render("on")
		}
		return lipgloss.NewStyle().
			Foreground(style.S.Subtle).
			Bold(true).
			Render("off")
	}

	if on {
		return RenderBadge("  ●", style.S.Success)
	}
	return RenderBadge("●  ", style.S.Subtle)
}

func RenderBadge(label string, fill color.Color) string {
	body := lipgloss.NewStyle().
		Background(fill).
		Foreground(style.S.Background).
		Bold(true).
		Render(label)

	cap := lipgloss.NewStyle().Foreground(fill).Background(style.S.Background)
	return cap.Render("") + body + cap.Render("")
}
