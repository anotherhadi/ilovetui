// @bun
import {
  ThemeConfigSchema2,
  layerConfigs2,
  resolveConfig2
} from "./chunk-njvj322f.js";
import {
  notify2
} from "./chunk-6ggk5fg7.js";
import {
  buildPresets2
} from "./chunk-yyg967q8.js";
import {
  configPath2
} from "./chunk-xykm10y9.js";

// src/index.ts
import { existsSync, readFileSync } from "fs";
import { join } from "path";
import { parse as parseYAML } from "yaml";
import { z } from "zod";
function reportThemeError(path, message) {
  const text = `${path}: ${message}`;
  console.error(`[ilovetui theme] ${text}`);
  notify2(text, { kind: "warning", duration: 0 });
}
function readThemeFile(path) {
  if (!existsSync(path))
    return {};
  let raw;
  try {
    raw = parseYAML(readFileSync(path, "utf8"));
  } catch (error) {
    reportThemeError(path, error instanceof Error ? error.message : String(error));
    return {};
  }
  const result = ThemeConfigSchema2.safeParse(raw ?? {});
  if (!result.success) {
    reportThemeError(path, z.prettifyError(result.error));
    return {};
  }
  return result.data;
}
var fileLayers = [readThemeFile(join(import.meta.dir, "default.yaml")), readThemeFile(configPath2())];
function buildTheme(overrides) {
  const resolved = resolveConfig2(layerConfigs2([...fileLayers, overrides]));
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
    borderStyle: resolved.border
  };
}
var theme2 = buildTheme({});
var presets2 = buildPresets2(theme2);
function configureTheme2(overrides = {}) {
  Object.assign(theme2, buildTheme(overrides));
  const next = buildPresets2(theme2);
  for (const key of Object.keys(next))
    Object.assign(presets2[key], next[key]);
}
function icon2(codepoint) {
  return theme2.nerdFonts ? `${codepoint} ` : "";
}

export { theme2, presets2, configureTheme2, icon2 };
