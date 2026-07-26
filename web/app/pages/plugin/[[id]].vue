<script setup lang="ts">
import dayjs from 'dayjs'
import {
  PhDownloadSimple,
  PhShare,
  PhStar,
} from '@phosphor-icons/vue'
import { getCommentsByPlugin, createComment, deleteComment, getPluginDownloadInfo, type ApiResponse, type CommentItem, type CommentListData } from '~/composables/api/plugin'
import { parseTargetCommentId, targetCommentElementId } from '~/utils/comment-target'
import { toggleStar } from '~/composables/api/public'
import { DEFAULT_PLUGIN_ICON } from '~/config'
import { getPluginSeoMeta } from '~/utils/plugin-seo'
import { getApiErrorMessage } from '~/composables/useApiFetch'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { configureMarkdownEditor } from '~/utils/md-editor-config'

configureMarkdownEditor()

type ReplyTarget = {
  id: number
  author?: { nick?: string }
}

const route = useRoute()
const userStore = useUserStore()
const assetUrl = useAssetUrl()
const { settings } = useSiteSettings()
const pluginId = computed(() => Number(route.params.id || 0))
const COMMENT_PAGE_SIZE = 20

const commentData = ref<{ list: CommentItem[]; total: number }>({ list: [], total: 0 })
const commentPage = ref(1)
const commentLoading = ref(false)
const commentText = ref('')
const replyingComment = ref<ReplyTarget | null>(null)
const submittingComment = ref(false)
const commentsLoaded = ref(false)
const commentTarget = ref<NonNullable<CommentListData['target']> | null>(null)
const downloadDialogVisible = ref(false)
const downloadLoading = ref(false)
const starLoading = ref(false)
const pluginStarCount = ref(0)
const pluginIsStarred = ref(false)
const pageMounted = ref(false)

const detailUrl = computed(() => (pluginId.value ? `/plugin/${pluginId.value}` : ''))
const { data: detailRes, pending: detailLoading, error: detailError } = await useAsyncData(
  () => `plugin-detail-${pluginId.value}`,
  () => useApiFetch<ApiResponse<any>>(detailUrl.value, { suppressErrorMessage: true }),
  {
    watch: [pluginId],
  }
)

const accessDeniedHandled = ref(false)
watch([detailRes, detailError], async ([response, error]) => {
  if (!import.meta.client) return
  if (accessDeniedHandled.value) return
  const message = response && response.code !== 200
    ? response.msg
    : error
      ? getApiErrorMessage(error)
      : ''
  if (!message) return
  accessDeniedHandled.value = true
  ElMessage.error(message)
  await navigateTo('/')
}, { immediate: true })

const pluginData = computed(() => detailRes.value?.data || null)
const pluginTitle = computed(() => pluginData.value?.name || '插件详情')
const pluginDesc = computed(() => pluginData.value?.desc_text || '暂无简介')
const pluginContent = computed(() => pluginData.value?.content || '')
const pluginIcon = computed(() => assetUrl(pluginData.value?.icon, DEFAULT_PLUGIN_ICON))
const pluginFrameworks = computed(() => pluginData.value?.frameworks || [])
const pluginTags = computed(() => pluginData.value?.tags || [])
const pluginCreated = computed(() => pluginData.value?.created ? dayjs(pluginData.value.created).format('YYYY年M月D日') : '-')
const pluginUpdated = computed(() => pluginData.value?.updated ? dayjs(pluginData.value.updated).format('YYYY年M月D日') : '-')
const pluginAuthor = computed(() => pluginData.value?.author)

const pluginSeo = computed(() => getPluginSeoMeta(pluginData.value, {
  title: settings.value.siteTitle,
  description: settings.value.siteDescription,
}))

useSeoMeta({
  title: () => pluginSeo.value.title,
  description: () => pluginSeo.value.description,
})

watch(pluginData, () => {
  pluginStarCount.value = Number(pluginData.value?.star || 0)
  pluginIsStarred.value = !!pluginData.value?.is_starred
}, { immediate: true })

