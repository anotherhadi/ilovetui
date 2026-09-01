import type { ColorInput } from "@opentui/core";
import "opentui-spinner/solid";
import type { ColorGenerator, SpinnerOptions } from "opentui-spinner";
import { theme } from "../index.ts";

export interface SpinnerProps {
  name?: SpinnerOptions["name"];
  frames?: string[];
  interval?: number;
  autoplay?: boolean;
  color?: ColorInput | ColorGenerator;
  backgroundColor?: ColorInput;
}

export function Spinner(props: SpinnerProps) {
  return (
    <spinner
      name={props.name}
      frames={props.frames}
      interval={props.interval}
      autoplay={props.autoplay ?? !theme.reducedMotion}
      color={props.color ?? theme.primary}
      backgroundColor={props.backgroundColor}
    />
  );
}
