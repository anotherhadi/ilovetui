import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { parse as parseYAML } from "yaml";
import { z } from "zod";
import {
  type ConfigYAML,
  type OpenTUIBorderStyle,
  layerConfigs,
  resolveConfig,
  ThemeConfigSchema,
} from "./config.ts";
import { notify } from "./context/notifications.ts";
import { buildPresets, type Presets } from "./presets.ts";
import { configPath } from "./yaml.ts";

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

function reportThemeError(path: string, message: string): void {
  const text = `${path}: ${message}`;
  console.error(`[ilovetui theme] ${text}`);
  notify(text, { kind: "warning", duration: 0 });
}

// A broken file is reported (as a sticky warning toast, since console output
// is hidden behind the TUI) and skipped as a whole, falling back to the
// layers below it.
function readThemeFile(path: string): ConfigYAML {
  if (!existsSync(path)) return {};
  let raw: unknown;
  try {
    raw = parseYAML(readFileSync(path, "utf8"));
  } catch (error) {
    reportThemeError(path, error instanceof Error ? error.message : String(error));
    return {};
  }
  const result = ThemeConfigSchema.safeParse(raw ?? {});
  if (!result.success) {
    reportThemeError(path, z.prettifyError(result.error));
    return {};
  }
  return result.data;
}

// Read once: the shipped defaults, then the user's shared
// ~/.config/ilovetui/config.yaml on top.
const fileLayers: ConfigYAML[] = [readThemeFile(join(import.meta.dir, "default.yaml")), readThemeFile(configPath())];

function buildTheme(overrides: ConfigYAML): Theme {
  const resolved = resolveConfig(layerConfigs([...fileLayers, overrides]));
  const colors = resolved.colors;

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

    nerdFonts: resolved.nerdFonts,
    emojiFonts: resolved.emojiFonts,
    reducedMotion: resolved.reducedMotion,
    mouse: resolved.mouse,
    borderStyle: resolved.border,
  };
}

export const theme: Theme = buildTheme({});

export const presets: Presets = buildPresets(theme);

/**
 * Layers app-specific overrides (typically a `theme:` section of the app's own
 * config, validated with `ThemeConfigSchema`) on top of the shared ilovetui
 * config. `theme` and `presets` are updated in place, so call this once at
 * startup, before the first render.
 */
export function configureTheme(overrides: ConfigYAML = {}): void {
  Object.assign(theme, buildTheme(overrides));
  const next = buildPresets(theme);
  for (const key of Object.keys(next) as (keyof Presets)[]) Object.assign(presets[key], next[key]);
}

/** A Nerd Font glyph followed by a space, or "" when `nerd_fonts` is off. */
export function icon(codepoint: string): string {
  return theme.nerdFonts ? `${codepoint} ` : "";
}

export { ThemeConfigSchema, type ConfigYAML as ThemeConfig, type OpenTUIBorderStyle } from "./config.ts";
export type { Presets, SelectPreset, TabSelectPreset, InputPreset, BoxPreset } from "./presets.ts";
