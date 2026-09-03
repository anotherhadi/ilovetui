// @bun
import {
  theme
} from "./chunk-9cbm4zqz.js";

// src/components/HelpBar.tsx
import { createTextNode as _$createTextNode } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { use as _$use } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { CliRenderEvents, TextAttributes } from "@opentui/core";
import { onResize, useRenderer } from "@opentui/solid";
import { formatCommandBindings } from "@opentui/keymap/extras";
import { useKeymapSelector } from "@opentui/keymap/solid";
import { createMemo, createSignal, For } from "solid-js";
var KEY_DISPLAY = {
  up: "\u2191",
  down: "\u2193",
  left: "\u2190",
  right: "\u2192",
  enter: "Enter",
  escape: "Esc"
};
var SEPARATOR = "   ";
var ELLIPSIS = "\u2026";
function HelpBar(props = {}) {
  const accent = () => props.accentColor ?? theme.primary;
  const muted = () => props.mutedColor ?? theme.muted;
  const renderer = useRenderer();
  const [width, setWidth] = createSignal(0);
  let container;
  const measure = () => {
    if (!container)
      return;
    setWidth(container.width);
  };
  const remeasureNextFrame = () => renderer.once(CliRenderEvents.FRAME, measure);
  onResize(remeasureNextFrame);
  const entries = useKeymapSelector((keymap) => keymap.getCommandEntries({
    visibility: "active"
  }).filter((entry) => entry.command.shortHelp === true).map((entry) => ({
    command: entry.command.name,
    keys: formatCommandBindings(entry.bindings, {
      keyNameAliases: KEY_DISPLAY
    }) ?? "",
    label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name
  })));
  const visible = createMemo(() => {
    const available = width() - 2;
    const shown = [];
    let used = 0;
    for (const entry of entries()) {
      const sepWidth = shown.length > 0 ? SEPARATOR.length : 0;
      const entryWidth = entry.keys.length + 1 + entry.label.length;
      if (used + sepWidth + entryWidth <= available) {
        used += sepWidth + entryWidth;
        shown.push(entry);
        continue;
      }
      const remaining = available - used - sepWidth;
      const prefix = `${entry.keys} `;
      if (remaining > prefix.length + ELLIPSIS.length) {
        const labelBudget = remaining - prefix.length - ELLIPSIS.length;
        shown.push({
          ...entry,
          label: `${entry.label.slice(0, labelBudget)}${ELLIPSIS}`
        });
      }
      break;
    }
    return shown;
  });
  return (() => {
    var _el$ = _$createElement("box");
    _$use((el) => {
      container = el;
      remeasureNextFrame();
    }, _el$);
    _$setProp(_el$, "flexDirection", "row");
    _$setProp(_el$, "paddingLeft", 1);
    _$setProp(_el$, "paddingRight", 1);
    _$setProp(_el$, "overflow", "hidden");
    _$insert(_el$, _$createComponent(For, {
      get each() {
        return visible();
      },
      children: (entry, i) => [_$memo(() => _$memo(() => i() > 0)() && (() => {
        var _el$4 = _$createElement("text");
        _$insertNode(_el$4, _$createTextNode(`   `));
        _$effect((_$p) => _$setProp(_el$4, "fg", muted(), _$p));
        return _el$4;
      })()), (() => {
        var _el$2 = _$createElement("text");
        _$insert(_el$2, () => entry.keys);
        _$effect((_p$) => {
          var _v$ = accent(), _v$2 = TextAttributes.BOLD;
          _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "fg", _v$, _p$.e));
          _v$2 !== _p$.t && (_p$.t = _$setProp(_el$2, "attributes", _v$2, _p$.t));
          return _p$;
        }, {
          e: undefined,
          t: undefined
        });
        return _el$2;
      })(), (() => {
        var _el$3 = _$createElement("text");
        _$insert(_el$3, () => ` ${entry.label}`);
        _$effect((_$p) => _$setProp(_el$3, "fg", muted(), _$p));
        return _el$3;
      })()]
    }));
    return _el$;
  })();
}

export { HelpBar };
