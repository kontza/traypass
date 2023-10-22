/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./src/**/*.js", "./dist/index.html"],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
    daisyui: { themes: ["dark"] }
}

