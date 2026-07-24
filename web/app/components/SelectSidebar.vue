<script setup lang="ts">
import {PhSquaresFour} from "@phosphor-icons/vue";
import { getFrameList, getTagList, type PluginListQuery } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'

interface FrameItem { id: number; name: string; icon?: string }
interface TagItem { id: number; name: string }

type SidebarValue = Required<PluginListQuery>

const props = defineProps<{
  modelValue?: SidebarValue
  responsive?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: SidebarValue): void
  (e: 'change', value: SidebarValue): void
}>()

const selectedValue = ref(props.modelValue?.type ?? -1)
const selectedFrameId = ref(props.modelValue?.frameId ?? -1)
const selectedTagId = ref(props.modelValue?.tagId ?? -1)

watch(() => props.modelValue, (value) => {
  selectedValue.value = value?.type ?? -1
  selectedFrameId.value = value?.frameId ?? -1
  selectedTagId.value = value?.tagId ?? -1
}, { deep: true })

const frameData: any = await getFrameList()
const tagData: any = await getTagList()
const frameList = computed<FrameItem[]>(() => frameData.data?.list ?? [])
const tagList = computed<TagItem[]>(() => tagData.data?.list ?? [])

const assetUrl = useAssetUrl()

const emitChange = () => {
  const value = {
    type: selectedValue.value,
    frameId: selectedFrameId.value,
    tagId: selectedTagId.value,
  }
  emit('update:modelValue', value)
  emit('change', value)
}

const selectFrame = (id: number) => {
  selectedFrameId.value = id
  emitChange()
}

const selectTag = (id: number) => {
  selectedTagId.value = selectedTagId.value === id ? -1 : id
  emitChange()
}
</script>

<template>
  <div :class="['select-sidebar', { 'select-sidebar--responsive': responsive }]">
    <div class="mobile-filters">
      <el-select v-model="selectedValue" placeholder="类型" aria-label="插件类型" @change="emitChange">
        <el-option :value="-1" label="全部类型"/>
        <el-option :value="0" label="免费"/>
        <el-option :value="1" label="收费"/>
      </el-select>
      <el-select v-model="selectedFrameId" placeholder="框架" aria-label="插件框架" @change="emitChange">
        <el-option :value="-1" label="全部框架"/>
        <el-option v-for="item in frameList" :key="item.id" :value="item.id" :label="item.name"/>
      </el-select>
      <el-select v-model="selectedTagId" placeholder="标签" aria-label="插件标签" @change="emitChange">
        <el-option :value="-1" label="全部标签"/>
        <el-option v-for="item in tagList" :key="item.id" :value="item.id" :label="item.name"/>
      </el-select>
    </div>
    <div class="desktop-filters">
    <div class="item-box">
      <div class="item-title">类型</div>
      <div class="item-content">
        <el-select v-model="selectedValue" placeholder="请选择" @change="emitChange">
          <el-option :value="-1" label="全部"/>
          <el-option :value="0" label="免费"/>
          <el-option :value="1" label="收费"/>
        </el-select>
      </div>
    </div>
    <div class="item-box">
      <div class="item-title">框架</div>
      <div class="item-content">
        <ul class="framework-list">
          <li :class="{ active: selectedFrameId === -1 }" @click="selectFrame(-1)">
            <PhSquaresFour :size="18"/>
            全部
          </li>
          <li v-for="item in frameList" :key="item.id" :class="{ active: selectedFrameId === item.id }" @click="selectFrame(item.id)">
            <img :src="assetUrl(item.icon, DEFAULT_PLUGIN_ICON)" alt="frame"> {{ item.name }}
          </li>
        </ul>
      </div>
    </div>
    <div class="item-box">
      <div class="item-title">标签</div>
      <div class="item-content">
        <div class="flex flex-wrap gap-1.5">
          <el-tag v-for="item in tagList" :key="item.id" :type="selectedTagId === item.id ? 'success' : 'primary'" @click="selectTag(item.id)">{{ item.name }}</el-tag>
        </div>
      </div>
    </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.select-sidebar {
  width: 240px;
  padding: 20px;
  flex-shrink: 0;
  border-radius: 20px;
  background: #fff;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.06);
}

.mobile-filters {
  display: none;
}

@media (max-width: 767px) {
  .select-sidebar--responsive {
    width: 100%;
    padding: 12px;
    border-radius: 16px;

    .desktop-filters {
      display: none;
    }

    .mobile-filters {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 8px;

      :deep(.el-select__wrapper) {
        min-height: 40px;
        border-radius: 12px;
      }
    }
  }
}

.item-box {
  margin-bottom: 24px;
  user-select: none;

  &:last-child {
    margin-bottom: 0;
  }


  .item-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    color: #1E293B;
    flex-shrink: 0;

    &:before {
      content: '';
      display: inline-block;
      width: 4px;
      height: 20px;
      background-color: #00BAAD;
      margin-right: 8px;
      border-radius: 2px;
    }
  }

  .item-content {
    color: #475569;

    :deep(.el-tag) {
      border-radius: 10px;
      cursor: pointer;
    }

    :deep(.el-select__wrapper) {
      border-radius: 20px;
    }

    .framework-list {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        padding: 8px 12px;
        border-radius: 7px;
        margin-bottom: 4px;
        cursor: pointer;
        font-size: 14px;
        display: flex;
        align-items: center;
        gap: 8px;
        color: #555;
        transition: all .3s;

        img {
          width: 18px;
          height: 18px;
          border-radius: 20%;
        }

        &:hover {
          background-color: #f5f6f8;

        }

        &.active {
          background-color: #5692F5;
          color: #FFF;
          box-shadow: 0 2px 8px rgba(66, 133, 244, 0.35);
        }
      }
    }
  }
}
</style>