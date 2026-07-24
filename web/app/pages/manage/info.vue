<script setup lang="ts">
import {
  PhCamera,
  PhChatCenteredText,
  PhEnvelopeSimple,
  PhIdentificationBadge,
  PhLock,
  PhUserFocus,
  PhUserList
} from '@phosphor-icons/vue'
import type { UploadFile, UploadRawFile } from 'element-plus'
import { DEFAULT_AVATAR, DEFAULT_USER_PROFILE } from '~/config'
import { completePasswordChange } from '~/utils/password-change'

interface ApiResponse<T = any> {
  code: number
  msg: string
  data?: T
}

interface UserInfoPayload {
  id: number
  username?: string
  nick?: string
  avatar?: string
  email?: string
  userdesc?: string
  power?: number
  titles?: any[]
  created?: string
  updated?: string
}

definePageMeta({
  layout: 'manage'
})

const userStore = useUserStore()
const assetUrl = useAssetUrl()

const isEditingInfo = ref(false)
const isEditingEmail = ref(false)
const isChangingPassword = ref(false)
const profileLoading = ref(false)
const infoSaving = ref(false)
const emailSaving = ref(false)
const avatarUploading = ref(false)
const passwordSaving = ref(false)
const avatarFile = ref<UploadRawFile | null>(null)
const avatarPreviewUrl = ref('')
const infoFormRef = ref()
const emailFormRef = ref()
const passwordFormRef = ref()

const infoForm = reactive({
  nick: '',
  userdesc: ''
})
const emailForm = reactive({
  email: ''
})
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const displayInfo = computed(() => userStore.userInfo)
const avatarUrl = computed(() => avatarPreviewUrl.value || assetUrl(displayInfo.value?.avatar, DEFAULT_AVATAR))
const userdescText = computed(() => displayInfo.value?.userdesc || DEFAULT_USER_PROFILE)

const infoRules = {
  nick: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 32, message: '昵称长度为 1-32 个字符', trigger: 'blur' }
  ],
  userdesc: [
    { max: 233, message: '简介最多 233 个字符', trigger: 'blur' }
  ]
}
const emailRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: ['blur', 'change'] }
  ]
}
const passwordRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度为 6-32 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== passwordForm.newPassword) callback(new Error('两次输入的新密码不一致'))
        else callback()
      },
      trigger: ['blur', 'change']
    }
  ]
}

const syncFormsFromStore = () => {
  const info = userStore.userInfo
  infoForm.nick = info?.nick || info?.username || ''
  infoForm.userdesc = info?.userdesc || ''
  emailForm.email = info?.email || ''
  avatarFile.value = null
  avatarPreviewUrl.value = ''
}

const loadUserInfo = async () => {
  profileLoading.value = true
  try {
    const res = await useApiFetch<ApiResponse<UserInfoPayload>>('/user/info')
    if (res.data) {
      userStore.setUserInfo(res.data)
      syncFormsFromStore()
    }
  } finally {
    profileLoading.value = false
  }
}

const startEditInfo = () => {
  syncFormsFromStore()
  isEditingInfo.value = true
  nextTick(() => infoFormRef.value?.clearValidate?.())
}

const cancelEditInfo = () => {
  syncFormsFromStore()
  isEditingInfo.value = false
}

const handleAvatarChange = async (file: UploadFile) => {
  if (!isEditingInfo.value || !file.raw) return
  if (!file.raw.type.startsWith('image/')) return ElMessage.warning('请选择图片文件')
  avatarFile.value = file.raw
  avatarPreviewUrl.value = URL.createObjectURL(file.raw)
}

const uploadAvatarIfNeeded = async () => {
  if (!avatarFile.value) return null
  avatarUploading.value = true
  try {
    const body = new FormData()
    body.append('file', avatarFile.value)
    const res = await useApiFetch<ApiResponse<{ avatar: string }>>('/user/avatar', { method: 'POST', body })
    return res.data?.avatar || null
  } finally {
    avatarUploading.value = false
  }
}

const saveInfo = async () => {
  const valid = await infoFormRef.value?.validate?.().catch(() => false)
  if (!valid) return

  infoSaving.value = true
  try {
    const avatar = await uploadAvatarIfNeeded()
    await useApiFetch('/user/info', {
      method: 'PUT',
      body: {
        nick: infoForm.nick.trim(),
        userdesc: infoForm.userdesc.trim()
      }
    })
    await loadUserInfo()
    if (avatar && userStore.userInfo) userStore.setUserInfo({ ...userStore.userInfo, avatar })
    isEditingInfo.value = false
    ElMessage.success('资料保存成功')
  } finally {
    infoSaving.value = false
  }
}

