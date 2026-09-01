import type { ColorInput } from "@opentui/core";
import { useTerminalDimensions } from "@opentui/solid";
import { Match, Switch, type ParentProps } from "solid-js";
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
    <Switch>
      <Match when={tooSmall()}>
        <box
          width={dimensions().width}
          height={dimensions().height}
          alignItems="center"
          justifyContent="center"
          backgroundColor={props.backgroundColor ?? theme.background}
        >
          <text fg={props.textColor ?? theme.muted}>
            {`Terminal too small\nMinimum size: ${minWidth()}x${minHeight()} — current: ${dimensions().width}x${dimensions().height}`}
          </text>
        </box>
      </Match>
      <Match when={!tooSmall()}>{props.children}</Match>
    </Switch>
  );
}
