package layout

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestModel(t *testing.T, root *Node) (Model, map[string]*stubPane) {
	t.Helper()
	panes := map[string]*stubPane{}
	walk(root, func(n *Node) {
		if s, ok := n.model.(*stubPane); ok {
			panes[n.id] = s
		}
	})
	m := New(root, AsRoot())
	m.Init()
	return m, panes
}

func TestInitFocusesFirstLeaf(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)

	if m.state.id != "a" {
		t.Fatalf("focusedID = %q, want %q", m.state.id, "a")
	}
	if panes["a"].focusN != 1 {
		t.Fatalf("a.focusN = %d, want 1", panes["a"].focusN)
	}
	if panes["b"].focusN != 0 {
		t.Fatalf("b.focusN = %d, want 0", panes["b"].focusN)
	}
}

func TestWindowSizeDispatchesSizeMsg(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	// newTestModel builds with AsRoot(), so the tree gets 40 minus the
	// 1-row help bar - see treeRect.
	if panes["a"].w != 50 || panes["a"].h != 39 {
		t.Fatalf("a size = %dx%d, want 50x39", panes["a"].w, panes["a"].h)
	}
	if panes["b"].w != 50 || panes["b"].h != 39 {
		t.Fatalf("b size = %dx%d, want 50x39", panes["b"].w, panes["b"].h)
	}
	if got := len(m.Leaves()); got != 2 {
		t.Fatalf("len(Leaves()) = %d, want 2", got)
	}
}

func TestCtrlLMovesFocusAndDispatchesBlurFocus(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyPressMsg{Code: 'l', Mod: tea.ModCtrl, Text: ""})
	m = updated.(Model)

	if m.state.id != "b" {
		t.Fatalf("focusedID = %q, want %q", m.state.id, "b")
	}
	if panes["a"].blurN != 1 {
		t.Fatalf("a.blurN = %d, want 1", panes["a"].blurN)
	}
	if panes["b"].focusN != 1 {
		t.Fatalf("b.focusN = %d, want 1", panes["b"].focusN)
	}
}

func TestKeyPressOnlyReachesFocusedPane(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyPressMsg{Text: "x"})
	_ = updated.(Model)

	if _, ok := panes["a"].last().(tea.KeyPressMsg); !ok {
		t.Fatalf("focused pane a should have received the key press, got %#v", panes["a"].last())
	}
	if _, ok := panes["b"].last().(tea.KeyPressMsg); ok {
		t.Fatalf("blurred pane b should not have received the key press")
	}
}

func TestSplitLeafAddsAndSizesNewPane(t *testing.T) {
	root := Leaf("a", newStub())
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	newPane := newStub()
	m, _ = m.SplitLeaf("a", Horizontal, "b", newPane)

	if got := len(m.Leaves()); got != 2 {
		t.Fatalf("len(Leaves()) = %d, want 2", got)
	}
	if newPane.w == 0 || newPane.h == 0 {
		t.Fatalf("new pane never received a SizeMsg: w=%d h=%d", newPane.w, newPane.h)
	}
	if panes["a"].w != 50 {
		t.Fatalf("a.w = %d, want 50 after 50/50 split", panes["a"].w)
	}
}

func TestSetPaneSwapsContentKeepingShape(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	// "a" is focused (first leaf).

	replacement := newStub()
	m, _ = m.SetPane("a", replacement)

	if got := len(m.Leaves()); got != 2 {
		t.Fatalf("len(Leaves()) = %d, want 2 (SetPane must not reshape the tree)", got)
	}
	if replacement.w != 50 || replacement.h != 39 {
		t.Fatalf("replacement size = %dx%d, want 50x39 (a's existing Rect)", replacement.w, replacement.h)
	}
	if replacement.focusN != 1 {
		t.Fatalf("replacement.focusN = %d, want 1: \"a\" currently holds focus", replacement.focusN)
	}
	if panes["b"].w != 50 {
		t.Fatalf("b.w = %d, want unchanged 50", panes["b"].w)
	}
}

