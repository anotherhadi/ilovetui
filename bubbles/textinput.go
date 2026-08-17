package bubbles

import (
	"charm.land/bubbles/v2/textinput"

	"github.com/anotherhadi/ilovetui/style"
)

// NewTextInput returns a textinput.Model styled with the active theme.
func NewTextInput() textinput.Model {
	t := textinput.New()
	t.SetStyles(themedTextInputStyles())
	return t
}

// themedTextInputStyles builds textinput.Styles from the active theme.
// Shared with NewList, which themes its filter input the same way.
func themedTextInputStyles() textinput.Styles {
	// isDark only affects textinput.DefaultStyles' Blurred.Text color, which
	// we override below regardless.
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
