const preset = require('../packages/ui/tailwind.preset')

/** @type {import('tailwindcss').Config} */
module.exports = {
  presets: [preset],
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
    '../packages/ui/src/**/*.{vue,js,ts}',
  ],
  plugins: [require('tailwindcss-animate')],
}
