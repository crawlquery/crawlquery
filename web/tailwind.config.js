/** @type {import('tailwindcss').Config} */
export default {
  purge: {
    content: [
      './src/**/*.html',
      './src/**/*.vue',
      './src/**/*.jsx',
    ],
  },
  darkMode: false, // or 'media' or 'class'
  content: [],
  theme: {
    extend: {},
  },
  daisy: {
    themes: ["light"],
  },
  plugins: [
    require('daisyui'),
  ],
}

