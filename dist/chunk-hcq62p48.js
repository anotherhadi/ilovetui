// @bun
// src/config.ts
var COLOR_KEYS = [
  "base00",
  "base01",
  "base02",
  "base03",
  "base04",
  "base05",
  "base06",
  "base07",
  "base08",
  "base09",
  "base0a",
  "base0b",
  "base0c",
  "base0d",
  "base0e",
  "base0f"
];
function mergeColors(base, user) {
  const merged = {};
  for (const key of COLOR_KEYS) {
    merged[key] = user[key] || base[key] || "";
  }
  return merged;
}
function mergeConfig(base, user) {
  return {
    nerdFonts: user.nerd_fonts ?? base.nerd_fonts ?? false,
    emojiFonts: user.emoji_fonts ?? base.emoji_fonts ?? false,
    reducedMotion: user.reduced_motion ?? base.reduced_motion ?? false,
    mouse: user.mouse ?? base.mouse ?? true,
    border: user.border || base.border || "rounded",
    colors: mergeColors(base.colors ?? {}, user.colors ?? {})
  };
}
function normalizeColor(hex) {
  const trimmed = hex.trim();
  if (trimmed && !trimmed.startsWith("#"))
    return `#${trimmed}`;
  return trimmed;
}
function normalizeColors(colors) {
  const out = {};
  for (const key of COLOR_KEYS)
    out[key] = normalizeColor(colors[key]);
  return out;
}
var OPENTUI_BORDER_STYLES = [
  "single",
  "double",
  "rounded",
  "heavy"
];
function resolveBorderStyle(name) {
  return OPENTUI_BORDER_STYLES.includes(name) ? name : "rounded";
}

export { COLOR_KEYS, mergeColors, mergeConfig, normalizeColor, normalizeColors, resolveBorderStyle };
