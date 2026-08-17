package layout

import tea "charm.land/bubbletea/v2"

// SizeMsg tells a pane the exact width/height it's been allocated, and its
// own id (so it can identify itself as the Source of a later
// RequestFocusMsg). Sent to every leaf whose Rect changed, whenever the tree
// is resized, split, closed or reshaped - never assume a pane can deduce its
// size any other way.
type SizeMsg struct {
	ID            string
	Width, Height int
}

// FocusMsg tells a pane it just gained keyboard focus. Named distinctly from
// bubbletea's own tea.FocusMsg (which is about terminal focus, not pane
// focus) to avoid any confusion between the two.
type FocusMsg struct{}

// BlurMsg tells a pane it just lost keyboard focus. See FocusMsg.
type BlurMsg struct{}

// SendMsg delivers Msg to the pane identified by Target, wherever it lives
// in the tree (including inside a nested layout.Model), regardless of which
// pane currently has focus. A pane never holds a reference to another, so
// this is how they talk: return
//
//	func() tea.Msg { return layout.SendMsg{Target: "editor", Msg: myMsg{}} }
//
// as a tea.Cmd from Update. An unknown Target is silently ignored.
type SendMsg struct {
	Target string
	Msg    tea.Msg
}

// RequestFocusMsg asks the layout to move keyboard focus to Target. Only
// honored when Source is the id of the pane that currently holds focus (a
// blurred pane can send other panes messages via SendMsg, but can't move
// focus itself or on anyone's behalf) - fill Source from the ID a pane
// learned via SizeMsg:
//
//	func() tea.Msg { return layout.RequestFocusMsg{Source: p.id, Target: "content"} }
//
// A mismatched Source, or an unknown Target, is silently ignored.
type RequestFocusMsg struct {
	Source string
	Target string
}

// SetPaneMsg is the message form of Model.SetPane, for a pane that wants to
// swap another leaf's content without holding a reference to its Model
// (which it never does).
type SetPaneMsg struct {
	ID      string
	NewPane Pane
}

// SplitLeafMsg is the message form of Model.SplitLeaf, for a pane that wants
// to reshape the tree without holding a reference to its Model (which it
// never does). See Model.SplitLeaf for parameters.
type SplitLeafMsg struct {
	ID       string
	Dir      Direction
	NewID    string
	NewModel Pane
	Opts     []SplitOption
}

// CloseLeafMsg is the message form of Model.CloseLeaf.
type CloseLeafMsg struct {
	ID string
}

// ResizeMsg is the message form of Model.Resize: sets the ratio of the
// first child of the Split identified by SplitID (a Split only reachable by
// id if it was given one via WithID or WithSplitID).
type ResizeMsg struct {
	SplitID string
	Ratio   float64
}
