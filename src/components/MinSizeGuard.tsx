import type { ColorInput } from "@opentui/core";
import { useTerminalDimensions } from "@opentui/solid";
import { Show, type ParentProps } from "solid-js";
import { theme } from "../index.ts";

export interface MinSizeGuardProps extends ParentProps {
  minWidth?: number;
  minHeight?: number;
  backgroundColor?: ColorInput;
  textColor?: ColorInput;
}

export function MinSizeGuard(props: MinSizeGuardProps) {
  const dimensions = useTerminalDimensions();
  const minWidth = () => props.minWidth ?? 80;
  const minHeight = () => props.minHeight ?? 24;
  const tooSmall = () => dimensions().width < minWidth() || dimensions().height < minHeight();

  return (
    <>
      {props.children}
      <Show when={tooSmall()}>
        <box
          position="absolute"
          top={0}
          left={0}
          width={dimensions().width}
          height={dimensions().height}
          alignItems="center"
          justifyContent="center"
          backgroundColor={props.backgroundColor ?? theme.background}
          zIndex={9999}
        >
          <text fg={props.textColor ?? theme.muted}>
            {`Terminal too small\nMinimum size: ${minWidth()}x${minHeight()} — current: ${dimensions().width}x${dimensions().height}`}
          </text>
        </box>
      </Show>
    </>
  );
}
