import { join } from "node:path";
import {
  type ConfigYAML,
  type OpenTUIBorderStyle,
  mergeConfig,
  normalizeColors,
  resolveBorderStyle,
} from "./config.ts";
import { buildPresets } from "./presets.ts";
import { configPath, readYamlFile } from "./yaml.ts";

export interface Theme {
  base00: string;
  base01: string;
  base02: string;
  base03: string;
  base04: string;
  base05: string;
  base06: string;
  base07: string;
  base08: string;
  base09: string;
  base0a: string;
  base0b: string;
  base0c: string;
  base0d: string;
  base0e: string;
  base0f: string;

  background: string;
  subtleBg: string;
  selection: string;
  subtle: string;
  muted: string;
  text: string;
  primary: string;
  success: string;
  warning: string;
  error: string;

  nerdFonts: boolean;
  emojiFonts: boolean;
  reducedMotion: boolean;
  mouse: boolean;
  borderStyle: OpenTUIBorderStyle;
}

export { configPath } from "./yaml.ts";

function buildTheme(): Theme {
  const defaultConfigPath = join(import.meta.dir, "default.yaml");
  const base = readYamlFile<ConfigYAML>(defaultConfigPath) ?? {};
  const user = readYamlFile<ConfigYAML>(configPath()) ?? {};
  const merged = mergeConfig(base, user);
  const colors = normalizeColors(merged.colors);

  return {
    ...colors,

    background: colors.base00,
    subtleBg: colors.base01,
    selection: colors.base02,
    subtle: colors.base03,
    muted: colors.base04,
    text: colors.base05,
    primary: colors.base0d,
    success: colors.base0b,
    warning: colors.base09,
    error: colors.base08,

    nerdFonts: merged.nerdFonts,
    emojiFonts: merged.emojiFonts,
    reducedMotion: merged.reducedMotion,
    mouse: merged.mouse,
    borderStyle: resolveBorderStyle(merged.border),
  };
}

export const theme: Theme = buildTheme();

export const presets = buildPresets(theme);

export type { OpenTUIBorderStyle } from "./config.ts";
export type { Presets, SelectPreset, TabSelectPreset, InputPreset, BoxPreset } from "./presets.ts";
