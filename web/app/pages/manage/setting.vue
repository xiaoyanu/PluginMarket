<script setup lang="ts">
import {PhBroom, PhCamera, PhEye, PhFilePlus, PhHash, PhNote, PhSelection, PhTextT} from "@phosphor-icons/vue";
import type { UploadRequestOptions } from 'element-plus'
import {
  DEFAULT_ALLOWED_IMAGE_EXTENSIONS,
  DEFAULT_ALLOW_COMMENT,
  DEFAULT_ALLOW_REGISTER,
  DEFAULT_ALLOW_UPLOAD,
  DEFAULT_ANTI_BRUSH_PLUGIN_DATA,
  DEFAULT_EMAIL_TEMPLATES,
  DEFAULT_MAX_IMAGE_SIZE,
  DEFAULT_NOTIFICATION_TEMPLATES,
  DEFAULT_SITE_LOGO,
  DEFAULT_SKIP_AUDIT,
} from '~/config'

definePageMeta({
  layout: 'manage'
})

const { settings: siteSettings, load: loadPublicSiteSettings } = useSiteSettings()
const assetUrl = useAssetUrl()
const basicSaving = ref(false)
const logoClearing = ref(false)
const functionSaving = ref(false)
const uploadSaving = ref(false)
const cleanupRunning = ref(false)
const emailSaving = ref(false)
const notificationSaving = ref(false)
const notificationSnapshot = ref<Record<string, string>>({})
const testEmailVisible = ref(false)
const testEmailAddress = ref('')
const testEmailSending = ref(false)
const emailSnapshot = ref<Record<string, string>>({})
const basicSnapshot = ref({
  site_logo: '',
  site_title: '',
  site_description: '',
  site_keywords: '',
})
const uploadSnapshot = ref({
  image_max_size: DEFAULT_MAX_IMAGE_SIZE,
  image_allowed_ext: DEFAULT_ALLOWED_IMAGE_EXTENSIONS,
})

// 预览模板变量示例数据
const previewExamples: Record<string, string> = {
  site_title: 'Plugin Market',
  plugin_name: 'Vue DevTools Pro',
  comment_user: '张三',
  comment_content: '这个插件非常好用，大大提高了我的开发效率！',
  reply_user: '李四',
  reply_content: '谢谢支持！我们还会持续更新优化。',
  content: '插件介绍不完整，请补充详细的安装和使用说明。',
  user_name: '王五',
  link: 'https://example.com/plugin/123'
}

// 预览浮层状态
const previewState = reactive({
  visible: false,
  title: '',
  content: '',
  x: 0,
  y: 0
})
let hideTimer: ReturnType<typeof setTimeout> | null = null

// 替换变量
const replaceVars = (template: string): string => {
  let result = template
  for (const [k, value] of Object.entries(previewExamples)) {
    result = result.replace(new RegExp(`\\{\\{${k}\\}\\}`, 'g'), value)
  }
  return result
}

// 显示预览
const showPreview = (e: MouseEvent, titleKey: string, bodyKey: string) => {
  if (hideTimer) {
    clearTimeout(hideTimer)
    hideTimer = null
  }
  const target = e.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  previewState.title = replaceVars(settings.email[titleKey as keyof typeof settings.email] as string)
  previewState.content = replaceVars(settings.email[bodyKey as keyof typeof settings.email] as string)
  previewState.x = rect.right + 12
  previewState.y = rect.top
  previewState.visible = true
}

// 隐藏预览
const hidePreview = () => {
  hideTimer = setTimeout(() => {
    previewState.visible = false
  }, 100)
}

// 鼠标进入预览DIV
const onPreviewEnter = () => {
  if (hideTimer) {
    clearTimeout(hideTimer)
    hideTimer = null
  }
}

