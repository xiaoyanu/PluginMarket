import { useApiFetch } from '~/composables/useApiFetch'
import type { NotificationItem } from './notification'

export interface AdminNotificationListResponse {
  list: NotificationItem[]
  total: number
}

export interface NotificationForm {
  type: string
  audience_type: 'all' | 'user' | 'admin' | 'normal'
  receiver_id?: number | null
  title: string
  content: string
  icon_name: string
  icon_color: string
  action_url: string
  include_future_users: boolean
  status: number
  published_at?: string | null
  expires_at?: string | null
}

export const getAdminNotifications = (query: Record<string, string | number>) => {
  return useApiFetch<{ code: number; msg: string; data: AdminNotificationListResponse }>('/admin/notification/list', {
    method: 'GET',
    query,
  })
}

export const getAdminNotification = (id: number) =>
  useApiFetch<{ code: number; msg: string; data: NotificationItem }>(`/admin/notification/${id}`)

export const createAdminNotification = (body: NotificationForm) =>
  useApiFetch('/admin/notification', { method: 'POST', body })

export const updateAdminNotification = (id: number, body: NotificationForm) =>
  useApiFetch(`/admin/notification/${id}`, { method: 'PUT', body })

export const deleteAdminNotification = (id: number) =>
  useApiFetch(`/admin/notification/${id}`, { method: 'DELETE' })
