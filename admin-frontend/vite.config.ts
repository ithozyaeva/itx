import path from 'node:path'

import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'
import Icons from 'unplugin-icons/vite'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  base: '/admin/',
  plugins: [
    vue(),
    // vue-devtools раздувает prod-bundle на ~50–100 KB и в проде не нужен
    ...(command === 'serve' ? [vueDevTools()] : []),
    Icons({
      autoInstall: true,
      compiler: 'vue3',
    }),
  ],
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@itx/ui': path.resolve(__dirname, '../packages/ui/src'),
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
