<script setup lang="ts">
import {PhCalendarMinus, PhHouse, PhIdentificationBadge} from "@phosphor-icons/vue";
import { DEFAULT_AVATAR, DEFAULT_USER_PROFILE } from '~/config'

type UserCardInfo = {
  id?: number
  username?: string
  nick?: string
  avatar?: string
  userdesc?: string
  titles?: Array<{ id?: number; name?: string; description?: string; icon?: string }>
  created?: string
}

const props = defineProps<{
  user?: UserCardInfo | null
  responsive?: boolean
}>()

const route = useRoute()
const showHomeButton = computed(() => !route.path.startsWith('/home/'))

const userStore = useUserStore()
const assetUrl = useAssetUrl()
// 显式传入 null 表示目标用户仍在加载，不能回退到当前登录用户。
const userInfo = computed(() => props.user !== undefined ? props.user : userStore.userInfo)

const displayName = computed(() => userInfo.value?.nick || userInfo.value?.username || '未命名用户')
const userDesc = computed(() => userInfo.value?.userdesc || DEFAULT_USER_PROFILE)
const avatarUrl = computed(() => assetUrl(userInfo.value?.avatar, DEFAULT_AVATAR))
const titles = computed(() => userInfo.value?.titles || [])

const formatDate = (value?: string) => {
  if (!value) return '未知'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const goHome = () => {
  navigateTo(userInfo.value?.id ? `/home/${userInfo.value.id}` : '/')
}
</script>

<template>
  <div :class="['user-card flex-1 shadow-pmbox bg-white rounded-[20px] p-5 select-none', { 'user-card--responsive': responsive }]">
    <div class="user-main flex gap-5">
      <div class="flex items-center shrink-0">
        <img class="avatar" :src="avatarUrl" alt="avatar" draggable="false"/>
      </div>
      <div class="user-content flex flex-col min-w-0">
        <div class="flex items-center gap-3 flex-wrap">
          <div class="font-bold text-[#1E293B] text-[20px]">{{ displayName }}</div>
          <div class="flex gap-2">
            <el-tooltip v-for="title in titles" :key="title.id || title.name" placement="bottom" effect="light" popper-class="title-description-tooltip">
              <template #content><div class="title-description-content" v-html="title.description || title.name || '称号'" /></template>
              <img class="titleIcon" :src="assetUrl(title.icon)" draggable="false" alt="">
            </el-tooltip>
          </div>
        </div>
        <div class="text-sm line-clamp-4 text-[#64748B] select-text whitespace-pre-wrap break-words">
          {{ userDesc }}
        </div>
      </div>
    </div>
    <div class="user-meta mt-5 flex items-center justify-between">
      <div class="user-tags flex gap-2">
        <span class="tag">
          <PhIdentificationBadge weight="duotone"/>
          UID：{{ userInfo?.id || '未知' }}
        </span>
        <span class="tag">
          <PhCalendarMinus weight="duotone"/>
          注册时间：{{ formatDate(userInfo?.created) }}
        </span>
      </div>
      <div v-if="showHomeButton">
        <pm-button text="查看主页" @click="goHome">
          <template #icon>
            <PhHouse/>
          </template>
        </pm-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.avatar {
  width: 100px;
  height: 100px;
  border-radius: 20px;
  background-color: #fff;
  border: 4px solid #FFF;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.titleIcon {
  height: 20px;
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  background-color: #F4F4F5;
  color: #909399;
  border: 1px solid rgba(0, 0, 0, 0.1);
  font-size: 12px;
  padding: 0 9px;
  border-radius: 5px;
  height: 22px;
  line-height: 1;
  flex-shrink: 0;
}

@media (max-width: 767px) {
  .user-card--responsive {
    width: 100%;
    padding: 16px;
    border-radius: 16px;

    .user-main {
      gap: 14px;
      align-items: flex-start;
    }

    .avatar {
      width: 72px;
      height: 72px;
      border-width: 2px;
      border-radius: 16px;
    }

    .user-content {
      flex: 1;
    }

    .user-meta {
      align-items: flex-start;
    }

    .user-tags {
      min-width: 0;
      flex-wrap: wrap;
    }

    .tag {
      max-width: 100%;
      height: auto;
      min-height: 24px;
      padding: 4px 8px;
      white-space: normal;
      line-height: 1.3;
    }
  }
}
</style>