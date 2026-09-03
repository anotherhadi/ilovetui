// @bun
import {
  Modal
} from "./chunk-3t6nhsa8.js";
import {
  helpOpen,
  toggleHelp
} from "./chunk-cyvpvpws.js";

// src/components/HelpModal.tsx
import { createComponent as _$createComponent } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { createTextNode as _$createTextNode } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { formatCommandBindings } from "@opentui/keymap/extras";
import { useKeymapSelector } from "@opentui/keymap/solid";
var KEY_DISPLAY = {
  up: "\u2191",
  down: "\u2193",
  left: "\u2190",
  right: "\u2192",
  enter: "Enter",
  escape: "Esc"
};
function HelpModal(props = {}) {
  const entries = useKeymapSelector((keymap) => keymap.getCommandEntries({
    visibility: "active"
  }).map((entry) => ({
    command: entry.command.name,
    keys: formatCommandBindings(entry.bindings, {
      keyNameAliases: KEY_DISPLAY
    }) ?? "",
    label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name
  })));
  return _$createComponent(Modal, {
    get open() {
      return helpOpen();
    },
    onDismiss: toggleHelp,
    get accentColor() {
      return props.accentColor;
    },
    get backgroundColor() {
      return props.backgroundColor;
    },
    get backdropColor() {
      return props.backdropColor;
    },
    get width() {
      return props.width;
    },
    get children() {
      return [(() => {
        var _el$ = _$createElement("text");
        _$insertNode(_el$, _$createTextNode(`Keybindings`));
        _$effect((_$p) => _$setProp(_el$, "attributes", TextAttributes.BOLD, _$p));
        return _el$;
      })(), (() => {
        var _el$3 = _$createElement("text");
        _$insertNode(_el$3, _$createTextNode(` `));
        return _el$3;
      })(), (() => {
        var _el$5 = _$createElement("text");
        _$insert(_el$5, () => entries().map((entry) => `${entry.keys.padEnd(10)}${entry.label}`).join(`
`));
        return _el$5;
      })()];
    }
  });
}

export { HelpModal };
