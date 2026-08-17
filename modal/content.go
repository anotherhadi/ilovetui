package modal

import tea "charm.land/bubbletea/v2"

// Modal is one popup. Build it via Show's opts rather than a literal: ID
// gets a default (see WithID) that a bare literal would silently skip.
type Modal struct {
	ID      string
	Title   string
	Content string
	// Style, if non-nil, overrides the Model's default Styles for this
	// modal alone.
	Style *Styles
}

// ModalOption configures a Modal built by Show.
type ModalOption func(*Modal)

// WithID gives the modal a stable id, so a later Show reusing the same id
// replaces it in place instead of pushing a duplicate on the stack, and so
// it can be targeted by Dismiss.
func WithID(id string) ModalOption {
	return func(mo *Modal) { mo.ID = id }
}

// WithModalStyle overrides the Model's default Styles for this modal alone,
// for a one-off custom look instead of the shared theme.
func WithModalStyle(s Styles) ModalOption {
	return func(mo *Modal) { mo.Style = &s }
}

func newModal(title, content string, opts ...ModalOption) Modal {
	mo := Modal{Title: title, Content: content}
	for _, opt := range opts {
		opt(&mo)
	}
	return mo
}

// ShowMsg tells a modal.Model to display Modal, pushing it on top of the
// stack. Any component in the same bubbletea program can trigger one via
// Show, without holding a reference to the modal.Model that will actually
// render it - that Model just needs to see every tea.Msg the program
// produces, same as any other child model.
type ShowMsg struct{ Modal Modal }

// Show returns a tea.Cmd that opens a new modal on top of the stack:
//
//	return m, modal.Show("Delete file?", "This can't be undone.")
func Show(title, content string, opts ...ModalOption) tea.Cmd {
	mo := newModal(title, content, opts...)
	return func() tea.Msg { return ShowMsg{Modal: mo} }
}

// DismissMsg closes the modal identified by ID, or the topmost modal if ID
// is empty.
type DismissMsg struct{ ID string }

// Dismiss returns a tea.Cmd that closes the modal identified by id. Only
// useful for modals shown with WithID.
func Dismiss(id string) tea.Cmd {
	return func() tea.Msg { return DismissMsg{ID: id} }
}

// Close returns a tea.Cmd that closes the topmost modal, whatever its id -
// the common case of a host handling esc/"cancel" without needing to know
// which modal is currently on top.
func Close() tea.Cmd {
	return func() tea.Msg { return DismissMsg{} }
}
