import tailwindcss from "@tailwindcss/vite";
import { API_BASE, ASSET_BASE, DEFAULT_SITE_FAVICON } from './app/config'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},
    app: {
        head: {
            link: [
                { rel: 'icon', type: 'image/x-icon', href: DEFAULT_SITE_FAVICON }
            ]
        }
    },
    runtimeConfig: {
        public: {
            // 浏览器从 Windows 访问 WSL 的 Nuxt 时，127.0.0.1 会指向 Windows 本机，不能指向 WSL 后端。
            // 这里使用 WSL 当前局域网地址，保证前端页面能直接请求 Go 服务。
            apiBase: API_BASE,
            assetBase: ASSET_BASE
        }
    },
    modules: [
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        '@element-plus/nuxt'
    ],
    css:[
        '~/assets/css/main.css'
    ],
    vite: {
        plugins: [
            tailwindcss()
        ],
    }
})
