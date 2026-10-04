import { z } from "zod";

export const COLOR_KEYS = [
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
  "base0f",
] as const;

export type ColorKey = (typeof COLOR_KEYS)[number];

export const OPENTUI_BORDER_STYLES = ["single", "double", "rounded", "heavy"] as const;

export type OpenTUIBorderStyle = (typeof OPENTUI_BORDER_STYLES)[number];

const hexColor = z.string().regex(/^#?[0-9a-fA-F]{6}$/, 'expected a hex color like "#1e1e2e"');

const colorsShape = {} as Record<ColorKey, z.ZodOptional<z.ZodString>>;
for (const key of COLOR_KEYS) colorsShape[key] = hexColor.optional();

// Not .strict(): ~/.config/ilovetui/config.yaml is shared by every
// ilovetui-based TUI, and some of them read extra keys of their own from it
// (e.g. `layout_border`). Unknown keys are stripped, known ones validated.
export const ThemeConfigSchema = z.object({
  nerd_fonts: z.boolean().optional(),
  emoji_fonts: z.boolean().optional(),
  reduced_motion: z.boolean().optional(),
  mouse: z.boolean().optional(),
  border: z.enum(OPENTUI_BORDER_STYLES).optional(),
  colors: z.object(colorsShape).optional(),
});

export type ConfigYAML = z.infer<typeof ThemeConfigSchema>;
export type ColorsYAML = NonNullable<ConfigYAML["colors"]>;

export type ResolvedColors = Record<ColorKey, string>;

export interface ResolvedConfig {
  nerdFonts: boolean;
  emojiFonts: boolean;
  reducedMotion: boolean;
  mouse: boolean;
  border: OpenTUIBorderStyle;
  colors: ResolvedColors;
}

// Later layers win, key by key (colors included).
export function layerConfigs(layers: ConfigYAML[]): ConfigYAML {
  const out: ConfigYAML = {};
  for (const layer of layers) {
    const { colors, ...rest } = layer;
    Object.assign(out, rest);
    if (colors) out.colors = { ...out.colors, ...colors };
  }
  return out;
}

export function resolveConfig(config: ConfigYAML): ResolvedConfig {
  const colors = {} as ResolvedColors;
  for (const key of COLOR_KEYS) colors[key] = normalizeColor(config.colors?.[key] ?? "");
  return {
    nerdFonts: config.nerd_fonts ?? false,
    emojiFonts: config.emoji_fonts ?? false,
    reducedMotion: config.reduced_motion ?? false,
    mouse: config.mouse ?? true,
    border: config.border ?? "rounded",
    colors,
  };
}

export function normalizeColor(hex: string): string {
  const trimmed = hex.trim();
  if (trimmed && !trimmed.startsWith("#")) return `#${trimmed}`;
  return trimmed;
}
