// @bun
import {
  theme2,
  presets2
} from "./chunk-p6jvscj3.js";

// src/components/ListRow.tsx
import { effect as _$effect } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
function ListRow2(props) {
  const tone = (fallback) => props.selected ? presets2.select.selectedTextColor : fallback;
  return (() => {
    var _el$ = _$createElement("box"), _el$2 = _$createElement("box"), _el$3 = _$createElement("box"), _el$4 = _$createElement("text");
    _$insertNode(_el$, _el$2);
    _$setProp(_el$, "flexShrink", 0);
    _$setProp(_el$, "flexDirection", "column");
    _$setProp(_el$, "paddingLeft", 1);
    _$setProp(_el$, "paddingRight", 1);
    _$insertNode(_el$2, _el$3);
    _$setProp(_el$2, "flexDirection", "row");
    _$setProp(_el$2, "justifyContent", "space-between");
    _$setProp(_el$2, "alignItems", "center");
    _$setProp(_el$2, "height", 1);
    _$insertNode(_el$3, _el$4);
    _$setProp(_el$3, "flexGrow", 1);
    _$setProp(_el$3, "flexShrink", 1);
    _$setProp(_el$3, "minWidth", 0);
    _$setProp(_el$3, "marginRight", 1);
    _$setProp(_el$3, "overflow", "hidden");
    _$setProp(_el$4, "wrapMode", "none");
    _$insert(_el$4, () => props.label);
    _$insert(_el$2, _$createComponent(Show, {
      get when() {
        return props.right;
      },
      get children() {
        var _el$5 = _$createElement("box");
        _$setProp(_el$5, "flexShrink", 0);
        _$setProp(_el$5, "flexDirection", "row");
        _$setProp(_el$5, "alignItems", "center");
        _$insert(_el$5, () => props.right?.(tone));
        return _el$5;
      }
    }), null);
    _$insert(_el$, () => props.below?.(tone), null);
    _$effect((_p$) => {
      var _v$ = props.id, _v$2 = props.selected ? presets2.select.selectedBackgroundColor : "transparent", _v$3 = theme2.mouse && props.onClick ? () => props.onClick?.() : undefined, _v$4 = tone(), _v$5 = props.bold ? TextAttributes.BOLD : undefined;
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "id", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "backgroundColor", _v$2, _p$.t));
      _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "onMouseDown", _v$3, _p$.a));
      _v$4 !== _p$.o && (_p$.o = _$setProp(_el$4, "fg", _v$4, _p$.o));
      _v$5 !== _p$.i && (_p$.i = _$setProp(_el$4, "attributes", _v$5, _p$.i));
      return _p$;
    }, {
      e: undefined,
      t: undefined,
      a: undefined,
      o: undefined,
      i: undefined
    });
    return _el$;
  })();
}

export { ListRow2 };
