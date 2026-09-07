// @bun
import {
  theme
} from "./chunk-x7m82z3z.js";

// src/components/ProgressBar.tsx
import { effect as _$effect } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { parseColor, RGBA } from "@opentui/core";
import { createMemo, For } from "solid-js";
function sampleGradient(stops, t) {
  if (stops.length === 1)
    return stops[0];
  const clamped = Math.min(Math.max(t, 0), 1);
  const segment = clamped * (stops.length - 1);
  const i = Math.min(Math.floor(segment), stops.length - 2);
  const localT = segment - i;
  const [ar, ag, ab, aa] = parseColor(stops[i]).toInts();
  const [br, bg, bb, ba] = parseColor(stops[i + 1]).toInts();
  return RGBA.fromInts(Math.round(ar + (br - ar) * localT), Math.round(ag + (bg - ag) * localT), Math.round(ab + (bb - ab) * localT), Math.round(aa + (ba - aa) * localT));
}
function ProgressBar(props) {
  const width = () => props.width ?? 20;
  const stops = createMemo(() => Array.isArray(props.color) ? props.color : [props.color ?? theme.primary]);
  const filled = createMemo(() => Math.round(width() * Math.min(Math.max(props.value, 0), 100) / 100));
  const cells = createMemo(() => Array.from({
    length: width()
  }, (_, i) => i));
  return (() => {
    var _el$ = _$createElement("box");
    _$setProp(_el$, "flexDirection", "row");
    _$insert(_el$, _$createComponent(For, {
      get each() {
        return cells();
      },
      children: (i) => {
        const isFilled = () => i < filled();
        const t = width() > 1 ? i / (width() - 1) : 0;
        return (() => {
          var _el$2 = _$createElement("text");
          _$insert(_el$2, (() => {
            var _c$ = _$memo(() => !!isFilled());
            return () => _c$() ? props.fillChar ?? "\u2588" : props.trackChar ?? "\u2591";
          })());
          _$effect((_$p) => _$setProp(_el$2, "fg", isFilled() ? sampleGradient(stops(), t) : props.trackColor ?? theme.muted, _$p));
          return _el$2;
        })();
      }
    }));
    return _el$;
  })();
}

export { ProgressBar };
