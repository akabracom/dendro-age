import { copyFileSync, writeFileSync } from "node:fs";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

function pagesFallback() {
  return {
    name: "pages-fallback",
    apply: "build" as const,
    closeBundle() {
      copyFileSync("dist/index.html", "dist/404.html");
      writeFileSync("dist/.nojekyll", "");
    },
  };
}

export default defineConfig({
  base: "/dendro-age/",
  plugins: [react(), pagesFallback()],
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8088",
        changeOrigin: true,
      },
      "/static": {
        target: "http://localhost:8088",
        changeOrigin: true,
      },
      "/minio": {
        target: "http://localhost:9000",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/minio/, ""),
      },
    },
  },
});
