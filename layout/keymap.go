package layout

import "charm.land/bubbles/v2/key"

// KeyMap holds the bindings Model itself reacts to: directional focus
// movement and toggling the help bar. A pane's own bindings are separate
// (see HelpProvider) - these are only the ones layout intercepts before a
// key ever reaches the focused pane.
type KeyMap struct {
	FocusLeft  key.Binding
	FocusRight key.Binding
	FocusUp    key.Binding
	FocusDown  key.Binding
	ToggleHelp key.Binding

	// ShowFocusInShortHelp also lists ctrl+hjkl on the short help line, not
	// just the full one. Off by default (see DefaultKeyMap): the four
	// bindings crowd a single line for little gain, since they're always
	// one '?' away in the full view regardless. Set it on a KeyMap passed
	// to WithKeyMap to opt back in.
	ShowFocusInShortHelp bool
}

// DefaultKeyMap returns the standard tmux/vim-style bindings: ctrl+h/j/k/l
// to move focus, ? to toggle the help bar.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		FocusLeft: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "focus left"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "focus right"),
		),
		FocusUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "focus up"),
		),
		FocusDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "focus down"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// focusBindings returns the four directional bindings in reading order.
func (k KeyMap) focusBindings() []key.Binding {
	return []key.Binding{k.FocusLeft, k.FocusDown, k.FocusUp, k.FocusRight}
}

// ShortHelp implements help.KeyMap so KeyMap can be fed to bubbles/help
// directly for layout's own controls. ToggleHelp always leads - see
// helpKeyMap.ShortHelp, which relies on that to put '?' first in the
// composed bar too.
func (k KeyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{k.ToggleHelp}
	if k.ShowFocusInShortHelp {
		bindings = append(bindings, k.focusBindings()...)
	}
	return bindings
}

// FullHelp implements help.KeyMap. ToggleHelp always leads, same reasoning
// as ShortHelp; the focus bindings show here unconditionally, regardless of
// ShowFocusInShortHelp.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleHelp},
		{k.FocusLeft, k.FocusRight},
		{k.FocusUp, k.FocusDown},
	}
}
