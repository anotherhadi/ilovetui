import { createSolidTransformPlugin } from "@opentui/solid/bun-plugin";

const entrypoints = await Array.fromAsync(new Bun.Glob("**/*.{ts,tsx}").scan({ cwd: "src" }));

await Bun.$`rm -rf dist`;

const result = await Bun.build({
  entrypoints: entrypoints.map((path) => `src/${path}`),
  outdir: "dist",
  root: "src",
  target: "bun",
  format: "esm",
  splitting: true,
  plugins: [createSolidTransformPlugin()],
  external: ["@opentui/core", "@opentui/solid", "@opentui/keymap", "solid-js", "opentui-spinner", "yaml", "zod"],
});

if (!result.success) {
  for (const message of result.logs) console.error(message);
  process.exit(1);
}

await Bun.write("dist/default.yaml", Bun.file("src/default.yaml"));

console.log(`Built ${result.outputs.length} files to dist/`);
