<script setup lang="ts">
import {PhUser, PhLock, PhShieldCheck, PhSignIn, PhUserCirclePlus} from '@phosphor-icons/vue'

const route = useRoute()
const router = useRouter()
const mode = ref<'login' | 'register'>((route.query.mode as 'login' | 'register') || 'login')
const userStore = useUserStore()
const loginNoticeShown = ref(false)
const passwordChangedNoticeShown = ref(false)
const redirectTarget = computed(() => {
  const redirect = route.query.redirect
  return typeof redirect === 'string' && redirect.startsWith('/') ? redirect : '/'
})

type ApiResponse<T = any> = { code: number; msg: string; data?: T }
type CaptchaData = { captchaId: string; image: string }
type LoginData = { token: string; user: { id: number; power: number } }

const authForm = ref({ username: '', password: '', confirmPassword: '', captcha: '', captchaId: '' })
const captchaImage = ref('')
const loading = ref(false)
const captchaLoading = ref(false)

const forgotPasswordVisible = ref(false)
const forgotForm = ref({ username: '', captcha: '', captchaId: '' })
const forgotCaptchaImage = ref('')
const forgotLoading = ref(false)

const fetchCaptcha = async (target: 'auth' | 'forgot' = 'auth') => {
  captchaLoading.value = target === 'auth'
  try {
    const res = await useApiFetch<ApiResponse<CaptchaData>>('/auth/captcha')
    const data = res.data
    if (!data?.captchaId || !data?.image) throw new Error('验证码返回数据异常')
    if (target === 'auth') {
      authForm.value.captchaId = data.captchaId
      authForm.value.captcha = ''
      captchaImage.value = data.image
    } else {
      forgotForm.value.captchaId = data.captchaId
      forgotForm.value.captcha = ''
      forgotCaptchaImage.value = data.image
    }
  } catch (error: any) {
    showApiError(error, '验证码加载失败')
  } finally {
    captchaLoading.value = false
  }
}

const handleForgotPassword = async () => {
  forgotPasswordVisible.value = true
  await nextTick()
  await fetchCaptcha('forgot')
}

const sendResetLink = async () => {
  if (!forgotForm.value.username.trim()) return ElMessage.warning('请输入账号')
  if (!forgotForm.value.captcha.trim()) return ElMessage.warning('请输入验证码')
  forgotLoading.value = true
  try {
    const res = await useApiFetch<ApiResponse>('/auth/forgot-password', {
      method: 'POST',
      body: { username: forgotForm.value.username.trim(), captcha: forgotForm.value.captcha.trim(), captchaId: forgotForm.value.captchaId }
    })
    ElMessage.success(res.msg || '重置邮件已发送')
    forgotPasswordVisible.value = false
    forgotForm.value = { username: '', captcha: '', captchaId: '' }
  } catch (error: any) {
    showApiError(error, '发送失败')
    await fetchCaptcha('forgot')
  } finally {
    forgotLoading.value = false
  }
}

const toggleMode = (newMode: 'login' | 'register') => {
  mode.value = newMode
  router.replace({query: {...route.query, mode: newMode}})
  authForm.value = { username: '', password: '', confirmPassword: '', captcha: '', captchaId: authForm.value.captchaId }
}

const submitAuth = async () => {
  const username = authForm.value.username.trim()
  const password = authForm.value.password
  const captcha = authForm.value.captcha.trim()
  if (!username) return ElMessage.warning('请输入账号')
  if (!password) return ElMessage.warning('请输入密码')
  if (mode.value === 'register' && username.length < 4) return ElMessage.warning('账号至少需要4位')
  if (mode.value === 'register' && password.length < 6) return ElMessage.warning('密码至少需要6位')
  if (mode.value === 'register' && password !== authForm.value.confirmPassword) return ElMessage.warning('两次输入的密码不一致')
  if (!captcha) return ElMessage.warning('请输入验证码')
  loading.value = true
  try {
    if (mode.value === 'login') {
      const res = await useApiFetch<ApiResponse<LoginData>>('/auth/login', {
        method: 'POST',
        body: { username, password, captcha, captchaId: authForm.value.captchaId }
      })
      const data = res.data
      if (!data?.token) throw new Error('登录接口未返回 token')
      userStore.login(data.token, data.user)
      try {
        const infoRes = await useApiFetch<ApiResponse>('/user/info', { suppressErrorMessage: true })
        userStore.setUserInfo(infoRes.data || data.user)
      } catch {
        userStore.setUserInfo(data.user)
      }
      ElMessage.success('登录成功')
      await navigateTo(redirectTarget.value)
    } else {
      const res = await useApiFetch<ApiResponse>('/auth/register', {
        method: 'POST',
        body: { username, password, captcha, captchaId: authForm.value.captchaId }
      })
      ElMessage.success(res.msg || '注册成功，请登录')
      toggleMode('login')
    }
  } catch (error: any) {
    showApiError(error, mode.value === 'login' ? '登录失败' : '注册失败')
  } finally {
    loading.value = false
    await fetchCaptcha('auth')
  }
}

