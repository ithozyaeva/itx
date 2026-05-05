import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import Icons from 'unplugin-icons/vite'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'
import svgLoader from 'vite-svg-loader'

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  plugins: [
    vue(),
    // vue-devtools раздувает prod-bundle на ~50–100 KB и в проде не нужен
    ...(command === 'serve' ? [vueDevTools()] : []),
    svgLoader({ svgo: false, defaultImport: 'component' }),
    Icons({
      autoInstall: true,
      compiler: 'vue3',
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/regl/'))
            return 'regl'
          if (id.includes('node_modules/@yeger/vue-masonry-wall/'))
            return 'masonry'
          if (id.includes('node_modules/@fontsource-variable/'))
            return 'fonts'
          if (
            id.includes('node_modules/vue/')
            || id.includes('node_modules/@vue/')
            || id.includes('node_modules/vue-router/')
            || id.includes('node_modules/@vueuse/')
          ) {
            return 'vue-vendor'
          }
          if (id.includes('node_modules/@tanstack/'))
            return 'query'
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
