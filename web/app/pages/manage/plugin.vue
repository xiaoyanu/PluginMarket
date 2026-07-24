<script setup lang="ts">
import {PhStackPlus} from "@phosphor-icons/vue";
import {Delete, Edit, Plus, WarningFilled} from "@element-plus/icons-vue";
import type {UploadFile, UploadRawFile, FormItemRule} from 'element-plus'
import {getMyPlugins} from "~/composables/api/public";
import {isValidPluginDownloadUrl} from "~/utils/plugin-download-url";


definePageMeta({
  layout: 'manage'
})

interface FrameOption {
  id: number
  name: string
  icon?: string
}

interface TagOption {
  id: number
  name: string
}

const assetUrl = useAssetUrl()

const imgUrl = ref('');
const iconFile = ref<UploadRawFile | null>(null)
const showDrawer = ref(false);
const submitLoading = ref(false)
const drawerMode = ref<'publish' | 'edit'>('publish')
const activePluginId = ref<number | null>(null)
const frameOptions = ref<FrameOption[]>([])
const tagOptions = ref<TagOption[]>([])
const drawerTitle = computed(() => drawerMode.value === 'edit' ? '编辑插件' : '发布插件')
const submitText = computed(() => drawerMode.value === 'edit' ? '保存修改' : '发布插件')

const form = reactive({
  title: '',
  description: '',
  type: 0,
  frames: [] as number[],
  tags: [] as number[],
  downloadUrl: '',
  password: '',
  content: ''
});

const formRules = reactive({
  title: [{required: true, message: '请输入插件名称', trigger: 'blur'}],
  description: [{required: true, message: '请输入插件简述', trigger: 'blur'}],
  type: [{required: true, message: '请选择插件类型', trigger: 'change'}],
  frames: [{required: true, message: '请选择适用框架', trigger: 'change'}],
  downloadUrl: [
    {required: true, message: '请输入下载链接', trigger: 'blur'},
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (!isValidPluginDownloadUrl(value)) {
          callback(new Error('请输入有效的 HTTP 或 HTTPS 下载链接'))
          return
        }
        callback()
      },
      trigger: ['blur', 'change']
    } as FormItemRule
  ],
  content: [{required: true, message: '请输入正文内容', trigger: 'blur'}]
});

const formRef = ref();

const resetPublishForm = () => {
  activePluginId.value = null
  form.title = ''
  form.description = ''
  form.type = 0
  form.frames = []
  form.tags = []
  form.downloadUrl = ''
  form.password = ''
  form.content = ''
  iconFile.value = null
  imgUrl.value = ''
  nextTick(() => formRef.value?.clearValidate?.())
}

const openPublishDrawer = async () => {
  drawerMode.value = 'publish'
  resetPublishForm()
  showDrawer.value = true
  await loadPublishOptions()
}

const openEditDrawer = async (row: any) => {
  drawerMode.value = 'edit'
  resetPublishForm()
  activePluginId.value = row.id
  showDrawer.value = true
  await loadPublishOptions()
  const res: any = await useApiFetch(`/plugin/${row.id}`)
  const detail = res.data || {}
  form.title = detail.name || ''
  form.description = detail.desc_text || ''
  form.type = detail.type ?? 0
  form.frames = (detail.frameworks || []).map((item: FrameOption) => item.id)
  form.tags = (detail.tags || []).map((item: TagOption) => item.id)
  form.downloadUrl = detail.url || ''
  form.password = detail.url_code || ''
  form.content = detail.content || ''
  imgUrl.value = assetUrl(detail.icon)
  nextTick(() => formRef.value?.clearValidate?.())
}

const loadPublishOptions = async () => {
  const [frameRes, tagRes]: any[] = await Promise.all([
    useApiFetch('/frame/list', { query: { page: 1, pageSize: 100 } }),
    useApiFetch('/tag/list', { query: { page: 1, pageSize: 100 } })
  ])
  frameOptions.value = frameRes.data?.list || []
  tagOptions.value = tagRes.data?.list || []
}

const handleIconChange = (file: UploadFile) => {
  if (!file.raw) return
  iconFile.value = file.raw
  imgUrl.value = URL.createObjectURL(file.raw)
}

const submitForm = async () => {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitLoading.value = true
  try {
    const body = new FormData()
    body.append('name', form.title.trim())
    body.append('desc_text', form.description.trim())
    body.append('type', String(form.type))
    body.append('content', form.content)
    body.append('url', form.downloadUrl.trim())
    body.append('url_code', form.password.trim())
    if (drawerMode.value === 'edit') {
      body.append('frameworkIds', JSON.stringify(form.frames))
      body.append('tagIds', JSON.stringify(form.tags))
    } else {
      form.frames.forEach(id => body.append('frameworkIds', String(id)))
      form.tags.forEach(id => body.append('tagIds', String(id)))
    }
    if (iconFile.value) body.append('icon', iconFile.value)
    if (drawerMode.value === 'edit' && activePluginId.value) {
      await useApiFetch(`/plugin/${activePluginId.value}`, {method: 'PUT', body})
      ElMessage.success('保存成功，插件已重新提交审核')
    } else {
      await useApiFetch('/plugin/publish', {method: 'POST', body})
      ElMessage.success('发布成功，等待审核')
    }
    showDrawer.value = false
    pagination.page = 1
    await loadTableData()
  } catch (error) {
    showApiError(error, '插件提交失败')
  } finally {
    submitLoading.value = false
  }
};

