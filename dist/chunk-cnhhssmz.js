// @bun
// src/presets.ts
function buildPresets(theme) {
  const select = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    selectedBackgroundColor: theme.selection,
    selectedTextColor: theme.primary,
    descriptionColor: theme.muted,
    selectedDescriptionColor: theme.text
  };
  const tabSelect = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    selectedBackgroundColor: theme.selection,
    selectedTextColor: theme.primary,
    selectedDescriptionColor: theme.text
  };
  const input = {
    backgroundColor: theme.background,
    textColor: theme.text,
    focusedBackgroundColor: theme.subtleBg,
    focusedTextColor: theme.text,
    placeholderColor: theme.muted
  };
  const box = {
    border: true,
    borderStyle: theme.borderStyle,
    borderColor: theme.subtle,
    focusedBorderColor: theme.primary,
    backgroundColor: theme.background
  };
  const slider = {
    backgroundColor: theme.subtleBg,
    foregroundColor: theme.primary
  };
  return { select, tabSelect, input, textarea: input, box, slider };
}

export { buildPresets };
