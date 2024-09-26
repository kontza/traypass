/** @type {import('tailwindcss').Config} */
import daisyui from 'daisyui'
import catppuccin from '@catppuccin/daisyui'
module.exports = {
  content: ['./src/**/*.{html,js}', 'index.html'],
  theme: {
    extend: {},
  },
  plugins: [daisyui],
  daisyui: { themes: [catppuccin('frappe')] },
}
