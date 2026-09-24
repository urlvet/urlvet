/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './src/**/*.{html,js,svelte,ts}',
    './src/**/*.{svelte,html,js,ts}',
    './src/routes/**/*.{svelte,html,js,ts}',
    './src/lib/**/*.{svelte,html,js,ts}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      // Warm neutrals: cream paper in light mode, warm near-black in dark mode.
      colors: {
        gray: {
          50: '#faf7f2',
          100: '#f3ede2',
          200: '#e6dfd2',
          300: '#d3cab9',
          400: '#a69c8c',
          500: '#7c7366',
          600: '#5c5449',
          700: '#3f3a33',
          800: '#2a2622',
          900: '#1a1714',
          950: '#0f0d0b',
        },
        // Brand accent (logo dot, tagline underline, bullets). Muted blue that sits well on
        // the warm palette; green stays reserved for the "Safe" meaning.
        accent: {
          light: '#3461c9',
          dark: '#8eb0ff',
        },
      },
      fontFamily: {
        sans: ['Geist', 'system-ui', '-apple-system', 'Segoe UI', 'Roboto', 'sans-serif'],
        serif: ['"Instrument Serif"', 'Georgia', 'serif'],
        mono: ['"Geist Mono"', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'monospace'],
      },
    },
  },
  plugins: [],
};
