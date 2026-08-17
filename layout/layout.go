// Package layout arranges Pane content in a binary split tree (BSP,
// tmux/i3-style), with spatial ctrl+hjkl focus navigation, message routing
// between panes by id, and a help bar that always reflects whatever pane is
// currently focused, however deep it's nested. It owns geometry, focus and
// routing only - it draws no border and imposes no style: each pane decides
// how to render itself for the size and focus state it's given (see SizeMsg,
// FocusMsg, BlurMsg).
package layout

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/style"
)

// Pane is a leaf's content. It's the same Init/Update/View shape used by
// every other custom component in this repo (see tabs.Tab): distinct from
// the real tea.Model, whose View returns tea.View rather than string -
// that's the top level's job (see Run), not a nested pane's.
type Pane interface {
	Init() tea.Cmd
	Update(tea.Msg) (Pane, tea.Cmd)
	View() string
}

// Navigable is what makes a Model composable: embed one layout.Model inside
// another (Leaf(id, innerModel)) and it works transparently, because
// layout.Model itself implements Navigable. ctrl+hjkl first tries
// MoveFocus on whatever's currently focused; SendMsg/RequestFocusMsg reach
// into nested trees via Route/Focus; the help bar drills in via
// FocusedHelp. A pane that isn't itself a layout.Model just doesn't
// implement this, and is treated as an ordinary leaf everywhere.
type Navigable interface {
	Pane
	Leaves() []LeafRect
	MoveFocus(dir FocusDirection) bool
	Route(target string, msg tea.Msg) (handled bool, cmd tea.Cmd)
	Focus(id string) (handled bool, cmd tea.Cmd)
	FocusedHelp() []key.Binding
}

// focusState holds the pieces of Model's state that Navigable's MoveFocus
// and Focus must be able to mutate despite having value receivers - a
// requirement of Model being usable by value as a Leaf's tea.Model and
// still satisfying Navigable when type-asserted back out of that interface.
// Boxed behind a pointer so the mutation persists across every copy of
// Model that shares it.
type focusState struct {
	id string
	// pendingCmd queues cmds produced by BlurMsg/FocusMsg dispatch that
	// happened inside MoveFocus, which - being bool-only, per Navigable -
	// has no return path for them. Drained by the nearest Update that
	// actually returns a tea.Cmd; delivery lags by at most one Update
	// cycle, never user-visible in practice.
	pendingCmd tea.Cmd
}

// Model is a running layout: a Node tree, focus, sizing, and (if AsRoot)
// the help bar. Build one with New.
type Model struct {
	root  *Node
	state *focusState

	leaves []LeafRect

	width, height int

	keyMap   KeyMap
	help     help.Model
	showHelp bool
	asRoot   bool
}

// Option configures a Model at construction. See AsRoot, WithKeyMap.
type Option func(*Model)

// AsRoot marks this Model as the outermost one: only a root Model renders
// its own help bar in View. Off by default, so an embedded Model (see
// Navigable) never shows a duplicate bar - only pass this to the Model
// actually handed to Run/tea.NewProgram.
func AsRoot() Option {
	return func(m *Model) { m.asRoot = true }
}

// WithKeyMap overrides the default ctrl+hjkl/? bindings.
func WithKeyMap(k KeyMap) Option {
	return func(m *Model) { m.keyMap = k }
}

