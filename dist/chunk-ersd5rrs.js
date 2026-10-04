// @bun
import {
  helpOpen2,
  toggleHelp2
} from "./chunk-ee2qaswx.js";
import {
  Modal2
} from "./chunk-t2cszvff.js";

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
import { createSignal } from "solid-js";
var KEY_DISPLAY = {
  up: "\u2191",
  down: "\u2193",
  left: "\u2190",
  right: "\u2192",
  enter: "Enter",
  escape: "Esc"
};
function HelpModal2(props = {}) {
  const [scope, setScope] = createSignal(null);
  const entries = useKeymapSelector((keymap) => keymap.getCommandEntries({
    visibility: "active",
    focused: scope()
  }).map((entry) => ({
    command: entry.command.name,
    keys: formatCommandBindings(entry.bindings, {
      keyNameAliases: KEY_DISPLAY
    }) ?? "",
    label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name
  })));
  return _$createComponent(Modal2, {
    get open() {
      return helpOpen2();
    },
    onDismiss: toggleHelp2,
    onOpen: (previous) => setScope(() => previous),
    get dismissKey() {
      return props.dismissKey;
    },
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
    get maxHeight() {
      return entries().length + 6;
    },
    scrollable: true,
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

export { HelpModal2 };
