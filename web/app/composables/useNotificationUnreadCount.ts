import { getNotificationUnreadCount } from '~/composables/api/notification'

const NOTIFICATION_POLL_INTERVAL = 30_000

export const useNotificationUnreadCount = () => {
  const count = useState<number>('notification-unread-count', () => 0)
  const loading = useState<boolean>('notification-unread-loading', () => false)
  const polling = useState<boolean>('notification-unread-polling', () => false)
  let timer: ReturnType<typeof setInterval> | null = null
  let requestInFlight = false

  const refresh = async () => {
    const userStore = useUserStore()
    if (!userStore.isLogin) {
      count.value = 0
      return
    }
    if (requestInFlight) return

    requestInFlight = true
    loading.value = true
    try {
      const response = await getNotificationUnreadCount()
      count.value = Number(response?.data?.unread || 0)
    } catch {
      // 轮询失败时保留上一次未读数，避免短暂断网导致角标闪烁或清零。
    } finally {
      loading.value = false
      requestInFlight = false
    }
  }

  const handleVisibilityChange = () => {
    if (document.visibilityState === 'visible') refresh()
  }

  const startPolling = () => {
    if (!import.meta.client || polling.value) return
    polling.value = true
    refresh()
    timer = window.setInterval(() => {
      if (document.visibilityState === 'visible') refresh()
    }, NOTIFICATION_POLL_INTERVAL)
    document.addEventListener('visibilitychange', handleVisibilityChange)
    window.addEventListener('focus', refresh)
  }

  const stopPolling = () => {
    if (!import.meta.client) return
    if (timer) window.clearInterval(timer)
    timer = null
    polling.value = false
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    window.removeEventListener('focus', refresh)
  }

  const decrement = (amount = 1) => {
    count.value = Math.max(0, count.value - amount)
  }

  return { count, loading, refresh, decrement, startPolling, stopPolling }
}
