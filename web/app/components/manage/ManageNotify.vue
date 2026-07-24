<script setup lang="ts">
import { PhArrowFatLinesUp } from '@phosphor-icons/vue'

const props = withDefaults(defineProps<{
  title: string
  content: string
  time?: string
  buttonText?: string
  showButton?: boolean
  disabled?: boolean
}>(), {
  time: '',
  buttonText: '前往处理',
  showButton: true,
  disabled: false,
})

defineEmits<{
  process: []
}>()
</script>

<template>
  <div class="main-box flex items-center gap-3 flex-1 shadow-[0_2px_10px_rgba(0,0,0,0.03)] bg-white rounded-[20px] p-5 border border-gray-200">
    <div class="text-[40px] shrink-0">
      <slot name="icon" />
    </div>
    <div class="flex flex-1 items-center gap-3 justify-between min-w-0">
      <div class="min-w-0 w-full max-w-170">
        <div class="font-bold truncate" :title="title">{{ title }}</div>
        <div class="text-sm truncate" :title="content">{{ content }}</div>
        <time v-if="props.time" class="block text-xs text-[#94A3B8]">{{ props.time }}</time>
      </div>
      <div v-if="props.showButton" class="shrink-0">
        <pm-button :text="props.buttonText" :disabled="props.disabled" @click="$emit('process')">
          <template #icon>
            <PhArrowFatLinesUp />
          </template>
        </pm-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.main-box {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 15px rgba(0, 0, 0, 0.06);
    border-color: rgba(0, 0, 0, 0.08);
  }
}

@media (max-width: 767px) {
  .main-box {
    align-items: flex-start;
    padding: 14px;
    border-radius: 16px;

    > .flex-1 {
      flex-direction: column;
      align-items: stretch;
    }

    .truncate {
      overflow: visible;
      white-space: normal;
      text-overflow: initial;
      word-break: break-word;
    }

    :deep(button) {
      width: 100%;
      min-height: 44px;
      justify-content: center;
    }
  }
}
</style>
