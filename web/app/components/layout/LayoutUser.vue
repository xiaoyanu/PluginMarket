<script setup lang="ts">
import {
  PhCardsThree,
  PhHouse,
  PhIdentificationCard,
  PhPuzzlePiece,
  PhSignOut,
  PhStar,
} from "@phosphor-icons/vue";
import { DEFAULT_AVATAR } from '~/config'

const userStore = useUserStore()
const assetUrl = useAssetUrl()
const userInfo = computed(() => userStore.userInfo)
const displayName = computed(() => userInfo.value?.nick || userInfo.value?.username || '已登录用户')
const avatarUrl = computed(() => assetUrl(userInfo.value?.avatar, DEFAULT_AVATAR))
const homeUrl = computed(() => userInfo.value?.id ? `/home/${userInfo.value.id}` : '/')

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  navigateTo('/')
}
</script>

<template>
  <div class="flex flex-col p-2 select-none gap-3">
    <div class="flex gap-3 items-center">
      <div>
        <img class="rounded-full w-12.5 h-12.5"
             :src="avatarUrl"
             alt="avatar" draggable="false"/>
      </div>
      <div class="flex flex-col gap-1">
        <span class="text-[#1E293B] font-bold text-[16px]">欢迎来到 插件市场</span>
        <span>{{ displayName }}</span>
      </div>
    </div>
    <hr class="border-gray-200">
    <div class="flex items-center gap-2 px-2 flex-wrap">
      <pm-button text="我的主页" @click="navigateTo(homeUrl)">
        <template #icon>
          <PhHouse/>
        </template>
      </pm-button>
      <pm-button text="管理首页" @click="navigateTo('/manage')">
        <template #icon>
          <PhCardsThree/>
        </template>
      </pm-button>
      <pm-button text="我的插件" @click="navigateTo('/manage/plugin')">
        <template #icon>
          <PhPuzzlePiece/>
        </template>
      </pm-button>
      <pm-button text="我的信息" @click="navigateTo('/manage/info')">
        <template #icon>
          <PhIdentificationCard/>
        </template>
      </pm-button>
      <pm-button text="我的收藏" @click="navigateTo('/manage/star')">
        <template #icon>
          <PhStar/>
        </template>
      </pm-button>
      <pm-button text="退出登录" @click="handleLogout" color="red">
        <template #icon>
          <PhSignOut/>
        </template>
      </pm-button>
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>
