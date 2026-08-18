# style

Shared Base16 theming for bubbletea/lipgloss TUIs. Loaded automatically on import from
`~/.config/ilovetui/config.yaml` (embedded default as fallback), exposed as the package-level
`style.S`.

- `S.NerdFonts` and `S.BorderType` come from the same config as the colors. This package has no
  icon registry: components that want icons read `NerdFonts` and pick their own glyphs.
- `RenderWithTitle` renders a bordered box with a title embedded in the top border, following
  `S.BorderType`.
- Never imports `charm.land/bubbles/v2/*`. Anything that needs to know about a specific component
  belongs in `bubbles/` instead.

## Config

Copy [`default.yaml`](default.yaml) to `~/.config/ilovetui/config.yaml` and edit it.

See `examples/style`.