// 设置数据
const settings = reactive({
  basic: {
    site_logo: siteSettings.value.siteLogo,
    site_title: siteSettings.value.siteTitle,
    site_description: siteSettings.value.siteDescription,
    site_keywords: siteSettings.value.siteKeywords,
    isEditing: false
  },
  function: {
    allow_register: DEFAULT_ALLOW_REGISTER,
    allow_comment: DEFAULT_ALLOW_COMMENT,
    plugin_no_audit: DEFAULT_SKIP_AUDIT,
    allow_upload_image: DEFAULT_ALLOW_UPLOAD,
    plugin_anti_brush: DEFAULT_ANTI_BRUSH_PLUGIN_DATA,
  },
  upload: {
    image_max_size: DEFAULT_MAX_IMAGE_SIZE,
    image_allowed_ext: DEFAULT_ALLOWED_IMAGE_EXTENSIONS,
    isEditing: false
  },
  notification: {
    new_comment_title: DEFAULT_NOTIFICATION_TEMPLATES.newComment.title,
    new_comment_content: DEFAULT_NOTIFICATION_TEMPLATES.newComment.content,
    reply_comment_title: DEFAULT_NOTIFICATION_TEMPLATES.replyComment.title,
    reply_comment_content: DEFAULT_NOTIFICATION_TEMPLATES.replyComment.content,
    plugin_approved_title: DEFAULT_NOTIFICATION_TEMPLATES.pluginApproved.title,
    plugin_approved_content: DEFAULT_NOTIFICATION_TEMPLATES.pluginApproved.content,
    plugin_rejected_title: DEFAULT_NOTIFICATION_TEMPLATES.pluginRejected.title,
    plugin_rejected_content: DEFAULT_NOTIFICATION_TEMPLATES.pluginRejected.content,
    isEditing: false,
  },
  email: {
    smtp_server: '',
    smtp_port: '465',
    smtp_user: '',
    smtp_pass: '',
    smtp_from: '',
    smtp_from_name: '',
    tpl_new_comment_title: DEFAULT_EMAIL_TEMPLATES.newComment.title,
    tpl_new_comment_body: DEFAULT_EMAIL_TEMPLATES.newComment.body,
    tpl_reply_comment_title: DEFAULT_EMAIL_TEMPLATES.replyComment.title,
    tpl_reply_comment_body: DEFAULT_EMAIL_TEMPLATES.replyComment.body,
    tpl_pending_plugin_review_title: DEFAULT_EMAIL_TEMPLATES.pendingPluginReview.title,
    tpl_pending_plugin_review_body: DEFAULT_EMAIL_TEMPLATES.pendingPluginReview.body,
    tpl_plugin_approved_title: DEFAULT_EMAIL_TEMPLATES.pluginApproved.title,
    tpl_plugin_approved_body: DEFAULT_EMAIL_TEMPLATES.pluginApproved.body,
    tpl_plugin_rejected_title: DEFAULT_EMAIL_TEMPLATES.pluginRejected.title,
    tpl_plugin_rejected_body: DEFAULT_EMAIL_TEMPLATES.pluginRejected.body,
    tpl_bind_email_title: DEFAULT_EMAIL_TEMPLATES.emailVerify.title,
    tpl_bind_email_body: DEFAULT_EMAIL_TEMPLATES.emailVerify.body,
    tpl_reset_pwd_title: DEFAULT_EMAIL_TEMPLATES.resetPassword.title,
    tpl_reset_pwd_body: DEFAULT_EMAIL_TEMPLATES.resetPassword.body,
    isEditing: false
  }
});

const notificationSettingMap = {
  notificationNewCommentTitle: 'new_comment_title',
  notificationNewCommentContent: 'new_comment_content',
  notificationReplyCommentTitle: 'reply_comment_title',
  notificationReplyCommentContent: 'reply_comment_content',
  notificationPluginApprovedTitle: 'plugin_approved_title',
  notificationPluginApprovedContent: 'plugin_approved_content',
  notificationPluginRejectedTitle: 'plugin_rejected_title',
  notificationPluginRejectedContent: 'plugin_rejected_content',
} as const

const notificationDefaults: Record<string, string> = {
  notificationNewCommentTitle: DEFAULT_NOTIFICATION_TEMPLATES.newComment.title,
  notificationNewCommentContent: DEFAULT_NOTIFICATION_TEMPLATES.newComment.content,
  notificationReplyCommentTitle: DEFAULT_NOTIFICATION_TEMPLATES.replyComment.title,
  notificationReplyCommentContent: DEFAULT_NOTIFICATION_TEMPLATES.replyComment.content,
  notificationPluginApprovedTitle: DEFAULT_NOTIFICATION_TEMPLATES.pluginApproved.title,
  notificationPluginApprovedContent: DEFAULT_NOTIFICATION_TEMPLATES.pluginApproved.content,
  notificationPluginRejectedTitle: DEFAULT_NOTIFICATION_TEMPLATES.pluginRejected.title,
  notificationPluginRejectedContent: DEFAULT_NOTIFICATION_TEMPLATES.pluginRejected.content,
}

const notificationTemplates = [
  { name: '收到评论模板', title: 'new_comment_title', content: 'new_comment_content', variables: [['user_name', '用户昵称'], ['content', '评论内容'], ['plugin_name', '插件名称']] },
  { name: '评论被回复模板', title: 'reply_comment_title', content: 'reply_comment_content', variables: [['user_name', '回复人昵称'], ['content', '评论内容'], ['plugin_name', '插件名称']] },
  { name: '插件审核通过模板', title: 'plugin_approved_title', content: 'plugin_approved_content', variables: [['plugin_name', '插件名称']] },
  { name: '插件审核拒绝模板', title: 'plugin_rejected_title', content: 'plugin_rejected_content', variables: [['content', '拒绝理由'], ['plugin_name', '插件名称']] },
] as const

const notificationPayload = () => Object.fromEntries(
  Object.entries(notificationSettingMap).map(([backendKey, frontendKey]) => [backendKey, String(settings.notification[frontendKey])]),
)

const applyNotificationSettings = async (data: Record<string, string>, persistMissing = true) => {
  const missing: Record<string, string> = {}
  for (const [backendKey, frontendKey] of Object.entries(notificationSettingMap)) {
    settings.notification[frontendKey] = data[backendKey] || notificationDefaults[backendKey] || ''
    if (!data[backendKey]) missing[backendKey] = notificationDefaults[backendKey] || ''
  }
  notificationSnapshot.value = notificationPayload()
  if (persistMissing && Object.keys(missing).length) {
    try {
      await useApiFetch('/admin/setting', { method: 'PUT', body: missing })
      Object.assign(notificationSnapshot.value, missing)
    } catch {
      ElMessage.warning('默认通知模板初始化失败，请稍后重试')
    }
  }
}

const copyBasicSettings = () => ({
  site_logo: settings.basic.site_logo,
  site_title: settings.basic.site_title,
  site_description: settings.basic.site_description,
  site_keywords: settings.basic.site_keywords,
})

