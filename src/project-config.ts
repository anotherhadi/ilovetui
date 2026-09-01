import { existsSync, readFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";
import { parse as parseYAML } from "yaml";
import { z } from "zod";
import { notify } from "./context/notifications.ts";

export interface ProjectConfigOptions<TDefault, TUser> {
  projectName: string;
  defaultConfigPath: string;
  defaultSchema: z.ZodType<TDefault>;
  userSchema: z.ZodType<TUser>;
}

export interface ProjectConfig<TDefault, TUser> {
  defaults: TDefault;
  user: TUser;
}

function userConfigPath(projectName: string): string {
  const configHome = process.env.XDG_CONFIG_HOME || join(homedir(), ".config");
  return join(configHome, projectName, "config.yaml");
}

function readYaml(path: string): unknown {
  return parseYAML(readFileSync(path, "utf-8"));
}

function warn(projectName: string, path: string, message: string): void {
  const text = `${path}: ${message}`;
  console.error(`[${projectName} config] ${text}`);
  notify(text, { kind: "warning", duration: 0 });
}

export function loadProjectConfig<TDefault, TUser>(
  options: ProjectConfigOptions<TDefault, TUser>,
): ProjectConfig<TDefault, TUser> {
  const { projectName, defaultConfigPath, defaultSchema, userSchema } = options;

  let rawDefaults: unknown;
  try {
    rawDefaults = readYaml(defaultConfigPath);
  } catch (error) {
    throw new Error(`Failed to read ${defaultConfigPath}: ${error instanceof Error ? error.message : error}`);
  }
  const defaultsResult = defaultSchema.safeParse(rawDefaults);
  if (!defaultsResult.success) {
    throw new Error(`${defaultConfigPath} is invalid:\n${z.prettifyError(defaultsResult.error)}`);
  }

  const path = userConfigPath(projectName);
  const user = ((): TUser => {
    if (!existsSync(path)) return {} as TUser;

    let raw: unknown;
    try {
      raw = readYaml(path);
    } catch (error) {
      warn(projectName, path, error instanceof Error ? error.message : String(error));
      return {} as TUser;
    }

    const result = userSchema.safeParse(raw ?? {});
    if (!result.success) {
      warn(projectName, path, z.prettifyError(result.error));
      return {} as TUser;
    }
    return result.data;
  })();

  return { defaults: defaultsResult.data, user };
}

export function keybindsSchema<TName extends string>(
  names: readonly TName[],
  required: true,
): z.ZodObject<Record<TName, z.ZodString>>;
export function keybindsSchema<TName extends string>(
  names: readonly TName[],
  required: false,
): z.ZodObject<Record<TName, z.ZodOptional<z.ZodString>>>;
export function keybindsSchema<TName extends string>(names: readonly TName[], required: boolean) {
  const shape = {} as Record<TName, z.ZodTypeAny>;
  for (const name of names) {
    shape[name] = required ? z.string().min(1) : z.string().min(1).optional();
  }
  return z.object(shape).strict();
}

export function mergeKeybinds<TName extends string>(
  names: readonly TName[],
  defaults: Record<TName, string>,
  user: Partial<Record<TName, string>> | undefined,
): Record<TName, string> {
  const merged = {} as Record<TName, string>;
  for (const name of names) {
    merged[name] = user?.[name] ?? defaults[name];
  }
  return merged;
}
