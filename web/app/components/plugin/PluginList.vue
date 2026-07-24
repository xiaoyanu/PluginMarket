<script setup lang="ts">
type TypePluginCard = {
  desc: string;
  downloads: number;
  frame_icon: string;
  icon: string;
  id: number;
  stars: number;
  title: string;
  type: number;
  views: number;
}

defineProps({
  pluginList: {
    type: Array<TypePluginCard>,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  responsive: {
    type: Boolean,
    default: false,
  }
})

</script>

<template>
  <div>
    <div class="section-header">
      <span class="text-[20px] font-bold">{{ title }}</span>
      <el-tag size="small">{{ pluginList.length }}</el-tag>
    </div>
    <div v-if="pluginList.length>0" :class="['plugin-grid', { 'plugin-grid--responsive': responsive }]">
      <nuxt-link v-for="item in pluginList" :key="item.id"
                 :to="`/plugin/`+item.id" target="_blank" class="plugin-link">
        <plugin-card
            :responsive="responsive"
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
    <div v-else>
      <el-empty description="这里什么也没有💨"/>
    </div>
  </div>
</template>

<style scoped lang="scss">
.plugin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 340px);
  justify-content: center;
  gap: 20px;
}

@media (max-width: 767px) {
  .plugin-grid--responsive {
    grid-template-columns: minmax(0, 1fr);
    gap: 12px;

    .plugin-link {
      min-width: 0;
    }
  }
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 5px;

  &:before {
    content: "";
    display: inline-block;
    width: 4px;
    height: 20px;
    background-color: #ffaa00;
    margin-right: 8px;
    border-radius: 2px;
  }

  :deep(.el-tag) {
    border-radius: 10px;
  }
}

@media (max-width: 767px) {
  .section-header {
    align-items: flex-start;
    margin-bottom: 12px;

    span:first-of-type {
      min-width: 0;
      overflow-wrap: anywhere;
      font-size: 18px;
      line-height: 1.4;
    }
  }
}
</style>