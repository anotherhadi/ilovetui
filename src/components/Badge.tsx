import { type ColorInput, TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
import { theme } from "../index.ts";

const CAP_LEFT = "";
const CAP_RIGHT = "";

export interface BadgeProps {
  label: string;
  color: ColorInput;
  textColor?: ColorInput;
  backgroundColor?: ColorInput;
  bold?: boolean;
  withNerdfont?: boolean;
}

export function Badge(props: BadgeProps) {
  const textColor = () => props.textColor ?? theme.background;
  const backgroundColor = () => props.backgroundColor ?? theme.background;
  const attributes = () => (props.bold === false ? undefined : TextAttributes.BOLD);
  const nerdFonts = () => props.withNerdfont ?? theme.nerdFonts;

  return (
    <box flexDirection="row">
      <Show when={nerdFonts()}>
        <text fg={props.color} bg={backgroundColor()}>
          {CAP_LEFT}
        </text>
      </Show>
      <text fg={textColor()} bg={props.color} attributes={attributes()}>
        {nerdFonts() ? props.label : ` ${props.label} `}
      </text>
      <Show when={nerdFonts()}>
        <text fg={props.color} bg={backgroundColor()}>
          {CAP_RIGHT}
        </text>
      </Show>
    </box>
  );
}