const tableData = ref<any[]>([])
const tableLoading = ref(false)
const rejectDialogVisible = ref(false)
const rejectPluginName = ref('')
const rejectReason = ref('')
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const deleteTarget = ref<any>(null)
const searchForm = reactive({
  keywords: '',
  type: -1,
  status: -1
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
    const query: Record<string, any> = {}
    if (searchForm.keywords.trim()) query.keywords = searchForm.keywords.trim()
    if (searchForm.type !== -1) query.type = searchForm.type
    if (searchForm.status !== -1) query.status = searchForm.status
    const res: any = await getMyPlugins(pagination.pageSize, pagination.page, query)
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } finally {
    tableLoading.value = false
  }
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadTableData()
}

const handleFilter = () => {
  pagination.page = 1
  loadTableData()
}

const handleReset = () => {
  searchForm.keywords = ''
  searchForm.type = -1
  searchForm.status = -1
  pagination.page = 1
  loadTableData()
}

const handleDelete = async (row: any) => {
  deleteTarget.value = row
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await useApiFetch(`/plugin/${deleteTarget.value.id}`, { method: 'DELETE' })
    deleteDialogVisible.value = false
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && pagination.page > 1) pagination.page -= 1
    await loadTableData()
  } finally {
    deleteLoading.value = false
  }
}

const showRejectReason = (row: any) => {
  rejectPluginName.value = row.name || '未命名插件'
  rejectReason.value = row.reject_msg?.trim() || '管理员未填写拒绝理由'
  rejectDialogVisible.value = true
}

onMounted(() => loadTableData())
</script>

