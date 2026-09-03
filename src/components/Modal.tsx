import { type ColorInput, type MouseEvent, RGBA } from "@opentui/core";
import { Portal, useTerminalDimensions } from "@opentui/solid";
import { Show, type ParentProps } from "solid-js";
import { theme } from "../index.ts";

export interface ModalProps extends ParentProps {
  open: boolean;
  onDismiss?: () => void;
  accentColor?: ColorInput;
  backgroundColor?: ColorInput;
  backdropColor?: ColorInput;
  width?: number;
}

export function Modal(props: ModalProps) {
  const dimensions = useTerminalDimensions();
  // Leave a margin of screen rows above/below so the panel never touches the
  // terminal edges, then reserve the border + vertical padding from that
  // budget for the scrollable content itself. On a short terminal (or with
  // a lot of content, e.g. HelpModal's keybinding list) this keeps the panel
  // fully on-screen — with a scrollbar — instead of spilling past the bottom
  // edge with no way to see or reach the rest of it.
  const maxPanelHeight = () => Math.max(5, dimensions().height - 4);
  const maxContentHeight = () => Math.max(1, maxPanelHeight() - 4);

  return (
    <Portal>
      <Show when={props.open}>
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
            width={props.width ?? 50}
            maxHeight={maxPanelHeight()}
            border
            borderStyle={theme.borderStyle}
            borderColor={props.accentColor ?? theme.primary}
            backgroundColor={props.backgroundColor ?? theme.background}
            flexDirection="column"
            paddingLeft={2}
            paddingRight={2}
            paddingTop={1}
            paddingBottom={1}
            onMouseDown={theme.mouse ? (event: MouseEvent) => event.stopPropagation() : undefined}
          >
            <scrollbox maxHeight={maxContentHeight()} flexShrink={1} scrollY scrollX={false}>
              {props.children}
            </scrollbox>
          </box>
        </box>
      </Show>
    </Portal>
  );
}
