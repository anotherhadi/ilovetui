// @bun
import {
  theme2
} from "./chunk-p6jvscj3.js";

// src/components/KeyValue.tsx
import { memo as _$memo } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { Show } from "solid-js";
function KeyValue2(props) {
  return (() => {
    var _el$ = _$createElement("box"), _el$2 = _$createElement("text");
    _$insertNode(_el$, _el$2);
    _$setProp(_el$, "flexShrink", 0);
    _$setProp(_el$, "flexDirection", "row");
    _$setProp(_el$, "justifyContent", "space-between");
    _$insert(_el$2, () => props.label);
    _$insert(_el$, _$createComponent(Show, {
      get when() {
        return props.children;
      },
      get fallback() {
        return (() => {
          var _el$3 = _$createElement("text");
          _$insert(_el$3, () => props.value ?? "");
          _$effect((_$p) => _$setProp(_el$3, "fg", props.valueColor, _$p));
          return _el$3;
        })();
      },
      get children() {
        return props.children;
      }
    }), null);
    _$effect((_$p) => _$setProp(_el$2, "fg", theme2.muted, _$p));
    return _el$;
  })();
}

export { KeyValue2 };
