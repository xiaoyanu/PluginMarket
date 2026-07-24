import tailwindcss from "@tailwindcss/vite";
import {API_BASE, ASSET_BASE, DEFAULT_SITE_FAVICON} from './app/config'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},
    app: {
        head: {
            link: [
                {rel: 'icon', type: 'image/x-icon', href: DEFAULT_SITE_FAVICON},
                {
                    rel: 'stylesheet',
                    href: 'https://cdn.jsdmirror.com/gh/xiaoyanu/file-test@26.7/more/AppleColorEmoji/AppleColorEmoji.css'
                    // href: 'https://cdn.jsdelivr.net/gh/xiaoyanu/file-test@26.7/more/AppleColorEmoji/AppleColorEmoji.css'
                }
            ]
        }
    },
    runtimeConfig: {
        public: {
            apiBase: API_BASE,
            assetBase: ASSET_BASE
        }
    },
    modules: [
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        '@element-plus/nuxt'
    ],
    css: [
        '~/assets/css/main.css'
    ],
    vite: {
        plugins: [
            tailwindcss()
        ],
    }
})
