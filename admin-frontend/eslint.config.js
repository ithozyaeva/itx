import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  typescript: true,
  ignores: ['.github/**', 'dist/**', 'node_modules/**', '**/__tests__/**'],
})
