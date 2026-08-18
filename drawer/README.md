# drawer

A full-height panel flush against the left or right edge of an already-rendered background, on top
of it dimmed - the sidebar/drawer equivalent of [`modal`](../modal/README.md), which this package
otherwise mirrors closely: same stack of panels triggered from anywhere via an exported `tea.Msg`
(`ShowMsg`/`Show`) rather than a direct reference to the `Model` that ends up rendering it, same
composite-over-an-already-rendered-string `Render`, no assumption about how the host builds that
string.

## Quick start

```go
import (
	"github.com/anotherhadi/ilovetui/drawer"
)

type model struct {
	d             drawer.Model
	width, height int
}

func newModel() model {
	return model{d: drawer.New()}
}

func (m model) Init() tea.Cmd { return m.d.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "l" {
			return m, drawer.Show("Nav", drawer.Text("Home\nProjects\nSettings"))
		}
		if msg.String() == "esc" && m.d.Open() {
			return m, drawer.Close()
		}
	}

	var cmd tea.Cmd
	m.d, cmd = m.d.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := renderYourUI(m.width, m.height)
	view := tea.NewView(m.d.Render(background))
	view.AltScreen = true
	return view
}
```

Any component in the same bubbletea program can trigger a drawer via `drawer.Show`, without
holding a reference to the `drawer.Model` that will actually render it - that `Model` just needs
to see every `tea.Msg` the program produces (i.e. get its `Update` called from the top-level
`Update`), same as any other child model.

## Content is a model

A drawer's body is a `tea.Model`, not a string. While a drawer is on top of the stack it gets
every message the `drawer.Model` receives, its `Init` runs when it opens, and its commands come
back out - so it can hold a file list, a filter form, or a picker that reports its choice with a
`tea.Msg` of its own, which the component that opened it listens for:

```go
type pickedMsg struct{ file string }

// somewhere else
return m, drawer.Show("Files", newFileList(dir), drawer.WithSide(drawer.Right))
```

Only the topmost drawer is updated: everything beneath it is dimmed and frozen until the drawers
above it close.

For a drawer with nothing to interact with, `drawer.Text` wraps a plain string:

```go
return m, drawer.Show("Nav", drawer.Text("Home\nProjects\nSettings"))
```

The box shrinks to fit whatever the content draws (unless `WithWidth` fixes it), so a content
that wants a specific size sets it on itself - the drawer only ever sees the rendered result.

## Side and width

```go
return m, drawer.Show("Nav", drawer.Text("Home\nProjects\nSettings"),
	drawer.WithSide(drawer.Right), drawer.WithWidth(24))
```

`WithSide` anchors the drawer to `drawer.Left` (the default) or `drawer.Right`. `WithWidth` fixes
the drawer's total width instead of shrinking to fit its content, still capped by whatever
actually fits the background - same unit as the `Model`-level `WithMaxWidth`. Either way, the
drawer always spans the background's full height, flush top to bottom.

## Showing and closing

```go
return m, drawer.Show("Files", newFileList(dir), drawer.WithSide(drawer.Right))

return m, drawer.Close() // close the topmost drawer
```

The stack is a plain LIFO, with no identity: a drawer is closed by being on top, never by being
named. There is nothing to tag a drawer with, and nothing that can target one in the middle of the
stack - the topmost is both the only one that receives messages and the only one `Close` can
reach, so the two rules never disagree.

- `drawer.Open()` reports whether at least one drawer is currently shown - handy for a host that
  wants to route key presses to the drawer instead of its normal UI while one is open. Note that a
  host doing this also makes it impossible for a key to open a second drawer, which is what keeps
  the stack shallow without any bookkeeping.
- A content model closes its own drawer by returning `drawer.Close()`, since it only ever runs
  while it is the topmost one.

## Stacking

Drawers stack: showing a second one while the first is still open pushes it on top, dimming both
the background and the first drawer to the same flat color, same as `modal`. Opening a left drawer
and a right drawer at once still stacks - if you want both visible at full color simultaneously,
render them as two ordinary panels via `layout` instead; this package is for transient
sidebars, not permanent chrome.

## Styling

```go
d := drawer.New(drawer.WithMaxWidth(30), drawer.WithStyles(myStyles))

return m, drawer.Show("Nav", drawer.Text("content"), drawer.WithDrawerStyle(oneOffStyles))
```

`WithMaxWidth` caps how wide a drawer can grow (border and padding included) before wrapping; a
drawer narrower than the cap shrinks to fit its content instead of padding out to it, unless shown
with `WithWidth`. A drawer can also never overflow past the edge of whatever background it's
rendered on. `WithStyles` sets the default look for every drawer shown by this `Model`;
`WithDrawerStyle` (a `Show` option) overrides it for one drawer alone. `DefaultStyles()` builds
from `style.S`, same palette roles as `modal.DefaultStyles()`.

## Examples

- `examples/drawer` - a left nav drawer and a right inspector drawer, either one at a time.
