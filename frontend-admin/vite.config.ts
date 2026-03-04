// File: vite.config.ts
// Purpose: Configure Vite build and dev behavior for the admin frontend.
// Module: frontend-admin, tooling/build layer.
// Related: Vue plugin, tsconfig path aliases, npm scripts.
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { fileURLToPath, URL } from 'node:url';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5174,
  },
});
