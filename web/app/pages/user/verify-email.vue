<script setup lang="ts">
import { PhCheckCircle, PhWarningCircle } from '@phosphor-icons/vue'

const route = useRoute()
const status = ref<'loading' | 'success' | 'error'>('loading')
const message = ref('正在验证邮箱，请稍候…')

onMounted(async () => {
  const token = typeof route.query.token === 'string' ? route.query.token : ''
  if (!token) {
    status.value = 'error'
    message.value = '验证链接无效'
    return
  }
  try {
    const response: any = await useApiFetch('/user/email/verify', {
      query: { token },
      suppressErrorMessage: true,
    })
    status.value = 'success'
    message.value = response?.msg || '邮箱绑定成功'
  } catch (error) {
    status.value = 'error'
    message.value = getApiErrorMessage(error, '邮箱验证失败')
  }
})
</script>

<template>
  <div class="min-h-[calc(100vh-120px)] flex items-center justify-center p-4">
    <div class="w-full max-w-110 bg-white rounded-3xl shadow-pmbox p-10 text-center">
      <div v-if="status === 'loading'" v-loading="true" class="h-20"></div>
      <PhCheckCircle v-else-if="status === 'success'" :size="64" class="mx-auto text-green-500 mb-5" weight="fill"/>
      <PhWarningCircle v-else :size="64" class="mx-auto text-red-500 mb-5" weight="fill"/>
      <h1 class="text-2xl font-bold text-[#1e293b] mb-3">{{ status === 'success' ? '验证成功' : status === 'error' ? '验证失败' : '正在验证' }}</h1>
      <p class="text-sm text-[#64748b] mb-7">{{ message }}</p>
      <pm-button v-if="status !== 'loading'" text="返回个人中心" @click="navigateTo('/manage/info')"/>
    </div>
  </div>
</template>
