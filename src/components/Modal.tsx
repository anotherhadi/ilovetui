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
          paddingTop={Math.floor(dimensions().height / 4)}
          backgroundColor={props.backdropColor ?? RGBA.fromInts(0, 0, 0, 150)}
          zIndex={2000}
          onMouseDown={theme.mouse ? props.onDismiss : undefined}
        >
          <box
            width={props.width ?? 50}
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
            {props.children}
          </box>
        </box>
      </Show>
    </Portal>
  );
}
