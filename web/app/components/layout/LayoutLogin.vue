<script setup lang="ts">
import {PhChatCircleText, PhEnvelope, PhPaperPlaneTilt, PhStar} from "@phosphor-icons/vue";
import { DEFAULT_AVATAR } from '~/config'
</script>

<template>
  <div class="flex flex-col p-2 select-none gap-3">
    <div class="flex gap-3 items-center">
      <div>
        <img class="rounded-full w-12.5 h-12.5"
             :src="DEFAULT_AVATAR"
             alt="avatar" draggable="false"/>
      </div>
      <div class="flex flex-col gap-1">
        <span class="text-[#1E293B] font-bold text-[16px]">欢迎来到 插件市场</span>
        <span>登录解锁更多功能</span>
      </div>
    </div>
    <hr class="border-gray-200">
    <div class="flex gap-2 justify-center flex-wrap">
      <div class="flex gap-1 items-center bg-[#F2F6F9] py-2 px-4 rounded-lg">
        <PhStar color="#0F7B72"/>
        收藏插件
      </div>
      <div class="flex gap-1 items-center bg-[#F2F6F9] py-2 px-4 rounded-lg">
        <PhPaperPlaneTilt color="#0F7B72"/>
        发布插件
      </div>
      <div class="flex gap-1 items-center bg-[#F2F6F9] py-2 px-4 rounded-lg">
        <PhEnvelope color="#0F7B72"/>
        消息通知
      </div>
      <div class="flex gap-1 items-center bg-[#F2F6F9] py-2 px-4 rounded-lg">
        <PhChatCircleText color="#0F7B72"/>
        回复评论
      </div>
    </div>
    <div class="flex items-center gap-2 px-2 justify-center">
      <pm-button text="登录" class="w-full" @click="navigateTo('/user/auth/?mode=login')"/>
      <pm-button text="注册" color="green" class="w-full" @click="navigateTo('/user/auth/?mode=register')"/>
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>