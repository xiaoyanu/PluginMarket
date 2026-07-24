<script setup lang="ts">
import { PhLock, PhCheckCircle } from '@phosphor-icons/vue'

const route = useRoute()
const router = useRouter()
const token = computed(() => typeof route.query.token === 'string' ? route.query.token : '')
const resetting = ref(false)

const resetForm = ref({
  password: '',
  confirmPassword: ''
})

const handleReset = async () => {
  if (!token.value) {
    ElMessage.error('重置链接无效')
    return
  }
  if (!resetForm.value.password || !resetForm.value.confirmPassword) {
    ElMessage.warning('请输入新密码并确认')
    return
  }
  if (resetForm.value.password.length < 6) {
    ElMessage.warning('新密码至少需要6位')
    return
  }
  if (resetForm.value.password !== resetForm.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  resetting.value = true
  try {
    await useApiFetch('/auth/reset-password', {
      method: 'POST',
      body: { token: token.value, password: resetForm.value.password },
    })
    ElMessage.success('密码重置成功')
    await router.push('/user/auth?mode=login')
  } finally {
    resetting.value = false
  }
}

onMounted(() => {
  if (!token.value) navigateTo('/user/auth?mode=login')
})
</script>

<template>
  <div class="min-h-[calc(100vh-120px)] flex items-center justify-center p-4 select-none">
    <div class="w-full max-w-110 bg-white rounded-3xl shadow-pmbox p-8 md:p-10">
      <div class="text-center mb-10">
        <h1 class="text-[28px] font-bold text-[#1e293b] mb-2">重置密码</h1>
        <p class="text-[#64748b] text-sm">请设置您的新密码并妥善保管</p>
      </div>

      <el-form :model="resetForm" label-position="top" @submit.prevent="handleReset">
        <el-form-item label="新密码">
          <el-input v-model="resetForm.password" type="password" placeholder="请输入新密码" show-password class="pm-input">
            <template #prefix><PhLock :size="20"/></template>
          </el-input>
        </el-form-item>

        <el-form-item label="确认新密码">
          <el-input v-model="resetForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password class="pm-input">
            <template #prefix><PhLock :size="20"/></template>
          </el-input>
        </el-form-item>

        <pm-button text="确认修改" color="blue" :loading="resetting" class="w-full h-12 text-base mt-4" @click="handleReset">
          <template #icon><PhCheckCircle :size="18" weight="bold"/></template>
        </pm-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
:deep(.el-form-item) {
  margin-bottom: 22px;
  .el-form-item__label {
    padding-bottom: 6px;
    font-size: 14px;
    font-weight: 500;
    color: #475569;
  }
}

:deep(.pm-input) {
  .el-input__wrapper {
    background-color: #F1F5F9;
    box-shadow: none;
    border: 1.5px solid transparent;
    border-radius: 12px;
    height: 44px;
    padding: 0 16px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &.is-focus {
      background-color: #fff;
      border-color: #60a5fa !important;
      box-shadow: 0 0 0 4px rgba(96, 165, 250, 0.1) !important;
    }

    &:hover:not(.is-focus) {
      border-color: #cbd5e1;
    }
  }
}
</style>
