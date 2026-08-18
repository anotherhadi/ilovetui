# modal

A centered popup box on top of a dimmed background, triggered from anywhere in a bubbletea
program via an exported `tea.Msg` (`ShowMsg`/`Show`) rather than a direct reference to the
`Model` that ends up rendering it - standard Elm architecture, no IPC between processes.

It composites over an already-rendered string, so it makes no assumption about how the host
builds that string: the same `Model` works whatever the host uses to lay out its main content.

## Quick start

```go
import (
	"github.com/anotherhadi/ilovetui/modal"
)

type model struct {
	m             modal.Model
	width, height int
}

func newModel() model {
	return model{m: modal.New()}
}

func (m model) Init() tea.Cmd { return m.m.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "d" {
			return m, modal.Show("Delete file?", modal.Text("This can't be undone.\n\ny: confirm  esc: cancel"))
		}
		if msg.String() == "esc" && m.m.Open() {
			return m, modal.Close()
		}
	}

	var cmd tea.Cmd
	m.m, cmd = m.m.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := renderYourUI(m.width, m.height)
	view := tea.NewView(m.m.Render(background))
	view.AltScreen = true
	return view
}
```

Any component in the same bubbletea program can trigger a modal via `modal.Show`, without holding
a reference to the `modal.Model` that will actually render it - that `Model` just needs to see
every `tea.Msg` the program produces (i.e. get its `Update` called from the top-level `Update`),
same as any other child model.

## Content is a model

A modal's body is a `tea.Model`, not a string. While a modal is on top of the stack it gets every
message the `modal.Model` receives, its `Init` runs when it opens, and its commands come back out
- so it can hold a form, a list, or a confirmation that reports its answer with a `tea.Msg` of
its own, which the component that opened it listens for:

```go
type confirmedMsg struct{ path string }

func (c confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "y" {
		path := c.path
		return c, func() tea.Msg { return confirmedMsg{path} }
	}
	return c, nil
}

// somewhere else
return m, modal.Show("Delete file?", confirm{path: p})
```

Only the topmost modal is updated: everything beneath it is dimmed and frozen until the modals
above it close.

For a modal with nothing to interact with, `modal.Text` wraps a plain string:

```go
return m, modal.Show("About", modal.Text("v1.0\n\nesc to close"))
```

The box shrinks to fit whatever the content draws, so a content that wants a specific size sets
it on itself - the modal only ever sees the rendered result.

## Showing and closing

```go
return m, modal.Show("Delete file?", confirm{})

return m, modal.Close() // close the topmost modal
```

The stack is a plain LIFO, with no identity: a modal is closed by being on top, never by being
named. There is nothing to tag a modal with, and nothing that can target one in the middle of the
stack - the topmost is both the only one that receives messages and the only one `Close` can
reach, so the two rules never disagree.

- `modal.Open()` reports whether at least one modal is currently shown - handy for a host that
  wants to route key presses to the modal instead of its normal UI while one is open. Note that a
  host doing this also makes it impossible for a key to open a second modal, which is what keeps
  the stack shallow without any bookkeeping.
- A content model closes its own modal by returning `modal.Close()`, since it only ever runs while
  it is the topmost one.

## Stacking

Modals stack: showing a second one while the first is still open pushes it on top, dimming both
the background and the first modal to the same flat color. Closing the top one reveals the one
beneath, still in full color. This is what lets a "delete?" confirmation open a nested "really
sure?" modal without any special-casing - the nested one is pushed by the content on top, the only
one able to act.

## Styling

```go
m := modal.New(modal.WithMaxWidth(60), modal.WithMaxHeight(20), modal.WithStyles(myStyles))

return m, modal.Show("Title", modal.Text("Message"), modal.WithModalStyle(oneOffStyles))
```

`WithMaxWidth`/`WithMaxHeight` cap how large a modal box can grow before wrapping/truncating; a
modal narrower than the cap shrinks to fit its content instead of padding out to it. A modal can
also never overflow past the edge of whatever background it's rendered on, regardless of these
caps. `WithStyles` sets the default look for every modal shown by this `Model`; `WithModalStyle`
(a `Show` option) overrides it for one modal alone. `DefaultStyles()` builds from `style.S`: the
box borrows `PanelFocused`'s border (the modal is what has focus while it's open), the dim color
reuses `Subtle`.

## Examples

- `examples/modal` - open/dismiss, nested modals, styled from theme colors.