const applyBasicSettings = (data: Record<string, string>) => {
  settings.basic.site_logo = data.siteLogo || siteSettings.value.siteLogo
  settings.basic.site_title = data.siteTitle || siteSettings.value.siteTitle
  settings.basic.site_description = data.siteDesc || siteSettings.value.siteDescription
  settings.basic.site_keywords = data.siteKeywords || siteSettings.value.siteKeywords
  basicSnapshot.value = copyBasicSettings()
}

const parseBooleanSetting = (value: string | undefined, fallback: boolean) => {
  if (value === 'true') return true
  if (value === 'false') return false
  return fallback
}

const applyFunctionSettings = (data: Record<string, string>) => {
  settings.function.allow_register = parseBooleanSetting(data.allowRegister, DEFAULT_ALLOW_REGISTER)
  settings.function.allow_comment = parseBooleanSetting(data.allowComment, DEFAULT_ALLOW_COMMENT)
  settings.function.plugin_no_audit = parseBooleanSetting(data.skipAudit, DEFAULT_SKIP_AUDIT)
  settings.function.allow_upload_image = parseBooleanSetting(data.allowUpload, DEFAULT_ALLOW_UPLOAD)
  settings.function.plugin_anti_brush = parseBooleanSetting(data.antiBrushPluginData, DEFAULT_ANTI_BRUSH_PLUGIN_DATA)
}

const copyUploadSettings = () => ({
  image_max_size: settings.upload.image_max_size,
  image_allowed_ext: settings.upload.image_allowed_ext,
})

const applyUploadSettings = (data: Record<string, string>) => {
  const maxSize = Number(data.maxFileSize)
  settings.upload.image_max_size = Number.isFinite(maxSize) && maxSize > 0
    ? maxSize
    : DEFAULT_MAX_IMAGE_SIZE
  settings.upload.image_allowed_ext = data.allowedExtensions || DEFAULT_ALLOWED_IMAGE_EXTENSIONS
  uploadSnapshot.value = copyUploadSettings()
}

const emailSettingMap = {
  emailHost: 'smtp_server',
  emailPort: 'smtp_port',
  emailUser: 'smtp_user',
  emailPass: 'smtp_pass',
  emailFrom: 'smtp_from',
  emailFromName: 'smtp_from_name',
  newCommentTitle: 'tpl_new_comment_title',
  newCommentBody: 'tpl_new_comment_body',
  replyCommentTitle: 'tpl_reply_comment_title',
  replyCommentBody: 'tpl_reply_comment_body',
  pendingPluginReviewTitle: 'tpl_pending_plugin_review_title',
  pendingPluginReviewBody: 'tpl_pending_plugin_review_body',
  pluginApprovedTitle: 'tpl_plugin_approved_title',
  pluginApprovedBody: 'tpl_plugin_approved_body',
  pluginRejectedTitle: 'tpl_plugin_rejected_title',
  pluginRejectedBody: 'tpl_plugin_rejected_body',
  emailVerifyTitle: 'tpl_bind_email_title',
  emailVerifyBody: 'tpl_bind_email_body',
  resetPasswordTitle: 'tpl_reset_pwd_title',
  resetPasswordBody: 'tpl_reset_pwd_body',
} as const

const emailPayload = () => Object.fromEntries(
  Object.entries(emailSettingMap).map(([backendKey, frontendKey]) => [backendKey, String(settings.email[frontendKey])]),
)

const applyEmailSettings = async (data: Record<string, string>, persistMissing = true) => {
  const defaults: Record<string, string> = {
    newCommentTitle: DEFAULT_EMAIL_TEMPLATES.newComment.title,
    newCommentBody: DEFAULT_EMAIL_TEMPLATES.newComment.body,
    replyCommentTitle: DEFAULT_EMAIL_TEMPLATES.replyComment.title,
    replyCommentBody: DEFAULT_EMAIL_TEMPLATES.replyComment.body,
    pendingPluginReviewTitle: DEFAULT_EMAIL_TEMPLATES.pendingPluginReview.title,
    pendingPluginReviewBody: DEFAULT_EMAIL_TEMPLATES.pendingPluginReview.body,
    pluginApprovedTitle: DEFAULT_EMAIL_TEMPLATES.pluginApproved.title,
    pluginApprovedBody: DEFAULT_EMAIL_TEMPLATES.pluginApproved.body,
    pluginRejectedTitle: DEFAULT_EMAIL_TEMPLATES.pluginRejected.title,
    pluginRejectedBody: DEFAULT_EMAIL_TEMPLATES.pluginRejected.body,
    emailVerifyTitle: DEFAULT_EMAIL_TEMPLATES.emailVerify.title,
    emailVerifyBody: DEFAULT_EMAIL_TEMPLATES.emailVerify.body,
    resetPasswordTitle: DEFAULT_EMAIL_TEMPLATES.resetPassword.title,
    resetPasswordBody: DEFAULT_EMAIL_TEMPLATES.resetPassword.body,
  }
  const missingTemplates: Record<string, string> = {}
  for (const [backendKey, frontendKey] of Object.entries(emailSettingMap)) {
    const fallback = defaults[backendKey] || (backendKey === 'emailPort' ? '465' : '')
    settings.email[frontendKey] = data[backendKey] || fallback
    if (backendKey in defaults && !data[backendKey]) missingTemplates[backendKey] = defaults[backendKey]
  }
  emailSnapshot.value = emailPayload()
  if (persistMissing && Object.keys(missingTemplates).length) {
    try {
      await useApiFetch('/admin/setting', { method: 'PUT', body: missingTemplates })
      Object.assign(emailSnapshot.value, missingTemplates)
    } catch {
      ElMessage.warning('默认邮件模板初始化失败，请稍后重试')
    }
  }
}

