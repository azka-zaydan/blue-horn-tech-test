import viteReact from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

import tanstackRouter from '@tanstack/router-plugin/vite'
import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter({
      autoCodeSplitting: true,
    }),
    viteReact(),
  ],

  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
})
