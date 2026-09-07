// @bun
import {
  theme
} from "./chunk-x7m82z3z.js";

// src/components/MinSizeGuard.tsx
import { createComponent as _$createComponent } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { useTerminalDimensions } from "@opentui/solid";
import { Show } from "solid-js";
function MinSizeGuard(props) {
  const dimensions = useTerminalDimensions();
  const minWidth = () => props.minWidth ?? 80;
  const minHeight = () => props.minHeight ?? 24;
  const tooSmall = () => dimensions().width < minWidth() || dimensions().height < minHeight();
  return [_$memo(() => props.children), _$createComponent(Show, {
    get when() {
      return tooSmall();
    },
    get children() {
      var _el$ = _$createElement("box"), _el$2 = _$createElement("text");
      _$insertNode(_el$, _el$2);
      _$setProp(_el$, "position", "absolute");
      _$setProp(_el$, "top", 0);
      _$setProp(_el$, "left", 0);
      _$setProp(_el$, "alignItems", "center");
      _$setProp(_el$, "justifyContent", "center");
      _$setProp(_el$, "zIndex", 9999);
      _$insert(_el$2, () => `Terminal too small
Minimum size: ${minWidth()}x${minHeight()} \u2014 current: ${dimensions().width}x${dimensions().height}`);
      _$effect((_p$) => {
        var _v$ = dimensions().width, _v$2 = dimensions().height, _v$3 = props.backgroundColor ?? theme.background, _v$4 = props.textColor ?? theme.muted;
        _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "width", _v$, _p$.e));
        _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "height", _v$2, _p$.t));
        _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "backgroundColor", _v$3, _p$.a));
        _v$4 !== _p$.o && (_p$.o = _$setProp(_el$2, "fg", _v$4, _p$.o));
        return _p$;
      }, {
        e: undefined,
        t: undefined,
        a: undefined,
        o: undefined
      });
      return _el$;
    }
  })];
}

export { MinSizeGuard };
