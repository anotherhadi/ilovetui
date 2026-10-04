import { type BoxRenderable, type ColorInput, type KeyEvent, type MouseEvent, type Renderable, RGBA } from "@opentui/core";
import type { Binding, Command } from "@opentui/keymap";
import { useBindings } from "@opentui/keymap/solid";
import { Portal, useRenderer, useTerminalDimensions } from "@opentui/solid";
import { createSignal, onCleanup, onMount, Show, type ParentProps } from "solid-js";
import { theme } from "../index.ts";

export type ModalCommand = Command<Renderable, KeyEvent>;
export type ModalBinding = Binding<Renderable, KeyEvent>;

export interface ModalProps extends ParentProps {
  open: boolean;
  onDismiss?: () => void;
  /**
   * Key that calls onDismiss while the modal has focus (default "escape").
   * `false` disables it. Needs a KeymapProvider unless disabled.
   */
  dismissKey?: string | false;
  /** Extra keymap commands/bindings, active only while the modal has focus. */
  commands?: ModalCommand[];
  bindings?: ModalBinding[];
  /** Called on open with whatever had focus before the modal took it. */
  onOpen?: (previousFocus: Renderable | null) => void;
  accentColor?: ColorInput;
  backgroundColor?: ColorInput;
  backdropColor?: ColorInput;
  width?: number;
  maxHeight?: number;
  scrollable?: boolean;
}

export function Modal(props: ModalProps) {
  return (
    <Portal>
      <Show when={props.open}>
        <ModalPanel {...props} />
      </Show>
    </Portal>
  );
}

// Mounted only while open: takes native focus on mount and hands it back to
// whatever had it on unmount, so keymap layers scoped to the page underneath
// go quiet while the modal is up and come back on their own afterwards.
function ModalPanel(props: ModalProps) {
  const renderer = useRenderer();
  const dimensions = useTerminalDimensions();
  const [panel, setPanel] = createSignal<BoxRenderable>();

  const screenCap = () => Math.max(5, dimensions().height - 4);
  const maxPanelHeight = () => Math.min(props.maxHeight ?? screenCap(), screenCap());
  const maxContentHeight = () => Math.max(1, maxPanelHeight() - 4);
  const maxPanelWidth = () => Math.max(10, dimensions().width - 4);
  const panelWidth = () => Math.min(props.width ?? 50, maxPanelWidth());

  let previousFocus: Renderable | null = null;
  onMount(() => {
    previousFocus = renderer.currentFocusedRenderable;
    props.onOpen?.(previousFocus);
    panel()?.focus();
  });
  onCleanup(() => {
    if (previousFocus && !previousFocus.isDestroyed) previousFocus.focus();
  });

  const dismissKey = props.dismissKey ?? "escape";
  if (dismissKey !== false || props.commands || props.bindings) {
    useBindings(() => ({
      target: panel,
      commands: [
        ...(dismissKey !== false
          ? [{ name: "modal_dismiss", desc: "close", shortHelp: true, run: () => props.onDismiss?.() }]
          : []),
        ...(props.commands ?? []),
      ],
      bindings: [...(dismissKey !== false ? [{ key: dismissKey, cmd: "modal_dismiss" }] : []), ...(props.bindings ?? [])],
    }));
  }

  return (
    <box
      width={dimensions().width}
      height={dimensions().height}
      position="absolute"
      top={0}
      left={0}
      alignItems="center"
      justifyContent="center"
      backgroundColor={props.backdropColor ?? RGBA.fromInts(0, 0, 0, 150)}
      zIndex={2000}
      onMouseDown={theme.mouse ? props.onDismiss : undefined}
    >
      <box
        ref={setPanel}
        focusable
        width={panelWidth()}
        maxHeight={maxPanelHeight()}
        border
        borderStyle={theme.borderStyle}
        borderColor={props.accentColor ?? theme.primary}
        focusedBorderColor={props.accentColor ?? theme.primary}
        backgroundColor={props.backgroundColor ?? theme.background}
        flexDirection="column"
        paddingLeft={2}
        paddingRight={2}
        paddingTop={1}
        paddingBottom={1}
        onMouseDown={theme.mouse ? (event: MouseEvent) => event.stopPropagation() : undefined}
      >
        <Show
          when={props.scrollable}
          fallback={
            <box flexDirection="column" maxHeight={maxContentHeight()} overflow="hidden">
              {props.children}
            </box>
          }
        >
          <scrollbox maxHeight={maxContentHeight()} flexShrink={1} scrollY scrollX={false} focusable={false}>
            {props.children}
          </scrollbox>
        </Show>
      </box>
    </box>
  );
}