const scrollToTargetComment = async () => {
  const target = commentTarget.value
  if (!target) return
  await nextTick()
  await nextTick()
  document.getElementById(targetCommentElementId(target.commentId))?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

const loadCommentList = async (page = commentPage.value, targetCommentId: number | null = null) => {
  if (!pluginId.value) return commentData.value
  commentLoading.value = true
  try {
    const res = await getCommentsByPlugin(pluginId.value, page, COMMENT_PAGE_SIZE, targetCommentId)
    commentData.value = {
      list: res.data?.list ?? [],
      total: res.data?.total ?? 0,
    }
    commentPage.value = res.data?.target?.page || page
    commentsLoaded.value = true
    commentTarget.value = res.data?.target || null
    if (commentTarget.value) await scrollToTargetComment()
    return commentData.value
  } catch (error) {
    if (targetCommentId) showApiError(error, '该评论已被删除、不存在或不属于当前插件')
    throw error
  } finally {
    commentLoading.value = false
  }
}

watch(pluginData, (detail) => {
  if (!pageMounted.value || !detail || commentsLoaded.value) return
  void loadCommentList(1, parseTargetCommentId(route.query.commentId as string | string[] | undefined)).catch(() => {})
}, { immediate: true })

onMounted(() => {
  pageMounted.value = true
  if (!pluginData.value || commentsLoaded.value) return
  void loadCommentList(1, parseTargetCommentId(route.query.commentId as string | string[] | undefined)).catch(() => {})
})

watch(() => route.query.commentId, (value, previous) => {
  if (!commentsLoaded.value || value === previous) return
  void loadCommentList(1, parseTargetCommentId(value as string | string[] | undefined)).catch(() => {})
})

watch(pluginId, () => {
  commentData.value = { list: [], total: 0 }
  commentPage.value = 1
  commentsLoaded.value = false
  commentTarget.value = null
  accessDeniedHandled.value = false
})

const handleCommentPageChange = (page: number) => {
  commentTarget.value = null
  void loadCommentList(page)
}

const handleToggleStar = async () => {
  if (!userStore.isLogin) {
    ElMessage.warning('登录后才可以收藏插件')
    await navigateTo('/user/auth?mode=login')
    return
  }
  if (!pluginId.value || starLoading.value) return
  starLoading.value = true
  try {
    const res = await toggleStar(pluginId.value)
    const isStarred = !!res.data?.is_starred
    pluginIsStarred.value = isStarred
    pluginStarCount.value = Math.max(0, pluginStarCount.value + (isStarred ? 1 : -1))
    if (detailRes.value?.data) {
      detailRes.value.data.is_starred = isStarred
      detailRes.value.data.star = pluginStarCount.value
    }
    ElMessage.success(isStarred ? '收藏成功' : '已取消收藏')
  } finally {
    starLoading.value = false
  }
}

const handleDownload = () => {
  if (!pluginData.value?.url) {
    ElMessage.warning('该插件未配置下载地址')
    return
  }
  downloadDialogVisible.value = true
}

const openDownloadUrl = async () => {
  if (!pluginId.value) return
  downloadLoading.value = true
  try {
    const res = await getPluginDownloadInfo(pluginId.value)
    if (!res.data?.url) {
      ElMessage.warning('该插件未配置下载地址')
      return
    }
    window.open(res.data.url, '_blank', 'noopener,noreferrer')
    downloadDialogVisible.value = false
  } finally {
    downloadLoading.value = false
  }
}

const closeDownloadDialog = () => {
  downloadDialogVisible.value = false
}

const copyShareTextToClipboard = async (text: string) => {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'true')
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  document.body.appendChild(textarea)
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
}

const handleCopyShare = async () => {
  const url = typeof window !== 'undefined' ? window.location.href : ''
  const content = `我在PluginMarket发现了一个宝藏插件《${pluginTitle.value}》，快去看看吧！${url}`
  await copyShareTextToClipboard(content)
  ElMessage.success('链接复制成功')
}

const handleCommentDisabledClick = () => {
  ElMessage.warning('登录后才可以发表评论')
}

