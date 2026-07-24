<script setup lang="ts">
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { PhPlus, PhSquaresFour } from '@phosphor-icons/vue'
import type { UploadFile, UploadRawFile } from 'element-plus'

definePageMeta({ layout: 'manage' })

interface FrameItem { id: number; name: string; icon: string; created?: string }

const assetUrl = useAssetUrl()
const formatDate = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}

const showDialog = ref(false)
const tableData = ref<FrameItem[]>([])
const tableLoading = ref(false)
const submitLoading = ref(false)
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const deleteTarget = ref<FrameItem | null>(null)
const searchKeyword = ref('')
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const activeId = ref<number | null>(null)
const iconFile = ref<UploadRawFile | null>(null)
const imgUrl = ref('')
const formRef = ref()
const form = reactive({ name: '' })
const formRules = { name: [{ required: true, message: '请输入框架名称', trigger: 'blur' }] }
const dialogTitle = computed(() => activeId.value ? '编辑框架' : '添加框架')
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
    const res: any = await useApiFetch('/frame/list', { query })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } finally {
    tableLoading.value = false
  }
}

const openDialog = (row?: FrameItem) => {
  activeId.value = row?.id || null
  form.name = row?.name || ''
  iconFile.value = null
  imgUrl.value = row?.icon ? assetUrl(row.icon) : ''
  showDialog.value = true
  nextTick(() => formRef.value?.clearValidate?.())
}

const handleIconChange = (file: UploadFile) => {
  if (!file.raw) return
  iconFile.value = file.raw
  imgUrl.value = URL.createObjectURL(file.raw)
}

const submitForm = async () => {
  const valid = await formRef.value?.validate?.().catch(() => false)
  if (!valid) return
  if (!activeId.value && !iconFile.value) return ElMessage.warning('请上传框架图标')
  submitLoading.value = true
  try {
    const body = new FormData()
    body.append('name', form.name.trim())
    if (iconFile.value) body.append('icon', iconFile.value)
    await useApiFetch(activeId.value ? `/frame/${activeId.value}` : '/frame', { method: activeId.value ? 'PUT' : 'POST', body })
    ElMessage.success(activeId.value ? '保存成功' : '添加成功')
    showDialog.value = false
    await loadTableData()
    adjustPageAfterDataChange()
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: FrameItem) => {
  deleteTarget.value = row
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await useApiFetch(`/frame/${deleteTarget.value.id}`, { method: 'DELETE' })
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
    <manage-box title="框架管理" :value="pagination.total">
      <template #header>
        <pm-button text="添加框架" @click="openDialog()">
          <template #icon><PhPlus /></template>
        </pm-button>
      </template>

      <el-form :inline="true" class="searchForm">
        <el-form-item><el-input v-model="searchKeyword" placeholder="框架名称关键字" /></el-form-item>
        <el-form-item>
          <div class="flex gap-3 items-center">
            <pm-button text="筛选" @click="handleFilter" />
            <pm-button text="重置" color="white" @click="handleReset" />
          </div>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" class="w-full" v-loading="tableLoading">
        <el-table-column label="ID" width="80" prop="id" />
        <el-table-column label="图标" width="100">
          <template #default="{ row }">
            <el-avatar :size="40" :src="assetUrl(row.icon)" shape="square"><el-icon><PhSquaresFour /></el-icon></el-avatar>
          </template>
        </el-table-column>
        <el-table-column label="框架名称" minWidth="200" prop="name" />
        <el-table-column label="创建日期" minWidth="150"><template #default="{ row }">{{ formatDate(row.created) }}</template></el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button :icon="Edit" circle plain type="primary" @click="openDialog(row)" />
            <el-button :icon="Delete" circle plain type="danger" @click="handleDelete(row)" />
          </template>
        </el-table-column>
        <template #empty><el-empty description="空空的什么也没有＞﹏＜" style="user-select: none" /></template>
      </el-table>
      <div class="mt-8 flex items-center justify-center">
        <el-pagination
            layout="prev, pager, next"
            :total="pagination.total"
            :page-size="pagination.pageSize"
            v-model:current-page="pagination.page"
            class="pageBox-diy"
            @current-change="loadTableData"
        />
      </div>
    </manage-box>

    <el-dialog v-model="showDialog" :title="dialogTitle" width="500px" align-center class="pm-manage-dialog">
      <div class="py-4 flex flex-col gap-6">
        <el-form ref="formRef" :model="form" :rules="formRules" label-position="top">
          <el-form-item label="框架名称" prop="name" required><el-input v-model="form.name" placeholder="请输入框架名称，例如：LiteLoaderQQNT" /></el-form-item>
          <el-form-item label="框架图标" required>
            <div class="flex items-center gap-4">
              <el-upload class="avatar-uploader" action="#" :show-file-list="false" :auto-upload="false" accept="image/png,image/jpeg,image/webp" :on-change="handleIconChange">
                <img v-if="imgUrl" :src="imgUrl" class="avatar" alt="" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
              <div class="text-xs text-slate-400"><p>建议尺寸: 200x200</p><p>支持 JPG, PNG, WEBP</p><p v-if="activeId">不选择新图标则保留原图标</p></div>
            </div>
          </el-form-item>
        </el-form>
      </div>
      <template #footer><div class="flex justify-end gap-3"><pm-button text="取消" color="white" @click="showDialog = false" /><pm-button :loading="submitLoading" :text="activeId ? '保存修改' : '确认添加'" @click="submitForm" /></div></template>
    </el-dialog>
    <manage-confirm-dialog v-model="deleteDialogVisible" title="删除框架" :message="`确认删除框架「${deleteTarget?.name || ''}」吗？`" confirm-text="确认删除" confirm-color="red" :loading="deleteLoading" @confirm="confirmDelete" />
  </div>
</template>

<style scoped lang="scss">
:deep(.pm-manage-dialog){border-radius:24px!important;padding:8px;.el-dialog__header{margin-right:0;padding:24px 32px 12px;.el-dialog__title{font-size:20px;font-weight:700;color:#0F172A}}.el-dialog__body{padding:12px 32px}.el-dialog__footer{padding:10px 24px 20px}.el-dialog__headerbtn{top:24px;right:24px}}
:deep(.el-select__wrapper),:deep(.el-input__wrapper){border-radius:7px!important}.avatar-uploader{.avatar{width:80px;height:80px;display:block;object-fit:contain;border-radius:12px}:deep(.el-upload){border:2px dashed #E2E8F0;border-radius:16px;cursor:pointer;position:relative;overflow:hidden;transition:all .2s;&:hover{border-color:var(--el-color-primary);background-color:#F8FAFC}}}:deep(.el-icon).avatar-uploader-icon{font-size:24px;color:#94A3B8;width:80px;height:80px;text-align:center}
</style>