func TestSetPaneOnBlurredLeafDoesNotFocus(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, _ := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	// "b" is blurred (only "a" is focused initially).

	replacement := newStub()
	m, _ = m.SetPane("b", replacement)

	if replacement.focusN != 0 {
		t.Fatalf("replacement.focusN = %d, want 0: \"b\" isn't focused", replacement.focusN)
	}
	if replacement.w != 50 || replacement.h != 39 {
		t.Fatalf("replacement size = %dx%d, want 50x39", replacement.w, replacement.h)
	}
}

func TestSetPaneUnknownIDIsNoop(t *testing.T) {
	root := Leaf("a", newStub())
	m, _ := newTestModel(t, root)

	m2, cmd := m.SetPane("nope", newStub())
	if cmd != nil {
		t.Fatalf("expected nil cmd for an unknown SetPane id, got %v", cmd)
	}
	if len(m2.Leaves()) != len(m.Leaves()) {
		t.Fatalf("tree changed after a no-op SetPane")
	}
}

func TestCloseLeafReassignsFocus(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, _ := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	// focus is on "a"; closing it must move focus to what remains.

	m, _ = m.CloseLeaf("a")

	if got := len(m.Leaves()); got != 1 {
		t.Fatalf("len(Leaves()) = %d, want 1", got)
	}
	if m.state.id != "b" {
		t.Fatalf("focusedID = %q, want %q after closing the focused leaf", m.state.id, "b")
	}
}

func TestCloseLeafRootIsNoop(t *testing.T) {
	root := Leaf("only", newStub())
	m, _ := newTestModel(t, root)

	m2, cmd := m.CloseLeaf("only")
	if cmd != nil {
		t.Fatalf("expected nil cmd closing the tree's only leaf, got %v", cmd)
	}
	if got := len(m2.Leaves()); got != len(m.Leaves()) {
		t.Fatalf("tree changed after a no-op close")
	}
}

func TestResizeChangesRatio(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub())).WithID("split")
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	m, _ = m.Resize("split", 0.8)

	if panes["a"].w != 80 {
		t.Fatalf("a.w = %d, want 80 after Resize to 0.8", panes["a"].w)
	}
	if panes["b"].w != 20 {
		t.Fatalf("b.w = %d, want 20 after Resize to 0.8", panes["b"].w)
	}
}

func TestSendMsgDeliversToTarget(t *testing.T) {
	type payload struct{ n int }
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)
	// a is focused via Init's initial FocusMsg; b never received anything yet.

	updated, _ := m.Update(SendMsg{Target: "b", Msg: payload{n: 42}})
	m = updated.(Model)

	if got, ok := panes["b"].last().(payload); !ok || got.n != 42 {
		t.Fatalf("b should have received payload{42}, got %#v", panes["b"].last())
	}
	for _, msg := range panes["a"].msgs {
		if _, ok := msg.(payload); ok {
			t.Fatalf("a should not have received the SendMsg meant for b")
		}
	}
}

func TestSendMsgUnknownTargetIsIgnored(t *testing.T) {
	root := Leaf("a", newStub())
	m, _ := newTestModel(t, root)

	updated, cmd := m.Update(SendMsg{Target: "nope", Msg: struct{}{}})
	_ = updated.(Model)
	if cmd != nil {
		t.Fatalf("expected nil cmd for an unknown SendMsg target, got %v", cmd)
	}
}

func TestRequestFocusFromFocusedPaneSucceeds(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	m, panes := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	updated, _ = m.Update(RequestFocusMsg{Source: "a", Target: "b"})
	m = updated.(Model)

	if m.state.id != "b" {
		t.Fatalf("focusedID = %q, want %q", m.state.id, "b")
	}
	if panes["b"].focusN != 1 {
		t.Fatalf("b.focusN = %d, want 1", panes["b"].focusN)
	}
}

func TestRequestFocusFromBlurredPaneIsIgnored(t *testing.T) {
	root := HSplit(0.5, Leaf("a", newStub()), Leaf("b", newStub()))
	third := Leaf("c", newStub())
	_ = third
	m, _ := newTestModel(t, root)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	// "a" is focused; "b" (blurred) tries to redirect focus to itself.

	updated, _ = m.Update(RequestFocusMsg{Source: "b", Target: "b"})
	m = updated.(Model)

	if m.state.id != "a" {
		t.Fatalf("focusedID = %q, want %q (unauthorized request must be ignored)", m.state.id, "a")
	}
}
