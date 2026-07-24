<script setup lang="ts">
import { Delete, Edit } from '@element-plus/icons-vue'
import { PhPlus } from '@phosphor-icons/vue'
import dayjs from 'dayjs'
import type { NotificationItem } from '~/composables/api/notification'
import {
  createAdminNotification,
  deleteAdminNotification,
  getAdminNotification,
  getAdminNotifications,
  updateAdminNotification,
  type NotificationForm,
} from '~/composables/api/admin-notification'

definePageMeta({ layout: 'manage' })

const typeOptions = [
  { value: 'system', label: '系统通知' },
  { value: 'custom', label: '自定义通知' },
  { value: 'new_comment', label: '新增评论' },
  { value: 'comment_reply', label: '评论回复' },
  { value: 'plugin_approved', label: '审核通过' },
  { value: 'plugin_rejected', label: '审核拒绝' },
]
const audienceOptions = [
  { value: 'all', label: '所有用户' },
  { value: 'normal', label: '普通用户' },
  { value: 'admin', label: '管理员' },
  { value: 'user', label: '指定用户' },
]
const predefinedIconColors = [
  '#FF8D1A', // 提醒橙
  '#1AB2FE', // 信息蓝
  '#4E89FF', // 回复蓝
  '#22C55E', // 成功绿
  '#EF4444', // 错误红
  '#A855F7', // 自定义通知紫
  '#EC4899', // 自定义通知粉
  '#14B8A6', // 自定义通知青
  '#64748B', // 默认灰蓝
]

const tableData = ref<NotificationItem[]>([])
const loading = ref(false)
const saving = ref(false)
const deleteLoading = ref(false)
const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const editId = ref<number | null>(null)
const deleteTarget = ref<NotificationItem | null>(null)
const formRef = ref()
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const filters = reactive({ keywords: '', type: '', status: -1 })

const form = reactive<NotificationForm>({
  type: 'system',
  audience_type: 'all',
  receiver_id: null,
  title: '',
  content: '',
  icon_name: 'PhBellRinging',
  icon_color: '#FF8D1A',
  action_url: '',
  include_future_users: false,
  status: 1,
  published_at: null,
  expires_at: null,
})

const dialogTitle = computed(() => editId.value ? '编辑通知' : '发布通知')
const formatDate = (value?: string | null) => value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-'
const statusLabel = (status: number) => ({ 0: '已撤回', 1: '已发布', 2: '草稿' }[status] || '未知')
const statusType = (status: number) => ({ 0: 'info', 1: 'success', 2: 'warning' }[status] || 'info')
const audienceLabel = (value: string) => audienceOptions.find(item => item.value === value)?.label || value
const typeLabel = (value: string) => typeOptions.find(item => item.value === value)?.label || value
const typeTagType = (value: string) => ({
  system: 'warning', custom: 'primary', new_comment: 'primary', comment_reply: 'info',
  plugin_approved: 'success', plugin_rejected: 'danger',
}[value] || 'info')
const stripLineBreaks = (value: string) => value.replace(/[\r\n]+/g, ' ')

const loadTable = async () => {
  loading.value = true
  try {
    const query: Record<string, string | number> = { page: pagination.page, pageSize: pagination.pageSize }
    if (filters.keywords.trim()) query.keywords = filters.keywords.trim()
    if (filters.type) query.type = filters.type
    if (filters.status !== -1) query.status = filters.status
    const response = await getAdminNotifications(query)
    tableData.value = response?.data?.list || []
    pagination.total = Number(response?.data?.total || 0)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    type: 'system', audience_type: 'all', receiver_id: null, title: '', content: '',
    icon_name: 'PhBellRinging', icon_color: '#FF8D1A', action_url: '',
    include_future_users: false, status: 1, published_at: null, expires_at: null,
  })
  nextTick(() => formRef.value?.clearValidate?.())
}

const openCreate = () => {
  editId.value = null
  resetForm()
  dialogVisible.value = true
}

