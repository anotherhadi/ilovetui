// @bun
import {
  mergeConfig,
  normalizeColors,
  resolveBorderStyle
} from "./chunk-hcq62p48.js";
import {
  buildPresets
} from "./chunk-cnhhssmz.js";
import {
  configPath,
  readYamlFile
} from "./chunk-7f8jagy5.js";

// src/index.ts
import { join } from "path";
function buildTheme() {
  const defaultConfigPath = join(import.meta.dir, "default.yaml");
  const base = readYamlFile(defaultConfigPath) ?? {};
  const user = readYamlFile(configPath()) ?? {};
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
    borderStyle: resolveBorderStyle(merged.border)
  };
}
var theme = buildTheme();
var presets = buildPresets(theme);

export { theme, presets };
