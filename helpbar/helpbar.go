package helpbar

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/bubbles"
)

type Model struct {
	ShowAll bool

	help   help.Model
	global []key.Binding
	toggle key.Binding
	width  int
}

type Option func(*Model)

func WithGlobal(bindings ...key.Binding) Option {
	return func(m *Model) { m.global = bindings }
}

func WithToggle(b key.Binding) Option {
	return func(m *Model) { m.toggle = b }
}

func WithStyles(s help.Styles) Option {
	return func(m *Model) { m.help.Styles = s }
}

func New(opts ...Option) Model {
	m := Model{help: bubbles.NewHelp()}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m *Model) SetWidth(w int) {
	m.width = w
	m.help.SetWidth(w)
}

func (m Model) Width() int { return m.width }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && key.Matches(keyMsg, m.toggle) {
		m.ShowAll = !m.ShowAll
	}
	return m, nil
}

func (m Model) View(contextual ...key.Binding) string {
	bindings := m.bindings(contextual)
	if len(bindings) == 0 || m.width <= 0 {
		return ""
	}
	var view string
	if !m.ShowAll {
		view = m.help.ShortHelpView(bindings)
	} else {
		view = m.help.FullHelpView(m.columns(bindings))
	}
	return m.clampWidth(view)
}

func (m Model) clampWidth(view string) string {
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, m.width, "")
	}
	return strings.Join(lines, "\n")
}

func (m Model) Height(contextual ...key.Binding) int {
	view := m.View(contextual...)
	if view == "" {
		return 0
	}
	return lipgloss.Height(view)
}

func (m Model) bindings(contextual []key.Binding) []key.Binding {
	all := make([]key.Binding, 0, 1+len(m.global)+len(contextual))
	if m.toggle.Enabled() {
		all = append(all, m.toggle)
	}
	for _, b := range append(append([]key.Binding{}, m.global...), contextual...) {
		if b.Enabled() {
			all = append(all, b)
		}
	}
	return all
}

func (m Model) columns(bindings []key.Binding) [][]key.Binding {
	if m.width <= 0 {
		return [][]key.Binding{bindings}
	}
	for rows := 1; rows < len(bindings); rows++ {
		groups := chunkColumns(bindings, rows)
		if m.renderedWidth(groups) <= m.width {
			return groups
		}
	}

	return chunkColumns(bindings, len(bindings))
}

func (m Model) renderedWidth(groups [][]key.Binding) int {
	unbounded := m.help
	unbounded.SetWidth(0)
	return lipgloss.Width(unbounded.FullHelpView(groups))
}

func chunkColumns(bindings []key.Binding, rows int) [][]key.Binding {
	if rows < 1 {
		rows = 1
	}
	groups := make([][]key.Binding, 0, (len(bindings)+rows-1)/rows)
	for i := 0; i < len(bindings); i += rows {
		groups = append(groups, bindings[i:min(i+rows, len(bindings))])
	}
	return groups
}
