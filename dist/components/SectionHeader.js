// @bun
import"../chunk-njvj322f.js";
import"../chunk-6ggk5fg7.js";
import"../chunk-yyg967q8.js";
import"../chunk-xykm10y9.js";
import {
  theme2
} from "../chunk-p6jvscj3.js";
import {
  Spinner2
} from "../chunk-kmqqzz9e.js";

// src/components/SectionHeader.tsx
import { createTextNode as _$createTextNode } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
function SectionHeader(props) {
  return (() => {
    var _el$ = _$createElement("box"), _el$2 = _$createElement("box"), _el$3 = _$createElement("text"), _el$4 = _$createElement("box");
    _$insertNode(_el$, _el$2);
    _$insertNode(_el$, _el$4);
    _$setProp(_el$, "flexShrink", 0);
    _$setProp(_el$, "flexDirection", "row");
    _$setProp(_el$, "justifyContent", "space-between");
    _$setProp(_el$, "alignItems", "center");
    _$setProp(_el$, "height", 1);
    _$insertNode(_el$2, _el$3);
    _$setProp(_el$2, "flexGrow", 1);
    _$setProp(_el$2, "flexShrink", 1);
    _$setProp(_el$2, "minWidth", 0);
    _$setProp(_el$2, "marginRight", 1);
    _$setProp(_el$2, "overflow", "hidden");
    _$setProp(_el$3, "wrapMode", "none");
    _$insert(_el$3, () => props.title ?? "");
    _$setProp(_el$4, "flexShrink", 0);
    _$insert(_el$4, _$createComponent(Show, {
      get when() {
        return !props.loading;
      },
      get fallback() {
        return (() => {
          var _el$6 = _$createElement("box"), _el$7 = _$createElement("text"), _el$8 = _$createTextNode(` `);
          _$insertNode(_el$6, _el$7);
          _$setProp(_el$6, "flexDirection", "row");
          _$setProp(_el$6, "alignItems", "center");
          _$insert(_el$6, _$createComponent(Spinner2, {}), _el$7);
          _$insertNode(_el$7, _el$8);
          _$insert(_el$7, () => props.loadingLabel ?? "Loading\u2026", null);
          _$effect((_$p) => _$setProp(_el$7, "fg", theme2.muted, _$p));
          return _el$6;
        })();
      },
      get children() {
        var _el$5 = _$createElement("text");
        _$insert(_el$5, () => props.hint ?? "");
        _$effect((_$p) => _$setProp(_el$5, "fg", theme2.muted, _$p));
        return _el$5;
      }
    }));
    _$effect((_p$) => {
      var _v$ = theme2.muted, _v$2 = TextAttributes.BOLD;
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$3, "fg", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$3, "attributes", _v$2, _p$.t));
      return _p$;
    }, {
      e: undefined,
      t: undefined
    });
    return _el$;
  })();
}
export {
  SectionHeader
};
