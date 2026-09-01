// @bun
// src/yaml.ts
import { existsSync, readFileSync } from "fs";
import { homedir } from "os";
import { join } from "path";
import { parse as parseYAML } from "yaml";
function configPath() {
  const configDir = process.env.XDG_CONFIG_HOME || join(homedir(), ".config");
  return join(configDir, "ilovetui", "config.yaml");
}
function readYamlFile(path) {
  if (!existsSync(path))
    return;
  try {
    return parseYAML(readFileSync(path, "utf8")) ?? undefined;
  } catch (err) {
    console.error(`ilovetui: failed to parse ${path}, ignoring it.`, err);
    return;
  }
}

export { configPath, readYamlFile };
