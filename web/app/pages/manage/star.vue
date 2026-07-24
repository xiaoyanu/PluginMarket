<script setup lang="ts">
import { getStarPlugins } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'
definePageMeta({
  layout: 'manage'
})

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
const tableLoading = ref(false)
const pagination = reactive({ page: 1, pageSize: 12, total: 0 })
const dataList = ref<any[]>([])

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

const loadData = async () => {
  tableLoading.value = true
  try {
    const res: any = await getStarPlugins(pagination.pageSize, pagination.page)
    const list = res.data?.list ?? []
    dataList.value = toPluginCardList(list)
    pagination.total = res.data?.total ?? 0
    const maxPage = Math.max(1, Math.ceil(pagination.total / pagination.pageSize))
    if (pagination.page > maxPage) {
      pagination.page = maxPage
      await loadData()
    }
  } finally {
    tableLoading.value = false
  }
}

const handlePageChange = async (page: number) => {
  pagination.page = page
  await loadData()
}

onMounted(() => {
  void loadData()
})
</script>

<template>
  <manage-box title="我的收藏" :value="pagination.total">
    <div v-loading="tableLoading">
      <div v-if="dataList.length > 0" class="star-grid grid grid-cols-[repeat(auto-fill,340px)] justify-center gap-5">
        <nuxt-link v-for="item in dataList" :key="item.id" :to="`/plugin/${item.id}`" target="_blank">
          <plugin-card
              responsive
              :title="item.title"
              :desc="item.desc"
              :type="item.type"
              :views="item.views"
              :download="item.downloads"
              :stars="item.stars"
              :icon="item.icon"
              :frameIcon="item.frame_icon"
          />
        </nuxt-link>
      </div>
      <el-empty v-else description="这里什么也没有💨" />
    </div>

    <div class="mt-5 flex items-center justify-center">
      <el-pagination
          layout="prev, pager, next"
          :total="pagination.total"
          :page-size="pagination.pageSize"
          v-model:current-page="pagination.page"
          class="pageBox-diy"
          @current-change="handlePageChange"
      />
    </div>
  </manage-box>
</template>

<style scoped lang="scss">
:deep(.section-header) {
  min-height: 32px;
  gap: 0;

  span {
    font-size: 16px;
  }

  &:before {
    background-color: #00BAAD;
  }
}

@media (max-width: 767px) {
  .star-grid {
    grid-template-columns: minmax(0, 1fr);
    gap: 12px;

    > a {
      min-width: 0;
    }
  }
}
</style>

