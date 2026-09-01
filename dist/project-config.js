// @bun
import {
  notify
} from "./chunk-y6txxzxk.js";

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
  notify(text, { kind: "warning", duration: 0 });
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
  mergeKeybinds,
  loadProjectConfig,
  keybindsSchema
};
