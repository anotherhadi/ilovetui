// @bun
import {
  theme
} from "../chunk-57mam43k.js";
import"../chunk-ergc99dh.js";
import"../chunk-7f8jagy5.js";
import"../chunk-hcq62p48.js";

// src/components/Spinner.tsx
import { setProp as _$setProp } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import"opentui-spinner/solid";
function Spinner(props) {
  return (() => {
    var _el$ = _$createElement("spinner");
    _$effect((_p$) => {
      var { name: _v$, frames: _v$2, interval: _v$3 } = props, _v$4 = props.autoplay ?? !theme.reducedMotion, _v$5 = props.color ?? theme.primary, _v$6 = props.backgroundColor;
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "name", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "frames", _v$2, _p$.t));
      _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "interval", _v$3, _p$.a));
      _v$4 !== _p$.o && (_p$.o = _$setProp(_el$, "autoplay", _v$4, _p$.o));
      _v$5 !== _p$.i && (_p$.i = _$setProp(_el$, "color", _v$5, _p$.i));
      _v$6 !== _p$.n && (_p$.n = _$setProp(_el$, "backgroundColor", _v$6, _p$.n));
      return _p$;
    }, {
      e: undefined,
      t: undefined,
      a: undefined,
      o: undefined,
      i: undefined,
      n: undefined
    });
    return _el$;
  })();
}
export {
  Spinner
};
