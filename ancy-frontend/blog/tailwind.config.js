/** @type {import('tailwindcss').Config} */
export default {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  content: [],
  theme: {
    extend: {
      height: {
        'screen-nh': 'calc(100vh - var(--ac-header-height))',
      },
    },
  },
  plugins: [],
}