const openEdit = async (row: NotificationItem) => {
  editId.value = row.id
  dialogVisible.value = true
  const response = await getAdminNotification(row.id)
  const item = response?.data
  if (!item) return
  Object.assign(form, {
    type: item.type,
    audience_type: item.audience_type,
    receiver_id: item.receiver_id ?? null,
    title: item.title,
    content: stripLineBreaks(item.content),
    icon_name: item.icon_name,
    icon_color: item.icon_color,
    action_url: item.action_url || '',
    include_future_users: Boolean(item.include_future_users),
    status: Number((item as any).status ?? 1),
    published_at: item.published_at || null,
    expires_at: item.expires_at || null,
  })
}

const submitForm = async () => {
  const valid = await formRef.value?.validate?.().catch(() => false)
  if (!valid) return
  if (form.audience_type !== 'user') form.receiver_id = null
  form.content = stripLineBreaks(form.content)
  saving.value = true
  try {
    if (editId.value) await updateAdminNotification(editId.value, form)
    else await createAdminNotification(form)
    ElMessage.success(editId.value ? '通知更新成功' : '通知发布成功')
    dialogVisible.value = false
    await loadTable()
  } finally {
    saving.value = false
  }
}

const requestDelete = (row: NotificationItem) => {
  deleteTarget.value = row
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await deleteAdminNotification(deleteTarget.value.id)
    deleteDialogVisible.value = false
    ElMessage.success('通知删除成功，关联状态记录已清理')
    if (tableData.value.length === 1 && pagination.page > 1) pagination.page--
    await loadTable()
  } finally {
    deleteLoading.value = false
  }
}

const applyFilter = () => {
  pagination.page = 1
  loadTable()
}
const resetFilter = () => {
  Object.assign(filters, { keywords: '', type: '', status: -1 })
  pagination.page = 1
  loadTable()
}

watch(() => form.audience_type, value => {
  if (value !== 'user') form.receiver_id = null
})
onMounted(loadTable)
</script>

