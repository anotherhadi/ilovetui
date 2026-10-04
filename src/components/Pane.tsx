import type { BoxRenderable } from "@opentui/core";
import type { BoxProps } from "@opentui/solid";
import { type Accessor, createContext, createSignal, splitProps, useContext } from "solid-js";
import { useFocusWithin } from "../context/focus.ts";

export interface PaneHandle {
  /** The pane's root renderable — use it as a keymap layer `target`. */
  target: Accessor<BoxRenderable | undefined>;
  /** True while native focus is anywhere inside the pane. */
  focused: Accessor<boolean>;
  focus(): void;
}

const PaneContext = createContext<PaneHandle>();

export interface PaneProps extends Omit<BoxProps, "ref"> {
  /** Receives the pane's handle, e.g. to focus it or target it from outside. */
  handle?: (handle: PaneHandle) => void;
}

/**
 * A focus region. Wrap each independently focusable area (sidebar, content,
 * …) in a Pane; anything inside can call usePane() to scope keymap layers to
 * it (`useBindings(() => ({ target: pane.target, … }))`) and to style itself
 * from `pane.focused()`. Pass `focusable` when the pane itself should be the
 * focus target (rather than e.g. a <select> inside it).
 */
export function Pane(props: PaneProps) {
  const [local, rest] = splitProps(props, ["handle", "children"]);
  const [el, setEl] = createSignal<BoxRenderable>();
  const focused = useFocusWithin(el);
  const handle: PaneHandle = { target: el, focused, focus: () => el()?.focus() };
  local.handle?.(handle);

  return (
    <PaneContext.Provider value={handle}>
      <box {...rest} ref={setEl}>
        {local.children}
      </box>
    </PaneContext.Provider>
  );
}

export function usePane(): PaneHandle {
  const pane = useContext(PaneContext);
  if (!pane) throw new Error("usePane() must be used inside a <Pane>");
  return pane;
}
