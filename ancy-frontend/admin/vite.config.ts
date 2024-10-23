/// <reference types="vite/client" />

import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    vue(),
    AutoImport({
      resolvers: [
        ElementPlusResolver({
          importStyle: 'sass',
        }),
        IconsResolver({
          prefix: 'Icon',
        }),
      ],
      eslintrc: { enabled: true },
    }),
    Components({
      resolvers: [
        ElementPlusResolver({
          importStyle: 'sass',
        }),
        IconsResolver({
          enabledCollections: ['ep'],
        }),
      ],
    }),
    Icons({
      autoInstall: true,
    }),
    // AutoImport({
    //   // targets to transform
    //   include: [/\.[tj]sx?$/],
    //   // global imports to register
    //   eslintrc: { enabled: true },
    // }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        javascriptEnabled: true,
        additionalData: `
          @use "@/styles/element.scss" as *;
          @use "@/styles/variable.scss" as *;
        `,
      },
    },
  },
  server: {
    proxy: {
      '/dev': {
        target: 'http://localhost:8889',
        // target: import.meta.env.VITE_APP_API_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/dev/, ''),
      },
    },
  },
})
