import type { OpenTUIBorderStyle, Theme } from "./index.ts";

export interface SelectPreset {
  backgroundColor: string;
  textColor: string;
  focusedBackgroundColor: string;
  focusedTextColor: string;
  selectedBackgroundColor: string;
  selectedTextColor: string;
  descriptionColor: string;
  selectedDescriptionColor: string;
}

export interface TabSelectPreset {
  backgroundColor: string;
  textColor: string;
  focusedBackgroundColor: string;
  focusedTextColor: string;
  selectedBackgroundColor: string;
  selectedTextColor: string;
  selectedDescriptionColor: string;
}

export interface InputPreset {
  backgroundColor: string;
  textColor: string;
  focusedBackgroundColor: string;
  focusedTextColor: string;
  placeholderColor: string;
}

export interface BoxPreset {
  border: true;
  borderStyle: OpenTUIBorderStyle;
  borderColor: string;
  focusedBorderColor: string;
  backgroundColor: string;
}

export interface Presets {
  select: SelectPreset;
  tabSelect: TabSelectPreset;
  input: InputPreset;
  textarea: InputPreset;
  box: BoxPreset;
}

export function buildPresets(theme: Theme): Presets {
  const select: SelectPreset = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    selectedBackgroundColor: theme.selection,
    selectedTextColor: theme.primary,
    descriptionColor: theme.muted,
    selectedDescriptionColor: theme.text,
  };

  const tabSelect: TabSelectPreset = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    selectedBackgroundColor: theme.selection,
    selectedTextColor: theme.primary,
    selectedDescriptionColor: theme.text,
  };

  const input: InputPreset = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    placeholderColor: theme.muted,
  };

  const box: BoxPreset = {
    border: true,
    borderStyle: theme.borderStyle,
    borderColor: theme.subtle,
    focusedBorderColor: theme.primary,
    backgroundColor: theme.background,
  };

  return { select, tabSelect, input, textarea: input, box };
}
