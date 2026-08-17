package notification

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

// Kind picks which of Styles' presets a toast renders with, unless
// overridden per-toast via WithToastStyle.
type Kind int

const (
	Info Kind = iota
	Success
	Warning
	Error
)

// DefaultDuration is how long a toast stays visible when WithDuration isn't
// used. Show has no reference to a Model (see ShowMsg's doc comment), so this
// lives as a package constant rather than a Model-level default.
const DefaultDuration = 3 * time.Second

// Toast is one notification. Build it via Show's opts rather than a literal:
// ID and Duration both get defaults (see WithID, DefaultDuration) that a bare
// literal would silently skip.
type Toast struct {
	ID      string
	Title   string
	Message string
	Kind    Kind
	// Duration is how long the toast stays up before auto-dismissing. 0
	// means sticky: it stays until DismissMsg/Dismiss(ID) removes it.
	Duration time.Duration
	// Style, if non-nil, overrides the Model's Kind-based preset for this
	// toast alone.
	Style *KindStyle
}

// ToastOption configures a Toast built by Show.
type ToastOption func(*Toast)

// WithID gives the toast a stable id, so a later Show reusing the same id
// replaces it in place (resetting its position and timer) instead of
// stacking a duplicate, and so it can be targeted by Dismiss.
func WithID(id string) ToastOption {
	return func(t *Toast) { t.ID = id }
}

// WithDuration overrides DefaultDuration. 0 makes the toast sticky: it never
// auto-dismisses, only Dismiss(ID) removes it.
func WithDuration(d time.Duration) ToastOption {
	return func(t *Toast) { t.Duration = d }
}

// WithToastStyle overrides the Model's Kind-based preset for this toast
// alone, for a one-off custom look instead of the type-based theme.
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

// ShowMsg tells a notification.Model to display Toast. Any component in the
// same bubbletea program can trigger one via Show, without holding a
// reference to the notification.Model that will actually render it - that
// Model just needs to see every tea.Msg the program produces, same as any
// other child model.
type ShowMsg struct{ Toast Toast }

// Show returns a tea.Cmd that shows a new toast of the given kind. Call it
// from any component's Update:
//
//	return m, notification.Show("Saved", "Config written to disk", notification.Success)
func Show(title, message string, kind Kind, opts ...ToastOption) tea.Cmd {
	t := newToast(title, message, kind, opts...)
	return func() tea.Msg { return ShowMsg{Toast: t} }
}

// DismissMsg removes the toast identified by ID, whether it's sticky or
// mid-countdown. A no-op if ID isn't currently shown (already expired, or
// never had an explicit id in the first place - see WithID).
type DismissMsg struct{ ID string }

// Dismiss returns a tea.Cmd that removes the toast identified by id. Only
// useful for toasts shown with WithID, since an auto-generated id is never
// exposed back to the caller.
func Dismiss(id string) tea.Cmd {
	return func() tea.Msg { return DismissMsg{ID: id} }
}

// expireMsg fires once a toast's Duration has elapsed, scheduled by
// Model.show via tea.Tick. Unexported: nothing outside the package should
// construct or match on it directly, that's what DismissMsg is for.
type expireMsg struct{ id string }
