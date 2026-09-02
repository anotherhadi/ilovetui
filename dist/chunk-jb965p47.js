// @bun
import {
  theme
} from "./chunk-57mam43k.js";

// src/components/Tabs.tsx
import { effect as _$effect } from "@opentui/solid";
import { insertNode as _$insertNode } from "@opentui/solid";
import { createComponent as _$createComponent } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { use as _$use } from "@opentui/solid";
import { setProp as _$setProp } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { BorderChars, CliRenderEvents, TextAttributes } from "@opentui/core";
import { onResize, useRenderer } from "@opentui/solid";
import { createMemo, createSignal, Index } from "solid-js";
function stepTabValue(items, value, delta) {
  const index = items.findIndex((t) => t.value === value);
  if (index === -1)
    return value;
  return items[(index + delta + items.length) % items.length].value;
}
function innerWidth(label) {
  return [...label].length + 2;
}
function segmentWidth(label) {
  return innerWidth(label) + 2;
}
function topEdge(chars, label) {
  return `${chars.topLeft}${chars.horizontal.repeat(innerWidth(label))}${chars.topRight}`;
}
function bottomEdge(chars, label, isFirst, isLast, isActive) {
  const left = isFirst ? isActive ? chars.vertical : chars.leftT : isActive ? chars.bottomRight : chars.bottomT;
  const right = isLast ? isActive ? chars.vertical : chars.rightT : isActive ? chars.bottomLeft : chars.bottomT;
  const fill = isActive ? " " : chars.horizontal;
  return left + fill.repeat(innerWidth(label)) + right;
}
function buildBottomRow(chars, segments, available) {
  const tabsWidth = segments.reduce((sum, s) => sum + segmentWidth(s.label), 0);
  const remainder = available - tabsWidth;
  const extending = remainder > 0;
  const tabs = segments.map((s, i) => bottomEdge(chars, s.label, i === 0, i === segments.length - 1 && !extending, s.isActive)).join("");
  return extending ? tabs + chars.horizontal.repeat(remainder - 1) + chars.topRight : tabs;
}
function fitCount(items, start, available) {
  const hiddenBefore = start;
  let used = 0;
  let count = 0;
  for (let i = start;i < items.length; i++) {
    const width = segmentWidth(items[i].label);
    const hiddenAfter = items.length - start - count - 1;
    const totalHidden = hiddenBefore + hiddenAfter;
    const reserve = totalHidden > 0 ? segmentWidth(`+${totalHidden}`) : 0;
    if (used + width + reserve > available)
      break;
    used += width;
    count++;
  }
  return Math.max(count, 1);
}
function Tabs(props) {
  const renderer = useRenderer();
  const [width, setWidth] = createSignal(0);
  let container;
  const measure = () => {
    if (!container)
      return;
    setWidth(container.width);
  };
  const remeasureNextFrame = () => renderer.once(CliRenderEvents.FRAME, measure);
  onResize(remeasureNextFrame);
  const segments = createMemo(() => {
    const items = props.items;
    const activeIndex = Math.max(0, items.findIndex((t) => t.value === props.value));
    let start = 0;
    let count = fitCount(items, start, width());
    if (activeIndex >= start + count) {
      start = activeIndex;
      count = fitCount(items, start, width());
    }
    const shownItems = items.slice(start, start + count);
    const shown = shownItems.map((item) => ({
      value: item.value,
      label: item.label,
      isActive: item.value === props.value
    }));
    const hidden = items.length - count;
    if (hidden > 0) {
      const used = shownItems.reduce((sum, item) => sum + segmentWidth(item.label), 0);
      const tag = `+${hidden}`;
      if (shown.length === 0 || used + segmentWidth(tag) <= width()) {
        shown.push({
          value: undefined,
          label: tag,
          isActive: false
        });
      }
    }
    return shown;
  });
  const accent = () => props.accentColor ?? theme.primary;
  const muted = () => props.mutedColor ?? theme.muted;
  const borderColor = () => props.focused ? accent() : muted();
  const chars = BorderChars[theme.borderStyle];
  return (() => {
    var _el$ = _$createElement("box"), _el$2 = _$createElement("text"), _el$3 = _$createElement("box"), _el$4 = _$createElement("text");
    _$insertNode(_el$, _el$2);
    _$insertNode(_el$, _el$3);
    _$insertNode(_el$, _el$4);
    _$use((el) => {
      container = el;
      remeasureNextFrame();
    }, _el$);
    _$setProp(_el$, "flexDirection", "column");
    _$setProp(_el$, "height", 3);
    _$setProp(_el$2, "selectable", false);
    _$insert(_el$2, () => segments().map((s) => topEdge(chars, s.label)).join(""));
    _$setProp(_el$3, "flexDirection", "row");
    _$insert(_el$3, _$createComponent(Index, {
      get each() {
        return segments();
      },
      children: (s) => [(() => {
        var _el$5 = _$createElement("text");
        _$setProp(_el$5, "selectable", false);
        _$insert(_el$5, () => chars.vertical);
        _$effect((_$p) => _$setProp(_el$5, "fg", borderColor(), _$p));
        return _el$5;
      })(), (() => {
        var _el$6 = _$createElement("text");
        _$setProp(_el$6, "selectable", false);
        _$insert(_el$6, () => ` ${s().label} `);
        _$effect((_p$) => {
          var _v$3 = s().isActive ? accent() : muted(), _v$4 = s().isActive ? TextAttributes.BOLD : undefined, _v$5 = theme.mouse && s().value !== undefined ? () => props.onChange(s().value) : undefined;
          _v$3 !== _p$.e && (_p$.e = _$setProp(_el$6, "fg", _v$3, _p$.e));
          _v$4 !== _p$.t && (_p$.t = _$setProp(_el$6, "attributes", _v$4, _p$.t));
          _v$5 !== _p$.a && (_p$.a = _$setProp(_el$6, "onMouseDown", _v$5, _p$.a));
          return _p$;
        }, {
          e: undefined,
          t: undefined,
          a: undefined
        });
        return _el$6;
      })(), (() => {
        var _el$7 = _$createElement("text");
        _$setProp(_el$7, "selectable", false);
        _$insert(_el$7, () => chars.vertical);
        _$effect((_$p) => _$setProp(_el$7, "fg", borderColor(), _$p));
        return _el$7;
      })()]
    }));
    _$setProp(_el$4, "selectable", false);
    _$insert(_el$4, () => buildBottomRow(chars, segments(), width()));
    _$effect((_p$) => {
      var _v$ = borderColor(), _v$2 = borderColor();
      _v$ !== _p$.e && (_p$.e = _$setProp(_el$2, "fg", _v$, _p$.e));
      _v$2 !== _p$.t && (_p$.t = _$setProp(_el$4, "fg", _v$2, _p$.t));
      return _p$;
    }, {
      e: undefined,
      t: undefined
    });
    return _el$;
  })();
}

export { stepTabValue, Tabs };
