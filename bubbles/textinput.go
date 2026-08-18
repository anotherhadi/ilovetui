package bubbles

import (
	"charm.land/bubbles/v2/textinput"

	"github.com/anotherhadi/ilovetui/style"
)

func NewTextInput() textinput.Model {
	t := textinput.New()
	t.SetStyles(themedTextInputStyles())
	return t
}

func themedTextInputStyles() textinput.Styles {
	s := textinput.DefaultStyles(true)
	s.Focused.Text = s.Focused.Text.Foreground(style.S.Text)
	s.Focused.Placeholder = s.Focused.Placeholder.Foreground(style.S.Subtle)
	s.Focused.Suggestion = s.Focused.Suggestion.Foreground(style.S.Subtle)
	s.Focused.Prompt = s.Focused.Prompt.Foreground(style.S.Primary)
	s.Blurred.Text = s.Blurred.Text.Foreground(style.S.Muted)
	s.Blurred.Placeholder = s.Blurred.Placeholder.Foreground(style.S.Subtle)
	s.Blurred.Suggestion = s.Blurred.Suggestion.Foreground(style.S.Subtle)
	s.Blurred.Prompt = s.Blurred.Prompt.Foreground(style.S.Subtle)
	s.Cursor.Color = style.S.Primary
	return s
}