const handleSubmitComment = async () => {
  if (!userStore.isLogin) {
    ElMessage.warning('登录后才可以发表评论')
    await navigateTo('/user/auth?mode=login')
    return
  }
  const content = commentText.value.trim()
  if (!content) {
    ElMessage.warning('请输入评论内容')
    return
  }
  submittingComment.value = true
  try {
    await createComment({ pluginId: pluginId.value, content, parentId: replyingComment.value?.id || 0 })
    commentText.value = ''
    replyingComment.value = null
    ElMessage.success('评论发布成功')
    await loadCommentList()
  } catch (error) {
    showApiError(error, '评论发布失败，请稍后重试')
  } finally {
    submittingComment.value = false
  }
}

const handleReply = (comment: ReplyTarget) => {
  if (!userStore.isLogin) {
    ElMessage.warning('登录后才可以回复评论')
    return
  }
  replyingComment.value = comment
}

const canDeleteComment = (comment: any) => {
  const currentUser = userStore.userInfo
  if (!userStore.isLogin || !currentUser?.id) return false
  if (currentUser.power === 1) return true
  if (currentUser.id === pluginAuthor.value?.id) return true
  return currentUser.id === comment.author?.id
}

const handleDeleteComment = async (comment: any) => {
  if (!canDeleteComment(comment)) {
    ElMessage.warning('你没有权限删除这条评论')
    return
  }
  await ElMessageBox.confirm('确认删除这条评论吗？', '删除确认', {
    type: 'warning',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
  await deleteComment(comment.id)
  ElMessage.success('删除成功')
  await loadCommentList()
}

const cancelReply = () => {
  replyingComment.value = null
}
</script>

<template>
  <div class="plugin-detail-layout flex flex-col lg:flex-row items-start gap-3 md:gap-5 max-w-350 m-auto px-3 md:px-5 mb-6 md:mb-10">
    <div class="flex w-full min-w-0 flex-1 gap-3 md:gap-5 flex-col select-none">
      <div class="plugin-hero flex flex-col sm:flex-row gap-4 md:gap-5 sm:items-center justify-between bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-4 md:p-5 md:px-8">
        <div class="flex gap-4 md:gap-5 min-w-0">
          <div class="flex items-center shrink-0">
            <img class="logo" :src="pluginIcon" alt="avatar" draggable="false"/>
          </div>
          <div class="flex flex-col min-w-0">
            <div class="flex items-center">
              <div class="font-bold text-[#1E293B] text-[20px] truncate">{{ pluginTitle }}</div>
            </div>
            <div class="text-sm line-clamp-2 text-[#64748B] wrap-break-word">
              {{ pluginDesc }}</div>
            <div class="mt-2">
              <plugin-card-tag :views="pluginData?.views || 0" :download="pluginData?.downloads || 0" :stars="pluginData?.star || 0" :type="pluginData?.type ?? 0"/>
            </div>
          </div>
        </div>
        <div class="download-action shrink-0">
          <pm-button text="前往下载" color="green" size="large" @click="handleDownload">
            <template #icon>
              <PhDownloadSimple weight="bold"/>
            </template>
          </pm-button>
        </div>
      </div>
      <div class="flex min-w-0 flex-col bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-4 md:p-5">
        <div class="item-title">概述</div>
        <div class="select-text plugin-preview">
          <MdPreview :modelValue="pluginContent" previewTheme="default"/>
        </div>
      </div>
      <div class="comments-panel flex min-w-0 flex-col bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-4 md:p-5 gap-4 md:gap-5">
        <div class="item-title">评论</div>
        <comment-input
            v-model="commentText"
            :disabled="!userStore.isLogin"
            :reply-target="replyingComment"
            :loading="submittingComment"
            @submit="handleSubmitComment"
            @cancel-reply="cancelReply"
            @disabled-click="handleCommentDisabledClick"
        />
        <div class="mt-5">
          <comment-box
              :comments="commentData.list"
              :total="commentData.total"
              :loading="commentLoading"
              :plugin-author-id="pluginAuthor?.id || 0"
              :comment-target="commentTarget"
              @reply="handleReply"
              @delete="handleDeleteComment"
          />
          <div class="mt-5 flex items-center justify-center">
            <el-pagination
                v-model:current-page="commentPage"
                layout="prev, pager, next"
                :total="commentData.total"
                :page-size="COMMENT_PAGE_SIZE"
                class="pageBox-diy"
                @current-change="handleCommentPageChange"
            />
          </div>
        </div>
      </div>
    </div>
    <div class="plugin-sidebar w-full lg:w-75 bg-white rounded-[16px] md:rounded-[20px] p-4 md:p-5 shadow-pmbox lg:sticky lg:top-20">
      <div class="item-box">
        <div class="item-title">作者</div>
        <plugin-author-card :author="pluginAuthor"/>
      </div>
      <div class="item-box">
        <div class="item-title">支持框架</div>
        <div class="item-content">
          <div class="flex flex-wrap gap-1.5">
            <el-tooltip v-for="frame in pluginFrameworks" :key="frame.id" :content="frame.name || '框架'" placement="bottom" effect="light">
              <img class="frame-img" :src="assetUrl(frame.icon, DEFAULT_PLUGIN_ICON)" alt="" draggable="false">
            </el-tooltip>
          </div>
        </div>
      </div>
      <div class="item-box">
        <div class="item-title">标签</div>
        <div class="item-content">
          <div class="flex flex-wrap gap-1.5">
            <el-tag v-for="tag in pluginTags" :key="tag.id" type="primary">{{ tag.name }}</el-tag>
          </div>
        </div>
      </div>
      <div class="item-box">
        <div class="item-title">其他</div>
        <div class="item-content flex flex-col gap-3">
          <div class="flex flex-col text-sm">
            <p>发布日期：{{ pluginCreated }}</p>
            <p>修改日期：{{ pluginUpdated }}</p>
          </div>
          <div class="flex justify-center gap-5">
            <pm-button :text="pluginIsStarred ? '已收藏' : '收藏'" color="orange" :loading="starLoading" @click="handleToggleStar">
              <template #icon>
                <PhStar :weight="pluginIsStarred ? 'fill' : 'bold'"/>
              </template>
            </pm-button>
            <pm-button text="分享" color="green" @click="handleCopyShare">
              <template #icon>
                <PhShare/>
              </template>
            </pm-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog
        v-model="downloadDialogVisible"
        title="提示"
        width="min(420px, calc(100vw - 24px))"
        align-center
        class="pm-manage-dialog"
        :close-on-click-modal="false"
    >
      <div class="flex flex-col gap-4 text-sm text-[#475569] leading-6">
        <div class="rounded-[12px] bg-[#F8FAFC] px-4 py-3 text-[#0F172A]">
          <span class="font-medium">密码 / 提取码：</span>{{ pluginData?.url_code || '无' }}
        </div>
        <div>
          下载地址由插件开发者提供，本站不保证其安全性，请自行判断。
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <pm-button text="取消" color="white" @click="closeDownloadDialog" />
          <pm-button text="前往下载" color="blue" :loading="downloadLoading" @click="openDownloadUrl" />
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
.frame-img {
  height: 32px;
  border: 2px solid #FFF;
  border-radius: 20%;
  background-color: #FFF;
}

.logo {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  border: 1px solid #FFF;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.item-box {
  margin-bottom: 24px;
  user-select: none;

  &:last-child {
    margin-bottom: 0;
  }
}

.item-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  color: #1E293B;

  &:before {
    content: '';
    display: inline-block;
    width: 4px;
    height: 20px;
    background-color: #00BAAD;
    margin-right: 8px;
    border-radius: 2px;
  }
}

.item-content {
  color: #475569;

  :deep(.el-tag) {
    border-radius: 10px;
    cursor: pointer;
  }

  :deep(.el-select__wrapper) {
    border-radius: 20px;
  }
}

.plugin-preview {
  min-width: 0;
  overflow-x: auto;

  :deep(.md-editor-preview) {
    padding: 0;
    background: transparent;
    overflow-wrap: anywhere;
  }

  :deep(img) {
    max-width: 100%;
    height: auto;
  }

  :deep(table) {
    display: block;
    max-width: 100%;
    overflow-x: auto;
  }
}

@media (max-width: 767px) {
  .plugin-hero {
    .logo {
      width: 64px;
      height: 64px;
    }

    .download-action,
    .download-action :deep(.pm-button) {
      width: 100%;
    }
  }

  .comments-panel {
    :deep(.comment-entry) {
      min-width: 0;
    }

    :deep(.comment-avatar) {
      width: 40px;
      height: 40px;
    }
  }

  .plugin-sidebar {
    position: static;
  }
}
</style>
