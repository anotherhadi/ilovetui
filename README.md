# I Love TUI

Shared OpenTUI/Solid theme, components, and config helpers for my TUIs.

## Install

```bash
bun add github:anotherhadi/ilovetui
```

Peer dependencies, install whichever you actually use: `@opentui/core`, `@opentui/solid`, `@opentui/keymap`, `solid-js`, `opentui-spinner`.

## Usage

```ts
import { theme, presets } from "ilovetui";
import "ilovetui/solid";
```

Entry points: `.`, `./solid`, `./context`, `./project-config`, `./keymap`, `./components`, `./components/help`, `./components/spinner`, `./components/section-header`. See the source for what each exports.

The theme is base16 and user-overridable via `~/.config/ilovetui/config.yaml`. Copy [`src/default.yaml`](https://github.com/anotherhadi/ilovetui/blob/main/src/default.yaml) there as a starting point.

An app can layer its own theme overrides on top of that file (e.g. a `theme:` section of its own config, validated with `ThemeConfigSchema`) by calling `configureTheme(overrides)` once at startup, before the first render. Invalid theme files are reported as warning toasts and skipped; unknown keys are ignored, so other TUIs can keep their own keys in the shared file.

## Focus and keymap scoping

Wrap each focus region in a `<Pane>` (`focusable` when the pane itself is the focus target). Inside it, `usePane()` gives `target` — pass it to `useBindings` so the layer is only active while focus is within the pane — and `focused()` for styling. `<Modal>` takes focus while open, restores it on close, and binds its `dismissKey` (default `escape`) plus any `commands`/`bindings` you pass, scoped to itself. Note that keymap layers are ordered by `priority` then registration order, not by tree depth: keep overlapping layers on disjoint targets, or give them explicit priorities.

