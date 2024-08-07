/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.html", "./frontend/**/*.ts"],
  theme: {
    extend: {},
  },
  plugins: [
      require('@tailwindcss/typography'),
      require('@tailwindcss/forms'),
  ],
  daisyui: {
    themes: [
       {
         "pique": {
			"primary": "#F5A9B8",
			"primary-content": "#150a0c",
			"secondary": "#5BCEFA",
			"secondary-content": "#030f15",
			"accent": "#A2E4B8",
			"accent-content": "#0a120c",
			"neutral": "#ff00ff",
			"neutral-content": "#160016",
			"base-100": "#ffffff",
			"base-200": "#dedede",
			"base-300": "#bebebe",
			"base-content": "#161616",
			"info": "#4f46e5",
			"info-content": "#d6ddfe",
			"success": "#0f766e",
			"success-content": "#d3e3e0",
			"warning": "#eab308",
			"warning-content": "#130c00",
			"error": "#ef4444",
			"error-content": "#140202",
         }
      },
    ],
  },
}

