package layout

import tea "charm.land/bubbletea/v2"

// Route implements Navigable: attempts to deliver msg to the leaf
// identified by target. Checks this Model's own direct leaves first, then
// recurses into any leaf whose model is itself Navigable (an embedded
// layout.Model), depth-first. Returns handled=false without touching
// anything if target isn't found anywhere in this (sub)tree - the caller
// (Model.Update, or a parent Model's own Route) is responsible for treating
// that as "silently ignore."
func (m Model) Route(target string, msg tea.Msg) (bool, tea.Cmd) {
	if n, ok := findNode(m.root, target); ok && n.leaf {
		updated, cmd := n.model.Update(msg)
		n.model = updated
		return true, cmd
	}

	var (
		handled bool
		cmd     tea.Cmd
	)
	walk(m.root, func(n *Node) {
		if handled {
			return
		}
		if nav, ok := n.model.(Navigable); ok {
			if h, c := nav.Route(target, msg); h {
				n.model = nav
				handled, cmd = true, c
			}
		}
	})
	return handled, cmd
}

// Focus implements Navigable: moves this (sub)tree's focus straight to id,
// wherever it is - a direct leaf, or nested inside a leaf's own Navigable.
// Unlike MoveFocus (one geometric step in a direction), this is "jump to
// this specific pane." When id lives inside a nested Navigable, that
// child's own internal focus is set first, and then this Model's own focus
// is brought to the leaf hosting it too, so the whole chain agrees on what's
// focused (required for FocusMsg/BlurMsg propagation and for MoveFocus's
// "ask the focused child first" rule to keep working afterward). That outer
// step re-notifies the leaf hosting the nested tree regardless of whether
// the nested Focus call already notified id directly - the two can't always
// be told apart cheaply (id might have already been that subtree's
// untouched default focus, which never got an initial FocusMsg at all, see
// Model.Init) - so id's pane may occasionally see FocusMsg twice for one
// real transition. Delivery is at-least-once, not exactly-once: a Pane
// should treat FocusMsg/BlurMsg as idempotent, the same way it would have
// to tolerate a redundant terminal focus event.
func (m Model) Focus(id string) (bool, tea.Cmd) {
	if _, ok := m.leafRect(id); ok {
		return true, m.setFocus(id)
	}

	var (
		handled  bool
		innerCmd tea.Cmd
		outerID  string
	)
	walk(m.root, func(n *Node) {
		if handled {
			return
		}
		if nav, ok := n.model.(Navigable); ok {
			if h, c := nav.Focus(id); h {
				n.model = nav
				handled, innerCmd, outerID = true, c, n.id
			}
		}
	})
	if !handled {
		return false, nil
	}
	return true, tea.Batch(m.setFocus(outerID), innerCmd)
}

// sourceIsFocused reports whether id is the leaf currently focused
// somewhere along this (sub)tree's active focus chain: either this Model's
// own focused leaf, or - recursively - whatever's focused inside that leaf
// if it's itself a nested Navigable. Used to authorize RequestFocusMsg: only
// the pane that genuinely holds focus right now, at whatever depth, is
// allowed to redirect focus elsewhere. A blurred pane reaching this code
// (e.g. reacting to a SendMsg while in the background) is correctly refused
// since it can never appear on the active chain.
func (m Model) sourceIsFocused(id string) bool {
	if m.state.id == id {
		return true
	}
	focused, ok := findNode(m.root, m.state.id)
	if !ok {
		return false
	}
	checker, ok := focused.model.(interface{ sourceIsFocused(string) bool })
	if !ok {
		return false
	}
	return checker.sourceIsFocused(id)
}

// handleRequestFocus honors a RequestFocusMsg only once sourceIsFocused
// clears its Source, then resolves Target the same way Focus does. An
// unauthorized Source, or an unknown Target, is silently ignored.
func (m Model) handleRequestFocus(msg RequestFocusMsg) tea.Cmd {
	if !m.sourceIsFocused(msg.Source) {
		return nil
	}
	_, cmd := m.Focus(msg.Target)
	return cmd
}
