<script setup lang="ts">
import { searchPlugins, type PluginListQuery } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'
const route = useRoute();
const keywords = route.params.keywords

if (!keywords) {
  await navigateTo('/')
}

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

const filters = ref<Required<PluginListQuery>>({ type: -1, frameId: -1, tagId: -1 })
const pagination = reactive({ page: 1, pageSize: 20 })

const data = ref<any>(await searchPlugins(String(keywords || ''), pagination.pageSize, pagination.page, filters.value))
const dataList = computed(() => toPluginCardList(data.value.data?.list ?? []))
const total = computed(() => data.value.data?.total ?? 0)

const loadData = async () => {
  data.value = await searchPlugins(String(keywords || ''), pagination.pageSize, pagination.page, filters.value)
}

const handlePageChange = async (page: number) => {
  pagination.page = page
  await loadData()
}

const handleSidebarChange = async (value: Required<PluginListQuery>) => {
  filters.value = value
  pagination.page = 1
  await loadData()
}
</script>
<template>
  <div class="flex flex-col md:flex-row items-start gap-3 md:gap-5 max-w-350 m-auto px-3 md:px-5 mb-6 md:mb-10">
    <select-sidebar v-model="filters" responsive @change="handleSidebarChange"/>
    <div class="w-full min-w-0 flex-1 select-none bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-3 md:p-5">
      <plugin-list responsive :plugin-list="dataList" :title="`搜索 “`+keywords+`” 的结果`"/>
      <div class="mt-5 flex items-center justify-center">
        <el-pagination
            layout="prev, pager, next"
            :total="total"
            :page-size="pagination.pageSize"
            v-model:current-page="pagination.page"
            class="pageBox-diy"
            @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>
