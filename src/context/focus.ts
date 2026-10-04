import { CliRenderEvents, type Renderable } from "@opentui/core";
import { useRenderer } from "@opentui/solid";
import { createEffect, createSignal, onCleanup, type Accessor } from "solid-js";

export function isInside(node: Renderable | null, root: Renderable | null): boolean {
  if (!root) return false;
  for (let current = node; current; current = current.parent) {
    if (current === root) return true;
  }
  return false;
}

/**
 * True while the renderer's native focus is `target` or one of its
 * descendants. Mouse clicks move native focus on their own (the renderer
 * focuses the nearest focusable ancestor of what was clicked), so this
 * follows both keyboard and mouse focus changes.
 */
export function useFocusWithin(target: Accessor<Renderable | null | undefined>): Accessor<boolean> {
  const renderer = useRenderer();
  const [focused, setFocused] = createSignal(false);

  const sync = () => setFocused(isInside(renderer.currentFocusedRenderable, target() ?? null));

  renderer.on(CliRenderEvents.FOCUSED_RENDERABLE, sync);
  onCleanup(() => renderer.off(CliRenderEvents.FOCUSED_RENDERABLE, sync));
  createEffect(sync);

  return focused;
}
