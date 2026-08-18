# drawer

A full-height panel flush against the left or right edge of an already-rendered background, on top
of it dimmed. The sidebar/drawer equivalent of [`modal`](../modal/README.md), which it otherwise
mirrors closely: same `tea.Msg`-triggered stack (`ShowMsg`/`Show`), same composite-over-a-string
`Render`, same rule that content is a `tea.Model`.

- `WithSide(Left|Right)` and `WithWidth` are per-drawer `Show` options; width otherwise shrinks to
  fit content, capped by the `Model`'s `WithMaxWidth`.
- The stack is a plain LIFO. Only the topmost drawer is updated; `Close()` closes it.
- `View(width, height)` renders on a blank background of that size, for use as a standalone pane.

See `examples/drawer`.
