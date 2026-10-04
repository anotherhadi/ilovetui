// @bun
// src/context/focus.ts
import { CliRenderEvents } from "@opentui/core";
import { useRenderer } from "@opentui/solid";
import { createEffect, createSignal, onCleanup } from "solid-js";
function isInside2(node, root) {
  if (!root)
    return false;
  for (let current = node;current; current = current.parent) {
    if (current === root)
      return true;
  }
  return false;
}
function useFocusWithin2(target) {
  const renderer = useRenderer();
  const [focused, setFocused] = createSignal(false);
  const sync = () => setFocused(isInside2(renderer.currentFocusedRenderable, target() ?? null));
  renderer.on(CliRenderEvents.FOCUSED_RENDERABLE, sync);
  onCleanup(() => renderer.off(CliRenderEvents.FOCUSED_RENDERABLE, sync));
  createEffect(sync);
  return focused;
}

export { isInside2, useFocusWithin2 };
