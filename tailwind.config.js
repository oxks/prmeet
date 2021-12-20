module.exports = {
  mode: 'jit',
  purge: [
    './views/layout/**/*.go.html',
    './views/layout/general/*.go.html',
    './views/layout/partials/*.go.html',
    './views/public/css/**/*.{js,jsx,ts,tsx,vue}',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
