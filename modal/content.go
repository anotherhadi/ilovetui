package modal

import tea "charm.land/bubbletea/v2"

// Modal is one popup on the stack.
type Modal struct {
	Title string
	// Content is the modal's body: a full model, updated and rendered by
	// the modal.Model while it's on top of the stack. Wrap a plain string
	// with Text for a modal with nothing to interact with.
	Content tea.Model
	// Style, if non-nil, overrides the Model's default Styles for this
	// modal alone.
	Style *Styles
}

// ModalOption configures a Modal built by Show.
type ModalOption func(*Modal)

// WithModalStyle overrides the Model's default Styles for this modal alone,
// for a one-off custom look instead of the shared theme.
func WithModalStyle(s Styles) ModalOption {
	return func(mo *Modal) { mo.Style = &s }
}

func newModal(title string, content tea.Model, opts ...ModalOption) Modal {
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

// Show returns a tea.Cmd that opens a new modal on top of the stack. content
// is a model, so a modal can hold anything a pane can - a form, a list, a
// confirmation that reports back with its own tea.Msg:
//
//	return m, modal.Show("Delete file?", modal.Text("This can't be undone."))
//	return m, modal.Show("Rename", newRenameForm(path))
//
// The content's Init runs when the modal opens, and it receives every message
// while it's the topmost modal (see Model.Update).
func Show(title string, content tea.Model, opts ...ModalOption) tea.Cmd {
	mo := newModal(title, content, opts...)
	return func() tea.Msg { return ShowMsg{Modal: mo} }
}

// DismissMsg closes the topmost modal. The stack is a plain LIFO: a modal is
// closed by being on top, never by being named.
type DismissMsg struct{}

// Close returns a tea.Cmd that closes the topmost modal - which is also the
// only one that can act (see Model.Update), so a content model closes itself
// by returning it:
//
//	return c, modal.Close()
func Close() tea.Cmd {
	return func() tea.Msg { return DismissMsg{} }
}

// text is a model wrapping a fixed string: a modal body with nothing to
// update.
type text string

func (t text) Init() tea.Cmd                       { return nil }
func (t text) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t text) View() tea.View                      { return tea.NewView(string(t)) }

// Text wraps a plain string as modal content, for the common modal that has
// nothing to interact with:
//
//	modal.Show("About", modal.Text("v1.0\n\nesc to close"))
func Text(s string) tea.Model { return text(s) }