const loadSettings = async () => {
  try {
    const response: any = await useApiFetch('/admin/setting')
    const data = response?.data || {}
    applyBasicSettings(data)
    applyFunctionSettings(data)
    applyUploadSettings(data)
    await applyNotificationSettings(data)
    await applyEmailSettings(data)
  } catch {
    applyBasicSettings({})
    applyFunctionSettings({})
    applyUploadSettings({})
    await applyNotificationSettings({}, false)
    await applyEmailSettings({}, false)
  }
}

const startBasicEdit = () => {
  basicSnapshot.value = copyBasicSettings()
  settings.basic.isEditing = true
}

const cancelBasicEdit = () => {
  Object.assign(settings.basic, basicSnapshot.value)
  settings.basic.isEditing = false
}

const uploadSiteLogo = async (options: UploadRequestOptions) => {
  const maxSize = Number(settings.upload.image_max_size) || DEFAULT_MAX_IMAGE_SIZE
  if (options.file.size > maxSize * 1024 * 1024) {
    const message = `图片大小超过限制，最大允许 ${maxSize}MB`
    ElMessage.error(message)
    options.onError(new Error(message))
    return
  }

  const body = new FormData()
  body.append('file', options.file)
  try {
    const response: any = await useApiFetch('/admin/setting/logo', { method: 'POST', body })
    settings.basic.site_logo = response?.data?.logo || settings.basic.site_logo
    options.onSuccess(response)
  } catch (error) {
    showApiError(error, '网站 Logo 上传失败')
    options.onError(error as Error)
  }
}

const clearSiteLogo = async () => {
  logoClearing.value = true
  try {
    await useApiFetch('/admin/setting/logo', { method: 'DELETE' })
    settings.basic.site_logo = DEFAULT_SITE_LOGO
    basicSnapshot.value.site_logo = DEFAULT_SITE_LOGO
    await loadPublicSiteSettings()
    ElMessage.success('已恢复默认 Logo')
  } finally {
    logoClearing.value = false
  }
}

const saveBasicSettings = async () => {
  if (!settings.basic.site_title.trim()) {
    ElMessage.warning('网站标题不能为空')
    return
  }

  basicSaving.value = true
  try {
    const basic = {
      siteLogo: settings.basic.site_logo,
      siteTitle: settings.basic.site_title.trim(),
      siteDesc: settings.basic.site_description.trim(),
      siteKeywords: settings.basic.site_keywords.trim(),
    }
    await useApiFetch('/admin/setting', { method: 'PUT', body: basic })
    await loadPublicSiteSettings()
    applyBasicSettings(basic)
    settings.basic.isEditing = false
    ElMessage.success('基础设置已保存')
  } finally {
    basicSaving.value = false
  }
}

const saveFunctionSettings = async () => {
  if (functionSaving.value) return
  functionSaving.value = true
  try {
    await useApiFetch('/admin/setting', {
      method: 'PUT',
      body: {
        allowRegister: String(settings.function.allow_register),
        allowComment: String(settings.function.allow_comment),
        skipAudit: String(settings.function.plugin_no_audit),
        allowUpload: String(settings.function.allow_upload_image),
        antiBrushPluginData: String(settings.function.plugin_anti_brush),
      },
    })
    ElMessage.success('功能设置已保存')
  } catch {
    await loadSettings()
  } finally {
    functionSaving.value = false
  }
}

const startUploadEdit = () => {
  uploadSnapshot.value = copyUploadSettings()
  settings.upload.isEditing = true
}

const cancelUploadEdit = () => {
  Object.assign(settings.upload, uploadSnapshot.value)
  settings.upload.isEditing = false
}

const saveUploadSettings = async () => {
  const maxSize = Number(settings.upload.image_max_size)
  const allowedExtensions = settings.upload.image_allowed_ext
    .split(',')
    .map(ext => ext.trim().toLowerCase())
    .filter(Boolean)
    .join(',')
  if (!Number.isFinite(maxSize) || maxSize <= 0) {
    ElMessage.warning('文件最大限制必须大于 0')
    return
  }
  if (!allowedExtensions) {
    ElMessage.warning('允许图片格式不能为空')
    return
  }

  uploadSaving.value = true
  try {
    await useApiFetch('/admin/setting', {
      method: 'PUT',
      body: {
        maxFileSize: String(maxSize),
        allowedExtensions,
      },
    })
    settings.upload.image_max_size = maxSize
    settings.upload.image_allowed_ext = allowedExtensions
    uploadSnapshot.value = copyUploadSettings()
    settings.upload.isEditing = false
    ElMessage.success('图片限制已保存')
  } finally {
    uploadSaving.value = false
  }
}

const startEmailEdit = () => {
  emailSnapshot.value = emailPayload()
  settings.email.isEditing = true
}

const startNotificationEdit = () => {
  notificationSnapshot.value = notificationPayload()
  settings.notification.isEditing = true
}

const cancelNotificationEdit = () => {
  for (const [backendKey, frontendKey] of Object.entries(notificationSettingMap)) {
    settings.notification[frontendKey] = notificationSnapshot.value[backendKey] || notificationDefaults[backendKey] || ''
  }
  settings.notification.isEditing = false
}

