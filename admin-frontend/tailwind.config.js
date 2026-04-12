const preset = require('../packages/ui/tailwind.preset')

/** @type {import('tailwindcss').Config} */
export default {
  presets: [preset],
  content: [
    './index.html',
    './src/**/*.{ts,js,vue}',
    '../packages/ui/src/**/*.{vue,js,ts}',
  ],
  plugins: [require('tailwindcss-animate')],
}
