package drawer

import tea "charm.land/bubbletea/v2"

type Model struct {
	drawers  []Drawer
	nextID   int
	maxWidth int
	styles   Styles
}

type Option func(*Model)

func WithMaxWidth(w int) Option {
	return func(m *Model) { m.maxWidth = w }
}

func WithStyles(s Styles) Option {
	return func(m *Model) { m.styles = s }
}

func New(opts ...Option) Model {
	m := Model{
		maxWidth: 30,
		styles:   DefaultStyles(),
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
		return m.show(msg.Drawer)
	case DismissMsg:
		return m.pop(), nil
	}
	return m.updateTop(msg)
}

func (m Model) updateTop(msg tea.Msg) (Model, tea.Cmd) {
	i := len(m.drawers) - 1
	if i < 0 || m.drawers[i].Content == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.drawers[i].Content, cmd = m.drawers[i].Content.Update(msg)
	return m, cmd
}

func (m Model) Open() bool { return len(m.drawers) > 0 }

func (m Model) show(d Drawer) (Model, tea.Cmd) {
	m.drawers = append(m.drawers, d)
	return m, initContent(d)
}

func initContent(d Drawer) tea.Cmd {
	if d.Content == nil {
		return nil
	}
	return d.Content.Init()
}

func (m Model) pop() Model {
	if len(m.drawers) == 0 {
		return m
	}
	m.drawers = m.drawers[:len(m.drawers)-1]
	return m
}
