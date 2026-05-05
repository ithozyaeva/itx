import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  base: '/platform/',
  plugins: [
    vue(),
    // vue-devtools раздувает prod-bundle на ~50–100 KB и в проде не нужен
    ...(command === 'serve' ? [vueDevTools()] : []),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@itx/ui': fileURLToPath(new URL('../packages/ui/src', import.meta.url)),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (
            id.includes('node_modules/vue/')
            || id.includes('node_modules/@vue/')
            || id.includes('node_modules/vue-router/')
            || id.includes('node_modules/pinia/')
          ) {
            return 'vue-vendor'
          }
          if (id.includes('node_modules/@tanstack/'))
            return 'query'
          if (id.includes('node_modules/@headlessui/'))
            return 'headless'
          return undefined
        },
      },
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:80',
      },
    },
  },
}))
