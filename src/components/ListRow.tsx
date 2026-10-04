import { type ColorInput, TextAttributes } from "@opentui/core";
import { Show, type JSX } from "solid-js";
import { presets, theme } from "../index.ts";

/** Returns the selected-row text color when the row is selected, `fallback` otherwise. */
export type RowTone = (fallback?: ColorInput) => ColorInput | undefined;

export interface ListRowProps {
  id?: string;
  label: string;
  selected: boolean;
  bold?: boolean;
  onClick?: () => void;
  /** Right-aligned content; use `tone` for its colors so it follows the selection. */
  right?: (tone: RowTone) => JSX.Element;
  /** Extra lines under the main one. */
  below?: (tone: RowTone) => JSX.Element;
}

/** One row of a keyboard-navigable list: truncated label on the left, details on the right. */
export function ListRow(props: ListRowProps) {
  const tone: RowTone = (fallback) => (props.selected ? presets.select.selectedTextColor : fallback);

  return (
    <box
      id={props.id}
      flexShrink={0}
      flexDirection="column"
      paddingLeft={1}
      paddingRight={1}
      backgroundColor={props.selected ? presets.select.selectedBackgroundColor : "transparent"}
      onMouseDown={theme.mouse && props.onClick ? () => props.onClick?.() : undefined}
    >
      <box flexDirection="row" justifyContent="space-between" alignItems="center" height={1}>
        <box flexGrow={1} flexShrink={1} minWidth={0} marginRight={1} overflow="hidden">
          <text fg={tone()} attributes={props.bold ? TextAttributes.BOLD : undefined} wrapMode="none">
            {props.label}
          </text>
        </box>
        <Show when={props.right}>
          <box flexShrink={0} flexDirection="row" alignItems="center">
            {props.right?.(tone)}
          </box>
        </Show>
      </box>
      {props.below?.(tone)}
    </box>
  );
}
