// @bun
import {
  Badge
} from "./chunk-wasbecqm.js";
import {
  theme
} from "./chunk-57mam43k.js";

// src/components/Toggle.tsx
import { setProp as _$setProp } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
function Toggle(props) {
  const onColor = () => props.onColor ?? theme.success;
  const offColor = () => props.offColor ?? theme.muted;
  const nerdFonts = () => props.withNerdfont ?? theme.nerdFonts;
  return (() => {
    var _el$ = _$createElement("box");
    _$insert(_el$, _$createComponent(Show, {
      get when() {
        return nerdFonts();
      },
      get fallback() {
        return (() => {
          var _el$2 = _$createElement("text");
          _$insert(_el$2, () => props.on ? "on" : "off");
          _$effect((_p$) => {
            var _v$ = props.on ? onColor() : offColor(), _v$2 = TextAttributes.BOLD;
            _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "fg", _v$, _p$.e));
            _v$2 !== _p$.t && (_p$.t = _$setProp(_el$2, "attributes", _v$2, _p$.t));
            return _p$;
          }, {
            e: undefined,
            t: undefined
          });
          return _el$2;
        })();
      },
      get children() {
        return _$createComponent(Badge, {
          get label() {
            return props.on ? "  \u25CF" : "\u25CF  ";
          },
          get color() {
            return _$memo(() => !!props.on)() ? onColor() : offColor();
          },
          withNerdfont: true
        });
      }
    }));
    _$effect((_$p) => _$setProp(_el$, "onMouseDown", theme.mouse ? props.onToggle : undefined, _$p));
    return _el$;
  })();
}

export { Toggle };
