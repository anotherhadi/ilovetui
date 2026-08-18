# I Love TUI

A shared [Base16](https://github.com/tinted-theming/home) theme, a themed wrapper around every official component, and a small set of custom Bubble Tea components, in one Go module, so every TUI built with it shares one config file and looks consistent.

## Install

```sh
go get github.com/anotherhadi/ilovetui
```

## Packages

- [`style`](style/README.md): the theme itself. Colors, pre-built panel styles, config loading.
- [`bubbles`](bubbles/README.md): themed constructors for official `bubbles/v2` components (`help`, `textarea`, `textinput`, `viewport`, ...).
- [`tabs`](tabs/README.md), [`modal`](modal/README.md), [`drawer`](drawer/README.md), [`notification`](notification/README.md), [`helpbar`](helpbar/README.md): custom components not found in the official `bubbles` library, styled from the same theme.

Each package has its own README and a runnable example under `examples/<package>`.

## Quick start

On import, `style` automatically loads the user's theme from `~/.config/ilovetui/config.yaml` (embedded default as fallback), exposed as the package-level `S`:

```go
import "github.com/anotherhadi/ilovetui/style"

s := lipgloss.NewStyle().Foreground(style.S.Primary)
box := style.RenderWithTitle(style.S.PanelFocused, "Title", content, w, h)
```

No setup required. See [`style/README.md`](style/README.md) for config details.

## Projects using ilovetui

- [anotherhadi/spilltea](https://github.com/anotherhadi/spilltea): A minimal, terminal-based HTTP(S) proxy for pentesters and CTF players. Think Burp Suite or Caido, but entirely in your terminal.
- [anotherhadi/usbguard-tui](https://github.com/anotherhadi/usbguard-tui): A terminal UI for managing USB devices via usbguard. TUI built with golang & bubbletea.
- [anotherhadi/jwt-tui](https://github.com/anotherhadi/jwt-tui): A terminal UI for inspecting, editing, and signing JSON Web Tokens (JWTs).
