/**
 * File: vite.config.ts
 * Purpose: Configure Vite build and development behavior for React admin app.
 * Module: frontend-admin-react/build, tooling configuration layer.
 * Related: TypeScript compiler settings and npm scripts.
 */
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
})
