import { type ColorInput, TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
import { theme } from "../index.ts";
import { Badge } from "./Badge.tsx";

export interface ToggleProps {
  on: boolean;
  onToggle?: () => void;
  onColor?: ColorInput;
  offColor?: ColorInput;
  withNerdfont?: boolean;
}

export function Toggle(props: ToggleProps) {
  const onColor = () => props.onColor ?? theme.success;
  const offColor = () => props.offColor ?? theme.muted;
  const nerdFonts = () => props.withNerdfont ?? theme.nerdFonts;

  return (
    <box onMouseDown={theme.mouse ? props.onToggle : undefined}>
      <Show
        when={nerdFonts()}
        fallback={
          <text fg={props.on ? onColor() : offColor()} attributes={TextAttributes.BOLD}>
            {props.on ? "on" : "off"}
          </text>
        }
      >
        <Badge label={props.on ? "  ●" : "●  "} color={props.on ? onColor() : offColor()} withNerdfont />
      </Show>
    </box>
  );
}
