export interface ColorsYAML {
  base00?: string;
  base01?: string;
  base02?: string;
  base03?: string;
  base04?: string;
  base05?: string;
  base06?: string;
  base07?: string;
  base08?: string;
  base09?: string;
  base0a?: string;
  base0b?: string;
  base0c?: string;
  base0d?: string;
  base0e?: string;
  base0f?: string;
}

export interface ConfigYAML {
  nerd_fonts?: boolean;
  emoji_fonts?: boolean;
  reduced_motion?: boolean;
  mouse?: boolean;
  border?: string;
  colors?: ColorsYAML;
}

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

export type ResolvedColors = Record<(typeof COLOR_KEYS)[number], string>;

export function mergeColors(
  base: ColorsYAML,
  user: ColorsYAML,
): ResolvedColors {
  const merged = {} as ResolvedColors;
  for (const key of COLOR_KEYS) {
    merged[key] = user[key] || base[key] || "";
  }
  return merged;
}

export interface ResolvedConfig {
  nerdFonts: boolean;
  emojiFonts: boolean;
  reducedMotion: boolean;
  mouse: boolean;
  border: string;
  colors: ResolvedColors;
}

export function mergeConfig(
  base: ConfigYAML,
  user: ConfigYAML,
): ResolvedConfig {
  return {
    nerdFonts: user.nerd_fonts ?? base.nerd_fonts ?? false,
    emojiFonts: user.emoji_fonts ?? base.emoji_fonts ?? false,
    reducedMotion: user.reduced_motion ?? base.reduced_motion ?? false,
    mouse: user.mouse ?? base.mouse ?? true,
    border: user.border || base.border || "rounded",
    colors: mergeColors(base.colors ?? {}, user.colors ?? {}),
  };
}

export function normalizeColor(hex: string): string {
  const trimmed = hex.trim();
  if (trimmed && !trimmed.startsWith("#")) return `#${trimmed}`;
  return trimmed;
}

export function normalizeColors(colors: ResolvedColors): ResolvedColors {
  const out = {} as ResolvedColors;
  for (const key of COLOR_KEYS) out[key] = normalizeColor(colors[key]);
  return out;
}

export type OpenTUIBorderStyle = "single" | "double" | "rounded" | "heavy";

const OPENTUI_BORDER_STYLES: readonly OpenTUIBorderStyle[] = [
  "single",
  "double",
  "rounded",
  "heavy",
];

export function resolveBorderStyle(name: string): OpenTUIBorderStyle {
  return (OPENTUI_BORDER_STYLES as readonly string[]).includes(name)
    ? (name as OpenTUIBorderStyle)
    : "rounded";
}
