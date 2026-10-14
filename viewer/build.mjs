// Сборка вьюера.
//   node build.mjs        → ../internal/htmlout/assets/viewer.js (минифицированный IIFE, его встраивает Go)
//   node build.mjs --dev  → dev-сервер с синтетической сценой на http://localhost:5173
import * as esbuild from "esbuild";
import { statSync } from "node:fs";

const common = {
  bundle: true,
  format: "iife",
  target: ["chrome100", "edge100", "firefox113", "safari16.4"],
  loader: { ".css": "text" },
  legalComments: "eof", // лицензия three.js (MIT) должна остаться в бандле
  logLevel: "info",
};

if (process.argv.includes("--dev")) {
  const ctx = await esbuild.context({
    ...common,
    entryPoints: ["dev/dev.ts"],
    outfile: "dev/out/dev.js",
    sourcemap: true,
  });
  await ctx.watch();
  const { port } = await ctx.serve({ servedir: "dev", port: 5173 });
  console.log(`dev: http://localhost:${port}`);
} else {
  const outfile = "../internal/htmlout/assets/viewer.js";
  await esbuild.build({ ...common, entryPoints: ["src/main.ts"], outfile, minify: true });
  console.log(`viewer.js: ${(statSync(outfile).size / 1024).toFixed(0)} КБ`);
}
