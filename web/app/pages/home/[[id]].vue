<script setup lang="ts">
import { getPluginsByUser, type PluginListQuery } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'
import { getApiErrorMessage } from '~/composables/useApiFetch'

type PluginListItem = {
  id: number
  name?: string
  title?: string
  desc_text?: string
  desc?: string
  icon?: string
  type?: number
  views?: number
  downloads?: number
  star?: number
  stars?: number
  frame_icon?: string
}

const route = useRoute()
const assetUrl = useAssetUrl()
const PAGE_SIZE = 9

const userId = computed(() => {
  const value = Number(route.params.id || 0)
  return Number.isInteger(value) && value > 0 ? value : 0
})

if (!userId.value) await navigateTo('/')

const toPluginCardList = (list: PluginListItem[] = []) => list.map((item) => ({
  id: item.id,
  title: item.title || item.name || '',
  desc: item.desc || item.desc_text || '',
  icon: assetUrl(item.icon, DEFAULT_PLUGIN_ICON),
  type: item.type ?? 0,
  views: item.views ?? 0,
  downloads: item.downloads ?? 0,
  stars: item.stars ?? item.star ?? 0,
  frame_icon: assetUrl(item.frame_icon, DEFAULT_PLUGIN_ICON),
}))

const filters = ref<Required<PluginListQuery>>({ type: -1, frameId: -1, tagId: -1 })
const pagination = reactive({ page: 1 })
const { data: pluginData, error: pluginError } = await useAsyncData(
  () => `user-plugin-list-${userId.value}`,
  () => getPluginsByUser(userId.value, PAGE_SIZE, pagination.page, filters.value),
  { watch: [userId], default: () => ({ data: { list: [], total: 0 } }),
  }
)
const userPluginList = computed(() => toPluginCardList(pluginData.value.data?.list ?? []))
const pluginTotal = computed(() => pluginData.value.data?.total ?? 0)
const pageUserInfo = ref<any>(null)
const userLoading = ref(true)

const loadUser = async () => {
  const requestedUserId = userId.value
  userLoading.value = true
  pageUserInfo.value = null
  try {
    const response = await useApiFetch(`/user/${userId.value}`, { suppressErrorMessage: true })
    if (response?.code !== 200) throw new Error(response?.msg)
    if (requestedUserId === userId.value) pageUserInfo.value = response.data || null
  } catch (error) {
    if (requestedUserId !== userId.value) return
    ElMessage.error(getApiErrorMessage(error))
    await new Promise((resolve) => window.setTimeout(resolve, 500))
    await navigateTo('/')
  } finally {
    if (requestedUserId === userId.value) userLoading.value = false
  }
}

onMounted(() => {
  if (userId.value) void loadUser()
})

const loadPlugins = async () => {
  pluginData.value = await getPluginsByUser(userId.value, PAGE_SIZE, pagination.page, filters.value)
}

const handlePageChange = async (page: number) => {
  pagination.page = page
  await loadPlugins()
}

const handleSidebarChange = async (value: Required<PluginListQuery>) => {
  filters.value = value
  pagination.page = 1
  await loadPlugins()
}

watch(userId, async (value) => {
  if (!value) return navigateTo('/')
  pagination.page = 1
  if (import.meta.client) await loadUser()
})
</script>

<template>
  <div>
    <div class="flex items-center justify-center mb-3 md:mb-5 max-w-350 px-3 md:px-5 m-auto">
      <div v-if="userLoading" class="w-full bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-4 md:p-5">
        <div class="flex gap-4 md:gap-5 items-center">
          <skeleton width="100" height="100" radius="20"/>
          <div class="flex-1 min-w-0 flex flex-col gap-3">
            <skeleton width="160" height="24"/>
            <skeleton width="100%" height="16"/>
            <skeleton width="70%" height="16"/>
          </div>
        </div>
        <div class="flex gap-2 mt-5">
          <skeleton width="88" height="22" radius="5"/>
          <skeleton width="180" height="22" radius="5"/>
        </div>
      </div>
      <user-card v-else-if="pageUserInfo" :user="pageUserInfo" responsive/>
    </div>
    <div class="flex flex-col md:flex-row items-start gap-3 md:gap-5 max-w-350 m-auto px-3 md:px-5 mb-6 md:mb-10">
      <select-sidebar v-model="filters" responsive @change="handleSidebarChange"/>
      <div class="w-full min-w-0 flex-1 select-none bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-3 md:p-5">
        <plugin-list responsive :plugin-list="userPluginList" title="插件 / 作品"/>
        <div class="mt-5 flex items-center justify-center">
          <el-pagination
              layout="prev, pager, next"
              :total="pluginTotal"
              :page-size="PAGE_SIZE"
              v-model:current-page="pagination.page"
              class="pageBox-diy"
              @current-change="handlePageChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>
