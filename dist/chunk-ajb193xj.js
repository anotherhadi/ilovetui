// @bun
import {
  notifications2,
  dismiss2
} from "./chunk-6ggk5fg7.js";
import {
  theme2
} from "./chunk-p6jvscj3.js";

// src/components/NotificationHost.tsx
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { Portal } from "@opentui/solid";
import { For } from "solid-js";
var KIND_LABEL = {
  info: "Info",
  success: "Success",
  warning: "Warning",
  error: "Error"
};
function NotificationHost2(props = {}) {
  const defaultColor = () => ({
    info: theme2.primary,
    success: theme2.success,
    warning: theme2.warning,
    error: theme2.error
  });
  const colorFor = (kind) => props.colors?.[kind] ?? defaultColor()[kind];
  return _$createComponent(Portal, {
    get children() {
      var _el$ = _$createElement("box");
      _$setProp(_el$, "position", "absolute");
      _$setProp(_el$, "top", 1);
      _$setProp(_el$, "right", 1);
      _$setProp(_el$, "flexDirection", "column");
      _$setProp(_el$, "zIndex", 1000);
      _$insert(_el$, _$createComponent(For, {
        get each() {
          return notifications2();
        },
        children: (toast) => (() => {
          var _el$2 = _$createElement("box"), _el$3 = _$createElement("text"), _el$4 = _$createElement("text");
          _$insertNode(_el$2, _el$3);
          _$insertNode(_el$2, _el$4);
          _$setProp(_el$2, "border", true);
          _$setProp(_el$2, "marginBottom", 1);
          _$insert(_el$3, () => toast.title ?? KIND_LABEL[toast.kind]);
          _$insert(_el$4, () => toast.message);
          _$effect((_p$) => {
            var _v$ = theme2.borderStyle, _v$2 = colorFor(toast.kind), _v$3 = props.backgroundColor ?? theme2.background, _v$4 = props.width ?? 32, _v$5 = theme2.mouse ? () => dismiss2(toast.id) : undefined, _v$6 = colorFor(toast.kind);
            _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "borderStyle", _v$, _p$.e));
            _v$2 !== _p$.t && (_p$.t = _$setProp(_el$2, "borderColor", _v$2, _p$.t));
            _v$3 !== _p$.a && (_p$.a = _$setProp(_el$2, "backgroundColor", _v$3, _p$.a));
            _v$4 !== _p$.o && (_p$.o = _$setProp(_el$2, "width", _v$4, _p$.o));
            _v$5 !== _p$.i && (_p$.i = _$setProp(_el$2, "onMouseDown", _v$5, _p$.i));
            _v$6 !== _p$.n && (_p$.n = _$setProp(_el$3, "fg", _v$6, _p$.n));
            return _p$;
          }, {
            e: undefined,
            t: undefined,
            a: undefined,
            o: undefined,
            i: undefined,
            n: undefined
          });
          return _el$2;
        })()
      }));
      return _el$;
    }
  });
}

export { NotificationHost2 };
