<script setup lang="ts">
import type { Component } from 'vue'
import {
  PhAirplaneTaxiing,
  PhBellRinging,
  PhChatText,
  PhChats,
  PhCheckCircle,
  PhImageSquare,
  PhMegaphone,
  PhPuzzlePiece,
  PhSealCheck,
  PhSealWarning,
  PhShieldWarning,
  PhWarningCircle,
  PhXCircle,
} from '@phosphor-icons/vue'

const props = withDefaults(defineProps<{
  name?: string
  color?: string
  isRead?: boolean
  size?: number
}>(), {
  name: 'PhBellRinging',
  color: '#64748B',
  isRead: false,
  size: 32,
})

// 使用显式注册表，避免生产构建时动态读取命名空间被 tree-shaking 掉。
const iconRegistry: Record<string, Component> = {
  PhAirplaneTaxiing,
  PhBellRinging,
  PhChatText,
  PhChats,
  PhCheckCircle,
  PhImageSquare,
  PhMegaphone,
  PhPuzzlePiece,
  PhSealCheck,
  PhSealWarning,
  PhShieldWarning,
  PhWarningCircle,
  PhXCircle,
}

const icon = computed<Component>(() => iconRegistry[props.name || ''] || PhBellRinging)
const iconColor = computed(() => props.isRead ? '#94A3B8' : (props.color || '#64748B'))
</script>

<template>
  <component :is="icon" :size="size" :color="iconColor" weight="duotone" />
</template>
