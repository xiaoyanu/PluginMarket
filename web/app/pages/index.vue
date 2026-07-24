<script setup lang="ts">
import { getHotPlugins, getLatestPlugins, type PluginListQuery } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'
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

const assetUrl = useAssetUrl()

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

const HOT_PAGE_SIZE = 3
const LATEST_PAGE_SIZE = 9

const filters = ref<Required<PluginListQuery>>({ type: -1, frameId: -1, tagId: -1 })
const pagination = reactive({ page: 1 })

const latestData = ref<any>(await getLatestPlugins(LATEST_PAGE_SIZE, pagination.page, filters.value))
const latestPluginList = computed(() => toPluginCardList(latestData.value.data?.list ?? []))
const latestTotal = computed(() => latestData.value.data?.total ?? 0)

const hotData = ref<any>(await getHotPlugins(HOT_PAGE_SIZE, 1, filters.value))
const hotPluginList = computed(() => toPluginCardList((hotData.value.data?.list ?? []).slice(0, HOT_PAGE_SIZE)))

const loadLatest = async () => {
  latestData.value = await getLatestPlugins(LATEST_PAGE_SIZE, pagination.page, filters.value)
}

const loadHot = async () => {
  hotData.value = await getHotPlugins(HOT_PAGE_SIZE, 1, filters.value)
}

const handlePageChange = async (page: number) => {
  pagination.page = page
  await loadLatest()
}

const handleSidebarChange = async (value: Required<PluginListQuery>) => {
  filters.value = value
  pagination.page = 1
  await Promise.all([loadHot(), loadLatest()])
}
</script>

<template>
  <div class="flex flex-col md:flex-row items-start gap-3 md:gap-5 max-w-350 m-auto px-3 md:px-5 mb-6 md:mb-10">
    <select-sidebar v-model="filters" responsive @change="handleSidebarChange"/>
    <div class="w-full min-w-0 flex-1 select-none bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-3 md:p-5">
      <plugin-list responsive :plugin-list="hotPluginList" title="热门"/>
      <plugin-list responsive class="mt-7 md:mt-10" :plugin-list="latestPluginList" title="最新"/>
      <div class="mt-5 flex items-center justify-center">
        <el-pagination
            layout="prev, pager, next"
            :total="latestTotal"
            :page-size="LATEST_PAGE_SIZE"
            v-model:current-page="pagination.page"
            class="pageBox-diy"
            @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>
