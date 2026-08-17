package layout

import (
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// buildNested returns an outer Model with two direct leaves: "sidebar" and
// "inner-root", the latter being a nested layout.Model (itself split into
// "inner-a"/"inner-b") embedded as an ordinary Pane. Both models are sized
// before being handed back so their leaf registries are populated.
func buildNested(t *testing.T) (outer Model, sidebar *stubPane, innerA, innerB *stubPane) {
	t.Helper()
	sidebar = newStub()
	innerA = newStub()
	innerB = newStub()

	// inner is deliberately built WITHOUT AsRoot(): it's embedded, so it
	// must not act focused on its own until outer actually focuses the
	// leaf that hosts it (see Model.Init's asRoot guard).
	inner := New(HSplit(0.5, Leaf("inner-a", innerA), Leaf("inner-b", innerB)))

	root := HSplit(0.3, Leaf("sidebar", sidebar), Leaf("inner-root", inner))
	outer = New(root, AsRoot())
	outer.Init()
	updated, _ := outer.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	outer = updated.(Model)
	return outer, sidebar, innerA, innerB
}

// Regression test: a nested Model used to fire its own initial FocusMsg to
// its first leaf unconditionally on Init, regardless of whether the outer
// tree's actual initial focus ever lands on the leaf hosting it - so a
// leaf buried in a subtree that isn't even initially focused would still
// show up as focused, alongside whatever the outer tree really focused.
func TestEmbeddedModelDoesNotSelfFocusOnInit(t *testing.T) {
	outer, sidebar, innerA, innerB := buildNested(t)

	if sidebar.focusN != 1 {
		t.Fatalf("sidebar.focusN = %d, want 1 (it's the outer tree's real initial focus)", sidebar.focusN)
	}
	if innerA.focusN != 0 || innerB.focusN != 0 {
		t.Fatalf("innerA.focusN=%d innerB.focusN=%d, want 0/0: the nested tree isn't focused yet", innerA.focusN, innerB.focusN)
	}
	if outer.state.id != "sidebar" {
		t.Fatalf("outer focusedID = %q, want %q", outer.state.id, "sidebar")
	}
}

func TestSendMsgReachesNestedLeaf(t *testing.T) {
	type payload struct{ n int }
	outer, _, innerA, _ := buildNested(t)

	updated, _ := outer.Update(SendMsg{Target: "inner-a", Msg: payload{n: 7}})
	_ = updated.(Model)

	if got, ok := innerA.last().(payload); !ok || got.n != 7 {
		t.Fatalf("inner-a should have received payload{7}, got %#v", innerA.last())
	}
}

func TestFocusJumpsIntoNestedSubtreeAndUpdatesOuterFocus(t *testing.T) {
	outer, _, _, innerB := buildNested(t)
	// Outer's own focus starts on "sidebar" (first leaf, depth-first).

	handled, _ := outer.Focus("inner-b")
	if !handled {
		t.Fatal("Focus(\"inner-b\") should have been handled")
	}
	if outer.state.id != "inner-root" {
		t.Fatalf("outer focusedID = %q, want %q (the leaf hosting the nested tree)", outer.state.id, "inner-root")
	}
	// At-least-once, not exactly-once (see Focus's doc comment): the outer
	// leaf's own re-notification can duplicate the nested tree's own
	// dispatch, so only assert inner-b actually got notified, not a count.
	if innerB.focusN < 1 {
		t.Fatalf("inner-b.focusN = %d, want at least 1", innerB.focusN)
	}
}

func TestRequestFocusAuthorizedThroughNestedChain(t *testing.T) {
	outer, _, innerA, _ := buildNested(t)

	// Move outer focus onto the nested subtree, and its own internal focus
	// onto inner-a, so inner-a is genuinely the focused leaf end-to-end.
	outer.Focus("inner-a")
	if innerA.focusN != 1 {
		t.Fatalf("inner-a.focusN = %d, want 1 before the request", innerA.focusN)
	}

	updated, _ := outer.Update(RequestFocusMsg{Source: "inner-a", Target: "sidebar"})
	outer = updated.(Model)

	if outer.state.id != "sidebar" {
		t.Fatalf("focusedID = %q, want %q: inner-a is genuinely focused, its request should be honored", outer.state.id, "sidebar")
	}
}

func TestRequestFocusFromNonFocusedNestedLeafIsIgnored(t *testing.T) {
	outer, _, _, innerB := buildNested(t)
	// Outer focus is on "sidebar"; the nested tree isn't even the focused
	// branch, so nothing inside it - including inner-b - is authorized.
	_ = innerB

	updated, _ := outer.Update(RequestFocusMsg{Source: "inner-b", Target: "sidebar"})
	outer = updated.(Model)

	if outer.state.id != "sidebar" {
		t.Fatalf("focusedID = %q, want unchanged %q", outer.state.id, "sidebar")
	}
}

func TestMoveFocusDelegatesToNestedBeforeGeometry(t *testing.T) {
	outer, _, _, innerB := buildNested(t)
	outer.Focus("inner-a")

	if ok := outer.MoveFocus(FocusRight); !ok {
		t.Fatal("MoveFocus(FocusRight) should have been handled by the nested tree (inner-a -> inner-b)")
	}
	if innerB.focusN != 1 {
		t.Fatalf("inner-b.focusN = %d, want 1 (moved within the nested tree)", innerB.focusN)
	}
	if outer.state.id != "inner-root" {
		t.Fatalf("outer focusedID changed to %q, should have stayed on the nested leaf", outer.state.id)
	}

	// Now at inner-b, the nested tree's own rightmost leaf: pressing right
	// again must bubble up and move the OUTER focus instead.
	if ok := outer.MoveFocus(FocusRight); ok {
		// buildNested's outer split is only sidebar|inner-root left-to-right,
		// so there's nothing further right at the outer level either -
		// MoveFocus should report false all the way up.
		t.Fatalf("expected no further neighbor to the right at either level")
	}
}

func TestHelpBindingsDelegatesThroughNesting(t *testing.T) {
	outer, _, innerA, _ := buildNested(t)
	outer.Focus("inner-a")

	want := key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "do x"))
	innerA.help = []key.Binding{want}

	got := outer.HelpBindings()
	if len(got) != 1 || got[0].Help().Key != "x" {
		t.Fatalf("HelpBindings() = %#v, want the focused inner leaf's bindings", got)
	}
}