const saveNotificationSettings = async () => {
  const payload = notificationPayload()
  if (Object.values(payload).some(value => !value.trim())) {
    ElMessage.warning('通知模板的标题和内容不能为空')
    return
  }
  notificationSaving.value = true
  try {
    await useApiFetch('/admin/setting', { method: 'PUT', body: payload })
    notificationSnapshot.value = payload
    settings.notification.isEditing = false
    ElMessage.success('通知配置已保存')
  } finally {
    notificationSaving.value = false
  }
}

const cancelEmailEdit = () => {
  for (const [backendKey, frontendKey] of Object.entries(emailSettingMap)) {
    settings.email[frontendKey] = emailSnapshot.value[backendKey] || ''
  }
  settings.email.isEditing = false
}

const saveEmailSettings = async () => {
  if (!settings.email.smtp_server.trim() || !settings.email.smtp_user.trim() || !settings.email.smtp_pass) {
    ElMessage.warning('请填写SMTP服务器、账号和授权码')
    return
  }
  const port = Number(settings.email.smtp_port)
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    ElMessage.warning('SMTP端口无效')
    return
  }
  emailSaving.value = true
  try {
    await useApiFetch('/admin/setting', { method: 'PUT', body: emailPayload() })
    emailSnapshot.value = emailPayload()
    settings.email.isEditing = false
    ElMessage.success('邮件配置已保存')
  } finally {
    emailSaving.value = false
  }
}

const openTestEmail = () => {
  testEmailAddress.value = ''
  testEmailVisible.value = true
}

const sendTestEmail = async () => {
  if (testEmailSending.value) return
  const email = testEmailAddress.value.trim()
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    ElMessage.warning('请输入正确的邮箱地址')
    return
  }
  testEmailSending.value = true
  try {
    await useApiFetch('/admin/setting/test-email', {
      method: 'POST',
      body: { email },
      errorFallback: '测试邮件发送失败，请检查SMTP配置',
    })
    testEmailVisible.value = false
    ElMessage.success('测试邮件已发送，请检查收件箱')
  } catch {
  } finally {
    testEmailSending.value = false
  }
}

