# layout

Arrange panes in a binary split tree (BSP, tmux/i3-style), navigate between them with
spatial `ctrl+hjkl`, and get a help bar that always reflects whatever's focused - no
matter how deep the tree is, or how many of those splits are themselves other `layout`
trees nested inside a pane.

`layout` only owns geometry, focus and routing. It draws no border and imposes no
style: every pane decides how to render itself for the size and focus state it's given.

## Concepts

Three things:

- **`Pane`** is your content: `Init() tea.Cmd`, `Update(tea.Msg) (Pane, tea.Cmd)`,
  `View() string` - the same shape used by every other custom component in this repo.
- **`Node`** is the shape of the tree. `Leaf(id, pane)` is a slot holding one `Pane`,
  addressed everywhere else (`SendMsg`, `RequestFocusMsg`, `SplitLeaf`, `CloseLeaf`,
  `Resize`) by that `id`. `Split`/`HSplit`/`VSplit` divide space between two child
  `Node`s.
- **`Model`** is the running layout: a `Node` tree plus focus, sizing and the help bar.
  Build one with `layout.New(root, layout.AsRoot())`.

## Quick start

```go
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

func main() {
	root := layout.HSplit(0.3,
		layout.Leaf("sidebar", newSidebarPane()),
		layout.Leaf("content", newContentPane()),
	)
	m := layout.New(root, layout.AsRoot())

	if err := layout.Run(m); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

That's a working two-pane app: `ctrl+h`/`ctrl+l` move focus between "sidebar" and
"content", `?` toggles the help bar, everything resizes with the terminal.
`layout.Run(m, opts ...tea.ProgramOption)` wraps `m` for `tea.NewProgram` and runs it,
alt-screen included. `layout` itself reserves no quit key - deciding how (and whether)
the app quits is the host's call, same as any other app-level policy; see the examples
for the usual `ctrl+c`/`q` pattern.

## Writing a pane

```go
type Pane interface {
	Init() tea.Cmd
	Update(tea.Msg) (Pane, tea.Cmd)
	View() string
}
```

`layout` tells a pane its size and focus state entirely through messages - never assume
either any other way:

- **`SizeMsg{ID, Width, Height}`** whenever the pane's allocated space changes (first
  layout, terminal resize, a sibling split/close/resize...). `ID` is the pane's own id,
  learned here so it can later fill `RequestFocusMsg.Source` (see below).
- **`FocusMsg{}`** / **`BlurMsg{}`** whenever the pane gains or loses keyboard focus.
  Named distinctly from bubbletea's own `tea.FocusMsg`/`tea.BlurMsg`, which are about
  terminal focus, not pane focus.

Two rules to actually get right:

- **`View()` must render exactly the last width/height you were told.** `layout`
  composes panes side by side with `lipgloss.JoinHorizontal`/`JoinVertical` - if a pane
  renders the wrong size, the whole layout visibly misaligns.
- **You own your own chrome.** `layout` never draws a border. `layout.Bordered(focused,
  w, h, content)` covers the common case (border color follows focus, via
  `style.S.Primary`/`style.S.Subtle`) as an optional helper - draw nothing, or draw
  something else entirely, if you want.

## Building the tree

```go
root := layout.HSplit(0.3,
	layout.Leaf("sidebar", newSidebarPane()),
	layout.Leaf("content", newContentPane()),
)
```

- `layout.Leaf(id, pane)` - one pane, addressed by `id` everywhere else in the API.
- `layout.HSplit(ratio, first, second)` / `layout.VSplit(ratio, first, second)` - divide
  space horizontally (side by side) or vertically (stacked), giving `ratio` (0 to 1) to
  `first` and the rest to `second`. `layout.Split(id, dir, ratio, first, second)` is the
  general form if the split itself needs an `id` (see "Resizing at runtime" below).
  Nest freely:

```go
root := layout.HSplit(0.25,
	sidebar,
	layout.VSplit(0.7,
		layout.HSplit(0.5, topLeft, topRight),
		bottom,
	),
)
```

### Sizing

```go
layout.HSplit(0.3, sidebar, content)                          // sidebar gets 30%, content the rest
layout.HSplit(0.3, sidebar, content).WithMinimum(20)           // 30%, but never below 20 cells
layout.HSplit(0.3, sidebar, content).WithMaximum(40)           // 30%, but never above 40 cells
layout.HSplit(0.3, sidebar, content).WithMinimum(20).WithMaximum(20) // fixed at 20 cells
```

`WithMinimum`/`WithMaximum` clamp the resolved size of the split's *first* child (in
cells: columns for a horizontal split, rows for a vertical one). Setting both to the
same value pins it regardless of ratio - the way to get an exact-width sidebar.

## Navigation and the help bar

`ctrl+h`/`ctrl+l`/`ctrl+j`/`ctrl+k` move focus spatially - whichever pane is actually
adjacent in that direction (tmux's `select-pane -L/-D/-U/-R`), not "next in the tree".
`?` toggles the help bar between its short and full form. Both work out of the box.

To customize the keys:

```go
km := layout.DefaultKeyMap()
km.FocusLeft = key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "focus left"))
m := layout.New(root, layout.AsRoot(), layout.WithKeyMap(km))
```

### Making a pane show up in the help bar

Implement `HelpProvider`:

```go
type HelpProvider interface {
	HelpBindings() []key.Binding
}
```

`layout` reads it fresh every render from whichever pane is currently focused, at
whatever depth (see "Composing bigger apps" below) - nothing to push or keep in sync. A
pane that doesn't implement it just contributes nothing; the bar still shows `layout`'s
own controls (`?`, `ctrl+hjkl`), it never disappears.

**Only the outermost `Model` should render a help bar.** Pass `layout.AsRoot()` only to
the one actually handed to `Run`/`tea.NewProgram` - an embedded `Model` (see "Composing
bigger apps") built without it contributes its focused pane's bindings to the outer bar
instead of drawing a second one of its own.

## Talking between panes

A pane never holds a reference to another - it returns a `tea.Cmd` and lets `layout`
deliver it, by id, wherever that id lives in the tree (including inside a nested
`layout.Model` - see "Composing bigger apps"):

```go
// deliver an arbitrary message to another pane's Update, regardless of focus
return p, func() tea.Msg {
	return layout.SendMsg{Target: "content", Msg: pageChangedMsg{page: selected}}
}
```

An unknown `Target` is silently ignored.

### Asking layout to move focus

```go
return p, func() tea.Msg {
	return layout.RequestFocusMsg{Source: p.id, Target: "content"}
}
```

`p.id` is whatever the pane last learned from `SizeMsg.ID`. **Only honored when
`Source` is the pane that currently, genuinely holds focus** - a blurred pane (say,
reacting to a `SendMsg` while in the background) can't redirect focus this way, for
itself or anyone else; only the pane actually focused right now can hand focus off to
another. An unauthorized `Source`, or an unknown `Target`, is silently ignored.

A concrete pattern: a sidebar list drives *and* jumps to a content pane, purely through
messages:

```go
func (p sidebarPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	prevIndex := p.list.Index()
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)

	if p.list.Index() != prevIndex {
		selected := p.list.SelectedItem().(page)
		cmd = tea.Batch(cmd,
			func() tea.Msg { return layout.SendMsg{Target: "content", Msg: pageChangedMsg{selected}} },
			func() tea.Msg { return layout.RequestFocusMsg{Source: p.id, Target: "content"} },
		)
	}
	return p, cmd
}
```

## Reshaping the tree at runtime

```go
m, cmd := m.SplitLeaf("editor", layout.Vertical, "terminal", newTerminalPane())
m, cmd = m.CloseLeaf("terminal")
m, cmd = m.Resize("main-split", 0.6)
m, cmd = m.SetPane("workspace", newSecondPagePane())
```

- **`SplitLeaf(id, dir, newID, newModel, opts ...SplitOption)`** splits the leaf `id`
  into two: `id` keeps its original pane on one side, `Leaf(newID, newModel)` takes the
  other, joined by a 50/50 split by default. Override with `WithSplitID`,
  `WithSplitRatio`, `WithSplitMinimum`, `WithSplitMaximum`.
- **`CloseLeaf(id)`** removes a leaf; its sibling takes the place of their parent split.
  Focus moves elsewhere automatically if `id` was focused. The tree's own last
  remaining leaf can't be closed this way.
- **`Resize(splitID, ratio)`** changes a split's ratio. Only reachable if the split was
  given an id, via `(*Node).WithID` (or `WithSplitID` when it was created by
  `SplitLeaf`) - `HSplit`/`VSplit` leave it unaddressable (`""`) by default.
- **`SetPane(id, newPane)`** swaps what's rendered at an existing leaf without touching
  the tree's shape - the way an app switches its content area between entirely
  different pages/sub-apps, each potentially its own package, as opposed to a pane
  updating its own internal state in response to a message. The new pane is `Init`'d and
  immediately told its size; it's told `FocusMsg` too if `id` currently holds focus,
  since whatever it's replacing never will be again.

A pane never holds a reference to the `Model` it lives in, so from inside a pane's own
`Update`, use the message forms instead - `SplitLeafMsg`, `CloseLeafMsg`, `ResizeMsg`,
`SetPaneMsg`:

```go
return p, func() tea.Msg {
	return layout.SplitLeafMsg{ID: "editor", Dir: layout.Vertical, NewID: "terminal", NewModel: newTerminalPane()}
}
```

## Composing bigger apps

A `layout.Model` is itself a `Pane` (and a `Navigable`, see below) - embed one inside
another directly, no wrapper needed:

```go
func newWorkspace() layout.Model {
	root := layout.VSplit(0.7,
		layout.Leaf("editor", newEditorPane()),
		layout.Leaf("terminal", newTerminalPane()),
	)
	return layout.New(root) // no AsRoot(): the outer Model already renders one help bar
}

