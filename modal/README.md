# modal

A centered popup box on top of a dimmed background, triggered from anywhere in a bubbletea
program via an exported `tea.Msg` (`ShowMsg`/`Show`) rather than a direct reference to the
`Model` that ends up rendering it - standard Elm architecture, no IPC between processes.

It composites over an already-rendered string, so it has no dependency on
`github.com/anotherhadi/ilovetui/layout`: the same `Model` works whether the host uses `layout`
for its main content or not.

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
			return m, modal.Show("Delete file?", "This can't be undone.\n\ny: confirm  esc: cancel")
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

## Showing and dismissing

```go
return m, modal.Show("Delete file?", "This can't be undone.", modal.WithID("confirm"))

return m, modal.Dismiss("confirm") // close a specific modal by id
return m, modal.Close()            // close whichever modal is on top, whatever its id
```

- `modal.Open()` reports whether at least one modal is currently shown - handy for a host that
  wants to route key presses to the modal (e.g. `esc` to dismiss, `enter` to confirm) instead of
  its normal UI while one is open.
- `modal.TopID()` returns the id of the topmost (currently interactive) modal, or `""` if none is
  open - useful to tell which modal a generic key like `enter` should act on.
- A modal shown without `WithID` gets an auto-generated id that's never returned to the caller -
  only a modal shown with `WithID` can be targeted by `Dismiss` later. Showing again with the same
  id replaces it in place instead of stacking a duplicate.

## Stacking

Modals stack: showing a second one while the first is still open pushes it on top, dimming both
the background and the first modal to the same flat color. Dismissing the top one reveals the one
beneath, still in full color. This is what lets a "delete?" confirmation open a nested "really
sure?" modal without any special-casing.

## Styling

```go
m := modal.New(modal.WithMaxWidth(60), modal.WithMaxHeight(20), modal.WithStyles(myStyles))

return m, modal.Show("Title", "Message", modal.WithModalStyle(oneOffStyles))
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
