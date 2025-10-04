import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    port: 5001,
    proxy: {
      '/api': {
        target: 'http://localhost:5082',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: '../gorev-mcpserver/cmd/gorev/web/dist',
    emptyOutDir: true,
    assetsDir: 'assets',
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          api: ['axios', '@tanstack/react-query']
        }
      }
    }
  }
})