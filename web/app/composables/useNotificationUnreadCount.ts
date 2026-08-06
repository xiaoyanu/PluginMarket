import { getNotifications } from '~/composables/api/notification'

export const useNotificationUnreadCount = () => {
  const count = useState<number>('notification-unread-count', () => 0)
  const loading = useState<boolean>('notification-unread-loading', () => false)

  const refresh = async () => {
    const userStore = useUserStore()
    if (!userStore.isLogin) {
      count.value = 0
      return
    }

    loading.value = true
    try {
      const response = await getNotifications(1, 1)
      count.value = Number(response?.data?.unread || 0)
    } catch {
      count.value = 0
    } finally {
      loading.value = false
    }
  }

  const decrement = (amount = 1) => {
    count.value = Math.max(0, count.value - amount)
  }

  return { count, loading, refresh, decrement }
}
