import tailwindcss from "@tailwindcss/vite";
import {DEFAULT_SITE_FAVICON} from './app/config'

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
            backendBase: ''
        }
    },
    modules: [
        '@pinia/nuxt',
        'pinia-plugin-persistedstate/nuxt',
        '@element-plus/nuxt'
    ],
    piniaPluginPersistedstate: {
        cookieOptions: {
            maxAge: 72 * 60 * 60,
            path: '/',
            sameSite: 'lax'
        }
    },
    css: [
        '~/assets/css/main.css'
    ],
    vite: {
        plugins: [
            tailwindcss()
        ],
    }
})
