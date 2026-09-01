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

Entry points: `.`, `./solid`, `./context`, `./project-config`, `./keymap`, `./components`, `./components/help`, `./components/spinner`. See the source for what each exports.

The theme is base16 and user-overridable via `~/.config/ilovetui/config.yaml`. Copy [`src/default.yaml`](https://github.com/anotherhadi/ilovetui/blob/main/src/default.yaml) there as a starting point.
