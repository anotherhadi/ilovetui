package app

import tea "charm.land/bubbletea/v2"

type QuitMsg struct{}

func Quit() tea.Cmd {
	return func() tea.Msg { return QuitMsg{} }
}

type FocusMsg struct{}

func Focus() tea.Cmd {
	return func() tea.Msg { return FocusMsg{} }
}

type BlurMsg struct{}

func Blur() tea.Cmd {
	return func() tea.Msg { return BlurMsg{} }
}

type RequestFocusMsg struct{}

func RequestFocus() tea.Cmd {
	return func() tea.Msg { return RequestFocusMsg{} }
}

type SetTitleMsg struct{ Title string }

func SetTitle(title string) tea.Cmd {
	return func() tea.Msg { return SetTitleMsg{Title: title} }
}
