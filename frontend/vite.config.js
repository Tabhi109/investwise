import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/register': 'http://localhost:8080',
      '/login': 'http://localhost:8080',
      '/portfolio': 'http://localhost:8080',
      '/trade': 'http://localhost:8080',
      '/tickers': 'http://localhost:8080',
      '/market': 'http://localhost:8080',
      '/ping': 'http://localhost:8080',
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true
      }
    }
  }
})
