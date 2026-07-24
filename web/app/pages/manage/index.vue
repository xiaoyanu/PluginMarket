<script setup lang="ts">
import { PhBellRinging, PhImageSquare } from '@phosphor-icons/vue'
import type { NotificationItem } from '~/composables/api/notification'
import { getNotifications, markAllNotificationsRead, markNotificationRead } from '~/composables/api/notification'
import { formatNotificationTime, getNotificationAction } from '~/utils/manage-notification'

const userStore = useUserStore()
const { count: navUnreadCount } = useNotificationUnreadCount()

const notifications = ref<NotificationItem[]>([])
const unreadCount = ref(0)
const loading = ref(false)
const handlingId = ref<number | null>(null)
const markingAllRead = ref(false)

const needsEmail = computed(() => !userStore.userInfo?.email?.trim())
const needsAvatar = computed(() => !userStore.userInfo?.avatar?.trim())
const notificationBadgeCount = computed(() => (
  unreadCount.value + Number(needsEmail.value) + Number(needsAvatar.value)
))

const loadNotifications = async () => {
  loading.value = true
  try {
    const response = await getNotifications(1, 50)
    notifications.value = response?.data?.list || []
    unreadCount.value = response?.data?.unread || 0
    navUnreadCount.value = unreadCount.value
  } finally {
    loading.value = false
  }
}

const goInfo = () => window.open('/manage/info', '_blank', 'noopener,noreferrer')

const processNotification = async (item: NotificationItem) => {
  const action = getNotificationAction(item)
  if (!action.visible) return
  handlingId.value = item.id
  if (action.target) {
    window.open(action.target, '_blank', 'noopener,noreferrer')
  }
  try {
    await markNotificationRead(item.id)
    item.is_read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
    navUnreadCount.value = unreadCount.value
  } finally {
    handlingId.value = null
  }
}

const markAllRead = async () => {
  if (unreadCount.value === 0) return
  markingAllRead.value = true
  try {
    await markAllNotificationsRead()
    notifications.value.forEach(item => { item.is_read = true })
    unreadCount.value = 0
    navUnreadCount.value = 0
    ElMessage.success('所有通知已标记为已读')
  } finally {
    markingAllRead.value = false
  }
}

onMounted(loadNotifications)

definePageMeta({
  layout: 'manage'
})
</script>

<template>
  <div class="flex flex-col gap-5">
    <div>
      <user-card responsive />
    </div>

    <manage-box title="通知" :value="notificationBadgeCount">
      <template #header>
        <pm-button
          v-if="unreadCount > 0"
          text="一键已读"
          color="white"
          :loading="markingAllRead"
          @click="markAllRead"
        />
      </template>

      <manage-notify
        v-if="needsEmail"
        title="邮箱绑定提醒"
        content="你还没有绑定邮箱，绑定邮箱后可以收到插件评论、审核结果等相关通知。"
        @process="goInfo"
      >
        <template #icon>
          <PhBellRinging color="#FF8D1A" weight="duotone" />
        </template>
      </manage-notify>

      <manage-notify
        v-if="needsAvatar"
        title="还没有设置头像"
        content="(｡･∀･)ﾉﾞ 嗨，来给自己设置一个个性头像吧！"
        @process="goInfo"
      >
        <template #icon>
          <PhImageSquare color="#1AB2FE" weight="duotone" />
        </template>
      </manage-notify>

      <template v-if="!loading">
        <manage-notify
          v-for="item in notifications"
          :key="item.id"
          :title="item.title"
          :content="item.content"
          :time="formatNotificationTime(item.published_at || item.created_at)"
          :button-text="getNotificationAction(item).text"
          :show-button="getNotificationAction(item).visible"
          :disabled="handlingId === item.id || markingAllRead"
          @process="processNotification(item)"
        >
          <template #icon>
            <manage-notification-icon
              :name="item.icon_name"
              :color="item.icon_color"
              :is-read="item.is_read"
              :size="40"
            />
          </template>
        </manage-notify>
      </template>

      <div
        v-if="!loading && !needsEmail && !needsAvatar && notifications.length === 0"
        class="rounded-[20px] bg-[#F8FAFC] py-8 text-center text-sm text-[#94A3B8]"
      >
        暂无通知
      </div>
    </manage-box>
  </div>
</template>
