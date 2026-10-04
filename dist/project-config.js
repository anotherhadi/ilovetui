// @bun
import {
  notify2
} from "./chunk-6ggk5fg7.js";

// src/project-config.ts
import { existsSync, readFileSync } from "fs";
import { homedir } from "os";
import { join } from "path";
import { parse as parseYAML } from "yaml";
import { z } from "zod";
function userConfigPath(projectName) {
  const configHome = process.env.XDG_CONFIG_HOME || join(homedir(), ".config");
  return join(configHome, projectName, "config.yaml");
}
function readYaml(path) {
  return parseYAML(readFileSync(path, "utf-8"));
}
function warn(projectName, path, message) {
  const text = `${path}: ${message}`;
  console.error(`[${projectName} config] ${text}`);
  notify2(text, { kind: "warning", duration: 0 });
}
function loadProjectConfig(options) {
  const { projectName, defaultConfigPath, defaultSchema, userSchema } = options;
  let rawDefaults;
  try {
    rawDefaults = readYaml(defaultConfigPath);
  } catch (error) {
    throw new Error(`Failed to read ${defaultConfigPath}: ${error instanceof Error ? error.message : error}`);
  }
  const defaultsResult = defaultSchema.safeParse(rawDefaults);
  if (!defaultsResult.success) {
    throw new Error(`${defaultConfigPath} is invalid:
${z.prettifyError(defaultsResult.error)}`);
  }
  const path = userConfigPath(projectName);
  const user = (() => {
    if (!existsSync(path))
      return {};
    let raw;
    try {
      raw = readYaml(path);
    } catch (error) {
      warn(projectName, path, error instanceof Error ? error.message : String(error));
      return {};
    }
    const result = userSchema.safeParse(raw ?? {});
    if (!result.success) {
      warn(projectName, path, z.prettifyError(result.error));
      return {};
    }
    return result.data;
  })();
  return { defaults: defaultsResult.data, user };
}
function isPlainObject(value) {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}
function deepMerge(base, over) {
  if (!isPlainObject(base) || !isPlainObject(over))
    return over === undefined ? base : over;
  const out = { ...base };
  for (const [key, value] of Object.entries(over))
    out[key] = deepMerge(base[key], value);
  return out;
}
function loadMergedProjectConfig(options) {
  const { projectName, defaultConfigPath, schema } = options;
  let rawDefaults;
  try {
    rawDefaults = readYaml(defaultConfigPath);
  } catch (error) {
    throw new Error(`Failed to read ${defaultConfigPath}: ${error instanceof Error ? error.message : error}`);
  }
  const defaults = schema.safeParse(rawDefaults);
  if (!defaults.success) {
    throw new Error(`${defaultConfigPath} is invalid:
${z.prettifyError(defaults.error)}`);
  }
  const path = userConfigPath(projectName);
  if (!existsSync(path))
    return defaults.data;
  let rawUser;
  try {
    rawUser = readYaml(path);
  } catch (error) {
    warn(projectName, path, error instanceof Error ? error.message : String(error));
    return defaults.data;
  }
  const merged = schema.safeParse(deepMerge(rawDefaults, rawUser ?? {}));
  if (!merged.success) {
    warn(projectName, path, z.prettifyError(merged.error));
    return defaults.data;
  }
  return merged.data;
}
function keybindsSchema(names, required) {
  const shape = {};
  for (const name of names) {
    shape[name] = required ? z.string().min(1) : z.string().min(1).optional();
  }
  return z.object(shape).strict();
}
function mergeKeybinds(names, defaults, user) {
  const merged = {};
  for (const name of names) {
    merged[name] = user?.[name] ?? defaults[name];
  }
  return merged;
}
export {
  keybindsSchema,
  loadMergedProjectConfig,
  loadProjectConfig,
  mergeKeybinds
};
