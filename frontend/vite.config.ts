import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, '..'), '')
  const publicSiteURL = (env.VITE_PUBLIC_SITE_URL || env.PUBLIC_SITE_URL || 'https://rostmebel.shop').replace(/\/+$/, '')
  const yandexMetrikaId = (env.VITE_YANDEX_METRIKA_ID || env.YANDEX_METRIKA_ID || '111598091').trim()

  return {
    plugins: [
      vue(),
      {
        name: 'public-site-url-html',
        transformIndexHtml(html) {
          return html
            .replace(/__PUBLIC_SITE_URL__/g, publicSiteURL)
            .replace(/__YANDEX_METRIKA_ID__/g, yandexMetrikaId)
        },
      },
    ],
    define: {
      __PUBLIC_SITE_URL__: JSON.stringify(publicSiteURL),
      __YANDEX_METRIKA_ID__: JSON.stringify(yandexMetrikaId),
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
