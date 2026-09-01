import { type ColorInput, type SelectOption, TextAttributes } from "@opentui/core";
import { theme } from "../index.ts";

export interface SidebarHandle {
  moveUp(): void;
  moveDown(): void;
  selectCurrent(): void;
}

export interface SidebarProps {
  title: string;
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
  return (
    <box
      width={props.width ?? 24}
      border
      borderStyle={theme.borderStyle}
      borderColor={props.focused ? (props.accentColor ?? theme.primary) : (props.mutedColor ?? theme.muted)}
      flexDirection="column"
    >
      <text attributes={TextAttributes.BOLD}>{props.title}</text>
      <select
        ref={(el) => {
          props.ref?.(
            el && {
              moveUp: () => el.moveUp(),
              moveDown: () => el.moveDown(),
              selectCurrent: () => el.selectCurrent(),
            },
          );
        }}
        flexGrow={1}
        focused={props.focused}
        showScrollIndicator
        options={props.items}
        backgroundColor={props.backgroundColor}
        textColor={props.textColor}
        focusedBackgroundColor={props.focusedBackgroundColor}
        focusedTextColor={props.focusedTextColor}
        selectedBackgroundColor={props.selectedBackgroundColor}
        selectedTextColor={props.selectedTextColor}
        descriptionColor={props.descriptionColor}
        selectedDescriptionColor={props.selectedDescriptionColor}
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
