# bubbles

Themed constructors for every official `bubbles/v2` component: `help`, `textarea`, `textinput`,
`list`, `table`, `filepicker`, `spinner`, `progress`, `paginator`, `viewport`. Each wraps the
official `New`, then applies `style.S` colors on top.

- `spinner`, `table`, `progress` forward any `opts` to the upstream `New` before the theme is
  applied, so the theme's colors always win over anything conflicting in `opts`.
- `ViewportView` renders a `viewport.Model` with a themed scrollbar thumb in the left gutter,
  shown only when content overflows.
- `NewList`/`NewDefaultDelegate` theme both the list chrome and the selection/filter colors of its
  default delegate.

Custom components in this repo (`tabs`, `modal`, `drawer`, `notification`, `helpbar`) build any
official component they need through here, never through `charm.land/bubbles/v2` directly.

See `examples/bubbles`.
