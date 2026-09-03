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
  maxHeight?: number;
  scrollable?: boolean;
}

export function Modal(props: ModalProps) {
  const dimensions = useTerminalDimensions();
  const screenCap = () => Math.max(5, dimensions().height - 4);
  const maxPanelHeight = () => Math.min(props.maxHeight ?? screenCap(), screenCap());
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
            <Show
              when={props.scrollable}
              fallback={
                <box flexDirection="column" maxHeight={maxContentHeight()} overflow="hidden">
                  {props.children}
                </box>
              }
            >
              <scrollbox maxHeight={maxContentHeight()} flexShrink={1} scrollY scrollX={false}>
                {props.children}
              </scrollbox>
            </Show>
          </box>
        </box>
      </Show>
    </Portal>
  );
}
