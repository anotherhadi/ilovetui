# Ilovetui

A minimal Go library that provides a shared [Base16](https://github.com/tinted-theming/home) color theme for terminal UIs built with [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss), plus a small collection of Bubble Tea v2 components on top of it.

The idea is simple: instead of every TUI app managing its own colors, they all share one theme file so the user customizes once and every app looks consistent.

## Packages

- `github.com/anotherhadi/ilovetui/style` — the theme itself: colors, pre-built panel styles, config loading.
- `github.com/anotherhadi/ilovetui/bubbles` — themed constructors for official `bubbles/v2` components (`help`, `textarea`, `textinput`, `list`, `table`, `filepicker`, `spinner`, `progress`, `paginator`, `viewport`).
- `github.com/anotherhadi/ilovetui/tabs`, `.../helpbar`, `.../modal`, `.../drawer`, `.../notification` — custom components not found in the official `bubbles` library, styled from the same theme.

## How it works

On import, `style` automatically loads the user's theme from `~/.config/ilovetui/config.yaml` (respecting `$XDG_CONFIG_HOME`).
If no config exists, it falls back to the embedded default. The active theme is exposed as the package-level variable `S`.

```go
import "github.com/anotherhadi/ilovetui/style"

// Use colors directly
s := lipgloss.NewStyle().Foreground(style.S.Primary)

// Use pre-built panel styles
box := style.RenderWithTitle(style.S.PanelFocused, "Title", content, w, h)
```

No setup required — just import and use.

## Installation

```sh
go get github.com/anotherhadi/ilovetui
```

## Theme

The theme follows the [Base16](https://github.com/tinted-theming/home) standard (16 colors). The library exposes both the raw palette and semantic aliases:

| Alias        | Base16 | Meaning                                 |
| ------------ | ------ | --------------------------------------- |
| `Background` | Base00 | Background                              |
| `SubtleBg`   | Base01 | Lighter Background / Status Bars        |
| `Selection`  | Base02 | Selection Background                    |
| `Subtle`     | Base03 | Comments / Invisibles                   |
| `Muted`      | Base04 | Dark Foreground / Status Bars           |
| `Text`       | Base05 | Default Foreground                      |
| `Primary`    | Base0D | Functions / Methods / Headings / Accent |
| `Success`    | Base0B | Strings / Success / Diff Inserted       |
| `Warning`    | Base09 | Integers / Constants / Booleans         |
| `Error`      | Base08 | Variables / Errors / Diff Deleted       |

The default theme is `style/default.yaml`. Copy it and edit to customize:

```sh
mkdir -p ~/.config/ilovetui
cp $(go env GOPATH)/pkg/mod/github.com/anotherhadi/ilovetui*/style/default.yaml ~/.config/ilovetui/config.yaml
```

Or let your app write it on first run:

```go
style.WriteDefaultConfig(style.DefaultConfigPath())
```

## Pre-built styles

`S` ships with a few ready-to-use lipgloss styles:

| Field            | Description                              |
| ---------------- | ---------------------------------------- |
| `S.Bold`         | Bold text                                |
| `S.Faint`        | Muted / dimmed text                      |
| `S.Panel`        | Rounded border, unfocused (Subtle color) |
| `S.PanelFocused` | Rounded border, focused (Primary color)  |

## Helpers

```go
// Inner usable height of a bordered panel with outer height h
inner := style.ContentHeight(h)

// Render a box with a title embedded in the top border
box := style.RenderWithTitle(style.S.PanelFocused, "Header", content, w, h)
```

## API

```go
style.Init()                      // Reload from default config path
style.InitFrom(path string)       // Reload from a custom path
style.InitFromBytes(data []byte)  // Parse raw YAML
style.DefaultConfigPath() string  // ~/.config/ilovetui/config.yaml
style.WriteDefaultConfig(path)    // Write default config if missing
```

## Themed official components

```go
import "github.com/anotherhadi/ilovetui/bubbles"

h  := bubbles.NewHelp()
ta := bubbles.NewTextarea(false)
ti := bubbles.NewTextInput()
l  := bubbles.NewList(items, width, height)
t  := bubbles.NewTable()
fp := bubbles.NewFilePicker()
sp := bubbles.NewSpinner()
pr := bubbles.NewProgress()
pg := bubbles.NewPaginator()
vp := bubbles.NewViewport()
```

Each constructor mirrors the official component's own `New`, then applies `style.S` on top. Where the
official `New` takes options (`spinner`, `table`, `progress`), they're forwarded before the theme is
applied, so you can still customize behavior; anything you pass that also sets colors will be overridden
by the theme afterward.

## Custom components

```go
import "github.com/anotherhadi/ilovetui/tabs"

t := tabs.New([]tabs.Item{{Title: "First", Model: firstPane}, {Title: "Second", Model: secondPane}})
```

`tabs` renders a horizontal tab bar styled from `style.S`. The host application renders the content
below it; see [`tabs/README.md`](tabs/README.md) and `examples/tabs` for a full example.

For the other custom components, see their own README: [`helpbar`](helpbar/README.md) (responsive
help bar that reflows into as many columns as fit), [`modal`](modal/README.md) (centered popup
dialogs), [`drawer`](drawer/README.md) (left/right sidebar panels, mirroring `modal`) and
[`notification`](notification/README.md) (toast notifications).

## Projects using ilovetui

- [anotherhadi/spilltea](https://github.com/anotherhadi/spilltea): A minimal, terminal-based HTTP(S) proxy for pentesters and CTF players. Think Burp Suite or Caido, but entirely in your terminal.
- [anotherhadi/usbguard-tui](https://github.com/anotherhadi/usbguard-tui): A terminal UI for managing USB devices via usbguard. TUI built with golang & bubbletea.
- [anotherhadi/jwt-tui](https://github.com/anotherhadi/jwt-tui): A terminal UI for inspecting, editing, and signing JSON Web Tokens (JWTs).