<template>
  <div class="flex flex-col gap-5">
    <div>
      <manage-box title="我的插件" :value="pagination.total">
        <template #header>
          <pm-button text="发布插件" @click="openPublishDrawer">
            <template #icon>
              <PhStackPlus/>
            </template>
          </pm-button>
        </template>
        <div>
          <el-form :inline="true" class="searchForm">
            <el-form-item>
              <el-input
                  v-model="searchForm.keywords"
                  placeholder="关键字"
                  @keydown.enter="handleFilter"
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

            <el-form-item label="状态">
              <el-select v-model="searchForm.status" placeholder="请选择">
                <el-option :value="-1" label="不限"/>
                <el-option :value="0" label="通过">
                  <el-tag size="small" type="success">通过</el-tag>
                </el-option>

                <el-option :value="1" label="拒绝">
                  <el-tag size="small" type="danger">拒绝</el-tag>
                </el-option>

                <el-option :value="2" label="审核中">
                  <el-tag size="small" type="warning">审核中</el-tag>
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
            <el-table-column label="创建日期" minWidth="150">
              <template #default="{ row }">
                {{ formatDate(row.created) }}
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 2 ? 'warning' : row.status === 1 ? 'danger' : 'success'">
                  {{ row.status === 2 ? '审核中' : row.status === 1 ? '拒绝' : '通过' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <el-button
                    :icon="Edit"
                    circle
                    plain
                    type="primary"
                    @click="openEditDrawer(row)"
                />
                <el-button
                    :icon="Delete"
                    circle
                    plain
                    type="danger"
                    @click="handleDelete(row)"
                />
                <el-button
                    v-if="row.status == 1"
                    :icon="WarningFilled"
                    circle
                    plain
                    type="warning"
                    @click="showRejectReason(row)"
                />
              </template>
            </el-table-column>
            <template #empty>
              <el-empty description="空空的什么也没有＞﹏＜" style="user-select: none"/>
            </template>
          </el-table>

          <div class="mt-8 flex items-center justify-center">
            <el-pagination
                layout="prev, pager, next"
                :total="pagination.total"
                :page-size="pagination.pageSize"
                v-model:current-page="pagination.page"
                class="pageBox-diy"
                @current-change="handlePageChange"
            />
          </div>
        </div>
      </manage-box>
    </div>

    <el-dialog
        v-model="rejectDialogVisible"
        title="审核拒绝理由"
        width="460px"
        align-center
        class="pm-manage-dialog"
    >
      <div class="flex flex-col gap-3">
        <div class="text-sm font-medium text-[#64748B]">
          插件「{{ rejectPluginName }}」未通过审核，理由如下：
        </div>
        <div class="min-h-24 rounded-xl bg-[#F8FAFC] px-4 py-3 text-sm leading-6 text-[#334155] whitespace-pre-wrap break-words select-text">
          {{ rejectReason }}
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <pm-button text="关闭" color="white" @click="rejectDialogVisible = false"/>
        </div>
      </template>
    </el-dialog>

    <manage-confirm-dialog
        v-model="deleteDialogVisible"
        title="删除插件"
        :message="`确认删除插件「${deleteTarget?.name || ''}」吗？删除后不可恢复。`"
        confirm-text="确认删除"
        confirm-color="red"
        :loading="deleteLoading"
        @confirm="confirmDelete"
    />

    <!-- 抽屉 -->
    <el-drawer
        v-model="showDrawer"
        :destroy-on-close="true"
        :resizable="true"
        size="80%"
        :title="drawerTitle"
        class="rounded-tl-[20px] rounded-tr-[20px] m-auto custom-btt-drawer"
        direction="btt"
        :close-on-press-escape="false"
    >
      <div class="plugin-drawer-content m-auto max-w-[90%]">
        <el-form ref="formRef" :model="form" :rules="formRules" label-width="auto">
          <el-form-item label="">
            <el-alert type="warning" :closable="false">
              <template #title>
                请务必阅读并遵守
                <el-popover placement="right-start" width="auto" popper-class="plugin-rules-popover">
                  <template #reference>
                    <span class="font-bold text-red-400 cursor-pointer">《插件规范》</span>
                  </template>
                  <template #default>
                    <div class="flex flex-col text-sm">
                      <p>1. 禁止发布任何违法违规插件，包括但不限于黄、赌、毒</p>
                      <p>2. 禁止发布任何含有恶意病毒的插件/应用</p>
                      <p>3. 插件禁止强制引流，例如：必须要加群才可以使用</p>
                      <p>4. 应明确说明插件收费/免费，避免含糊说明：半免费、半收费</p>
                      <p>5. 插件正文介绍至少满足3行，避免简略介绍</p>
                      <p>6. 禁止抄袭、盗用他人插件/代码，应尊重原创</p>
                      <p>7. 禁止未经用户同意收集个人隐私数据</p>
                      <p>8. 插件名称不得冒充官方或使用误导性名称</p>
                    </div>
                  </template>
                </el-popover>
              </template>
            </el-alert>
          </el-form-item>
          <el-form-item label="插件名称" prop="title" required>
            <el-input v-model="form.title" placeholder="请输入插件名称，例如：智能助手"/>
          </el-form-item>
          <el-form-item label="插件简述" prop="description" required>
            <el-input v-model="form.description" placeholder="请输入插件简述，例如：提供简单的天气查询功能"/>
          </el-form-item>
          <el-form-item label="适用框架" prop="frames" required>
            <el-select
                v-model="form.frames"
                multiple
                placeholder="请选择适用框架（可多选）"
                :loading="!frameOptions.length"
            >
              <el-option
                  v-for="item in frameOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
              >
                <div class="flex items-center gap-2">
                  <img v-if="item.icon" :src="assetUrl(item.icon)" class="w-5 h-5 object-contain" alt="" />
                  <span>{{ item.name }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="标签">
            <el-select
                v-model="form.tags"
                multiple
                placeholder="请选择标签（可多选）"
                :loading="!tagOptions.length"
            >
              <el-option
                  v-for="item in tagOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="插件类型" prop="type" required>
            <el-radio-group v-model="form.type">
              <el-radio-button :value="0" label="免费"/>
              <el-radio-button :value="1" label="收费"/>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="插件图标">
            <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :auto-upload="false"
                accept="image/png,image/jpeg,image/webp"
                :on-change="handleIconChange"
            >
              <img v-if="imgUrl" :src="imgUrl" class="avatar" alt=""/>
              <el-icon v-else class="avatar-uploader-icon">
                <Plus/>
              </el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item label="下载链接" prop="downloadUrl" required>
            <el-input
                v-model="form.downloadUrl"
                placeholder="请输入插件下载链接，例如：https://example.com/plugin.zip"
            />
          </el-form-item>
          <el-form-item label="提取密码">
            <el-input
                v-model="form.password"
                placeholder="下载所需 密码/提取码（没有可以不填写）"
            />
          </el-form-item>
          <el-form-item label="正文内容" prop="content" required>
            <div class="flex flex-col flex-1 gap-2">
              <el-alert
                  title="请使用Markdown语法书写，在当下环境非常建议你花一点点时间去了解一下Markdown语法（一点也不难）"
                  type="primary"/>
              <manage-editor v-model="form.content"/>
            </div>
          </el-form-item>
        </el-form>
        <div class="flex justify-center">
          <pm-button size="large" :text="submitText" :loading="submitLoading" @click="submitForm"/>
        </div>
      </div>

    </el-drawer>
  </div>
</template>

<style scoped lang="scss">
:deep(.custom-btt-drawer) {
  width: 80%;
}

.avatar-uploader {
  .avatar {
    width: 100px;
    height: 100px;
    display: block;
  }

  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 10px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);

    &:hover {
      border-color: var(--el-color-primary);
    }
  }
}

:deep(.el-icon) {
  &.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 100px;
    height: 100px;
    text-align: center;
  }
}

@media (max-width: 767px) {
  :deep(.custom-btt-drawer) {
    width: calc(100vw - 16px) !important;
  }

  .plugin-drawer-content {
    width: 100%;
    max-width: 100%;
  }

  :deep(.md-editor) {
    width: 100%;
    min-width: 0;
  }

  :deep(.md-editor-toolbar-wrapper) {
    overflow-x: auto;
  }
}
</style>
