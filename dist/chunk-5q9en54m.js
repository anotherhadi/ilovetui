// @bun
import {
  dismiss,
  notifications
} from "./chunk-y6txxzxk.js";
import {
  theme
} from "./chunk-9cbm4zqz.js";

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
function NotificationHost(props = {}) {
  const defaultColor = () => ({
    info: theme.primary,
    success: theme.success,
    warning: theme.warning,
    error: theme.error
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
          return notifications();
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
            var _v$ = theme.borderStyle, _v$2 = colorFor(toast.kind), _v$3 = props.width ?? 32, _v$4 = theme.mouse ? () => dismiss(toast.id) : undefined, _v$5 = colorFor(toast.kind);
            _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "borderStyle", _v$, _p$.e));
            _v$2 !== _p$.t && (_p$.t = _$setProp(_el$2, "borderColor", _v$2, _p$.t));
            _v$3 !== _p$.a && (_p$.a = _$setProp(_el$2, "width", _v$3, _p$.a));
            _v$4 !== _p$.o && (_p$.o = _$setProp(_el$2, "onMouseDown", _v$4, _p$.o));
            _v$5 !== _p$.i && (_p$.i = _$setProp(_el$3, "fg", _v$5, _p$.i));
            return _p$;
          }, {
            e: undefined,
            t: undefined,
            a: undefined,
            o: undefined,
            i: undefined
          });
          return _el$2;
        })()
      }));
      return _el$;
    }
  });
}

export { NotificationHost };
