// @bun
// src/yaml.ts
import { homedir } from "os";
import { join } from "path";
function configPath2() {
  const configDir = process.env.XDG_CONFIG_HOME || join(homedir(), ".config");
  return join(configDir, "ilovetui", "config.yaml");
}

export { configPath2 };
