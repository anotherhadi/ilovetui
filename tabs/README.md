# tabs

A horizontal tab bar styled from `style.S`. Renders the bar and the `Content` frame around the
active item; the content itself runs and renders through the host, same as any other custom
component in this repo.

- `Focused()` controls the frame's border color, independent of which tab is active (that's shown
  only by the title style).
- Tabs collapse into a trailing `+N` badge when they don't fit `Width`, keeping the active one
  always visible.
- `WithLoop(false)` clamps navigation at either end instead of wrapping.
- The frame follows `style.S.BorderType`: corners and junctions are derived from the border's own
  glyphs, not hardcoded.

See `examples/tabs`.