const startEditEmail = () => {
  emailForm.email = userStore.userInfo?.email || ''
  isEditingEmail.value = true
  nextTick(() => emailFormRef.value?.clearValidate?.())
}

const cancelEditEmail = () => {
  emailForm.email = userStore.userInfo?.email || ''
  isEditingEmail.value = false
}

const saveEmail = async () => {
  const valid = await emailFormRef.value?.validate?.().catch(() => false)
  if (!valid) return
  if (emailForm.email === userStore.userInfo?.email) {
    isEditingEmail.value = false
    return
  }

  emailSaving.value = true
  try {
    await useApiFetch('/user/email/send-verify', {
      method: 'POST',
      body: { email: emailForm.email.trim() }
    })
    isEditingEmail.value = false
    ElMessage.success('验证邮件已发送，请到邮箱完成绑定')
  } finally {
    emailSaving.value = false
  }
}

const openPasswordDialog = () => {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  isChangingPassword.value = true
  nextTick(() => passwordFormRef.value?.clearValidate?.())
}

const savePassword = async () => {
  const valid = await passwordFormRef.value?.validate?.().catch(() => false)
  if (!valid) return

  passwordSaving.value = true
  try {
    await useApiFetch('/user/password', {
      method: 'PUT',
      body: {
        oldPassword: passwordForm.oldPassword,
        newPassword: passwordForm.newPassword
      }
    })
    isChangingPassword.value = false
    await completePasswordChange({
      logout: userStore.logout,
      redirectToLogin: () => navigateTo({
        path: '/user/auth',
        query: {
          mode: 'login',
          redirect: '/manage/info',
          notice: 'password-changed'
        }
      })
    })
  } catch (error) {
    showApiError(error, '密码修改失败')
  } finally {
    passwordSaving.value = false
  }
}

onMounted(loadUserInfo)
</script>

<template>
  <div class="flex flex-col gap-5" v-loading="profileLoading">
    <manage-box title="我的信息">
      <template #header>
        <div v-if="!isEditingInfo">
          <pm-button text="编辑资料" @click="startEditInfo"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelEditInfo"/>
          <pm-button text="保存" :loading="infoSaving || avatarUploading" @click="saveInfo"/>
        </div>
      </template>
      <div class="profile-layout flex gap-10">
        <div class="shrink-0 text-center">
          <div>
            <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :disabled="!isEditingInfo"
                :auto-upload="false"
                accept="image/png,image/jpeg,image/gif,image/webp"
                :on-change="handleAvatarChange"
            >
              <div class="w-24 h-24 rounded-2xl overflow-hidden border-2 border-[#F1F5F9] relative group bg-[#F8FAFC]">
                <img :src="avatarUrl" class="w-full h-full object-cover" alt="头像"/>
                <div v-if="isEditingInfo"
                     class="absolute inset-0 bg-black/40 flex items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-opacity">
                  <PhCamera :size="24"/>
                </div>
              </div>
            </el-upload>
          </div>
          <span class="text-xs text-[#94A3B8]" v-if="isEditingInfo">点击更换头像</span>
        </div>
        <el-form ref="infoFormRef" :model="infoForm" :rules="infoRules" label-position="top" class="flex-1 w-full">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1">
            <el-form-item>
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhIdentificationBadge :size="18"/>
                  UID
                </div>
              </template>
              <div
                  class="w-full h-12 flex items-center bg-gray-100 rounded-xl px-3.5 text-[#94A3B8] select-text text-[14px]">
                {{ displayInfo?.id || '-' }}
              </div>
            </el-form-item>
            <el-form-item>
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhUserList :size="18"/>
                  账号
                </div>
              </template>
              <div
                  class="w-full h-12 flex items-center bg-gray-100 rounded-xl px-3.5 text-[#94A3B8] select-text text-[14px]">
                {{ displayInfo?.username || '-' }}
              </div>
            </el-form-item>
            <el-form-item prop="nick">
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhUserFocus :size="18"/>
                  昵称
                </div>
              </template>
              <div v-if="!isEditingInfo"
                   class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
                {{ displayInfo?.nick || displayInfo?.username || '-' }}
              </div>
              <el-input v-model="infoForm.nick" v-else placeholder="请输入昵称" class="pm-manage-input" maxlength="32"/>
            </el-form-item>
            <el-form-item class="md:col-span-2" prop="userdesc">
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhChatCenteredText :size="18"/>
                  简介
                </div>
              </template>
              <div v-if="!isEditingInfo"
                   class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 leading-normal py-2.5 text-[#606266] select-text text-[14px] whitespace-pre-wrap break-words overflow-y-auto">
                {{ userdescText }}
              </div>
              <el-input
                  v-else
                  type="textarea"
                  v-model="infoForm.userdesc"
                  placeholder="介绍一下自己吧..."
                  class="pm-manage-input"
                  maxlength="233"
                  show-word-limit
              />
            </el-form-item>
          </div>
        </el-form>
      </div>
    </manage-box>
    <manage-box title="消息通知">
      <template #header>
        <div v-if="!isEditingEmail">
          <pm-button text="编辑邮箱" @click="startEditEmail"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelEditEmail"/>
          <pm-button text="发送验证" :loading="emailSaving" @click="saveEmail"/>
        </div>
      </template>
      <div>
        <el-form ref="emailFormRef" :model="emailForm" :rules="emailRules" label-position="top" class="max-w-md">
          <el-form-item class="mb-0" prop="email">
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                <PhEnvelopeSimple :size="18"/>
                邮箱地址
              </div>
            </template>
            <div v-if="!isEditingEmail"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ displayInfo?.email || '暂未绑定邮箱' }}
            </div>
            <el-input v-else v-model="emailForm.email" placeholder="请输入邮箱地址" class="pm-manage-input"/>
            <p class="text-xs text-[#94A3B8] mt-2 ml-1 leading-normal w-full">用于接收插件审核通知、站内信提醒等。修改邮箱需要点击验证邮件后才会生效。</p>
          </el-form-item>
        </el-form>
      </div>
    </manage-box>
    <manage-box title="其他">
      <div class="flex flex-wrap gap-4">
        <pm-button text="修改密码" color="orange" @click="openPasswordDialog"/>
      </div>
    </manage-box>

    <el-dialog
        v-model="isChangingPassword"
        title="修改密码"
        width="460px"
        align-center
        class="pm-manage-dialog"
    >
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-position="top">
        <el-form-item prop="oldPassword">
          <template #label>
            <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
              <PhLock :size="18"/>
              旧密码
            </div>
          </template>
          <el-input
              v-model="passwordForm.oldPassword"
              type="password"
              show-password
              placeholder="请输入当前密码"
              class="pm-manage-input"
          />
        </el-form-item>
        <el-form-item prop="newPassword">
          <template #label>
            <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
              <PhLock :size="18"/>
              新密码
            </div>
          </template>
          <el-input
              v-model="passwordForm.newPassword"
              type="password"
              show-password
              placeholder="请输入新密码"
              class="pm-manage-input"
          />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <template #label>
            <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
              <PhLock :size="18"/>
              确认新密码
            </div>
          </template>
          <el-input
              v-model="passwordForm.confirmPassword"
              type="password"
              show-password
              placeholder="请再次输入新密码"
              class="pm-manage-input"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3 pb-2">
          <pm-button text="取消" color="white" @click="isChangingPassword = false"/>
          <pm-button text="确认修改" :loading="passwordSaving" @click="savePassword"/>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
