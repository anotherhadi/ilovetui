# tabs

A horizontal tab bar, styled from `style.S`. Switches between a set of items with
`left`/`right`/`h`/`l`/`tab`/`shift+tab`, wrapping around at either end by default. Draws its own
frame (tab bar + a `Content` box below it) that follows the theme's configured border family
(`style.S.BorderType`) and reads as one continuous box.

`tabs` only renders the bar and the frame around the active item's content - it never runs the
content itself; that's the host's job, same as any other custom component in this repo.

## Concepts

- **`Tab`** is what each tab shows: `Init() tea.Cmd`, `Update(tea.Msg) (Tab, tea.Cmd)`,
  `View() string` - the same shape used by every other custom component in this repo.
- **`Item`** pairs a `Tab` with the `Title` shown on its tab.
- **`Model`** is the running tab bar: the item list, which one is active, focus state, size, and
  styles. Build one with `tabs.New(items, opts...)`.

## Quick start

```go
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/tabs"
)

func main() {
	items := []tabs.Item{
		{Title: "First", Model: newPane("First")},
		{Title: "Second", Model: newPane("Second")},
		{Title: "Third", Model: newPane("Third")},
	}
	m := model{tabs: tabs.New(items)}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

See `examples/tabs` for the full `model`/`Tab` implementation, including sizing.

## Writing a Tab

```go
type Tab interface {
	Init() tea.Cmd
	Update(tea.Msg) (Tab, tea.Cmd)
	View() string
}
```

`tabs.Update` routes every message that isn't a Next/Prev key press straight to the active item's
`Update`, so a `Tab` behaves like any other bubbletea model - it just never sees messages while a
different tab is active.

## Sizing

`tabs` has no generic way to size an arbitrary `Tab` itself (the interface is intentionally
minimal), so a host building a fullscreen app sizes the whole component, then forwards the actual
content area back into it:

```go
m.tabs.SetSize(width, height)

var cmd tea.Cmd
m.tabs, cmd = m.tabs.Update(tea.WindowSizeMsg{
	Width:  m.tabs.ContentWidth(),
	Height: m.tabs.ContentHeight(),
})
```

The tab bar itself always keeps its intrinsic width (the sum of its tab labels); only `Content`
stretches to fill `Width`, so the bar never looks artificially stretched. `ContentWidth`/
`ContentHeight` report the usable inner area once `Width`/`Height` are set (`Content`'s box size
minus its own border and padding) - forward that to whatever `Tab` implementation needs to know
its own size, exactly as you'd size any other nested bubbles component.

When there are more tabs than fit `Width`, `tabs` collapses the overflow into a single trailing
`+N` badge, keeping a contiguous window around the active tab.

## Focus vs. active tab

Two independent things:

- **`Focused()`/`Focus()`/`Blur()`/`WithFocus(bool)`** (on by default) control the frame's border
  color: `style.S.Primary` when focused, `style.S.Subtle` when blurred. Meant for host apps with
  several panes that toggle focus between them (e.g. alongside `layout`) - the border color never
  depends on which tab is active, only on whether `tabs` itself currently has keyboard focus.
- **Which item is active** is shown only by the tab's title style (`Styles.ActiveTitle` vs.
  `InactiveTitle`), not by border color.

## Navigation

```go
km := tabs.DefaultKeyMap()
km.Next = key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "next"))
m := tabs.New(items, tabs.WithKeyMap(km))
```

`WithLoop(false)` clamps at either end instead of wrapping; `WithActive(i)` sets the initial tab.

## Examples

- `examples/tabs` - a full fullscreen app: sizing, per-tab independent state, `+`/counter demo.
