// @bun
import {
  theme
} from "./chunk-x7m82z3z.js";

// src/components/SignalBars.tsx
import { effect as _$effect } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { For } from "solid-js";
var BARS = ["\u2581", "\u2583", "\u2585", "\u2587"];
function SignalBars(props) {
  const level = () => Math.floor(props.signal * BARS.length / 100);
  return (() => {
    var _el$ = _$createElement("box");
    _$setProp(_el$, "flexDirection", "row");
    _$insert(_el$, _$createComponent(For, {
      each: BARS,
      children: (bar, i) => (() => {
        var _el$2 = _$createElement("text");
        _$insert(_el$2, bar);
        _$effect((_$p) => _$setProp(_el$2, "fg", i() < level() ? props.filledColor : props.emptyColor ?? theme.muted, _$p));
        return _el$2;
      })()
    }));
    return _el$;
  })();
}

export { SignalBars };
