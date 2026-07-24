<script setup lang="ts">
import {CloseBold, Delete, Select} from "@element-plus/icons-vue";
definePageMeta({
  layout: 'manage'
})
interface PendingPlugin {
  id: number
  name: string
  type: number
  updated: string
}

const tableData = ref<PendingPlugin[]>([])
const tableLoading = ref(false)
const actionPluginId = ref<number | null>(null)
const rejectDialogVisible = ref(false)
const rejectLoading = ref(false)
const rejectFormRef = ref()
const rejectTarget = ref<PendingPlugin | null>(null)
const confirmDialogVisible = ref(false)
const confirmTarget = ref<PendingPlugin | null>(null)
const confirmAction = ref<'approve' | 'delete'>('approve')
const rejectForm = reactive({reason: ''})
const rejectRules = {
  reason: [
    {required: true, message: '请输入拒绝理由', trigger: 'blur'},
    {max: 255, message: '拒绝理由不能超过255个字符', trigger: 'blur'}
  ]
}

const searchForm = reactive({
  keywords: '',
  type: -1
})

const appliedFilters = reactive({
  keywords: '',
  type: -1
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formatDate = (value?: string) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const loadTableData = async () => {
  tableLoading.value = true
  try {
    const query: Record<string, string | number> = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (appliedFilters.keywords) query.keywords = appliedFilters.keywords
    if (appliedFilters.type !== -1) query.type = appliedFilters.type
    const res: any = await useApiFetch('/admin/plugin/pending', {query})
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } finally {
    tableLoading.value = false
  }
}

const handleFilter = () => {
  appliedFilters.keywords = searchForm.keywords.trim()
  appliedFilters.type = searchForm.type
  pagination.page = 1
  loadTableData()
}

const handleReset = () => {
  searchForm.keywords = ''
  searchForm.type = -1
  handleFilter()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadTableData()
}

const refreshAfterAction = async () => {
  if (tableData.value.length === 1 && pagination.page > 1) pagination.page -= 1
  await loadTableData()
}

const handleApprove = async (row: PendingPlugin) => {
  confirmTarget.value = row
  confirmAction.value = 'approve'
  confirmDialogVisible.value = true
}

const submitConfirmAction = async () => {
  if (!confirmTarget.value) return
  const row = confirmTarget.value
  actionPluginId.value = row.id
  try {
    if (confirmAction.value === 'approve') {
      await useApiFetch(`/admin/plugin/${row.id}/approve`, {method: 'PUT'})
      ElMessage.success('审核通过')
    } else {
      await useApiFetch(`/admin/plugin/${row.id}`, {method: 'DELETE'})
      ElMessage.success('删除成功')
    }
    confirmDialogVisible.value = false
    await refreshAfterAction()
  } finally {
    actionPluginId.value = null
  }
}

const handleReject = async (row: PendingPlugin) => {
  rejectTarget.value = row
  rejectForm.reason = ''
  rejectDialogVisible.value = true
  await nextTick()
  rejectFormRef.value?.clearValidate?.()
}

const closeRejectDialog = () => {
  if (rejectLoading.value) return
  rejectDialogVisible.value = false
}

const submitReject = async () => {
  if (!rejectTarget.value || !rejectFormRef.value) return
  const valid = await rejectFormRef.value.validate().catch(() => false)
  const reason = rejectForm.reason.trim()
  if (!valid || !reason) return

  const pluginId = rejectTarget.value.id
  rejectLoading.value = true
  actionPluginId.value = pluginId
  try {
    await useApiFetch(`/admin/plugin/${pluginId}/reject`, {
      method: 'PUT',
      body: {reason}
    })
    rejectDialogVisible.value = false
    ElMessage.success('已拒绝')
    await refreshAfterAction()
  } finally {
    rejectLoading.value = false
    actionPluginId.value = null
  }
}

const handleDelete = async (row: PendingPlugin) => {
  confirmTarget.value = row
  confirmAction.value = 'delete'
  confirmDialogVisible.value = true
}

onMounted(loadTableData)
</script>

<template>
  <div class="flex flex-col gap-5">
    <div>
      <manage-box title="审核插件" :value="pagination.total">
        <div>
          <el-form :inline="true" class="searchForm">
            <el-form-item>
              <el-input
                  v-model="searchForm.keywords"
                  placeholder="关键字"
                  @keydown.enter.prevent="handleFilter"
              />
            </el-form-item>

            <el-form-item label="类型">
              <el-select v-model="searchForm.type" placeholder="请选择">
                <el-option :value="-1" label="不限"/>
                <el-option :value="0" label="免费">
                  <el-tag size="small" type="success">免费</el-tag>
                </el-option>
                <el-option :value="1" label="收费">
                  <el-tag size="small" type="warning">收费</el-tag>
                </el-option>
              </el-select>
            </el-form-item>

            <el-form-item>
              <div class="flex gap-3 items-center">
                <pm-button text="筛选" @click="handleFilter"/>
                <pm-button text="重置" color="white" @click="handleReset"/>
              </div>
            </el-form-item>
          </el-form>


          <el-table :data="tableData" class="w-full" v-loading="tableLoading">
            <el-table-column label="ID">
              <template #default="{row}">
                {{ row.id }}
              </template>
            </el-table-column>
            <el-table-column label="插件名称" minWidth="250">
              <template #default="{row}">
                <el-link :href="`/plugin/`+row.id" target="_blank">
                  {{ row.name }}
                </el-link>
              </template>
            </el-table-column>
            <el-table-column label="类型">
              <template #default="{ row }">
                <el-tag :type="row.type === 1 ? 'warning' : 'success'">
                  {{ row.type === 1 ? '收费' : '免费' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="提交时间" minWidth="150">
              <template #default="{ row }">
                {{ formatDate(row.updated) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <el-button
                    :icon="Select"
                    circle
                    plain
                    type="success"
                    :loading="actionPluginId === row.id"
                    :disabled="actionPluginId !== null"
                    @click="handleApprove(row)"
                />
                <el-button
                    :icon="CloseBold"
                    circle
                    plain
                    type="warning"
                    :disabled="actionPluginId !== null"
                    @click="handleReject(row)"
                />
                <el-button
                    :icon="Delete"
                    circle
                    plain
                    type="danger"
                    :disabled="actionPluginId !== null"
                    @click="handleDelete(row)"
                />
              </template>
            </el-table-column>
            <template #empty>
              <el-empty description="空空的什么也没有＞﹏＜" style="user-select: none"/>
            </template>
          </el-table>

          <div class="mt-8 flex items-center justify-center">
            <el-pagination
                v-model:current-page="pagination.page"
                layout="prev, pager, next"
                :total="pagination.total"
                :page-size="pagination.pageSize"
                class="pageBox-diy"
                @current-change="handlePageChange"
            />
          </div>
        </div>
      </manage-box>
    </div>

    <el-dialog
        v-model="rejectDialogVisible"
        title="审核拒绝"
        width="460px"
        align-center
        class="pm-manage-dialog"
        :close-on-click-modal="!rejectLoading"
        :close-on-press-escape="!rejectLoading"
        :show-close="!rejectLoading"
    >
      <el-form ref="rejectFormRef" :model="rejectForm" :rules="rejectRules" label-position="top">
        <el-form-item prop="reason">
          <template #label>
            <div class="text-[#64748B] font-medium">
              拒绝插件「{{ rejectTarget?.name }}」的理由
            </div>
          </template>
          <el-input
              v-model="rejectForm.reason"
              type="textarea"
              placeholder="请输入拒绝理由"
              maxlength="255"
              show-word-limit
              class="pm-manage-input"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-3 pb-2">
          <pm-button text="取消" color="white" :disabled="rejectLoading" @click="closeRejectDialog"/>
          <pm-button text="确认拒绝" color="orange" :loading="rejectLoading" @click="submitReject"/>
        </div>
      </template>
    </el-dialog>

    <manage-confirm-dialog
        v-model="confirmDialogVisible"
        :title="confirmAction === 'approve' ? '通过审核' : '删除插件'"
        :message="confirmAction === 'approve' ? `确认通过插件「${confirmTarget?.name || ''}」的审核吗？` : `确认删除插件「${confirmTarget?.name || ''}」吗？删除后不可恢复。`"
        :confirm-text="confirmAction === 'approve' ? '确认通过' : '确认删除'"
        :confirm-color="confirmAction === 'approve' ? 'green' : 'red'"
        :loading="actionPluginId === confirmTarget?.id"
        @confirm="submitConfirmAction"
    />
  </div>
</template>

<style scoped lang="scss">
:deep(.el-form-item.is-required > .el-form-item__label::before) {
  display: none;
}

:deep(.pm-manage-input) {
  width: 100%;

  .el-textarea__inner {
    font-size: 14px !important;
    background-color: #F8FAFC;
    border: none;
    box-shadow: inset 0 0 0 1.5px transparent;
    border-radius: 12px;
    padding: 10px 14px !important;
    height: 120px;
    resize: none;
    transition: background-color 0.2s, border-color 0.2s, box-shadow 0.2s;

    &:hover {
      box-shadow: inset 0 0 0 1.5px #E2E8F0;
    }

    &:focus {
      box-shadow: inset 0 0 0 1.5px #00BAAD !important;
    }
  }
}

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
    top: 18px;
    right: 20px;
  }
}
</style>