// Package helpbar is a responsive help bar: a single line of key bindings
// that expands, on demand, into a multi-column view reflowed to use as many
// columns as the available width allows.
//
// It has no dependency on any particular layout or container - it's just a
// component that takes a width and some bindings and returns a string, so it
// works as well under a plain lipgloss.JoinVertical as anywhere else.
//
// Bindings come from two places. Global ones (quit, toggle help, whatever
// your app reserves for itself) are set once via WithGlobal and always shown
// first. Contextual ones are passed to View at render time, so the bar can
// track whatever component currently has focus:
//
//	bar := m.help.View(m.focused().HelpBindings()...)
//	body := m.height - lipgloss.Height(bar)
package helpbar

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
)

// Model is a help bar. The zero value isn't usable - build one with New.
type Model struct {
	// ShowAll switches between the one-line short view and the full
	// multi-column view. Set it directly, or let WithToggle bind a key to
	// it and have Update flip it for you.
	ShowAll bool

	help   help.Model
	global []key.Binding
	toggle key.Binding
	width  int
}

// Option configures a Model at construction.
type Option func(*Model)

// WithGlobal sets the bindings shown before the contextual ones on every
// render - the keys your app reserves for itself regardless of what's
// focused.
func WithGlobal(bindings ...key.Binding) Option {
	return func(m *Model) { m.global = bindings }
}

// WithToggle makes Update flip ShowAll when b matches, and lists b ahead of
// every other binding (including WithGlobal's) so the way to expand the bar
// is always the first thing shown. Without it, Update ignores key presses
// and toggling ShowAll is entirely up to the caller.
func WithToggle(b key.Binding) Option {
	return func(m *Model) { m.toggle = b }
}

// WithStyles overrides the themed default styles.
func WithStyles(s help.Styles) Option {
	return func(m *Model) { m.help.Styles = s }
}

// New builds a help bar themed from the shared style package.
func New(opts ...Option) Model {
	m := Model{help: bubbles.NewHelp()}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// SetWidth sets the width the bar renders within. Nothing is shown until
// this is called with a positive value - typically from your
// tea.WindowSizeMsg handler.
func (m *Model) SetWidth(w int) {
	m.width = w
	m.help.SetWidth(w)
}

// Width returns the width last given to SetWidth.
func (m Model) Width() int { return m.width }

// Update flips ShowAll when a key press matches WithToggle's binding. It's
// optional: a Model built without WithToggle ignores every message, and you
// can always set ShowAll yourself instead.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && key.Matches(keyMsg, m.toggle) {
		m.ShowAll = !m.ShowAll
	}
	return m, nil
}

// View renders the bar: one line while ShowAll is false, otherwise the full
// multi-column view. contextual bindings follow the toggle and global ones,
// and are meant to change from render to render as focus moves.
//
// Returns "" when no width has been set or nothing is left to show, in which
// case the bar occupies no rows at all.
func (m Model) View(contextual ...key.Binding) string {
	bindings := m.bindings(contextual)
	if len(bindings) == 0 || m.width <= 0 {
		return ""
	}
	if !m.ShowAll {
		return m.help.ShortHelpView(bindings)
	}
	return m.help.FullHelpView(m.columns(bindings))
}

// Height is the number of rows View would take for the same bindings. Use it
// to work out how much room is left for the rest of your UI. View is the
// source of truth: this is exactly lipgloss.Height of its output, so the two
// can't disagree about where the bar begins.
func (m Model) Height(contextual ...key.Binding) int {
	view := m.View(contextual...)
	if view == "" {
		return 0
	}
	return lipgloss.Height(view)
}

// bindings is the full ordered list: the toggle first (so the way to expand
// the bar always leads), then global, then contextual. Disabled bindings are
// dropped here so column reflow never budgets width for something
// help.FullHelpView will skip anyway.
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

// columns reflows bindings into as many columns as fit within the bar's
// width, which is the same as using as few rows as possible. It walks row
// counts upward and returns the first arrangement that fits, so the result is
// the widest (fewest-rows) layout the width allows.
//
// help.FullHelpView fills each column top to bottom, so a group is a column,
// not a row.
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
	// Everything in one column: the narrowest arrangement possible. It may
	// still overflow, in which case help.FullHelpView truncates as usual.
	return chunkColumns(bindings, len(bindings))
}

// renderedWidth measures what FullHelpView would actually produce for
// groups, by asking it - rather than reimplementing its column and separator
// arithmetic here, which would silently drift the moment upstream changes a
// style or a separator.
//
// The width is zeroed first because that's what disables FullHelpView's own
// truncation (see its shouldAddItem): at width 0 it lays every column out in
// full, which is the untruncated width this needs to measure.
func (m Model) renderedWidth(groups [][]key.Binding) int {
	unbounded := m.help
	unbounded.SetWidth(0)
	return lipgloss.Width(unbounded.FullHelpView(groups))
}

// chunkColumns slices bindings into consecutive groups of at most rows each.
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
