import { type ColorInput, type MouseEvent, type SelectOption, type SelectRenderable, TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
import { presets, theme } from "../index.ts";

export interface SidebarHandle {
  moveUp(): void;
  moveDown(): void;
  selectCurrent(): void;
  focus(): void;
}

function itemIndexAtScreenY(el: SelectRenderable, screenY: number): number | null {
  const internals = el as unknown as { scrollOffset: number; linesPerItem: number };
  if (!internals.linesPerItem) return null;
  const localY = screenY - el.screenY;
  if (localY < 0) return null;
  const index = internals.scrollOffset + Math.floor(localY / internals.linesPerItem);
  return index >= 0 && index < el.options.length ? index : null;
}

export interface SidebarProps {
  title?: string;
  items: SelectOption[];
  width?: number;
  focused?: boolean;
  onSelect: (option: SelectOption) => void;
  onConfirm?: (option: SelectOption) => void;
  ref?: (handle: SidebarHandle | null) => void;

  accentColor?: ColorInput;
  mutedColor?: ColorInput;

  backgroundColor?: ColorInput;
  textColor?: ColorInput;
  focusedBackgroundColor?: ColorInput;
  focusedTextColor?: ColorInput;
  selectedBackgroundColor?: ColorInput;
  selectedTextColor?: ColorInput;
  descriptionColor?: ColorInput;
  selectedDescriptionColor?: ColorInput;
}

export function Sidebar(props: SidebarProps) {
  let select: SelectRenderable | undefined;

  const handleMouseDown = (event: MouseEvent) => {
    if (!select) return;
    const index = itemIndexAtScreenY(select, event.y);
    if (index !== null) select.setSelectedIndex(index);
  };

  const handleMouseScroll = (event: MouseEvent) => {
    if (!select || !event.scroll) return;
    const steps = Math.max(1, event.scroll.delta);
    if (event.scroll.direction === "down") select.moveDown(steps);
    else if (event.scroll.direction === "up") select.moveUp(steps);
  };

  return (
    <box
      width={props.width ?? 24}
      border
      borderStyle={theme.borderStyle}
      borderColor={props.focused ? (props.accentColor ?? theme.primary) : (props.mutedColor ?? theme.muted)}
      flexDirection="column"
    >
      <Show when={props.title}>
        <text attributes={TextAttributes.BOLD}>{props.title}</text>
        <text> </text>
      </Show>
      <select
        ref={(el) => {
          select = el ?? undefined;
          props.ref?.(
            el && {
              moveUp: () => el.moveUp(),
              moveDown: () => el.moveDown(),
              selectCurrent: () => el.selectCurrent(),
              focus: () => el.focus(),
            },
          );
        }}
        flexGrow={1}
        focused={props.focused}
        showScrollIndicator
        options={props.items}
        onMouseDown={theme.mouse ? handleMouseDown : undefined}
        onMouseScroll={theme.mouse ? handleMouseScroll : undefined}
        backgroundColor={props.backgroundColor ?? "transparent"}
        textColor={props.textColor ?? presets.select.textColor}
        focusedBackgroundColor={props.focusedBackgroundColor ?? "transparent"}
        focusedTextColor={props.focusedTextColor ?? presets.select.focusedTextColor}
        selectedBackgroundColor={props.selectedBackgroundColor ?? presets.select.selectedBackgroundColor}
        selectedTextColor={props.selectedTextColor ?? presets.select.selectedTextColor}
        descriptionColor={props.descriptionColor ?? presets.select.descriptionColor}
        selectedDescriptionColor={props.selectedDescriptionColor ?? presets.select.selectedDescriptionColor}
        onChange={(_index, option) => {
          if (!option) return;
          props.onSelect(option);
        }}
        onSelect={(_index, option) => {
          if (!option) return;
          props.onConfirm?.(option);
        }}
      />
    </box>
  );
}
