/** @type {import('tailwindcss').Config} */
export default {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  content: [],
  theme: {
    extend: {
      height: {
        'screen-nh': 'calc(100vh - var(--ac-header-height))',
      },
      colors: {
        gray: 'var(--gray)',
        'gray-bg': 'var(--gray-bg)',
        primary: 'var(--primary-color)',
        'primary-bg': 'var(--primary-background)',
      },
    },
  },
  plugins: [],
}
