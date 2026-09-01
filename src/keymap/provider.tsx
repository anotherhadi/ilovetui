import { createDefaultOpenTuiKeymap } from "@opentui/keymap/opentui";
import { registerCommaBindings } from "@opentui/keymap/addons";
import { KeymapProvider as BaseKeymapProvider } from "@opentui/keymap/solid";
import { useRenderer } from "@opentui/solid";
import { onCleanup, type ParentProps } from "solid-js";

export { useKeymap, useBindings, useKeymapSelector, reactiveMatcherFromSignal } from "@opentui/keymap/solid";

export function KeymapProvider(props: ParentProps) {
  const renderer = useRenderer();
  const keymap = createDefaultOpenTuiKeymap(renderer);
  onCleanup(registerCommaBindings(keymap));

  return <BaseKeymapProvider keymap={keymap}>{props.children}</BaseKeymapProvider>;
}
