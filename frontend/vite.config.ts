import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { fileURLToPath } from 'url'
import tailwindcss from '@tailwindcss/vite'
import ui from '@nuxt/ui/vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [
      vue(),
      vueDevTools(),
      tailwindcss(),
      ui(),
    ],
    resolve: {
      tsconfigPaths: true,
      alias: {
        '@types': fileURLToPath(new URL('./src/types.ts', import.meta.url)),
        '@api': fileURLToPath(new URL('./src/utils/api.ts', import.meta.url)),
        '@utils': fileURLToPath(new URL('./src/utils', import.meta.url)),
        '@router': fileURLToPath(new URL('./src/router/index.ts', import.meta.url)),
        '@components': fileURLToPath(new URL('./src/components', import.meta.url)),
        '@layouts': fileURLToPath(new URL('./src/layouts', import.meta.url)),
        '@ui': fileURLToPath(new URL('./src/components/ui', import.meta.url)),
        '@views': fileURLToPath(new URL('./src/views', import.meta.url)),
        '@assets': fileURLToPath(new URL('./src/assets', import.meta.url)),
        '@stores': fileURLToPath(new URL('./src/stores', import.meta.url)),
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      }
    },
    server: {
      host: "0.0.0.0",
      port: parseInt(env.FRONTEND_PORT || '8080'),
      proxy: {
        '/api': {
          target: 'http://localhost:8081',
          changeOrigin: true,
          secure: false,
        },
      }
    }
  }
})
