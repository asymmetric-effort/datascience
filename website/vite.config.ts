import { defineConfig } from "vite";

export default defineConfig({
  root: ".",
  base: "/",
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  esbuild: {
    jsxFactory: "createElement",
    jsxFragment: "Fragment",
    jsxImportSource: "@asymmetric-effort/specifyjs",
  },
});
