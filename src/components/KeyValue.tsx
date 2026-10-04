import type { ColorInput } from "@opentui/core";
import { Show, type JSX } from "solid-js";
import { theme } from "../index.ts";

export interface KeyValueProps {
  label: string;
  /** Plain value; ignored when children are given. */
  value?: string | number;
  valueColor?: ColorInput;
  children?: JSX.Element;
}

/** A "Label ........ value" line, as used in details panels and modals. */
export function KeyValue(props: KeyValueProps) {
  return (
    <box flexShrink={0} flexDirection="row" justifyContent="space-between">
      <text fg={theme.muted}>{props.label}</text>
      <Show when={props.children} fallback={<text fg={props.valueColor}>{props.value ?? ""}</text>}>
        {props.children}
      </Show>
    </box>
  );
}