// New builds a Model from root. The first leaf (depth-first, first child
// before second) starts focused.
func New(root *Node, opts ...Option) Model {
	m := Model{
		root:   root,
		state:  &focusState{id: firstLeafID(root)},
		keyMap: DefaultKeyMap(),
		help:   bubbles.NewHelp(),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// program adapts a Model (a Pane, like any other layout leaf content) into
// a real tea.Model for tea.NewProgram: the only place a Model's View needs
// to become a tea.View instead of a string (see Pane's doc comment).
type program struct{ m Model }

func (p program) Init() tea.Cmd { return p.m.Init() }

func (p program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := p.m.Update(msg)
	p.m = updated.(Model)
	return p, cmd
}

func (p program) View() tea.View {
	view := tea.NewView(p.m.View())
	view.AltScreen = true
	return view
}

// Run builds and starts a tea.Program for m (which should have been
// constructed with AsRoot). A convenience for the common standalone-binary
// case; an app assembling layout into a bigger tea.Program of its own can
// wrap m the same way program does above instead.
func Run(m Model, opts ...tea.ProgramOption) error {
	_, err := tea.NewProgram(program{m: m}, opts...).Run()
	return err
}

func firstLeafID(n *Node) string {
	for n != nil && !n.leaf {
		n = n.first
	}
	if n == nil {
		return ""
	}
	return n.id
}

// walk visits every leaf in the tree rooted at n, depth-first, first child
// before second - the same order computeLayout produces, so it's safe to
// rely on for anything that should stay in step with the flat registry.
func walk(n *Node, fn func(*Node)) {
	if n == nil {
		return
	}
	if n.leaf {
		fn(n)
		return
	}
	walk(n.first, fn)
	walk(n.second, fn)
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	walk(m.root, func(n *Node) {
		if cmd := n.model.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	})

	// Only the actual root originates the initial FocusMsg. An embedded
	// Model's own state.id already defaults to its first leaf (see New),
	// but it must stay quiet about it until its parent actually focuses the
	// leaf hosting it - which happens naturally through the ordinary
	// FocusMsg/BlurMsg case in Update, cascading down as deep as needed.
	// Without this guard, every nested Model fires its own initial
	// FocusMsg independently, so a leaf that isn't even the outer tree's
	// initial focus still shows as focused until the first real move.
	if m.asRoot {
		if n, ok := findNode(m.root, m.state.id); ok {
			updated, cmd := n.model.Update(FocusMsg{})
			n.model = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m.resize()

	case SizeMsg:
		// A nested Model receiving its own allocation from a parent
		// layout.Model - equivalent to tea.WindowSizeMsg at the root.
		m.width, m.height = msg.Width, msg.Height
		return m.resize()

	case FocusMsg, BlurMsg:
		// This whole (sub)tree just gained/lost focus at the parent's
		// level: redistribute to whichever of our own leaves is focused,
		// not to the tree "globally" (see Navigable doc).
		return m, m.deliverToFocused(msg)

	case SendMsg:
		_, cmd := m.Route(msg.Target, msg.Msg)
		return m, cmd

	case RequestFocusMsg:
		return m, m.handleRequestFocus(msg)

	case SetPaneMsg:
		return m.SetPane(msg.ID, msg.NewPane)

	case SplitLeafMsg:
		return m.SplitLeaf(msg.ID, msg.Dir, msg.NewID, msg.NewModel, msg.Opts...)

	case CloseLeafMsg:
		return m.CloseLeaf(msg.ID)

	case ResizeMsg:
		return m.Resize(msg.SplitID, msg.Ratio)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.ToggleHelp):
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp
			return m.resize()
		case key.Matches(msg, m.keyMap.FocusLeft):
			m.MoveFocus(FocusLeft)
			return m, m.drainCmd()
		case key.Matches(msg, m.keyMap.FocusRight):
			m.MoveFocus(FocusRight)
			return m, m.drainCmd()
		case key.Matches(msg, m.keyMap.FocusUp):
			m.MoveFocus(FocusUp)
			return m, m.drainCmd()
		case key.Matches(msg, m.keyMap.FocusDown):
			m.MoveFocus(FocusDown)
			return m, m.drainCmd()
		default:
			return m, m.deliverToFocused(msg)
		}

	default:
		return m, m.broadcast(msg)
	}
}

func (m Model) drainCmd() tea.Cmd {
	cmd := m.state.pendingCmd
	m.state.pendingCmd = nil
	return cmd
}

// deliverToFocused sends msg to the currently focused leaf's model only.
// Used for ordinary key presses (only the focused pane should react to
// keyboard input) and for relaying FocusMsg/BlurMsg into a nested subtree.
func (m Model) deliverToFocused(msg tea.Msg) tea.Cmd {
	n, ok := findNode(m.root, m.state.id)
	if !ok {
		return nil
	}
	updated, cmd := n.model.Update(msg)
	n.model = updated
	return cmd
}

// broadcast sends msg to every leaf's model, focused or not - the default
// for anything that isn't one of layout's own reserved message types or a
// key press, so a blurred pane can still receive its own async messages
// (a tick driving a spinner, an HTTP response, ...).
func (m Model) broadcast(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	walk(m.root, func(n *Node) {
		updated, cmd := n.model.Update(msg)
		n.model = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	})
	return tea.Batch(cmds...)
}

// setFocus moves focus to id (assumed already validated as an existing
// leaf), dispatching BlurMsg to the old focus and FocusMsg to the new one,
// and returns the resulting batched cmd. A no-op (nil cmd) if id is already
// focused.
func (m Model) setFocus(id string) tea.Cmd {
	if id == m.state.id {
		return nil
	}
	var cmds []tea.Cmd
	if old, ok := findNode(m.root, m.state.id); ok {
		updated, cmd := old.model.Update(BlurMsg{})
		old.model = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	m.state.id = id
	if n, ok := findNode(m.root, id); ok {
		updated, cmd := n.model.Update(FocusMsg{})
		n.model = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// SplitLeaf splits the leaf identified by id into two: id keeps its
// original pane on one side, a new Leaf(newID, newModel) takes the other,
// joined by a 50/50 Split (override via WithSplitRatio, WithSplitID,
// WithSplitMinimum, WithSplitMaximum). A no-op (m unchanged, nil cmd) if id
// doesn't identify an existing leaf.
func (m Model) SplitLeaf(id string, dir Direction, newID string, newModel Pane, opts ...SplitOption) (Model, tea.Cmd) {
	newRoot, ok := splitLeaf(m.root, id, dir, newID, newModel, opts...)
	if !ok {
		return m, nil
	}
	m.root = newRoot

	var cmds []tea.Cmd
	if cmd := newModel.Init(); cmd != nil {
		cmds = append(cmds, cmd)
	}
	resized, cmd := m.resize()
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return resized, tea.Batch(cmds...)
}

// SetPane replaces the Pane at leaf id with newPane, keeping its place and
// shape in the tree unchanged - unlike SplitLeaf/CloseLeaf, which reshape
// the tree, this only swaps what's rendered at an existing slot (the way an
// app switches its content area between entirely different sub-apps/pages,
// each its own package). newPane is Init'd and immediately told its size
// via SizeMsg (using id's current Rect, which by definition hasn't changed);
// it's also told FocusMsg if id currently holds focus, since the pane it's
// replacing never will. A no-op if id doesn't identify an existing leaf.
func (m Model) SetPane(id string, newPane Pane) (Model, tea.Cmd) {
	n, ok := findNode(m.root, id)
	if !ok || !n.leaf {
		return m, nil
	}
	n.model = newPane

	var cmds []tea.Cmd
	if cmd := newPane.Init(); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if lr, ok := m.leafRect(id); ok {
		updated, cmd := n.model.Update(SizeMsg{ID: id, Width: lr.Rect.W, Height: lr.Rect.H})
		n.model = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if m.state.id == id {
		updated, cmd := n.model.Update(FocusMsg{})
		n.model = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// CloseLeaf removes the leaf identified by id, promoting its sibling to take
// the place of their parent Split. If id currently has focus, focus moves to
// the tree's new first leaf. A no-op if id is the tree's own root (the last
// remaining pane can't be closed this way) or doesn't exist.
func (m Model) CloseLeaf(id string) (Model, tea.Cmd) {
	newRoot, ok := closeLeaf(m.root, id)
	if !ok {
		return m, nil
	}
	m.root = newRoot

	var focusCmd tea.Cmd
	if m.state.id == id {
		focusCmd = m.setFocus(firstLeafID(m.root))
	}

	resized, resizeCmd := m.resize()
	return resized, tea.Batch(focusCmd, resizeCmd)
}

// Resize sets the ratio of the first child of the Split identified by
// splitID (only reachable if it was given one, via (*Node).WithID or
// WithSplitID). A no-op if splitID isn't found or identifies a Leaf.
func (m Model) Resize(splitID string, ratio float64) (Model, tea.Cmd) {
	n, ok := findNode(m.root, splitID)
	if !ok || n.leaf {
		return m, nil
	}
	n.ratio = ratio
	return m.resize()
}

// resize recomputes the flat leaf registry from the current root/width/
// height and dispatches SizeMsg to every leaf whose Rect actually changed
// (not just the ones directly touched by whatever triggered this - a
// sibling's size can shift too). Shared by every path that can change
// geometry: tea.WindowSizeMsg, SizeMsg (nested), SplitLeaf, CloseLeaf,
// Resize, and toggling the help bar (which changes how much height the tree
// itself gets).
func (m Model) resize() (Model, tea.Cmd) {
	if m.width <= 0 || m.height <= 0 {
		return m, nil
	}
	m.help.SetWidth(m.width)

	rect := m.treeRect()
	newLeaves := computeLayout(m.root, rect)

	old := make(map[string]Rect, len(m.leaves))
	for _, lr := range m.leaves {
		old[lr.ID] = lr.Rect
	}

	var cmds []tea.Cmd
	for _, lr := range newLeaves {
		if prev, ok := old[lr.ID]; ok && prev == lr.Rect {
			continue
		}
		if n, ok := findNode(m.root, lr.ID); ok {
			updated, cmd := n.model.Update(SizeMsg{ID: lr.ID, Width: lr.Rect.W, Height: lr.Rect.H})
			n.model = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	m.leaves = newLeaves

	return m, tea.Batch(cmds...)
}

// treeRect is the region left for the split tree once the help bar (if
// AsRoot) has taken its share of the height. resize and View both go
// through this so they can never disagree about where the tree ends and
// the help bar begins.
func (m Model) treeRect() Rect {
	h := m.height - m.helpHeight()
	if h < 0 {
		h = 0
	}
	return Rect{W: m.width, H: h}
}

func (m Model) helpHeight() int {
	if !m.asRoot {
		return 0
	}
	if rendered := m.renderHelp(); rendered != "" {
		return lipgloss.Height(rendered)
	}
	return 0
}

func (m Model) renderHelp() string {
	return m.help.View(helpKeyMap{pane: m.FocusedHelp(), own: m.keyMap, width: m.width})
}

// Leaves implements Navigable.
func (m Model) Leaves() []LeafRect {
	return m.leaves
}

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	rect := m.treeRect()
	tree := renderNode(m.root, rect.W, rect.H)

	if !m.asRoot {
		return tree
	}
	help := m.renderHelp()
	if help == "" {
		return tree
	}
	return lipgloss.JoinVertical(lipgloss.Left, tree, help)
}

// renderNode mirrors computeLayout's own allocation (same resolveSize calls
// on the same w/h at each level), so what's rendered here always matches
// the SizeMsg values leaves were already told via resize.
func renderNode(n *Node, w, h int) string {
	if n.leaf {
		// A misbehaving Pane that renders wider/taller than the SizeMsg it
		// was given would otherwise desync every ancestor Join*, so clip it
		// here rather than trusting the contract to hold. MaxWidth/MaxHeight
		// truncate via ansi.Truncate internally, so this stays escape-code
		// safe instead of mangling a Pane's own styling mid-sequence.
		return lipgloss.NewStyle().MaxWidth(w).MaxHeight(h).Render(n.model.View())
	}
	if n.dir == Horizontal {
		w1 := resolveSize(n, w)
		return lipgloss.JoinHorizontal(lipgloss.Top,
			renderNode(n.first, w1, h),
			renderNode(n.second, w-w1, h),
		)
	}
	h1 := resolveSize(n, h)
	return lipgloss.JoinVertical(lipgloss.Left,
		renderNode(n.first, w, h1),
		renderNode(n.second, w, h-h1),
	)
}

// Bordered is an optional helper for panes that want the common look: a
// border that follows focus (style.S.Primary focused, style.S.Subtle
// blurred) drawn with the configured BorderType. Not required - a pane
// that wants something else, or nothing, just doesn't call this. Renders
// to exactly w by h, border included, as View() must (see FocusMsg/BlurMsg
// and SizeMsg docs).
func Bordered(focused bool, w, h int, content string) string {
	color := style.S.Subtle
	if focused {
		color = style.S.Primary
	}
	// lipgloss's Width/Height already count the border as part of the box
	// (they subtract its size internally before sizing the content), so w
	// and h go straight through - no manual -2 here, or the box comes out
	// two cells smaller than asked in both dimensions.
	return lipgloss.NewStyle().
		Border(style.S.BorderType).
		BorderForeground(color).
		Width(max(w, 0)).
		Height(max(h, 0)).
		Render(content)
}
