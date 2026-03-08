/**
 * File: vite.config.ts
 * Purpose: Configure Vite build and development behavior for React admin app.
 * Module: frontend-admin-react/build, tooling configuration layer.
 * Related: TypeScript compiler settings and npm scripts.
 */
import path from 'path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const devApiTarget = process.env.VITE_DEV_API_TARGET || 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: devApiTarget,
        changeOrigin: true,
      },
    },
  },
})