/* 隐藏必填星号，保留表单校验 */
:deep(.el-form-item.is-required > .el-form-item__label::before) {
  display: none;
}

/* 表单输入框深度定制 */
:deep(.pm-manage-input) {
  width: 100%;

  .el-input__wrapper, .el-textarea__inner {
    font-size: 14px !important;
    background-color: #F8FAFC;
    border: none;
    box-shadow: inset 0 0 0 1.5px transparent;
    border-radius: 12px;
    padding: 8px 14px;
    transition: background-color 0.2s, border-color 0.2s, box-shadow 0.2s;

    &:hover {
      box-shadow: inset 0 0 0 1.5px #E2E8F0;
    }

    &.is-focus, &:focus {
      box-shadow: inset 0 0 0 1.5px #00BAAD !important;
    }
  }

  .el-input__wrapper {
    height: 48px !important;
  }
  .el-textarea__inner {
    padding: 10px 14px !important;
    height: 120px;
    resize: none;
  }
}

.avatar-uploader {
  :deep(.el-upload) {
    cursor: pointer;
    border: none;
    background: none;

    &:hover {
      border: none;
    }
  }
}

/* 对话框深度定制 */
:deep(.pm-manage-dialog) {
  border-radius: 20px !important;
  padding: 10px;
  overflow: hidden;

  .el-dialog__header {
    margin-right: 0;
    padding: 20px 24px 10px;

    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #1E293B;
    }
  }

  .el-dialog__body {
    padding: 10px 24px;
  }

  .el-dialog__footer {
    padding: 10px 24px 20px;
  }

  .el-dialog__headerbtn {
    top: 20px;
    right: 20px;

    &:hover .el-dialog__close {
      color: #00BAAD;
    }
  }
}

@media (max-width: 767px) {
  .profile-layout {
    flex-direction: column;
    gap: 20px;
    align-items: center;

    > :deep(.el-form) {
      min-width: 0;
    }
  }
}

</style>
