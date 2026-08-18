package bubbles

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

func NewTextarea(showLineNumbers bool) textarea.Model {
	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = showLineNumbers
	ta.CharLimit = 0
	ta.EndOfBufferCharacter = '~'
	ts := ta.Styles()
	ts.Focused.Base = lipgloss.NewStyle()
	ts.Blurred.Base = lipgloss.NewStyle()
	ts.Focused.Text = lipgloss.NewStyle().Foreground(style.S.Text)
	ts.Focused.CursorLine = lipgloss.NewStyle().Background(style.S.Selection).Foreground(style.S.Text)
	ts.Focused.CursorLineNumber = lipgloss.NewStyle().Background(style.S.Selection).Foreground(style.S.Primary).Bold(true)
	ts.Focused.LineNumber = lipgloss.NewStyle().Foreground(style.S.Subtle)
	ts.Focused.Placeholder = lipgloss.NewStyle().Foreground(style.S.Subtle)
	ts.Focused.EndOfBuffer = lipgloss.NewStyle().Foreground(style.S.SubtleBg)
	ts.Blurred.Text = lipgloss.NewStyle().Foreground(style.S.Muted)
	ts.Blurred.LineNumber = lipgloss.NewStyle().Foreground(style.S.SubtleBg)
	ts.Blurred.Placeholder = lipgloss.NewStyle().Foreground(style.S.Subtle)
	ts.Blurred.EndOfBuffer = lipgloss.NewStyle().Foreground(style.S.SubtleBg)
	ta.SetStyles(ts)
	return ta
}
