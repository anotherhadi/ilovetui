package layout

import (
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

// HelpProvider is how a pane opts into the help bar: implement it and
// return whatever bindings you want shown while you're focused. A pane
// that doesn't implement it just contributes nothing - the help bar still
// shows layout's own controls (ctrl+hjkl, ?), it never disappears entirely.
type HelpProvider interface {
	HelpBindings() []key.Binding
}

// HelpBindings implements HelpProvider by delegating to FocusedHelp, so a
// Model embedded as a HelpProvider behaves identically to one consulted as
// a Navigable.
func (m Model) HelpBindings() []key.Binding {
	return m.FocusedHelp()
}

// FocusedHelp implements Navigable: the bindings for whatever pane
// currently has focus, drilling into a nested Navigable automatically until
// it reaches the real pane at the bottom. Returns nil if the focused pane
// (at any depth) implements neither Navigable nor HelpProvider.
func (m Model) FocusedHelp() []key.Binding {
	focused, ok := findNode(m.root, m.state.id)
	if !ok {
		return nil
	}
	if nav, ok := focused.model.(Navigable); ok {
		return nav.FocusedHelp()
	}
	if hp, ok := focused.model.(HelpProvider); ok {
		return hp.HelpBindings()
	}
	return nil
}

// helpKeyMap adapts a focused pane's flat HelpProvider bindings plus
// layout's own KeyMap into the shape bubbles/help.KeyMap expects.
type helpKeyMap struct {
	pane  []key.Binding
	own   KeyMap
	width int
}

// ShortHelp implements help.KeyMap. own leads (its ToggleHelp is always
// first, see KeyMap.ShortHelp), the focused pane's own bindings follow.
func (h helpKeyMap) ShortHelp() []key.Binding {
	return append(append([]key.Binding{}, h.own.ShortHelp()...), h.pane...)
}

// FullHelp implements help.KeyMap. Rather than the fixed grouping ShortHelp
// mirrors, it flattens every binding into one ordered list - own first, so
// ToggleHelp lands in the first column's first row, then the pane's - and
// re-flows it into as many columns as fit within width, maximizing columns
// to minimize the number of rows the full help view takes.
func (h helpKeyMap) FullHelp() [][]key.Binding {
	all := append(append([]key.Binding{}, flattenGroups(h.own.FullHelp())...), h.pane...)
	return flowColumns(all, h.width)
}

func flattenGroups(groups [][]key.Binding) []key.Binding {
	var flat []key.Binding
	for _, g := range groups {
		flat = append(flat, g...)
	}
	return flat
}

// flowColumns arranges bindings into as many columns as fit within width
// without overflowing. It fills each column top-to-bottom before moving to
// the next, which is what bubbles/help.FullHelpView expects: one inner
// slice per column, rendered as a vertical stack.
func flowColumns(bindings []key.Binding, width int) [][]key.Binding {
	enabled := make([]key.Binding, 0, len(bindings))
	for _, kb := range bindings {
		if kb.Enabled() {
			enabled = append(enabled, kb)
		}
	}
	if len(enabled) == 0 {
		return nil
	}
	if width <= 0 {
		return [][]key.Binding{enabled}
	}

	for rows := 1; rows <= len(enabled); rows++ {
		groups := chunkRows(enabled, rows)
		if columnsWidth(groups) <= width {
			return groups
		}
	}
	return [][]key.Binding{enabled}
}

// chunkRows splits bindings into groups of at most rows items each, filling
// each group before moving to the next - the column-major order help's
// FullHelpView renders (first group is the leftmost column).
func chunkRows(bindings []key.Binding, rows int) [][]key.Binding {
	var groups [][]key.Binding
	for i := 0; i < len(bindings); i += rows {
		end := min(i+rows, len(bindings))
		groups = append(groups, bindings[i:end])
	}
	return groups
}

// columnsWidth mirrors bubbles/help.FullHelpView's own width accounting:
// each column is as wide as its longest key plus a space plus its longest
// description, columns separated by FullSeparator's width (4 cells, "    ").
func columnsWidth(groups [][]key.Binding) int {
	const separatorWidth = 4
	total := 0
	for i, group := range groups {
		if i > 0 {
			total += separatorWidth
		}
		var keyWidth, descWidth int
		for _, kb := range group {
			keyWidth = max(keyWidth, lipgloss.Width(kb.Help().Key))
			descWidth = max(descWidth, lipgloss.Width(kb.Help().Desc))
		}
		total += keyWidth + 1 + descWidth
	}
	return total
}
