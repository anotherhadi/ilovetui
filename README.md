# Ilovetui

A minimal Go library that provides a shared [Base16](https://github.com/tinted-theming/home) color theme for terminal UIs built with [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss).

The idea is simple: instead of every TUI app managing its own colors, they all share one theme file so the user customizes once and every app looks consistent.

## How it works

On import, `ilovetui` automatically loads the user's theme from `~/.config/ilovetui/config.yaml` (respecting `$XDG_CONFIG_HOME`).
If no config exists, it falls back to the embedded default. The active theme is exposed as the package-level variable `S`.

```go
import "github.com/anotherhadi/ilovetui"

// Use colors directly
style := lipgloss.NewStyle().Foreground(ilovetui.S.Primary)

// Use pre-built panel styles
box := ilovetui.RenderWithTitle(ilovetui.S.PanelFocused, "Title", content, w, h)
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

The default theme is `./default.yaml`. Copy it and edit to customize:

```sh
mkdir -p ~/.config/ilovetui
cp $(go env GOPATH)/pkg/mod/github.com/anotherhadi/ilovetui*/default.yaml ~/.config/ilovetui/config.yaml
```

Or let your app write it on first run:

```go
ilovetui.WriteDefaultConfig(ilovetui.DefaultConfigPath())
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
inner := ilovetui.ContentHeight(h)

// Render a box with a title embedded in the top border
box := ilovetui.RenderWithTitle(ilovetui.S.PanelFocused, "Header", content, w, h)
```

## API

```go
ilovetui.Init()                      // Reload from default config path
ilovetui.InitFrom(path string)       // Reload from a custom path
ilovetui.InitFromBytes(data []byte)  // Parse raw YAML
ilovetui.DefaultConfigPath() string  // ~/.config/ilovetui/config.yaml
ilovetui.WriteDefaultConfig(path)    // Write default config if missing
```

## Projects using ilovetui

- [anotherhadi/spilltea](https://github.com/anotherhadi/spilltea): A minimal, terminal-based HTTP(S) proxy for pentesters and CTF players. Think Burp Suite or Caido, but entirely in your terminal.
- [anotherhadi/usbguard-tui](https://github.com/anotherhadi/usbguard-tui): A terminal UI for managing USB devices via usbguard. TUI built with golang & bubbletea.
- [anotherhadi/jwt-tui](https://github.com/anotherhadi/jwt-tui): A terminal UI for inspecting, editing, and signing JSON Web Tokens (JWTs).
