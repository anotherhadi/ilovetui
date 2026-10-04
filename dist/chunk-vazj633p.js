// @bun
import {
  useFocusWithin2
} from "./chunk-59dzyjgd.js";

// src/components/Pane.tsx
import { createComponent as _$createComponent } from "@opentui/solid";
import { insert as _$insert } from "@opentui/solid";
import { use as _$use } from "@opentui/solid";
import { spread as _$spread } from "@opentui/solid";
import { createElement as _$createElement } from "@opentui/solid";
import { createContext, createSignal, splitProps, useContext } from "solid-js";
var PaneContext = createContext();
function Pane2(props) {
  const [local, rest] = splitProps(props, ["handle", "children"]);
  const [el, setEl] = createSignal();
  const focused = useFocusWithin2(el);
  const handle = {
    target: el,
    focused,
    focus: () => el()?.focus()
  };
  local.handle?.(handle);
  return _$createComponent(PaneContext.Provider, {
    value: handle,
    get children() {
      var _el$ = _$createElement("box");
      _$use(setEl, _el$);
      _$spread(_el$, rest, true);
      _$insert(_el$, () => local.children);
      return _el$;
    }
  });
}
function usePane2() {
  const pane = useContext(PaneContext);
  if (!pane)
    throw new Error("usePane() must be used inside a <Pane>");
  return pane;
}

export { Pane2, usePane2 };
