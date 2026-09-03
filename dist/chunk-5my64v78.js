// @bun
import {
  theme
} from "./chunk-9cbm4zqz.js";

// src/components/Badge.tsx
import { insert as _$insert } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { createTextNode as _$createTextNode } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
function Badge(props) {
  const textColor = () => props.textColor ?? theme.background;
  const backgroundColor = () => props.backgroundColor ?? theme.background;
  const attributes = () => props.bold === false ? undefined : TextAttributes.BOLD;
  const nerdFonts = () => props.withNerdfont ?? theme.nerdFonts;
  return (() => {
    var _el$ = _$createElement("box"), _el$4 = _$createElement("text");
    _$insertNode(_el$, _el$4);
    _$setProp(_el$, "flexDirection", "row");
    _$insert(_el$, _$createComponent(Show, {
      get when() {
        return nerdFonts();
      },
      get children() {
        var _el$2 = _$createElement("text");
        _$insertNode(_el$2, _$createTextNode(`\uE0B6`));
        _$effect((_p$) => {
          var _v$ = props.color, _v$2 = backgroundColor();
          _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "fg", _v$, _p$.e));
          _v$2 !== _p$.t && (_p$.t = _$setProp(_el$2, "bg", _v$2, _p$.t));
          return _p$;
        }, {
          e: undefined,
          t: undefined
        });
        return _el$2;
      }
    }), _el$4);
    _$insert(_el$4, (() => {
      var _c$ = _$memo(() => !!nerdFonts());
      return () => _c$() ? props.label : ` ${props.label} `;
    })());
    _$insert(_el$, _$createComponent(Show, {
      get when() {
        return nerdFonts();
      },
      get children() {
        var _el$5 = _$createElement("text");
        _$insertNode(_el$5, _$createTextNode(`\uE0B4`));
        _$effect((_p$) => {
          var _v$3 = props.color, _v$4 = backgroundColor();
          _v$3 !== _p$.e && (_p$.e = _$setProp(_el$5, "fg", _v$3, _p$.e));
          _v$4 !== _p$.t && (_p$.t = _$setProp(_el$5, "bg", _v$4, _p$.t));
          return _p$;
        }, {
          e: undefined,
          t: undefined
        });
        return _el$5;
      }
    }), null);
    _$effect((_p$) => {
      var _v$5 = textColor(), _v$6 = props.color, _v$7 = attributes();
      _v$5 !== _p$.e && (_p$.e = _$setProp(_el$4, "fg", _v$5, _p$.e));
      _v$6 !== _p$.t && (_p$.t = _$setProp(_el$4, "bg", _v$6, _p$.t));
      _v$7 !== _p$.a && (_p$.a = _$setProp(_el$4, "attributes", _v$7, _p$.a));
      return _p$;
    }, {
      e: undefined,
      t: undefined,
      a: undefined
    });
    return _el$;
  })();
}

export { Badge };
