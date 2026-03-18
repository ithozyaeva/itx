import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  typescript: true,
  ignores: ['.github/**', 'dist/**', 'node_modules/**', '**/__tests__/**', '.omc/**', '.claude/**'],
}, {
  files: ['**/*.vue'],
  rules: {
    'vue/max-attributes-per-line': ['error', {
      singleline: 1,
      multiline: 1,
    }],
  },
})
