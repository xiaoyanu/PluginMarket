<script setup lang="ts">
import type { Component } from 'vue'
import {
  PhAirplaneTaxiing,
  PhApplePodcastsLogo,
  PhArrowsCounterClockwise,
  PhBalloon,
  PhBarricade,
  PhBellRinging,
  PhBracketsCurly,
  PhBroom,
  PhBug,
  PhCake,
  PhCalendar,
  PhCalendarStar,
  PhCarrot,
  PhChatText,
  PhChats,
  PhChatsTeardrop,
  PhCheckCircle,
  PhCloudArrowDown,
  PhCloudArrowUp,
  PhCoffee,
  PhConfetti,
  PhCrown,
  PhCursorClick,
  PhDeviceMobileSpeaker,
  PhFire,
  PhHeart,
  PhImageSquare,
  PhLightbulb,
  PhMegaphone,
  PhNotification,
  PhPuzzlePiece,
  PhQrCode,
  PhSealCheck,
  PhSealWarning,
  PhShareNetwork,
  PhShieldWarning,
  PhShootingStar,
  PhSignOut,
  PhSketchLogo,
  PhSnapchatLogo,
  PhStar,
  PhStarFour,
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
  PhApplePodcastsLogo,
  PhArrowsCounterClockwise,
  PhBalloon,
  PhBarricade,
  PhBellRinging,
  PhBracketsCurly,
  PhBroom,
  PhBug,
  PhCake,
  PhCalendar,
  PhCalendarStar,
  PhCarrot,
  PhChatText,
  PhChats,
  PhChatsTeardrop,
  PhCheckCircle,
  PhCloudArrowDown,
  PhCloudArrowUp,
  PhCoffee,
  PhConfetti,
  PhCrown,
  PhCursorClick,
  PhDeviceMobileSpeaker,
  PhFire,
  PhHeart,
  PhImageSquare,
  PhLightbulb,
  PhMegaphone,
  PhNotification,
  PhPuzzlePiece,
  PhQrCode,
  PhSealCheck,
  PhSealWarning,
  PhShareNetwork,
  PhShieldWarning,
  PhShootingStar,
  PhSignOut,
  PhSketchLogo,
  PhSnapchatLogo,
  PhStar,
  PhStarFour,
  PhWarningCircle,
  PhXCircle,
}

const icon = computed<Component>(() => iconRegistry[props.name || ''] || PhBellRinging)
const iconColor = computed(() => props.isRead ? '#94A3B8' : (props.color || '#64748B'))
</script>

<template>
  <component :is="icon" :size="size" :color="iconColor" weight="duotone" />
</template>
