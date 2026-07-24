<script setup lang="ts">
import { DEFAULT_AVATAR, DEFAULT_USER_PROFILE } from '~/config'

const props = defineProps<{
  author?: {
    id: number
    nick?: string
    avatar?: string
    userdesc?: string
    titles?: Array<{ id: number; name?: string; icon?: string }>
  } | null
}>()

const assetUrl = useAssetUrl()
const author = computed(() => props.author || null)
const authorAvatar = computed(() => assetUrl(author.value?.avatar, DEFAULT_AVATAR))
const authorName = computed(() => author.value?.nick || '未命名作者')
const authorDesc = computed(() => author.value?.userdesc || DEFAULT_USER_PROFILE)
const authorTitles = computed(() => author.value?.titles || [])
const authorHomeUrl = computed(() => author.value?.id ? `/home/${author.value.id}` : '')
</script>

<template>
  <NuxtLink
    v-if="authorHomeUrl"
    :to="authorHomeUrl"
    target="_blank"
    rel="noopener noreferrer"
    class="block"
  >
    <div class="flex gap-2 flex-col hover:opacity-90 transition-opacity">
      <div class="flex gap-3">
        <div class="flex items-start shrink-0">
          <img class="avatar" :src="authorAvatar" alt="avatar" draggable="false"/>
        </div>
        <div class="flex flex-col gap-0.5">
          <div class="font-bold text-[#1E293B] text-[16px]">{{ authorName }}</div>
          <div class="flex gap-1 flex-wrap">
            <img v-for="title in authorTitles" :key="title.id" class="titleIcon" :src="assetUrl(title.icon)" draggable="false" alt="">
          </div>
        </div>
      </div>
      <div class="text-sm line-clamp-4 text-[#64748B] select-text whitespace-pre-wrap break-words">
        {{ authorDesc }}
      </div>
    </div>
  </NuxtLink>
  <div v-else class="flex gap-2 flex-col">
    <div class="flex gap-3">
      <div class="flex items-start shrink-0">
        <img class="avatar" :src="authorAvatar" alt="avatar" draggable="false"/>
      </div>
      <div class="flex flex-col gap-0.5">
        <div class="font-bold text-[#1E293B] text-[16px]">{{ authorName }}</div>
        <div class="flex gap-1 flex-wrap">
          <img v-for="title in authorTitles" :key="title.id" class="titleIcon" :src="assetUrl(title.icon)" draggable="false" alt="">
        </div>
      </div>
    </div>
    <div class="text-sm line-clamp-4 text-[#64748B] select-text whitespace-pre-wrap break-words">
      {{ authorDesc }}
    </div>
  </div>
</template>

<style scoped lang="scss">
.avatar {
  width: 50px;
  height: 50px;
  border-radius: 10px;
  background-color: #fff;
  border: 1px solid #FFF;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.titleIcon {
  height: 20px;;
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
</style>