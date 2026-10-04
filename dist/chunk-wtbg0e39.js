// @bun
import {
  theme2
} from "./chunk-p6jvscj3.js";
import {
  Box2
} from "./chunk-8r4rajaz.js";
import {
  Tabs2
} from "./chunk-vf2ybhtx.js";

// src/components/TabbedPanel.tsx
import { insert as _$insert } from "@opentui/solid";
import { memo as _$memo } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
function TabbedPanel2(props) {
  return (() => {
    var _el$ = _$createElement("box");
    _$setProp(_el$, "flexDirection", "column");
    _$setProp(_el$, "flexGrow", 1);
    _$insert(_el$, _$createComponent(Tabs2, {
      get items() {
        return props.items;
      },
      get value() {
        return props.value;
      },
      get onChange() {
        return props.onChange;
      },
      get focused() {
        return props.focused;
      }
    }), null);
    _$insert(_el$, _$createComponent(Box2, {
      flexGrow: 1,
      border: ["left", "right", "bottom"],
      get borderColor() {
        return _$memo(() => !!props.focused)() ? theme2.primary : theme2.muted;
      },
      flexDirection: "column",
      get children() {
        return props.children;
      }
    }), null);
    return _el$;
  })();
}

export { TabbedPanel2 };
