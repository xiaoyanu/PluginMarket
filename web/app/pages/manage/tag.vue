<script setup lang="ts">
import { Delete, Edit } from '@element-plus/icons-vue'
import { PhPlus } from '@phosphor-icons/vue'

definePageMeta({ layout: 'manage' })

interface TagItem { id: number; name: string; created?: string }

const showDialog = ref(false)
const tableData = ref<TagItem[]>([])
const tableLoading = ref(false)
const submitLoading = ref(false)
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const deleteTarget = ref<TagItem | null>(null)
const searchKeyword = ref('')
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const activeId = ref<number | null>(null)
const formRef = ref()
const form = reactive({ name: '' })
const formRules = { name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }] }
const dialogTitle = computed(() => activeId.value ? '编辑标签' : '添加标签')
const formatDate = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}
const handleFilter = () => {
  pagination.page = 1
  loadTableData()
}
const handleReset = () => {
  searchKeyword.value = ''
  pagination.page = 1
  loadTableData()
}
const adjustPageAfterDataChange = () => {
  const maxPage = Math.max(1, Math.ceil(pagination.total / pagination.pageSize))
  if (pagination.page > maxPage) pagination.page = maxPage
}
const loadTableData = async () => {
  tableLoading.value = true
  try {
    const query: Record<string, any> = { page: pagination.page, pageSize: pagination.pageSize }
    if (searchKeyword.value.trim()) query.keywords = searchKeyword.value.trim()
    const res: any = await useApiFetch('/tag/list', { query })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } finally {
    tableLoading.value = false
  }
}
const openDialog = (row?: TagItem) => {
  activeId.value = row?.id || null
  form.name = row?.name || ''
  showDialog.value = true
  nextTick(() => formRef.value?.clearValidate?.())
}
const submitForm = async () => {
  const valid = await formRef.value?.validate?.().catch(() => false)
  if (!valid) return
  submitLoading.value = true
  try {
    await useApiFetch(activeId.value ? `/tag/${activeId.value}` : '/tag', { method: activeId.value ? 'PUT' : 'POST', body: { name: form.name.trim() } })
    ElMessage.success(activeId.value ? '保存成功' : '添加成功')
    showDialog.value = false
    await loadTableData()
    adjustPageAfterDataChange()
  } finally {
    submitLoading.value = false
  }
}
const handleDelete = async (row: TagItem) => {
  deleteTarget.value = row
  deleteDialogVisible.value = true
}
const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await useApiFetch(`/tag/${deleteTarget.value.id}`, { method: 'DELETE' })
    deleteDialogVisible.value = false
    ElMessage.success('删除成功')
    await loadTableData()
    adjustPageAfterDataChange()
  } finally {
    deleteLoading.value = false
  }
}
onMounted(loadTableData)
</script>

<template>
  <div class="flex flex-col gap-5">
    <manage-box title="标签管理" :value="pagination.total">
      <template #header><pm-button text="添加标签" @click="openDialog()"><template #icon><PhPlus /></template></pm-button></template>
      <el-form :inline="true" class="searchForm"><el-form-item><el-input v-model="searchKeyword" placeholder="标签名称关键字" /></el-form-item><el-form-item><div class="flex gap-3 items-center"><pm-button text="筛选" @click="handleFilter" /><pm-button text="重置" color="white" @click="handleReset" /></div></el-form-item></el-form>
      <el-table :data="tableData" class="w-full" v-loading="tableLoading">
        <el-table-column label="ID" width="100" prop="id" />
        <el-table-column label="标签名称" minWidth="200"><template #default="{ row }"><el-tag size="small" effect="plain">{{ row.name }}</el-tag></template></el-table-column>
        <el-table-column label="创建日期" minWidth="150"><template #default="{ row }">{{ formatDate(row.created) }}</template></el-table-column>
        <el-table-column label="操作" width="150"><template #default="{ row }"><el-button :icon="Edit" circle plain type="primary" @click="openDialog(row)" /><el-button :icon="Delete" circle plain type="danger" @click="handleDelete(row)" /></template></el-table-column>
        <template #empty><el-empty description="还没有任何标签哦＞﹏＜" style="user-select: none" /></template>
      </el-table>
      <div class="mt-8 flex items-center justify-center">
        <el-pagination layout="prev, pager, next" :total="pagination.total" :page-size="pagination.pageSize" v-model:current-page="pagination.page" class="pageBox-diy" @current-change="loadTableData" />
      </div>
    </manage-box>
    <el-dialog v-model="showDialog" :title="dialogTitle" width="450px" align-center class="pm-manage-dialog"><div class="py-4 flex flex-col gap-6"><el-form ref="formRef" :model="form" :rules="formRules" label-position="top"><el-form-item label="标签名称" prop="name" required><el-input v-model="form.name" placeholder="请输入标签名称，例如：开发工具" /></el-form-item></el-form></div><template #footer><div class="flex justify-end gap-3"><pm-button text="取消" color="white" @click="showDialog = false" /><pm-button :loading="submitLoading" :text="activeId ? '保存修改' : '确认添加'" @click="submitForm" /></div></template></el-dialog>
    <manage-confirm-dialog v-model="deleteDialogVisible" title="删除标签" :message="`确认删除标签「${deleteTarget?.name || ''}」吗？`" confirm-text="确认删除" confirm-color="red" :loading="deleteLoading" @confirm="confirmDelete" />
  </div>
</template>

<style scoped lang="scss">
:deep(.pm-manage-dialog){border-radius:24px!important;padding:8px;.el-dialog__header{margin-right:0;padding:24px 32px 12px;.el-dialog__title{font-size:20px;font-weight:700;color:#0F172A}}.el-dialog__body{padding:12px 32px}.el-dialog__footer{padding:10px 24px 20px}.el-dialog__headerbtn{top:24px;right:24px}}:deep(.el-input__wrapper){border-radius:7px!important}
</style>