<template>
  <div class="flex flex-col gap-5">
    <manage-box title="通知管理" :value="pagination.total">
      <template #header>
        <pm-button text="发布通知" @click="openCreate">
          <template #icon><PhPlus /></template>
        </pm-button>
      </template>

      <el-form :inline="true" class="searchForm" @submit.prevent="applyFilter">
        <el-form-item>
          <el-input v-model="filters.keywords" placeholder="搜索标题或内容" clearable @keyup.enter="applyFilter" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.type" placeholder="请选择">
            <el-option value="" label="不限" />
            <el-option v-for="item in typeOptions" :key="item.value" :value="item.value" :label="item.label">
              <el-tag size="small" :type="typeTagType(item.value)">{{ item.label }}</el-tag>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="请选择">
            <el-option :value="-1" label="不限" />
            <el-option :value="0" label="已撤回"><el-tag size="small" type="info">已撤回</el-tag></el-option>
            <el-option :value="1" label="已发布"><el-tag size="small" type="success">已发布</el-tag></el-option>
            <el-option :value="2" label="草稿"><el-tag size="small" type="warning">草稿</el-tag></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <div class="flex gap-3">
            <pm-button text="筛选" @click="applyFilter" />
            <pm-button text="重置" color="white" @click="resetFilter" />
          </div>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" class="w-full">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="通知标题" min-width="220">
          <template #default="{ row }">
            <div class="flex items-center gap-2 min-w-0">
              <manage-notification-icon class="shrink-0" :name="row.icon_name" :color="row.icon_color" :size="24" />
              <span class="truncate">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="typeTagType(row.type)">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="接收范围" width="120">
          <template #default="{ row }">{{ audienceLabel(row.audience_type) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><el-tag :type="statusType(row.status)">{{ statusLabel(row.status) }}</el-tag></template>
        </el-table-column>
        <el-table-column label="发布时间" min-width="170">
          <template #default="{ row }">{{ formatDate(row.published_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-tooltip content="编辑">
              <el-button :icon="Edit" circle plain @click="openEdit(row)" />
            </el-tooltip>
            <el-tooltip content="删除">
              <el-button :icon="Delete" circle plain type="danger" @click="requestDelete(row)" />
            </el-tooltip>
          </template>
        </el-table-column>
        <template #empty><el-empty description="暂无通知" /></template>
      </el-table>

      <div class="mt-8 flex justify-center">
        <el-pagination
          v-model:current-page="pagination.page"
          layout="prev, pager, next"
          :total="pagination.total"
          :page-size="pagination.pageSize"
          class="pageBox-diy"
          @current-change="loadTable"
        />
      </div>
    </manage-box>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" align-center class="pm-manage-dialog" :close-on-click-modal="!saving" :close-on-press-escape="!saving">
      <el-form ref="formRef" :model="form" label-position="top" class="notification-form">
        <div class="notification-form-grid grid grid-cols-2 gap-4">
          <el-form-item label="通知类型" prop="type" :rules="[{ required: true, message: '请选择通知类型' }]">
            <el-select v-model="form.type" class="w-full"><el-option v-for="item in typeOptions" :key="item.value" :value="item.value" :label="item.label" /></el-select>
          </el-form-item>
          <el-form-item label="接收范围" prop="audience_type" :rules="[{ required: true, message: '请选择接收范围' }]">
            <el-select v-model="form.audience_type" class="w-full"><el-option v-for="item in audienceOptions" :key="item.value" :value="item.value" :label="item.label" /></el-select>
          </el-form-item>
        </div>
        <el-form-item v-if="form.audience_type === 'user'" label="接收用户ID" prop="receiver_id" :rules="[{ required: true, message: '请输入接收用户ID' }]">
          <el-input-number v-model="form.receiver_id" :min="1" class="w-full" />
        </el-form-item>
        <el-form-item label="标题" prop="title" :rules="[{ required: true, message: '请输入通知标题' }]">
          <el-input v-model="form.title" maxlength="255" show-word-limit />
        </el-form-item>
        <el-form-item label="内容" prop="content" :rules="[{ required: true, message: '请输入通知内容' }]">
          <el-input v-model="form.content" maxlength="2000" show-word-limit @input="form.content = stripLineBreaks(form.content)" />
        </el-form-item>
        <div class="notification-form-grid grid grid-cols-2 gap-4">
          <el-form-item prop="icon_name" :rules="[{ required: true, message: '请输入图标名称' }]">
            <template #label>
              <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                <span>图标名称</span>
                <el-link
                  href="https://phosphoricons.com/"
                  target="_blank"
                  type="primary"
                  :underline="false"
                >
                  前往 Phosphor Icons 查找图标
                </el-link>
              </div>
            </template>
            <el-input v-model="form.icon_name" placeholder="例如 PhBellRinging" />
          </el-form-item>
          <el-form-item label="图标颜色" prop="icon_color" :rules="[{ required: true, message: '请输入图标颜色' }]">
            <div class="flex items-center gap-3">
              <el-color-picker
                v-model="form.icon_color"
                color-format="hex"
                :predefine="predefinedIconColors"
              />
              <span class="text-xs text-[#94A3B8]">自定义颜色主要用于自定义通知</span>
            </div>
          </el-form-item>
        </div>
        <el-form-item label="处理地址">
          <el-input v-model="form.action_url" placeholder="例如 /manage/info，可留空" />
        </el-form-item>
        <el-form-item v-if="form.audience_type !== 'user'">
          <el-checkbox v-model="form.include_future_users">允许未来注册用户看到此通知</el-checkbox>
        </el-form-item>
        <el-form-item label="发布状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">立即发布</el-radio>
            <el-radio :value="2">保存草稿</el-radio>
            <el-radio :value="0">撤回</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3 pb-2">
          <pm-button text="取消" color="white" :disabled="saving" @click="dialogVisible = false" />
          <pm-button text="保存" :loading="saving" @click="submitForm" />
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="deleteDialogVisible" title="删除通知" width="460px" align-center class="pm-manage-dialog">
      <div class="text-[#475569]">确定删除通知“{{ deleteTarget?.title }}”吗？删除后会同时清理所有用户的已读/隐藏状态记录。</div>
      <template #footer>
        <div class="flex justify-end gap-3 pb-2">
          <pm-button text="取消" color="white" :disabled="deleteLoading" @click="deleteDialogVisible = false" />
          <pm-button text="删除" color="red" :loading="deleteLoading" @click="confirmDelete" />
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
@media (max-width: 767px) {
  .notification-form-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  :deep(.el-radio-group) {
    display: flex;
    flex-wrap: wrap;
  }
}
</style>
