// @bun
// src/keymap/provider.tsx
import { createComponent as _$createComponent } from "@opentui/solid";
import { createDefaultOpenTuiKeymap } from "@opentui/keymap/opentui";
import { registerCommaBindings } from "@opentui/keymap/addons";
import { KeymapProvider as BaseKeymapProvider } from "@opentui/keymap/solid";
import { useRenderer } from "@opentui/solid";
import { onCleanup } from "solid-js";
import { useKeymap, useBindings, useKeymapSelector, reactiveMatcherFromSignal } from "@opentui/keymap/solid";
function KeymapProvider(props) {
  const renderer = useRenderer();
  const keymap = createDefaultOpenTuiKeymap(renderer);
  onCleanup(registerCommaBindings(keymap));
  return _$createComponent(BaseKeymapProvider, {
    keymap,
    get children() {
      return props.children;
    }
  });
}
export {
  useKeymapSelector,
  useKeymap,
  useBindings,
  reactiveMatcherFromSignal,
  KeymapProvider
};
