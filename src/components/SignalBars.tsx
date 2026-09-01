import type { ColorInput } from "@opentui/core";
import { For } from "solid-js";
import { theme } from "../index.ts";

const BARS = ["▁", "▃", "▅", "▇"] as const;

export interface SignalBarsProps {
  signal: number;
  filledColor?: ColorInput;
  emptyColor?: ColorInput;
}

export function SignalBars(props: SignalBarsProps) {
  const level = () => Math.floor((props.signal * BARS.length) / 100);

  return (
    <box flexDirection="row">
      <For each={BARS}>
        {(bar, i) => <text fg={i() < level() ? props.filledColor : (props.emptyColor ?? theme.muted)}>{bar}</text>}
      </For>
    </box>
  );
}
