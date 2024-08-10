/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.html", "./frontend/**/*.ts"],
  theme: {
    extend: {
      boxShadow: {
        'trans': '6px 3px 0px rgb(91, 206, 250), 3px 6px 0px rgb(245, 169, 184)',
        'thicc-pink': '4px 4px 0px rgb(245, 169, 184)',
      },
    },
  },
  plugins: [
      require('@tailwindcss/typography'),
      require('@tailwindcss/forms'),
  ],
}
