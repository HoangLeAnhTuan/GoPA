import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          "vendor-data": ["@tanstack/react-query", "axios", "zustand"],
          "vendor-forms": ["@hookform/resolvers", "react-hook-form", "zod"],
          "vendor-i18n": ["i18next", "react-i18next"],
          "vendor-motion": ["framer-motion"],
          "vendor-react": ["react", "react-dom", "react-router-dom"],
        },
      },
    },
  },
  test: {
    environment: "jsdom",
    setupFiles: "./src/test/setup.ts",
  },
});
