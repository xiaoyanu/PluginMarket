import {defineStore} from "pinia";

export interface UserInfo {
    id: number
    username?: string
    nick?: string
    avatar?: string
    email?: string
    userdesc?: string
    power?: number
    titles?: any[]
    created?: string
    updated?: string
}

export const useUserStore = defineStore('PluginMarketNuxt', () => {
    const token = ref<string | null>(null)
    const userInfo = ref<UserInfo | null>(null)
    const isLogin = computed(() => !!token.value)

    const setToken = (newToken: string) => {
        token.value = newToken
    }

    const setUserInfo = (info: UserInfo | null) => {
        userInfo.value = info
    }

    const login = (newToken: string, info: UserInfo | null = null) => {
        token.value = newToken
        userInfo.value = info
    }

    const logout = () => {
        token.value = null
        userInfo.value = null
    }

    return {
        token,
        userInfo,
        isLogin,
        setToken,
        setUserInfo,
        login,
        logout
    }
}, {
    persist: true,
})
