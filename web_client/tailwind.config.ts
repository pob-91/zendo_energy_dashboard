import type { Config } from "tailwindcss";
//const defaultTheme = require('tailwindcss/defaultTheme');

const config: Config = {
  content: [
    "./src/routes/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    // fontFamily: {
    //   montserrat: ["Montserrat", "sans-serif"],
    //   ptmono: ["PT Mono", "mono"],
    //   poppins: ["Poppins", "sans-serif"],
    // },
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {},
      screens: {
        smh: {
          raw: "(max-height: 667px)",
        },
        mdh: {
          raw: "(max-height: 1000px)",
        },
        lgh: {
          raw: "(max-height: 1200px)",
        },
        xlh: {
          raw: "(max-height: 1500px)",
        },
        "3xl": "1920px",
        "4xl": "2560px",
      },
      transitionProperty: {
        width: "width",
      },
    },
  },
  plugins: [],
};
module.exports = config;
