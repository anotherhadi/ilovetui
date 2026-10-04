// @bun
// src/config.ts
import { z } from "zod";
var COLOR_KEYS2 = [
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
var OPENTUI_BORDER_STYLES2 = ["single", "double", "rounded", "heavy"];
var hexColor = z.string().regex(/^#?[0-9a-fA-F]{6}$/, 'expected a hex color like "#1e1e2e"');
var colorsShape = {};
for (const key of COLOR_KEYS2)
  colorsShape[key] = hexColor.optional();
var ThemeConfigSchema2 = z.object({
  nerd_fonts: z.boolean().optional(),
  emoji_fonts: z.boolean().optional(),
  reduced_motion: z.boolean().optional(),
  mouse: z.boolean().optional(),
  border: z.enum(OPENTUI_BORDER_STYLES2).optional(),
  colors: z.object(colorsShape).optional()
});
function layerConfigs2(layers) {
  const out = {};
  for (const layer of layers) {
    const { colors, ...rest } = layer;
    Object.assign(out, rest);
    if (colors)
      out.colors = { ...out.colors, ...colors };
  }
  return out;
}
function resolveConfig2(config) {
  const colors = {};
  for (const key of COLOR_KEYS2)
    colors[key] = normalizeColor2(config.colors?.[key] ?? "");
  return {
    nerdFonts: config.nerd_fonts ?? false,
    emojiFonts: config.emoji_fonts ?? false,
    reducedMotion: config.reduced_motion ?? false,
    mouse: config.mouse ?? true,
    border: config.border ?? "rounded",
    colors
  };
}
function normalizeColor2(hex) {
  const trimmed = hex.trim();
  if (trimmed && !trimmed.startsWith("#"))
    return `#${trimmed}`;
  return trimmed;
}

export { COLOR_KEYS2, OPENTUI_BORDER_STYLES2, ThemeConfigSchema2, layerConfigs2, resolveConfig2, normalizeColor2 };
