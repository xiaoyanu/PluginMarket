<script setup>
import {computed} from 'vue'

const props = defineProps({
  width: {
    type: [String, Number],
    default: '100%'
  },
  height: {
    type: [String, Number],
    default: '32px'
  },
  radius: {
    type: [String, Number],
    default: '6px'
  }
})

const formatUnit = (val) => {
  if (typeof val === 'number') return `${val}px`
  if (typeof val === 'string' && !isNaN(Number(val))) return `${val}px`
  return val
}

const skeletonStyle = computed(() => ({
  width: formatUnit(props.width),
  height: formatUnit(props.height),
  borderRadius: formatUnit(props.radius)
}))
</script>

<template>
  <div class="skeleton-shimmer" :style="skeletonStyle"></div>
</template>

<style scoped>
.skeleton-shimmer {
  /* 使用极浅的灰色，适合纯白背景 */
  background-color: #fafafa;
  /* 渐变过渡配合纯白高光 */
  background-image: linear-gradient(90deg, #fafafa 25%, #ffffff 50%, #fafafa 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite linear;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>