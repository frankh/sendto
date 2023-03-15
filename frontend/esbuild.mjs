import * as esbuild from 'esbuild'
import stylePlugin from 'esbuild-style-plugin'

import tailwindcss from 'tailwindcss'
import autoprefixer from 'autoprefixer'

let result = await esbuild.build({
  entryPoints: ['src/send.jsx', 'src/receive.jsx', 'src/main.css'],
  bundle: true,
  outdir: 'public/dist',
  minify: true,
  sourcemap: true,
  plugins: [
    stylePlugin({
      postcss: {
        plugins: [tailwindcss, autoprefixer],
      },
    }),
  ],
})

console.log(result)
