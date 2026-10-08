// Bundles src/ into dist/ (load dist/ in chrome://extensions) and copies public/.
// Shared wording and link parsing are imported from the website's source, so the
// extension and url.vet always say the same thing.
import * as esbuild from "esbuild";
import { cpSync, rmSync, readFileSync, writeFileSync } from "node:fs";

const watch = process.argv.includes("--watch");
const pkg = JSON.parse(readFileSync("package.json", "utf8"));

rmSync("dist", { recursive: true, force: true });
cpSync("public", "dist", { recursive: true });

// One version number, kept in package.json.
const manifest = JSON.parse(readFileSync("dist/manifest.json", "utf8"));
manifest.version = pkg.version;
writeFileSync("dist/manifest.json", JSON.stringify(manifest, null, 2) + "\n");

const common = {
  bundle: true,
  target: "chrome120",
  logLevel: "info",
  legalComments: "none",
};

const builds = [
  // Background worker and extension pages: ES modules.
  {
    ...common,
    entryPoints: {
      background: "src/background.js",
      popup: "src/popup.js",
      blocked: "src/blocked.js",
    },
    outdir: "dist",
    format: "esm",
  },
  // Content scripts can't be modules: one self-contained file each.
  // Their CSS goes in as text, applied as a constructed stylesheet (see card.js).
  {
    ...common,
    entryPoints: {
      card: "src/content/card-entry.js",
      page: "src/content/page-entry.js",
    },
    outdir: "dist",
    format: "iife",
    loader: { ".css": "text" },
  },
];

if (watch) {
  for (const b of builds) await (await esbuild.context(b)).watch();
  console.log(
    "Watching src/ for changes. public/ is copied once; restart after changing it.",
  );
} else {
  await Promise.all(builds.map((b) => esbuild.build(b)));
}
