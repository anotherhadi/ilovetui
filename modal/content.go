package modal

import tea "charm.land/bubbletea/v2"

type Modal struct {
	Title string

	Content tea.Model

	Style *Styles
}

type ModalOption func(*Modal)

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

type ShowMsg struct{ Modal Modal }

func Show(title string, content tea.Model, opts ...ModalOption) tea.Cmd {
	mo := newModal(title, content, opts...)
	return func() tea.Msg { return ShowMsg{Modal: mo} }
}

type DismissMsg struct{}

func Close() tea.Cmd {
	return func() tea.Msg { return DismissMsg{} }
}

type text string

func (t text) Init() tea.Cmd                       { return nil }
func (t text) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t text) View() tea.View                      { return tea.NewView(string(t)) }

func Text(s string) tea.Model { return text(s) }
