# notification

Toast notifications, triggered from anywhere in a bubbletea program via an exported `tea.Msg`
(`ShowMsg`/`Show`), not a direct reference to the `Model` that renders them. Composites over an
already-rendered string, so it makes no assumption about how the host builds that string.

- Four kinds: `Info`, `Success`, `Warning`, `Error`, each with its own `style.S` color preset.
- Auto-dismisses after `DefaultDuration` (3s) unless shown with `WithDuration(0)`, which makes it
  sticky; a sticky toast needs `WithID` so `Dismiss(id)` can remove it later.
- Six anchors (`Top`, `TopLeft`, `TopRight`, `Bottom`, `BottomLeft`, `BottomRight`). Toasts stack
  along the anchored edge, newest closest to it.
- `WithMaxWidth` caps growth; a toast narrower than the cap shrinks to fit instead.

See `examples/notification`.
