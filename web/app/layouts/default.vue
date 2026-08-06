<script setup>
import {PhMagnifyingGlass,} from '@phosphor-icons/vue'
import { DEFAULT_AVATAR } from '~/config'

const keywords = ref();
const userStore = useUserStore()
const assetUrl = useAssetUrl()
const avatarUrl = computed(() => assetUrl(userStore.userInfo?.avatar, DEFAULT_AVATAR))
const avatarTargetUrl = computed(() => userStore.isLogin ? '/manage' : '/user/auth/?mode=login')
const isTouchDevice = ref(false)
const avatarPopoverTrigger = computed(() => isTouchDevice.value ? 'click' : 'hover')
const { count: notificationUnreadCount, startPolling: startNotificationPolling, stopPolling: stopNotificationPolling } = useNotificationUnreadCount()
const { settings: siteSettings } = useSiteSettings()
const logoUrl = computed(() => {
  const logo = siteSettings.value.siteLogo
  return logo.startsWith('/uploads/') ? assetUrl(logo, logo) : logo
})
const goSearch = (keywords) => {
  if (!keywords || !keywords.trim()) {
    return navigateTo('/');
  }
  return navigateTo(`/search/${encodeURIComponent(keywords)}`);
}

const handleAvatarClick = () => {
  if (isTouchDevice.value) return
  return navigateTo(avatarTargetUrl.value)
}

onMounted(() => {
  isTouchDevice.value = !window.matchMedia('(hover: hover) and (pointer: fine)').matches
})

watch(() => userStore.token, (token) => {
  if (token) startNotificationPolling()
  else {
    stopNotificationPolling()
    notificationUnreadCount.value = 0
  }
}, { immediate: true })

onBeforeUnmount(stopNotificationPolling)
</script>

<template>
  <div class="flex z-50 justify-between items-center bg-white/80 backdrop-blur-md h-15 w-full px-4 navbox navbox--responsive mb-10 select-none sticky top-0">
    <div class="brand-box">
      <nuxt-link to="/">
        <div class="flex items-center cursor-pointer">
          <img class="w-10 h-10" :src="logoUrl" alt="logo" draggable="false"/>
          <span class="brand-title ml-3.75 mr-1.25 text-[24px]">插件市场</span>
          <span class="brand-en text-[#a6a6a6] font-bold opacity-30">PluginMarket</span>
        </div>
      </nuxt-link>
    </div>
    <div class="nav-actions flex items-center gap-5">
      <div class=
               "search-box flex items-center bg-[#F1F5F9] h-8.75 rounded-[20px] text-[#64748b] pl-4
                border border-transparent transition-all
                hover: hover:border-[#cbd5e1] cursor-text
                focus-within:bg-white focus-within:border-blue-400 focus-within:ring-2 focus-within:ring-blue-100">
        <input v-model="keywords"
               @keydown.enter="goSearch(keywords)"
               class="search-input text-[14px] w-37.5 mr-2 bg-transparent outline-none" type="text"
               placeholder="搜索"/>
        <div @click="goSearch(keywords)" class="flex items-center pr-4 cursor-pointer">
          <PhMagnifyingGlass :size="18"/>
        </div>
      </div>
      <div class="shrink-0">
        <el-popover
            placement="bottom-start"
            transition="dropdown-fade"
            :width="280"
            :trigger="avatarPopoverTrigger"
        >
          <template #reference>
            <button type="button" class="relative block" @click="handleAvatarClick">
              <img
                  draggable="false"
                  :src="avatarUrl"
                  alt="User"
                  class="w-10 h-10 rounded-full avatar cursor-pointer transition-all hover:scale-110 hover:shadow-lg active:scale-95"
              >
              <span v-if="userStore.isLogin && notificationUnreadCount > 0" class="notification-badge">
                {{ notificationUnreadCount > 99 ? '99+' : notificationUnreadCount }}
              </span>
            </button>
          </template>
          <template #default>
            <layout-user v-if="userStore.isLogin"/>
            <layout-login v-else/>
          </template>
        </el-popover>
      </div>
    </div>
  </div>
  <slot></slot>
</template>

<style scoped lang="scss">
.navbox {
  border-bottom: 1px solid rgba(0, 0, 0, .04);
  box-shadow: 0 1px 3px rgba(0, 0, 0, .06), 0 2px 8px rgba(0, 0, 0, .02);

  .avatar {
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06), 0 2px 8px rgba(0, 0, 0, 0.02);
  }

  .notification-badge {
    position: absolute;
    top: -5px;
    right: -7px;
    display: flex;
    min-width: 19px;
    height: 19px;
    padding: 0 5px;
    align-items: center;
    justify-content: center;
    border: 1px solid #FECACA;
    border-radius: 10px;
    background: #FEF2F2;
    color: #DC2626;
    font-size: 10px;
    font-weight: 600;
    line-height: 1;
    pointer-events: none;
  }
}

@media (max-width: 767px) {
  .navbox--responsive {
    height: auto;
    min-height: 116px;
    margin-bottom: 16px;
    padding: 10px 12px 12px;
    flex-wrap: wrap;
    align-content: space-between;

    .brand-box {
      max-width: calc(100% - 52px);

      img {
        width: 36px;
        height: 36px;
      }
    }

    .brand-title {
      margin-left: 10px;
      font-size: 20px;
    }

    .brand-en {
      display: none;
    }

    .nav-actions {
      display: contents;
    }

    .search-box {
      order: 3;
      width: 100%;
      height: 40px;
    }

    .search-input {
      width: 100%;
      min-width: 0;
    }
  }
}
</style>
