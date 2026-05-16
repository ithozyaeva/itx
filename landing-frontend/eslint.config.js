import antfu from '@antfu/eslint-config'
import astroPlugin from 'eslint-plugin-astro'

export default antfu(
  {
    typescript: true,
    astro: true,
    formatters: false,
    stylistic: {
      indent: 2,
      quotes: 'single',
    },
    ignores: [
      'dist/**',
      '.astro/**',
      'node_modules/**',
      'public/**',
    ],
  },
  ...astroPlugin.configs.recommended,
)