watch(() => route.query.mode, (newMode) => {
  if (newMode === 'login' || newMode === 'register') mode.value = newMode
})

watch(() => route.query.notice, (notice) => {
  if (!import.meta.client) return
  if (notice === 'login-required' && !loginNoticeShown.value) {
    loginNoticeShown.value = true
    ElMessage.warning('请先登录后访问')
  }
  if (notice === 'password-changed' && !passwordChangedNoticeShown.value) {
    passwordChangedNoticeShown.value = true
    ElMessage.success('密码修改成功，请重新登录')
  }
}, { immediate: true })

onMounted(() => fetchCaptcha())
</script>

<template>
  <div class="auth-page min-h-[calc(100vh-120px)] flex items-start sm:items-center justify-center px-3 pb-6 sm:p-4 select-none">
    <div class="auth-card w-full max-w-110 bg-white rounded-[16px] sm:rounded-3xl shadow-pmbox p-5 py-6 sm:p-8 md:p-10 transition-all">
      <div class="text-center mb-7 sm:mb-10">
        <h1 class="text-[24px] sm:text-[28px] font-bold text-[#1e293b] mb-2">{{ mode === 'login' ? '欢迎回来' : '创建账号' }}</h1>
        <p class="text-[#64748b] text-sm">{{ mode === 'login' ? '输入账号密码以登录' : '加入插件市场，发现更多可能' }}</p>
      </div>

      <el-form :model="authForm" label-position="top" hide-required-asterisk @submit.prevent="submitAuth">
        <el-form-item label="账号">
          <el-input v-model="authForm.username" :placeholder="mode === 'login' ? '请输入账号' : '仅支持字母或数字'" class="pm-input" clearable>
            <template #prefix><PhUser :size="20"/></template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <template #label>
            <div class="flex justify-between items-center w-full pr-1">
              <span class="text-sm font-medium text-[#475569]">密码</span>
              <span v-if="mode === 'login'" class="text-xs text-blue-500 hover:text-blue-600 cursor-pointer font-normal" @click="handleForgotPassword">忘记密码？</span>
            </div>
          </template>
          <el-input v-model="authForm.password" type="password" placeholder="请输入密码" show-password class="pm-input" clearable @keydown.enter="submitAuth">
            <template #prefix><PhLock :size="20"/></template>
          </el-input>
        </el-form-item>

        <transition name="fade-slide">
          <el-form-item v-if="mode === 'register'" label="确认密码">
            <el-input v-model="authForm.confirmPassword" type="password" placeholder="请再次输入密码" show-password class="pm-input" clearable @keydown.enter="submitAuth">
              <template #prefix><PhLock :size="20"/></template>
            </el-input>
          </el-form-item>
        </transition>

        <el-form-item label="验证码">
          <div class="flex gap-2 sm:gap-3 w-full min-w-0">
            <el-input v-model="authForm.captcha" placeholder="验证码" maxlength="4" class="pm-input min-w-0 flex-1" @keydown.enter="submitAuth">
              <template #prefix><PhShieldCheck :size="20"/></template>
            </el-input>
            <div v-loading="captchaLoading" class="w-26 sm:w-30 h-11 bg-[#F1F5F9] rounded-xl overflow-hidden cursor-pointer border border-transparent hover:border-[#cbd5e1] transition-all active:scale-95 shrink-0" @click="fetchCaptcha('auth')">
              <img v-if="captchaImage" :src="captchaImage" alt="captcha" class="w-full h-full object-contain bg-white"/>
            </div>
          </div>
        </el-form-item>

        <pm-button :text="mode === 'login' ? '登录' : '注册'" :color="mode === 'login' ? 'blue' : 'green'" :loading="loading" :class="['w-full h-12 text-base mt-4']" @click="submitAuth">
          <template #icon>
            <PhSignIn v-if="mode === 'login'" :size="18" weight="bold"/>
            <PhUserCirclePlus v-else :size="18" weight="bold"/>
          </template>
        </pm-button>

        <p class="text-center text-sm text-[#64748b] mt-6">
          {{ mode === 'login' ? '还没有账号？' : '已有账号？' }}
          <span class="text-blue-500 font-medium hover:underline cursor-pointer" @click="toggleMode(mode === 'login' ? 'register' : 'login')">
            {{ mode === 'login' ? '立即注册' : '返回登录' }}
          </span>
        </p>
      </el-form>
    </div>

    <el-dialog v-model="forgotPasswordVisible" title="找回密码" width="min(440px, calc(100vw - 24px))" class="pm-manage-dialog" align-center>
      <el-form :model="forgotForm" label-position="top">
        <el-form-item label="账号">
          <el-input v-model="forgotForm.username" placeholder="请输入绑定的账号" class="pm-input"><template #prefix><PhUser :size="20"/></template></el-input>
        </el-form-item>
        <el-form-item label="验证码">
          <div class="flex gap-2 sm:gap-3 w-full min-w-0">
            <el-input v-model="forgotForm.captcha" placeholder="验证码" maxlength="4" class="pm-input min-w-0 flex-1"><template #prefix><PhShieldCheck :size="20"/></template></el-input>
            <div class="w-26 sm:w-30 h-11 bg-[#F1F5F9] rounded-xl overflow-hidden cursor-pointer border border-transparent hover:border-[#cbd5e1] transition-all active:scale-95 shrink-0" @click="fetchCaptcha('forgot')">
              <img v-if="forgotCaptchaImage" :src="forgotCaptchaImage" alt="captcha" class="w-full h-full object-contain bg-white"/>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3 pb-2">
          <pm-button text="取消" color="white" @click="forgotPasswordVisible = false" />
          <pm-button text="发送重置链接" color="blue" :loading="forgotLoading" @click="sendResetLink" />
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
.fade-slide-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-12px);
  max-height: 0;
  margin-bottom: 0 !important;
  overflow: hidden;
}

