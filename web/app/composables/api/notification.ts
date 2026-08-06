import { useApiFetch } from '~/composables/useApiFetch'

export interface NotificationItem {
  id: number
  type: string
  audience_type: string
  receiver_id?: number | null
  sender_id?: number | null
  title: string
  content: string
  icon_name: string
  icon_color: string
  source_type?: string | null
  source_id?: number | null
  action_url?: string | null
  extra_data?: string | null
  include_future_users: boolean
  status: number
  is_read: boolean
  published_at?: string | null
  expires_at?: string | null
  created_at?: string
  updated_at?: string
}

export interface NotificationListResponse {
  list: NotificationItem[]
  total: number
  unread: number
}

export const getNotifications = (page = 1, pageSize = 20) => {
  return useApiFetch<{ code: number; msg: string; data: NotificationListResponse }>('/notification/list', {
    method: 'GET',
    query: { page, pageSize },
  })
}

export const getNotificationUnreadCount = () => {
  return useApiFetch<{ code: number; msg: string; data: { unread: number } }>('/notification/unread-count', {
    method: 'GET',
    suppressErrorMessage: true,
  })
}

export const markNotificationRead = (id: number) => {
  return useApiFetch(`/notification/${id}/read`, { method: 'PUT' })
}

export const markAllNotificationsRead = () => {
  return useApiFetch('/notification/read-all', { method: 'PUT' })
}
