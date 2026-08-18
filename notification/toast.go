package notification

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type Kind int

const (
	Info Kind = iota
	Success
	Warning
	Error
)

const DefaultDuration = 3 * time.Second

type Toast struct {
	ID      string
	Title   string
	Message string
	Kind    Kind

	Duration time.Duration

	Style *KindStyle
}

type ToastOption func(*Toast)

func WithID(id string) ToastOption {
	return func(t *Toast) { t.ID = id }
}

func WithDuration(d time.Duration) ToastOption {
	return func(t *Toast) { t.Duration = d }
}

func WithToastStyle(s KindStyle) ToastOption {
	return func(t *Toast) { t.Style = &s }
}

func newToast(title, message string, kind Kind, opts ...ToastOption) Toast {
	t := Toast{
		Title:    title,
		Message:  message,
		Kind:     kind,
		Duration: DefaultDuration,
	}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

type ShowMsg struct{ Toast Toast }

func Show(title, message string, kind Kind, opts ...ToastOption) tea.Cmd {
	t := newToast(title, message, kind, opts...)
	return func() tea.Msg { return ShowMsg{Toast: t} }
}

type DismissMsg struct{ ID string }

func Dismiss(id string) tea.Cmd {
	return func() tea.Msg { return DismissMsg{ID: id} }
}

type expireMsg struct{ id string }
