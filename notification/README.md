# notification

Toast-style notifications, triggered from anywhere in a bubbletea program via an exported
`tea.Msg` (`ShowMsg`/`Show`) rather than a direct reference to the `Model` that ends up rendering
them - standard Elm architecture, no IPC between processes.

It composites over an already-rendered string, so it has no dependency on
`github.com/anotherhadi/ilovetui/layout`: the same `Model` works whether the host uses `layout`
for its main content or not.

## Quick start

```go
import (
	"github.com/anotherhadi/ilovetui/notification"
)

type model struct {
	notif         notification.Model
	width, height int
}

func newModel() model {
	return model{notif: notification.New()}
}

func (m model) Init() tea.Cmd { return m.notif.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "s" {
			return m, notification.Show("Saved", "Config written to disk", notification.Success)
		}
	}

	var cmd tea.Cmd
	m.notif, cmd = m.notif.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	background := renderYourUI(m.width, m.height)
	view := tea.NewView(m.notif.Render(background))
	view.AltScreen = true
	return view
}
```

Any component in the same bubbletea program can trigger a toast via `notification.Show`, without
holding a reference to the `notification.Model` that will actually render it - that `Model` just
needs to see every `tea.Msg` the program produces (i.e. get its `Update` called from the top-level
`Update`), same as any other child model.

## Showing and dismissing

```go
return m, notification.Show("Saved", "Config written to disk", notification.Success)

return m, notification.Show("Sticky", "Stays until dismissed",
	notification.Info, notification.WithID("sticky-demo"), notification.WithDuration(0))
return m, notification.Dismiss("sticky-demo")
```

Four kinds: `Info`, `Success`, `Warning`, `Error`, each with its own color preset (see Styling
below). By default a toast auto-dismisses after `notification.DefaultDuration` (3s);
`WithDuration(0)` makes it sticky - it stays until `Dismiss(id)` removes it, so a sticky toast
needs `WithID` to be dismissable later (an auto-generated id is never returned to the caller).
Showing again with the same id replaces the toast in place, resetting its position and timer,
instead of stacking a duplicate.

## Position and stacking

```go
n := notification.New(notification.WithPosition(notification.TopRight))
```

Six anchors: `Top`, `TopLeft`, `TopRight`, `Bottom`, `BottomLeft`, `BottomRight` - toasts always
hug an edge or corner, never the middle of the screen. Multiple toasts stack along the anchored
edge, newest closest to it; a stack that overflows the background's height clips the oldest
toasts first, so the newest ones stay visible.

## Styling

```go
n := notification.New(notification.WithMaxWidth(40), notification.WithStyles(myStyles))

return m, notification.Show("Title", "Message", notification.Success,
	notification.WithToastStyle(oneOffStyle))
```

`WithMaxWidth` caps how wide a toast box can grow before its message wraps; a toast narrower than
the cap shrinks to fit its content instead of padding out to it. A toast can also never overflow
past the edge of whatever background it's rendered on, regardless of this cap. `WithStyles` sets
the default per-`Kind` look for every toast shown by this `Model`; `WithToastStyle` (a `Show`
option) overrides it for one toast alone. `DefaultStyles()` builds from `style.S`: `Info` uses
`Primary` (no dedicated "info" color in the theme), `Success`/`Warning`/`Error` use their matching
`style.S` alias.

## Examples

- `examples/notification` - all four kinds, a sticky toast with manual dismiss, cycling through
  all six positions.
