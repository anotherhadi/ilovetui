// @bun
import {
  theme
} from "./chunk-9cbm4zqz.js";

// src/components/Modal.tsx
import { insertNode as _$insertNode } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { RGBA } from "@opentui/core";
import { Portal, useTerminalDimensions } from "@opentui/solid";
import { Show } from "solid-js";
function Modal(props) {
  const dimensions = useTerminalDimensions();
  const screenCap = () => Math.max(5, dimensions().height - 4);
  const maxPanelHeight = () => Math.min(props.maxHeight ?? screenCap(), screenCap());
  const maxContentHeight = () => Math.max(1, maxPanelHeight() - 4);
  const maxPanelWidth = () => Math.max(10, dimensions().width - 4);
  const panelWidth = () => Math.min(props.width ?? 50, maxPanelWidth());
  return _$createComponent(Portal, {
    get children() {
      return _$createComponent(Show, {
        get when() {
          return props.open;
        },
        get children() {
          var _el$ = _$createElement("box"), _el$2 = _$createElement("box");
          _$insertNode(_el$, _el$2);
          _$setProp(_el$, "position", "absolute");
          _$setProp(_el$, "top", 0);
          _$setProp(_el$, "left", 0);
          _$setProp(_el$, "alignItems", "center");
          _$setProp(_el$, "justifyContent", "center");
          _$setProp(_el$, "zIndex", 2000);
          _$setProp(_el$2, "border", true);
          _$setProp(_el$2, "flexDirection", "column");
          _$setProp(_el$2, "paddingLeft", 2);
          _$setProp(_el$2, "paddingRight", 2);
          _$setProp(_el$2, "paddingTop", 1);
          _$setProp(_el$2, "paddingBottom", 1);
          _$insert(_el$2, _$createComponent(Show, {
            get when() {
              return props.scrollable;
            },
            get fallback() {
              return (() => {
                var _el$4 = _$createElement("box");
                _$setProp(_el$4, "flexDirection", "column");
                _$setProp(_el$4, "overflow", "hidden");
                _$insert(_el$4, () => props.children);
                _$effect((_$p) => _$setProp(_el$4, "maxHeight", maxContentHeight(), _$p));
                return _el$4;
              })();
            },
            get children() {
              var _el$3 = _$createElement("scrollbox");
              _$setProp(_el$3, "flexShrink", 1);
              _$setProp(_el$3, "scrollY", true);
              _$setProp(_el$3, "scrollX", false);
              _$insert(_el$3, () => props.children);
              _$effect((_$p) => _$setProp(_el$3, "maxHeight", maxContentHeight(), _$p));
              return _el$3;
            }
          }));
          _$effect((_p$) => {
            var _v$ = dimensions().width, _v$2 = dimensions().height, _v$3 = props.backdropColor ?? RGBA.fromInts(0, 0, 0, 150), _v$4 = theme.mouse ? props.onDismiss : undefined, _v$5 = panelWidth(), _v$6 = maxPanelHeight(), _v$7 = theme.borderStyle, _v$8 = props.accentColor ?? theme.primary, _v$9 = props.backgroundColor ?? theme.background, _v$0 = theme.mouse ? (event) => event.stopPropagation() : undefined;
            _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "width", _v$, _p$.e));
            _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "height", _v$2, _p$.t));
            _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "backgroundColor", _v$3, _p$.a));
            _v$4 !== _p$.o && (_p$.o = _$setProp(_el$, "onMouseDown", _v$4, _p$.o));
            _v$5 !== _p$.i && (_p$.i = _$setProp(_el$2, "width", _v$5, _p$.i));
            _v$6 !== _p$.n && (_p$.n = _$setProp(_el$2, "maxHeight", _v$6, _p$.n));
            _v$7 !== _p$.s && (_p$.s = _$setProp(_el$2, "borderStyle", _v$7, _p$.s));
            _v$8 !== _p$.h && (_p$.h = _$setProp(_el$2, "borderColor", _v$8, _p$.h));
            _v$9 !== _p$.r && (_p$.r = _$setProp(_el$2, "backgroundColor", _v$9, _p$.r));
            _v$0 !== _p$.d && (_p$.d = _$setProp(_el$2, "onMouseDown", _v$0, _p$.d));
            return _p$;
          }, {
            e: undefined,
            t: undefined,
            a: undefined,
            o: undefined,
            i: undefined,
            n: undefined,
            s: undefined,
            h: undefined,
            r: undefined,
            d: undefined
          });
          return _el$;
        }
      });
    }
  });
}

export { Modal };