root := layout.HSplit(0.25,
	layout.Leaf("sidebar", newSidebarPane()),
	layout.Leaf("workspace", newWorkspace()),
)
m := layout.New(root, layout.AsRoot())
```

Once embedded this way, everything works transparently:

- `ctrl+hjkl` tries moving focus *inside* whatever's currently focused first; only once
  that reports being at its own edge does the level above move between its own direct
  children instead.
- The help bar keeps showing exactly one bar, reflecting whatever's focused anywhere in
  the nesting.
- `SendMsg`/`RequestFocusMsg` reach an id inside a nested tree automatically, without
  the outer tree needing to know it's there.
- The nested tree's own shape stays its own business - nothing from outside reaches
  into it structurally, only messages cross that boundary.

This all works because `layout.Model` implements `Navigable`:

```go
type Navigable interface {
	Pane
	Leaves() []LeafRect
	MoveFocus(dir FocusDirection) bool
	Route(target string, msg tea.Msg) (handled bool, cmd tea.Cmd)
	Focus(id string) (handled bool, cmd tea.Cmd)
	FocusedHelp() []key.Binding
}
```

A hand-rolled `Pane` never needs to implement this itself - it's what lets one
`layout.Model` recognize *another* `layout.Model` sitting in one of its leaves and
delegate to it, at arbitrary nesting depth. You'll only reach for it directly if you're
building something that itself wants to compose with `layout` the same way `layout`
composes with itself.

## Examples

- `examples/layout/basic` - a 2x2 grid, no custom border, the minimum to get started.
- `examples/layout/bordered` - each pane draws its own border via `layout.Bordered`,
  following focus.
- `examples/layout/nested` - a whole `layout.Model` embedded as a pane, ctrl+hjkl and
  the help bar both working transparently across the boundary.
- `examples/layout/messaging` - a "control" pane driving and focusing an "editor" pane
  by id via `SendMsg`/`RequestFocusMsg`.
- `examples/layout/help` - the help bar changing to match whatever's focused, including
  a pane that implements no bindings at all.
