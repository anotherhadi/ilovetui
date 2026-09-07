// @bun
import {
  presets,
  theme
} from "./chunk-x7m82z3z.js";

// src/components/Sidebar.tsx
import { use as _$use } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { createTextNode as _$createTextNode } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { effect as _$effect } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
function itemIndexAtScreenY(el, screenY) {
  const internals = el;
  if (!internals.linesPerItem)
    return null;
  const localY = screenY - el.screenY;
  if (localY < 0)
    return null;
  const index = internals.scrollOffset + Math.floor(localY / internals.linesPerItem);
  return index >= 0 && index < el.options.length ? index : null;
}
function Sidebar(props) {
  let select;
  const handleMouseDown = (event) => {
    if (!select)
      return;
    const index = itemIndexAtScreenY(select, event.y);
    if (index !== null) {
      select.setSelectedIndex(index);
      select.selectCurrent();
    }
  };
  const handleMouseScroll = (event) => {
    if (!select || !event.scroll)
      return;
    const steps = Math.max(1, event.scroll.delta);
    if (event.scroll.direction === "down")
      select.moveDown(steps);
    else if (event.scroll.direction === "up")
      select.moveUp(steps);
  };
  return (() => {
    var _el$ = _$createElement("box"), _el$5 = _$createElement("select");
    _$insertNode(_el$, _el$5);
    _$setProp(_el$, "border", true);
    _$setProp(_el$, "flexDirection", "column");
    _$insert(_el$, _$createComponent(Show, {
      get when() {
        return props.title;
      },
      get children() {
        return [(() => {
          var _el$2 = _$createElement("text");
          _$insert(_el$2, () => props.title);
          _$effect((_$p) => _$setProp(_el$2, "attributes", TextAttributes.BOLD, _$p));
          return _el$2;
        })(), (() => {
          var _el$3 = _$createElement("text");
          _$insertNode(_el$3, _$createTextNode(` `));
          return _el$3;
        })()];
      }
    }), _el$5);
    _$use((el) => {
      select = el ?? undefined;
      props.ref?.(el && {
        moveUp: () => el.moveUp(),
        moveDown: () => el.moveDown(),
        selectCurrent: () => el.selectCurrent(),
        focus: () => el.focus()
      });
    }, _el$5);
    _$setProp(_el$5, "flexGrow", 1);
    _$setProp(_el$5, "showScrollIndicator", true);
    _$setProp(_el$5, "onChange", (_index, option) => {
      if (!option)
        return;
      props.onSelect(option);
    });
    _$setProp(_el$5, "onSelect", (_index, option) => {
      if (!option)
        return;
      props.onConfirm?.(option);
    });
    _$effect((_p$) => {
      var _v$ = props.width ?? 24, _v$2 = theme.borderStyle, _v$3 = props.focused ? props.accentColor ?? theme.primary : props.mutedColor ?? theme.muted, _v$4 = props.focused, _v$5 = props.items, _v$6 = theme.mouse ? handleMouseDown : undefined, _v$7 = theme.mouse ? handleMouseScroll : undefined, _v$8 = props.backgroundColor ?? "transparent", _v$9 = props.textColor ?? presets.select.textColor, _v$0 = props.focusedBackgroundColor ?? "transparent", _v$1 = props.focusedTextColor ?? presets.select.focusedTextColor, _v$10 = props.selectedBackgroundColor ?? presets.select.selectedBackgroundColor, _v$11 = props.selectedTextColor ?? presets.select.selectedTextColor, _v$12 = props.descriptionColor ?? presets.select.descriptionColor, _v$13 = props.selectedDescriptionColor ?? presets.select.selectedDescriptionColor;
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$, "width", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$, "borderStyle", _v$2, _p$.t));
      _v$3 !== _p$.a && (_p$.a = _$setProp(_el$, "borderColor", _v$3, _p$.a));
      _v$4 !== _p$.o && (_p$.o = _$setProp(_el$5, "focused", _v$4, _p$.o));
      _v$5 !== _p$.i && (_p$.i = _$setProp(_el$5, "options", _v$5, _p$.i));
      _v$6 !== _p$.n && (_p$.n = _$setProp(_el$5, "onMouseDown", _v$6, _p$.n));
      _v$7 !== _p$.s && (_p$.s = _$setProp(_el$5, "onMouseScroll", _v$7, _p$.s));
      _v$8 !== _p$.h && (_p$.h = _$setProp(_el$5, "backgroundColor", _v$8, _p$.h));
      _v$9 !== _p$.r && (_p$.r = _$setProp(_el$5, "textColor", _v$9, _p$.r));
      _v$0 !== _p$.d && (_p$.d = _$setProp(_el$5, "focusedBackgroundColor", _v$0, _p$.d));
      _v$1 !== _p$.l && (_p$.l = _$setProp(_el$5, "focusedTextColor", _v$1, _p$.l));
      _v$10 !== _p$.u && (_p$.u = _$setProp(_el$5, "selectedBackgroundColor", _v$10, _p$.u));
      _v$11 !== _p$.c && (_p$.c = _$setProp(_el$5, "selectedTextColor", _v$11, _p$.c));
      _v$12 !== _p$.w && (_p$.w = _$setProp(_el$5, "descriptionColor", _v$12, _p$.w));
      _v$13 !== _p$.m && (_p$.m = _$setProp(_el$5, "selectedDescriptionColor", _v$13, _p$.m));
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
      w: undefined,
      m: undefined
    });
    return _el$;
  })();
}

export { Sidebar };
