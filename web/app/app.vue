<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'

const { settings, apply } = useSiteSettings()

// 通过 Nuxt 数据层在服务端读取公共设置，并把结果写入 SSR Payload。
// 客户端水合时直接复用 Payload，不再重复请求基础设置接口。
const { data: publicSiteSettings } = await useAsyncData(
  'public-site-settings',
  () => useApiFetch<{ data?: Record<string, string> }>('/setting/public'),
)
apply(publicSiteSettings.value?.data)

useSeoMeta({
  title: () => settings.value.siteTitle,
  keywords: () => settings.value.siteKeywords,
  description: () => settings.value.siteDescription,
})
</script>

<template>
  <el-config-provider :locale="zhCn">
    <NuxtLayout>
      <NuxtPage/>
    </NuxtLayout>
  </el-config-provider>
</template>
