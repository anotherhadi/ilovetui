// @bun
import {
  theme
} from "./chunk-57mam43k.js";

// src/components/Sidebar.tsx
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { use as _$use } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
function Sidebar(props) {
  return (() => {
    var _el$ = _$createElement("box"), _el$2 = _$createElement("text"), _el$3 = _$createElement("select");
    _$insertNode(_el$, _el$2);
    _$insertNode(_el$, _el$3);
    _$setProp(_el$, "border", true);
    _$setProp(_el$, "flexDirection", "column");
    _$insert(_el$2, () => props.title);
    _$use((el) => {
      props.ref?.(el && {
        moveUp: () => el.moveUp(),
        moveDown: () => el.moveDown(),
        selectCurrent: () => el.selectCurrent()
      });
    }, _el$3);
    _$setProp(_el$3, "flexGrow", 1);
    _$setProp(_el$3, "showScrollIndicator", true);
    _$setProp(_el$3, "onChange", (_index, option) => {
      if (!option)
        return;
      props.onSelect(option);
    });
    _$setProp(_el$3, "onSelect", (_index, option) => {
      if (!option)
        return;
      props.onConfirm?.(option);
    });
    _$effect((_p$) => {
      var _v$ = props.width ?? 24, _v$2 = theme.borderStyle, _v$3 = props.focused ? props.accentColor ?? theme.primary : props.mutedColor ?? theme.muted, _v$4 = TextAttributes.BOLD, _v$5 = props.focused, _v$6 = props.items, _v$7 = props.backgroundColor, _v$8 = props.textColor, _v$9 = props.focusedBackgroundColor, _v$0 = props.focusedTextColor, _v$1 = props.selectedBackgroundColor, _v$10 = props.selectedTextColor, _v$11 = props.descriptionColor, _v$12 = props.selectedDescriptionColor;
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "width", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "borderStyle", _v$2, _p$.t));
      _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "borderColor", _v$3, _p$.a));
      _v$4 !== _p$.o && (_p$.o = _$setProp(_el$2, "attributes", _v$4, _p$.o));
      _v$5 !== _p$.i && (_p$.i = _$setProp(_el$3, "focused", _v$5, _p$.i));
      _v$6 !== _p$.n && (_p$.n = _$setProp(_el$3, "options", _v$6, _p$.n));
      _v$7 !== _p$.s && (_p$.s = _$setProp(_el$3, "backgroundColor", _v$7, _p$.s));
      _v$8 !== _p$.h && (_p$.h = _$setProp(_el$3, "textColor", _v$8, _p$.h));
      _v$9 !== _p$.r && (_p$.r = _$setProp(_el$3, "focusedBackgroundColor", _v$9, _p$.r));
      _v$0 !== _p$.d && (_p$.d = _$setProp(_el$3, "focusedTextColor", _v$0, _p$.d));
      _v$1 !== _p$.l && (_p$.l = _$setProp(_el$3, "selectedBackgroundColor", _v$1, _p$.l));
      _v$10 !== _p$.u && (_p$.u = _$setProp(_el$3, "selectedTextColor", _v$10, _p$.u));
      _v$11 !== _p$.c && (_p$.c = _$setProp(_el$3, "descriptionColor", _v$11, _p$.c));
      _v$12 !== _p$.w && (_p$.w = _$setProp(_el$3, "selectedDescriptionColor", _v$12, _p$.w));
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
      d: undefined,
      l: undefined,
      u: undefined,
      c: undefined,
      w: undefined
    });
    return _el$;
  })();
}

export { Sidebar };
