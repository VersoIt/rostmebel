import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, '..'), '')
  const publicSiteURL = (env.VITE_PUBLIC_SITE_URL || env.PUBLIC_SITE_URL || 'https://rostmebel.shop').replace(/\/+$/, '')

  return {
    plugins: [
      vue(),
      {
        name: 'public-site-url-html',
        transformIndexHtml(html) {
          return html.replace(/__PUBLIC_SITE_URL__/g, publicSiteURL)
        },
      },
    ],
    define: {
      __PUBLIC_SITE_URL__: JSON.stringify(publicSiteURL),
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
  }
})