.fade-slide-enter-to,
.fade-slide-leave-from {
  max-height: 90px; // 足够容纳一个 form-item 的高度
  opacity: 1;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

  .el-form-item__label {
    padding-bottom: 6px;
    font-size: 14px;
    font-weight: 500;
    color: #475569;
    line-height: 1;
    margin-left: 4px;
    display: block;
    width: 100%;
  }

  .el-form-item__error {
    padding-top: 4px;
    margin-left: 4px;
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

      .el-input__prefix-inner {
        color: #3b82f6;
      }
    }

    &:hover:not(.is-focus) {
      border-color: #cbd5e1;
    }
  }

  .el-input__inner {
    font-size: 14px;
    color: #1e293b;

    &::placeholder {
      color: #94a3b8;
    }
  }

  .el-input__prefix-inner {
    color: #64748b;
    margin-right: 8px;
    transition: color 0.3s;
  }
}

/* 即使去掉了校验逻辑，依旧保留这段关于报错的 UI 样式控制以备后续需要 */
:deep(.el-form-item) {
  &.is-error {
    .pm-input .el-input__wrapper {
      border-color: #f87171 !important;
      background-color: #fef2f2;
      box-shadow: none !important;

      .el-input__prefix-inner {
        color: #ef4444;
      }
    }
  }
}

/* 对话框深度定制 - 匹配 pm-manage-dialog 规范 */
:deep(.pm-manage-dialog) {
  border-radius: 24px !important;
  padding: 8px;

  .el-dialog__header {
    margin-right: 0;
    padding: 24px 32px 12px;

    .el-dialog__title {
      font-size: 20px;
      font-weight: 700;
      color: #0F172A;
    }
  }

  .el-dialog__body {
    padding: 12px 32px;
  }

  .el-dialog__footer {
    padding: 10px 24px 20px;
  }

  .el-dialog__headerbtn {
    top: 24px;
    right: 24px;
  }
}

@media (max-width: 639px) {
  .auth-page {
    min-height: calc(100vh - 132px);
  }

  :deep(.el-form-item) {
    margin-bottom: 18px;
  }

  :deep(.pm-input) {
    .el-input__wrapper {
      padding: 0 12px;
    }
  }

  :deep(.pm-manage-dialog) {
    border-radius: 16px !important;
    padding: 4px;

    .el-dialog__header {
      padding: 20px 20px 10px;
    }

    .el-dialog__body {
      padding: 10px 20px;
    }

    .el-dialog__footer {
      padding: 8px 16px 16px;

      > div {
        display: grid;
        grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
        gap: 8px;
      }

      .pm-button {
        width: 100%;
      }
    }

    .el-dialog__headerbtn {
      top: 18px;
      right: 16px;
    }
  }
}
</style>
