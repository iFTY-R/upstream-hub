import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, 'VITE_')
  const BACKEND_TARGET = env.VITE_BACKEND_URL ?? 'http://localhost:8418'

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, '.'),
      },
    },
    server: {
      port: 3010,
      strictPort: true,
      proxy: {
        '/api': { target: BACKEND_TARGET, changeOrigin: true },
        '/healthz': { target: BACKEND_TARGET, changeOrigin: true },
      },
    },
  }
})
