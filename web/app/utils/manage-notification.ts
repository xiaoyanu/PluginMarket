import dayjs from 'dayjs'

export interface NotificationActionInput {
  is_read: boolean
  action_url?: string | null
}

export interface NotificationAction {
  visible: boolean
  text: string
  target: string | null
}

export const getNotificationAction = (item: NotificationActionInput): NotificationAction => {
  if (item.is_read) {
    return { visible: false, text: '', target: null }
  }
  const target = item.action_url?.trim() || null
  return {
    visible: true,
    text: target ? '前往处理' : '标记已读',
    target,
  }
}

export const formatNotificationTime = (value?: string | null) => (
  value ? dayjs(value).format('YYYY-MM-DD HH:mm') : ''
)
