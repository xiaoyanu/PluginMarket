import { API_BASE, ASSET_BASE } from '~/config'
import { appendForwardedClientIpHeaders } from '~/utils/forwarded-headers'

/** 统一响应数据结构 */
interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

type FetchOptions = NonNullable<Parameters<typeof $fetch>[1]>
type ApiFetchOptions = FetchOptions & {
  /** 使用页面内嵌错误状态时可关闭全局消息，默认始终显示。 */
  suppressErrorMessage?: boolean
  /** 后端没有返回可读消息时使用的业务兜底文案。 */
  errorFallback?: string
}

const API_ERROR_NOTIFIED = Symbol('api-error-notified')
type NotifiedError = Error & { [API_ERROR_NOTIFIED]?: boolean }

class ApiRequestError extends Error {
  constructor(message: string, public readonly code: number) {
    super(message)
    this.name = 'ApiRequestError'
  }
}

const runHooks = async (hooks: any, context: any) => {
  if (!hooks) return
  for (const hook of Array.isArray(hooks) ? hooks : [hooks]) await hook(context)
}

const getBackendBases = () => {
  const config = useRuntimeConfig()
  const backendBase = String(config.public.backendBase || '').replace(/\/$/, '')
  return {
    apiBase: backendBase ? `${backendBase}/api` : API_BASE,
    assetBase: backendBase || ASSET_BASE,
  }
}

const getApiErrorCode = (error: any) =>
  error?.code ?? error?.data?.code ?? error?.response?._data?.code ?? error?.response?.status ?? error?.statusCode ?? error?.status

export const getApiErrorMessage = (error: unknown, fallback = '网络请求发生错误') => {
  const fetchError = error as any
  const message = fetchError?.data?.msg
    || fetchError?.data?.message
    || fetchError?.response?._data?.msg
    || fetchError?.response?._data?.message
  if (typeof message === 'string' && message.trim()) return message
  if (error instanceof Error && error.message && !fetchError?.request) return error.message
  return fallback
}

/** 页面 catch 里调用也只会显示一次，避免和请求器的保底消息重复。 */
export const showApiError = (error: unknown, fallback?: string) => {
  if (!import.meta.client) return
  const notifiedError = error as NotifiedError
  if (notifiedError?.[API_ERROR_NOTIFIED]) return
  ElMessage.error(getApiErrorMessage(error, fallback))
  if (error && (typeof error === 'object' || typeof error === 'function')) {
    notifiedError[API_ERROR_NOTIFIED] = true
  }
}

export const useApiFetch = <T = any>(
  url: string,
  options: ApiFetchOptions = {}
) => {
  const userStore = useUserStore()
  const { apiBase } = getBackendBases()
  const {
    suppressErrorMessage = false,
    errorFallback,
    onRequest,
    onRequestError,
    onResponse,
    onResponseError,
    ...requestOptions
  } = options

  const fetchOptions: FetchOptions = {
    ...requestOptions,
    baseURL: apiBase,

    async onRequest(context) {
      const headers = new Headers(context.options.headers)
      if (import.meta.server) {
        appendForwardedClientIpHeaders(headers, useRequestHeaders(['x-forwarded-for', 'x-real-ip']))
      }
      if (userStore.token) headers.set('Authorization', `Bearer ${userStore.token}`)
      context.options.headers = headers
      await runHooks(onRequest, context)
    },

    async onRequestError(context) {
      await runHooks(onRequestError, context)
    },

    async onResponse(context) {
      const body = context.response._data as ApiResponse
      if (body && typeof body.code === 'number' && body.code !== 200) {
        throw new ApiRequestError(body.msg || '业务请求失败', body.code)
      }
      await runHooks(onResponse, context)
    },

    async onResponseError(context) {
      await runHooks(onResponseError, context)
    }
  }

  // 非 2xx、HTTP 200 业务失败、断网/CORS/超时最终都由这里保底通知。
  return $fetch<T>(url, fetchOptions).catch((error: unknown) => {
    if (getApiErrorCode(error) === 401) {
      userStore.logout()
      showApiError(error, '登录状态已失效，请重新登录')
      if (import.meta.client) navigateTo('/user/auth')
    } else if (!suppressErrorMessage) {
      showApiError(error, errorFallback)
    }
    throw error
  })
}

export const useAssetUrl = () => {
  const { assetBase } = getBackendBases()

  return (url?: string, fallback = '') => {
    if (!url) return fallback
    if (/^https?:\/\//.test(url)) return url
    return `${assetBase}${url.startsWith('/') ? url : `/${url}`}`
  }
}
