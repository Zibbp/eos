/** @type {import('tailwindcss').Config} */
module.exports = {
 	content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  theme: {
    extend: {
      fontFamily: {
        'inter': ['Inter', 'sans-serif'],
        'roboto': ['Roboto', 'sans-serif']
      },
    },
  }
}

