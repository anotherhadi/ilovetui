import { existsSync, readFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";
import { parse as parseYAML } from "yaml";

export function configPath(): string {
  const configDir = process.env.XDG_CONFIG_HOME || join(homedir(), ".config");
  return join(configDir, "ilovetui", "config.yaml");
}

export function readYamlFile<T>(path: string): T | undefined {
  if (!existsSync(path)) return undefined;
  try {
    return (parseYAML(readFileSync(path, "utf8")) as T) ?? undefined;
  } catch (err) {
    console.error(`ilovetui: failed to parse ${path}, ignoring it.`, err);
    return undefined;
  }
}
