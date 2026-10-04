// @bun
import {
  theme2
} from "./chunk-p6jvscj3.js";

// src/components/Box.tsx
import { spread as _$spread } from "@opentui/solid";
import { mergeProps as _$mergeProps } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
function Box2(props) {
  return (() => {
    var _el$ = _$createElement("box");
    _$spread(_el$, _$mergeProps(props, {
      get borderStyle() {
        return props.borderStyle ?? theme2.borderStyle;
      }
    }), false);
    return _el$;
  })();
}

export { Box2 };
