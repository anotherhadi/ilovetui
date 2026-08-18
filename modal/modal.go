package modal

import tea "charm.land/bubbletea/v2"

type Model struct {
	modals    []Modal
	maxWidth  int
	maxHeight int
	styles    Styles
}

type Option func(*Model)

func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

func WithMaxHeight(h int) Option {
	return func(m *Model) { m.maxHeight = h }
}

func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

func New(opts ...Option) Model {
	m := Model{
		maxWidth:  60,
		maxHeight: 20,
		styles:    DefaultStyles(),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowMsg:
		return m.show(msg.Modal)
	case DismissMsg:
		return m.pop(), nil
	}
	return m.updateTop(msg)
}

func (m Model) updateTop(msg tea.Msg) (Model, tea.Cmd) {
	i := len(m.modals) - 1
	if i < 0 || m.modals[i].Content == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.modals[i].Content, cmd = m.modals[i].Content.Update(msg)
	return m, cmd
}

func (m Model) Open() bool { return len(m.modals) > 0 }

func (m Model) show(mo Modal) (Model, tea.Cmd) {
	m.modals = append(m.modals, mo)
	return m, initContent(mo)
}

func initContent(mo Modal) tea.Cmd {
	if mo.Content == nil {
		return nil
	}
	return mo.Content.Init()
}

func (m Model) pop() Model {
	if len(m.modals) == 0 {
		return m
	}
	m.modals = m.modals[:len(m.modals)-1]
	return m
}
