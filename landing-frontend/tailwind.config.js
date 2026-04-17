/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ['class'],
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2.5rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      fontFamily: {
        sans: ['"Inter Variable"', 'Inter', 'sans-serif'],
        mono: ['"JetBrains Mono Variable"', 'JetBrains Mono', 'ui-monospace', 'Menlo', 'monospace'],
        display: ['Unbounded', '"Space Grotesk Variable"', 'Space Grotesk', 'sans-serif'],
      },
      colors: {
        'term-amber': '#ffb547',
        'term-magenta': '#ff4d8b',
        'term-cyan': '#5eead4',
        'border': 'hsl(var(--border))',
        'input': 'hsl(var(--input))',
        'ring': 'hsl(var(--ring))',
        'background': 'hsl(var(--background))',
        'foreground': 'hsl(var(--foreground))',
        'primary': {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        'secondary': {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        'destructive': {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        'muted': {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        'accent': {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        'popover': {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        'card': {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      keyframes: {
        'marquee': {
          from: { transform: 'translateX(0)' },
          to: { transform: 'translateX(-50%)' },
        },
        'type-caret': {
          '50%': { opacity: 0 },
        },
        'glitch-x': {
          '0%, 100%': { transform: 'translateX(0)' },
          '20%': { transform: 'translateX(-1px)' },
          '40%': { transform: 'translateX(1px)' },
          '60%': { transform: 'translateX(-0.5px)' },
        },
        'scan': {
          '0%': { backgroundPosition: '0 0' },
          '100%': { backgroundPosition: '0 200%' },
        },
        'accordion-down': {
          from: { height: 0 },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: 0 },
        },
        'card-reveal': {
          from: {
            opacity: '0',
            filter: 'blur(4px)',
            transform: 'translateY(60px) scale(0.95)',
          },
          to: {
            opacity: '1',
            filter: 'blur(0)',
            transform: 'translateY(0) scale(1)',
          },
        },
        'card-reveal-highlight': {
          '0%': {
            opacity: '0',
            filter: 'blur(4px)',
            transform: 'translateY(60px) scale(0.95)',
            boxShadow: '0 0 0 0 hsl(var(--accent) / 0)',
          },
          '60%': {
            opacity: '1',
            filter: 'blur(0)',
            transform: 'translateY(0) scale(1)',
            boxShadow: '0 0 30px 5px hsl(var(--accent) / 0.3)',
          },
          '100%': {
            opacity: '1',
            filter: 'blur(0)',
            transform: 'translateY(0) scale(1)',
            boxShadow: '0 0 0 0 hsl(var(--accent) / 0)',
          },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'card-reveal': 'card-reveal 0.7s cubic-bezier(0.16, 1, 0.3, 1) both',
        'card-reveal-highlight': 'card-reveal-highlight 1s cubic-bezier(0.16, 1, 0.3, 1) both',
        'marquee': 'marquee 38s linear infinite',
        'marquee-slow': 'marquee 60s linear infinite',
        'type-caret': 'type-caret 1.1s steps(1) infinite',
        'glitch-x': 'glitch-x 3s steps(1) infinite',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
}
