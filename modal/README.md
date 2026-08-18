# modal

A centered popup box on top of a dimmed background, triggered from anywhere in a bubbletea program
via an exported `tea.Msg` (`ShowMsg`/`Show`), not a direct reference to the `Model` that renders it.
Composites over an already-rendered string, so it makes no assumption about how the host builds
that string.

- A modal's content is a `tea.Model`, not a string: it gets `Init`/`Update` while it's on top, and
  can report back with its own `tea.Msg` (see `modal.Text` for a plain, non-interactive content).
- The stack is a plain LIFO with no identity. Only the topmost modal is updated; `Close()` always
  closes it.
- Showing a second modal while one is open pushes it on top, nesting confirmations naturally.
- `WithMaxWidth`/`WithMaxHeight` cap growth; a modal narrower than the cap shrinks to fit instead.

See `examples/modal`.
