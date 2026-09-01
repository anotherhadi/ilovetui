import { type ColorInput, parseColor, RGBA } from "@opentui/core";
import { createMemo, For } from "solid-js";
import { theme } from "../index.ts";

function sampleGradient(stops: ColorInput[], t: number): ColorInput {
  if (stops.length === 1) return stops[0]!;
  const clamped = Math.min(Math.max(t, 0), 1);
  const segment = clamped * (stops.length - 1);
  const i = Math.min(Math.floor(segment), stops.length - 2);
  const localT = segment - i;
  const [ar, ag, ab, aa] = parseColor(stops[i]!).toInts();
  const [br, bg, bb, ba] = parseColor(stops[i + 1]!).toInts();
  return RGBA.fromInts(
    Math.round(ar + (br - ar) * localT),
    Math.round(ag + (bg - ag) * localT),
    Math.round(ab + (bb - ab) * localT),
    Math.round(aa + (ba - aa) * localT),
  );
}

export interface ProgressBarProps {
  value: number;
  width?: number;
  color?: ColorInput | ColorInput[];
  trackColor?: ColorInput;
  fillChar?: string;
  trackChar?: string;
}

export function ProgressBar(props: ProgressBarProps) {
  const width = () => props.width ?? 20;
  const stops = createMemo(() => (Array.isArray(props.color) ? props.color : [props.color ?? theme.primary]));
  const filled = createMemo(() => Math.round((width() * Math.min(Math.max(props.value, 0), 100)) / 100));
  const cells = createMemo(() => Array.from({ length: width() }, (_, i) => i));

  return (
    <box flexDirection="row">
      <For each={cells()}>
        {(i) => {
          const isFilled = () => i < filled();
          const t = width() > 1 ? i / (width() - 1) : 0;
          return (
            <text fg={isFilled() ? sampleGradient(stops(), t) : (props.trackColor ?? theme.muted)}>
              {isFilled() ? (props.fillChar ?? "█") : (props.trackChar ?? "░")}
            </text>
          );
        }}
      </For>
    </box>
  );
}