const formatBytes = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`
}

const handleCleanImages = async () => {
  cleanupRunning.value = true
  try {
    const response: any = await useApiFetch('/admin/setting/clean-images', { method: 'POST' })
    const data = response?.data || {}
    ElMessage.success(`清理成功，共清理 ${Number(data.cleanedCount) || 0} 个文件，释放 ${formatBytes(Number(data.cleanedBytes) || 0)}`)
  } finally {
    cleanupRunning.value = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <div class="flex flex-col gap-5">
    <manage-box title="基础设置">
      <template #header>
        <div v-if="!settings.basic.isEditing">
          <pm-button text="编辑" @click="startBasicEdit"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelBasicEdit"/>
          <pm-button text="保存" :loading="basicSaving" @click="saveBasicSettings"/>
        </div>
      </template>
      <div class="setting-basic-layout flex gap-10">
        <div class="shrink-0 text-center">
          <el-upload class="avatar-uploader" action="#" :show-file-list="false" accept="image/*"
                     :http-request="uploadSiteLogo" :disabled="!settings.basic.isEditing">
            <div class="w-24 h-24 rounded-2xl overflow-hidden border-2 border-[#F1F5F9] relative group">
              <img :src="settings.basic.site_logo.startsWith('/uploads/') ? assetUrl(settings.basic.site_logo, settings.basic.site_logo) : settings.basic.site_logo"
                   class="w-full h-full object-cover" alt="Logo"/>
              <div v-if="settings.basic.isEditing"
                   class="absolute inset-0 bg-black/40 flex items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-opacity">
                <PhCamera :size="24"/>
              </div>
            </div>
          </el-upload>
          <span class="text-xs text-[#94A3B8] block mt-2" v-if="settings.basic.isEditing">点击更换图标</span>
          <pm-button v-if="settings.basic.site_logo !== DEFAULT_SITE_LOGO"
                     class="mt-2" text="清除 Logo" color="white" size="small"
                     :loading="logoClearing" @click="clearSiteLogo"/>
        </div>
        <el-form label-position="top" class="flex-1 w-full">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1">
            <el-form-item>
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhTextT :size="18"/>
                  网站标题
                </div>
              </template>
              <div v-if="!settings.basic.isEditing"
                   class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
                {{ settings.basic.site_title }}
              </div>
              <el-input v-else v-model="settings.basic.site_title" placeholder="请输入网站标题" class="pm-manage-input"/>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhHash :size="18"/>
                  网站关键字
                </div>
              </template>
              <div v-if="!settings.basic.isEditing"
                   class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
                {{ settings.basic.site_keywords }}
              </div>
              <el-input v-else v-model="settings.basic.site_keywords" placeholder="请输入关键字，用逗号分隔"
                        class="pm-manage-input"/>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label>
                <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                  <PhNote :size="18"/>
                  网站描述
                </div>
              </template>
              <div v-if="!settings.basic.isEditing"
                   class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] select-text text-[14px]">
                {{ settings.basic.site_description }}
              </div>
              <el-input v-else type="textarea" v-model="settings.basic.site_description" placeholder="请输入网站描述"
                        class="pm-manage-input"/>
            </el-form-item>
          </div>
        </el-form>
      </div>
    </manage-box>

    <manage-box title="功能设置">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div
            class="flex items-center justify-between p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-[#1E293B]">允许注册</span>
            <span class="text-xs text-[#94A3B8]">开启后新用户可注册账号</span>
          </div>
          <el-switch v-model="settings.function.allow_register" :loading="functionSaving" @change="saveFunctionSettings"/>
        </div>
        <div
            class="flex items-center justify-between p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-[#1E293B]">允许评论</span>
            <span class="text-xs text-[#94A3B8]">开启后插件详情页可评论</span>
          </div>
          <el-switch v-model="settings.function.allow_comment" :loading="functionSaving" @change="saveFunctionSettings"/>
        </div>
        <div
            class="flex items-center justify-between p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-[#1E293B]">插件无需审核</span>
            <span class="text-xs text-[#94A3B8]">开启后提交插件直接发布</span>
          </div>
          <el-switch v-model="settings.function.plugin_no_audit" :loading="functionSaving" @change="saveFunctionSettings"/>
        </div>
        <div
            class="flex items-center justify-between p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-[#1E293B]">允许上传图片</span>
            <span class="text-xs text-[#94A3B8]">开启后则可以上传头像、插件图片等</span>
          </div>
          <el-switch v-model="settings.function.allow_upload_image" :loading="functionSaving" @change="saveFunctionSettings"/>
        </div>
        <div
            class="flex items-center justify-between p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-[#1E293B]">防刷插件数据</span>
            <span class="text-xs text-[#94A3B8]">可有效降低刷查看量和下载量</span>
          </div>
          <el-switch v-model="settings.function.plugin_anti_brush" :loading="functionSaving" @change="saveFunctionSettings"/>
        </div>
      </div>
    </manage-box>

    <manage-box title="图片限制">
      <template #header>
        <div v-if="!settings.upload.isEditing">
          <pm-button text="编辑" @click="startUploadEdit"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelUploadEdit"/>
          <pm-button text="保存" :loading="uploadSaving" @click="saveUploadSettings"/>
        </div>
      </template>
      <el-form label-position="top" class="max-w-2xl">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1">
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                <PhSelection :size="18"/>
                文件最大限制 (MB)
              </div>
            </template>
            <div v-if="!settings.upload.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.upload.image_max_size }}
            </div>
            <el-input v-else v-model="settings.upload.image_max_size" placeholder="请输入最大大小"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                <PhFilePlus :size="18"/>
                允许图片格式
              </div>
            </template>
            <div v-if="!settings.upload.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.upload.image_allowed_ext }}
            </div>
            <el-input v-else v-model="settings.upload.image_allowed_ext" placeholder="如: jpg,png,webp"
                      class="pm-manage-input"/>
          </el-form-item>
        </div>
      </el-form>
    </manage-box>
    <manage-box title="程序优化">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div
            class="flex flex-col gap-4 p-4 bg-[#F8FAFC] rounded-2xl border border-transparent hover:border-[#E2E8F0] transition-all">
          <div class="flex flex-col gap-1">
            <div class="flex items-center gap-1.5">
              <PhBroom :size="18" class="text-[#00BAAD]"/>
              <span class="text-sm font-medium text-[#1E293B]">清理未使用的图片</span>
            </div>
            <span class="text-xs text-[#94A3B8] leading-5">清理服务器上未引用的冗余图片，释放磁盘空间</span>
          </div>
          <div class="flex justify-end">
            <pm-button text="一键清理" size="small" :disabled="cleanupRunning" @click="handleCleanImages"/>
          </div>
        </div>
      </div>
    </manage-box>

    <manage-box title="通知配置">
      <template #header>
        <div v-if="!settings.notification.isEditing">
          <pm-button text="编辑" @click="startNotificationEdit"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelNotificationEdit"/>
          <pm-button text="保存" :loading="notificationSaving" @click="saveNotificationSettings"/>
        </div>
      </template>
      <el-form label-position="top" class="w-full">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1">
          <div v-for="template in notificationTemplates" :key="template.title" class="md:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-x-6">
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">{{ template.name }}</span>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">标题</span></template>
              <div v-if="!settings.notification.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">
                {{ settings.notification[template.title] }}
              </div>
              <el-input v-else v-model="settings.notification[template.title]" class="pm-manage-input"/>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">内容</span></template>
              <div v-if="!settings.notification.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">
                {{ settings.notification[template.content] }}
              </div>
              <el-input v-else type="textarea" v-model="settings.notification[template.content]" class="pm-manage-input"/>
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div v-for="variable in template.variables" :key="variable[0]" class="flex items-center gap-1">
                  <span class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">&#123;&#123;{{ variable[0] }}&#125;&#125;</span>
                  <span class="text-[10px] text-[#94A3B8]">{{ variable[1] }}</span>
                </div>
              </div>
            </el-form-item>
          </div>
        </div>
      </el-form>
    </manage-box>

    <manage-box title="邮件配置">
      <template #header>
        <div v-if="!settings.email.isEditing" class="flex gap-2">
          <pm-button text="测试可用性" color="white" @click="openTestEmail"/>
          <pm-button text="编辑" @click="startEmailEdit"/>
        </div>
        <div v-else class="flex gap-2">
          <pm-button text="取消" color="white" @click="cancelEmailEdit"/>
          <pm-button text="保存" :loading="emailSaving" @click="saveEmailSettings"/>
        </div>
      </template>
      <el-form label-position="top" class="w-full">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1">
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                邮件发送SMTP服务器地址
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_server }}
            </div>
            <el-input v-else v-model="settings.email.smtp_server" placeholder="例如: smtp.qq.com"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                SMTP端口
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_port }}
            </div>
            <el-input v-else v-model="settings.email.smtp_port" placeholder="例如: 465"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                SMTP账号
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_user || '未设置' }}
            </div>
            <el-input v-else v-model="settings.email.smtp_user" placeholder="通常为发件邮箱账号"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                发件邮箱
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_from || settings.email.smtp_user || '未设置' }}
            </div>
            <el-input v-else v-model="settings.email.smtp_from" placeholder="留空时使用SMTP账号"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                发件人名称
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_from_name || 'PluginMarket' }}
            </div>
            <el-input v-else v-model="settings.email.smtp_from_name" placeholder="例如: PluginMarket"
                      class="pm-manage-input"/>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="flex items-center gap-1.5 text-[#64748B] font-medium">
                密码 / 授权码
              </div>
            </template>
            <div v-if="!settings.email.isEditing"
                 class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] select-text text-[14px]">
              {{ settings.email.smtp_pass ? '••••••••' : '未设置' }}
            </div>
            <el-input v-else v-model="settings.email.smtp_pass" type="password" show-password
                      placeholder="请输入SMTP授权码" class="pm-manage-input"/>
          </el-form-item>

          <!-- 模板设置 -->
          <div class="md:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-x-6 mt-4">
            <!-- 新评论 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">新评论邮件模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_new_comment_title', 'tpl_new_comment_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_new_comment_title }}</div>
              <el-input v-else v-model="settings.email.tpl_new_comment_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap break-words overflow-x-hidden overflow-y-auto">{{ settings.email.tpl_new_comment_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_new_comment_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span>
                  <span class="text-[10px] text-[#94A3B8]">网站标题</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{plugin_name}}</span>
                  <span class="text-[10px] text-[#94A3B8]">插件名称</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{comment_user}}</span>
                  <span class="text-[10px] text-[#94A3B8]">评论者</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{comment_content}}</span>
                  <span class="text-[10px] text-[#94A3B8]">评论内容</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{time}}</span>
                  <span class="text-[10px] text-[#94A3B8]">评论时间</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{link}}</span>
                  <span class="text-[10px] text-[#94A3B8]">查看链接</span>
                </div>
              </div>
            </el-form-item>

            <!-- 评论回复 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">评论被回复邮件模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_reply_comment_title', 'tpl_reply_comment_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_reply_comment_title }}</div>
              <el-input v-else v-model="settings.email.tpl_reply_comment_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap break-words overflow-x-hidden overflow-y-auto">{{ settings.email.tpl_reply_comment_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_reply_comment_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span>
                  <span class="text-[10px] text-[#94A3B8]">网站标题</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{plugin_name}}</span>
                  <span class="text-[10px] text-[#94A3B8]">插件名称</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{reply_user}}</span>
                  <span class="text-[10px] text-[#94A3B8]">回复者</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{reply_content}}</span>
                  <span class="text-[10px] text-[#94A3B8]">回复内容</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{time}}</span>
                  <span class="text-[10px] text-[#94A3B8]">被回复时间</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{user}}</span>
                  <span class="text-[10px] text-[#94A3B8]">被回复人昵称</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{content}}</span>
                  <span class="text-[10px] text-[#94A3B8]">被回复人评论内容</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{link}}</span>
                  <span class="text-[10px] text-[#94A3B8]">查看链接</span>
                </div>
              </div>
            </el-form-item>

            <!-- 待审核插件 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">有插件待审核通知模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_pending_plugin_review_title', 'tpl_pending_plugin_review_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_pending_plugin_review_title }}</div>
              <el-input v-else v-model="settings.email.tpl_pending_plugin_review_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">{{ settings.email.tpl_pending_plugin_review_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_pending_plugin_review_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span><span class="text-[10px] text-[#94A3B8]">网站标题</span></div>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{plugin_name}}</span><span class="text-[10px] text-[#94A3B8]">插件名称</span></div>
              </div>
            </el-form-item>

            <!-- 插件审核通过 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">插件审核通过通知模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_plugin_approved_title', 'tpl_plugin_approved_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_plugin_approved_title }}</div>
              <el-input v-else v-model="settings.email.tpl_plugin_approved_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">{{ settings.email.tpl_plugin_approved_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_plugin_approved_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span><span class="text-[10px] text-[#94A3B8]">网站标题</span></div>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{plugin_name}}</span><span class="text-[10px] text-[#94A3B8]">插件名称</span></div>
              </div>
            </el-form-item>

            <!-- 插件审核未通过 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">插件审核未通过通知模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_plugin_rejected_title', 'tpl_plugin_rejected_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_plugin_rejected_title }}</div>
              <el-input v-else v-model="settings.email.tpl_plugin_rejected_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">{{ settings.email.tpl_plugin_rejected_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_plugin_rejected_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span><span class="text-[10px] text-[#94A3B8]">网站标题</span></div>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{plugin_name}}</span><span class="text-[10px] text-[#94A3B8]">插件名称</span></div>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{content}}</span><span class="text-[10px] text-[#94A3B8]">拒绝理由</span></div>
                <div class="flex items-center gap-1"><span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{link}}</span><span class="text-[10px] text-[#94A3B8]">插件链接</span></div>
              </div>
            </el-form-item>

            <!-- 绑定邮箱 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">绑定邮箱邮件模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_bind_email_title', 'tpl_bind_email_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_bind_email_title }}</div>
              <el-input v-else v-model="settings.email.tpl_bind_email_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">{{ settings.email.tpl_bind_email_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_bind_email_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span>
                  <span class="text-[10px] text-[#94A3B8]">网站标题</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{user_name}}</span>
                  <span class="text-[10px] text-[#94A3B8]">用户名称</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{link}}</span>
                  <span class="text-[10px] text-[#94A3B8]">验证链接</span>
                </div>
              </div>
            </el-form-item>

            <!-- 重置密码 -->
            <el-form-item class="md:col-span-2">
              <div class="flex items-center gap-2 mb-1">
                <div class="w-1 h-4 bg-[#00BAAD] rounded-xs"></div>
                <span class="text-sm font-bold text-[#1E293B]">重置密码邮件模板</span>
                <PhEye :size="16" class="text-[#00BAAD] cursor-pointer hover:text-[#00A396]"
                        @mouseenter="showPreview($event, 'tpl_reset_pwd_title', 'tpl_reset_pwd_body')"
                        @mouseleave="hidePreview"/>
              </div>
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件标题</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-12 flex items-center bg-[#F8FAFC] rounded-xl px-3.5 text-[#606266] text-[14px]">{{ settings.email.tpl_reset_pwd_title }}</div>
              <el-input v-else v-model="settings.email.tpl_reset_pwd_title" class="pm-manage-input" />
            </el-form-item>
            <el-form-item class="md:col-span-2">
              <template #label><span class="text-xs text-[#94A3B8]">邮件正文</span></template>
              <div v-if="!settings.email.isEditing" class="w-full h-30 bg-[#F8FAFC] rounded-xl px-3.5 py-2.5 leading-normal text-[#606266] text-[14px] whitespace-pre-wrap overflow-y-auto">{{ settings.email.tpl_reset_pwd_body }}</div>
              <el-input v-else type="textarea" v-model="settings.email.tpl_reset_pwd_body" class="pm-manage-input" />
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 items-center">
                <span class="text-[11px] text-[#94A3B8]">可用变量：</span>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{site_title}}</span>
                  <span class="text-[10px] text-[#94A3B8]">网站标题</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{user_name}}</span>
                  <span class="text-[10px] text-[#94A3B8]">用户名称</span>
                </div>
                <div class="flex items-center gap-1">
                  <span v-pre class="px-1.5 py-0.5 bg-[#F1F5F9] text-[#64748B] text-[10px] rounded select-all cursor-pointer">{{link}}</span>
                  <span class="text-[10px] text-[#94A3B8]">重置链接</span>
                </div>
              </div>
            </el-form-item>
          </div>
        </div>
      </el-form>
    </manage-box>

    <el-dialog v-model="testEmailVisible" title="测试SMTP邮件" width="440px" align-center class="pm-manage-dialog">
      <el-form label-position="top" @submit.prevent="sendTestEmail">
        <el-form-item label="收件邮箱">
          <el-input v-model="testEmailAddress" placeholder="请输入用于接收测试邮件的邮箱" class="pm-manage-input"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <pm-button text="取消" color="white" @click="testEmailVisible = false"/>
          <pm-button text="发送测试邮件" :disabled="testEmailSending" @click="sendTestEmail"/>
        </div>
      </template>
    </el-dialog>

    <!-- 统一悬浮预览 -->
    <Transition name="preview">
      <div v-if="previewState.visible"
           class="email-preview-popover fixed z-9999 bg-white rounded-xl shadow-xl border border-[#E2E8F0] overflow-hidden"
           :style="{ left: previewState.x + 'px', top: previewState.y + 'px' }"
           @mouseenter="onPreviewEnter"
           @mouseleave="hidePreview">
        <div class="px-4 py-3 bg-[#F8FAFC] border-b border-[#E2E8F0]">
          <span class="text-sm font-semibold text-[#1E293B]">{{ previewState.title }}</span>
        </div>
        <div class="p-4 max-h-100 overflow-y-auto">
          <iframe sandbox :srcdoc="previewState.content" class="email-preview-frame w-120 min-h-60 border-0" title="邮件模板预览"></iframe>
        </div>
      </div>
    </Transition>

  </div>
</template>

<style scoped lang="scss">

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

    &.is-focus,&:focus {
      box-shadow: inset 0 0 0 1.5px #00BAAD !important;
    }

  }

  .el-input__wrapper {
    height: 48px !important;
    padding-block: 0;
  }

  .el-input__inner {
    height: 100%;
    line-height: normal;
  }

  .el-textarea__inner {
    padding: 10px 14px !important;
    height: 120px;
    resize: none;
  }
}

.avatar-uploader :deep(.el-upload) {
  cursor: pointer;
  background: none;
  border: none;
}

// 预览悬浮动画
.preview-enter-active,
.preview-leave-active {
  transition: opacity 0.2s ease;
}

.preview-enter-from,
.preview-leave-to {
  opacity: 0;
}

// 邮件预览内容样式
.email-preview-content {
  font-size: 14px;
  line-height: 1.8;
  color: #334155;
  word-break: break-word;

  :deep(p) {
    margin: 8px 0;
  }

  :deep(blockquote) {
    border-left: 3px solid #00BAAD;
    padding-left: 12px;
    margin: 16px 0;
    color: #64748B;
  }

  :deep(a) {
    color: #00BAAD;
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }
}

@media (max-width: 767px) {
  .setting-basic-layout {
    flex-direction: column;
    gap: 20px;
    align-items: center;

    > :deep(.el-form) {
      min-width: 0;
    }
  }

  .email-preview-popover {
    left: 12px !important;
    right: 12px !important;
    top: 12px !important;
    max-height: calc(100dvh - 24px);
  }

  .email-preview-frame {
    width: 100%;
    min-width: 0;
  }
}

</style>
