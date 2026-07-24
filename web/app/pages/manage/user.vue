<script setup lang="ts">
import { CollectionTag, Delete, Key, Operation, Plus, View } from '@element-plus/icons-vue'
import { PhArticle, PhClock, PhEnvelope, PhFingerprint, PhLock, PhShieldCheck, PhUser } from '@phosphor-icons/vue'
import dayjs from 'dayjs'
import { DEFAULT_AVATAR } from '~/config'

definePageMeta({ layout: 'manage' })

interface TitleItem { id: number; name: string; icon?: string }
interface UserItem {
  id: number
  username: string
  nick: string
  email: string
  userdesc: string
  avatar: string
  power: number
  is_delete: boolean
  created: string
  updated: string
  titles?: TitleItem[]
}
interface ApiResponse<T> { code: number; msg: string; data?: T }
interface PageData<T> { list: T[]; total: number }

const userStore = useUserStore()
const assetUrl = useAssetUrl()
const loading = ref(false)
const tableData = ref<UserItem[]>([])
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const filters = reactive({ keywords: '', power: -1, isDelete: -1 })

const showViewDialog = ref(false)
const detailLoading = ref(false)
const userInfo = ref<UserItem | null>(null)

const showTitleDialog = ref(false)
const titleSaving = ref(false)
const titleTargetUser = ref<UserItem | null>(null)
const selectedTitleId = ref<number | null>(null)
const titleList = ref<TitleItem[]>([])
const titleData = ref<TitleItem[]>([])

const showResetDialog = ref(false)
const resetPassword = ref('')
const resetConfirmPassword = ref('')
const resetTargetUser = ref<UserItem | null>(null)
const resetSaving = ref(false)

const showPowerDialog = ref(false)
const currentPower = ref(0)
const powerTargetUser = ref<UserItem | null>(null)
const powerSaving = ref(false)

const deleteVisible = ref(false)
const deleteTarget = ref<UserItem | null>(null)
const deleteLoading = ref(false)

const formatDate = (value: string) => value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-'
const currentUserId = computed(() => Number(userStore.userInfo?.id || 0))

const loadUsers = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: String(pagination.page),
      pageSize: String(pagination.pageSize),
    })
    if (filters.keywords.trim()) params.set('keywords', filters.keywords.trim())
    if (filters.power !== -1) params.set('power', String(filters.power))
    if (filters.isDelete !== -1) params.set('isDelete', filters.isDelete === 1 ? 'true' : 'false')
    const res = await useApiFetch<ApiResponse<PageData<UserItem>>>(`/admin/user/list?${params}`)
    tableData.value = res.data?.list || []
    pagination.total = Number(res.data?.total || 0)
  } finally {
    loading.value = false
  }
}

const applyFilter = () => {
  pagination.page = 1
  loadUsers()
}
const resetFilter = () => {
  filters.keywords = ''
  filters.power = -1
  filters.isDelete = -1
  pagination.page = 1
  loadUsers()
}
const changePage = (page: number) => {
  pagination.page = page
  loadUsers()
}

const loadUserDetail = async (id: number) => {
  detailLoading.value = true
  try {
    const res = await useApiFetch<ApiResponse<UserItem>>(`/admin/user/${id}`)
    userInfo.value = res.data || null
    return res.data || null
  } finally {
    detailLoading.value = false
  }
}
const openView = async (row: UserItem) => {
  showViewDialog.value = true
  await loadUserDetail(row.id)
}

const loadTitles = async () => {
  const res = await useApiFetch<ApiResponse<PageData<TitleItem>>>('/title/list?page=1&pageSize=100')
  titleList.value = res.data?.list || []
}
const openTitles = async (row: UserItem) => {
  titleTargetUser.value = row
  selectedTitleId.value = null
  showTitleDialog.value = true
  const [, detail] = await Promise.all([loadTitles(), loadUserDetail(row.id)])
  titleData.value = [...(detail?.titles || [])]
}
const handleAddTitle = () => {
  if (!selectedTitleId.value || titleData.value.some(item => item.id === selectedTitleId.value)) return
  const title = titleList.value.find(item => item.id === selectedTitleId.value)
  if (title) titleData.value.push(title)
  selectedTitleId.value = null
}
const handleDeleteTitle = (index: number) => titleData.value.splice(index, 1)
const saveTitles = async () => {
  if (!titleTargetUser.value) return
  titleSaving.value = true
  try {
    await useApiFetch(`/admin/user/${titleTargetUser.value.id}/titles`, {
      method: 'PUT', body: { titleIds: titleData.value.map(item => item.id) },
    })
    showTitleDialog.value = false
    ElMessage.success('用户称号已更新')
  } finally {
    titleSaving.value = false
  }
}

const openReset = (row: UserItem) => {
  resetTargetUser.value = row
  resetPassword.value = ''
  resetConfirmPassword.value = ''
  showResetDialog.value = true
}
const saveResetPassword = async () => {
  if (!resetTargetUser.value) return
  if (resetPassword.value.length < 6 || resetPassword.value.length > 32) return ElMessage.warning('密码长度为 6-32 位')
  if (resetPassword.value !== resetConfirmPassword.value) return ElMessage.warning('两次输入的新密码不一致')
  resetSaving.value = true
  try {
    await useApiFetch(`/admin/user/${resetTargetUser.value.id}/reset-password`, {
      method: 'PUT', body: { newPassword: resetPassword.value },
    })
    showResetDialog.value = false
    ElMessage.success('密码重置成功')
  } finally {
    resetSaving.value = false
  }
}

const openPower = (row: UserItem) => {
  if (row.id === currentUserId.value) return ElMessage.warning('不能修改当前登录账号的权限')
  powerTargetUser.value = row
  currentPower.value = row.power
  showPowerDialog.value = true
}
const savePower = async () => {
  if (!powerTargetUser.value) return
  powerSaving.value = true
  try {
    await useApiFetch(`/admin/user/${powerTargetUser.value.id}/role`, {
      method: 'PUT', body: { power: currentPower.value },
    })
    showPowerDialog.value = false
    ElMessage.success('用户权限已更新')
    await loadUsers()
  } finally {
    powerSaving.value = false
  }
}

const requestDelete = (row: UserItem) => {
  if (row.id === currentUserId.value) return ElMessage.warning('不能删除当前登录账号')
  deleteTarget.value = row
  deleteVisible.value = true
}
const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await useApiFetch(`/admin/user/${deleteTarget.value.id}`, { method: 'DELETE' })
    deleteVisible.value = false
    if (tableData.value.length === 1 && pagination.page > 1) pagination.page--
    ElMessage.success('用户已删除')
    await loadUsers()
  } finally {
    deleteLoading.value = false
  }
}

onMounted(loadUsers)
</script>

<template>
  <div class="flex flex-col gap-5">
    <manage-box title="用户管理" :value="pagination.total">
      <el-form :inline="true" class="searchForm" @submit.prevent="applyFilter">
        <el-form-item>
          <el-input v-model="filters.keywords" placeholder="搜索账号或昵称" clearable @keyup.enter="applyFilter"/>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.power" placeholder="请选择">
            <el-option :value="-1" label="不限"/>
            <el-option :value="0" label="用户">
              <el-tag size="small" type="success">用户</el-tag>
            </el-option>
            <el-option :value="1" label="管理员">
              <el-tag size="small" type="danger">管理员</el-tag>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.isDelete" placeholder="请选择">
            <el-option :value="-1" label="不限"/>
            <el-option :value="0" label="正常">
              <el-tag size="small" type="success">正常</el-tag>
            </el-option>
            <el-option :value="1" label="已注销">
              <el-tag size="small" type="info">已注销</el-tag>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <div class="flex gap-3 items-center">
            <pm-button text="筛选" @click="applyFilter"/>
            <pm-button text="重置" color="white" @click="resetFilter"/>
          </div>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" class="w-full">
        <el-table-column prop="id" label="ID" width="80"/>
        <el-table-column label="账号" min-width="130">
          <template #default="{ row }"><el-link :href="`/home/${row.id}`" target="_blank">{{ row.username }}</el-link></template>
        </el-table-column>
        <el-table-column prop="nick" label="昵称" min-width="120"/>
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><el-tag :type="row.is_delete ? 'info' : 'success'">{{ row.is_delete ? '已注销' : '正常' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="权限" width="110">
          <template #default="{ row }"><el-tag :type="row.power === 1 ? 'danger' : 'success'">{{ row.power === 1 ? '管理员' : '用户' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="注册时间" min-width="170">
          <template #default="{ row }">{{ formatDate(row.created) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-tooltip content="查看信息"><el-button :icon="View" circle plain @click="openView(row)"/></el-tooltip>
            <el-tooltip content="称号管理"><el-button :icon="CollectionTag" circle plain @click="openTitles(row)"/></el-tooltip>
            <el-tooltip content="重置密码"><el-button :icon="Key" circle plain @click="openReset(row)"/></el-tooltip>
            <el-tooltip content="设置权限"><el-button :icon="Operation" circle plain :disabled="row.id === currentUserId" @click="openPower(row)"/></el-tooltip>
            <el-tooltip :content="row.is_delete ? '用户已注销' : '删除用户'"><el-button :icon="Delete" circle plain type="danger" :disabled="row.id === currentUserId || row.is_delete" @click="requestDelete(row)"/></el-tooltip>
          </template>
        </el-table-column>
        <template #empty><el-empty description="空空的什么也没有＞﹏＜"/></template>
      </el-table>

      <div class="mt-8 flex items-center justify-center">
        <el-pagination
            layout="prev, pager, next"
            :total="pagination.total"
            :page-size="pagination.pageSize"
            v-model:current-page="pagination.page"
            class="pageBox-diy"
            @current-change="changePage"
        />
      </div>
    </manage-box>

    <el-dialog v-model="showViewDialog" title="用户信息详情" width="800px" align-center class="pm-manage-dialog">
      <div v-loading="detailLoading" class="user-detail-container flex gap-8 py-4 min-h-60">
        <template v-if="userInfo">
          <div class="shrink-0 flex flex-col items-center gap-4">
            <div class="p-1 rounded-full border-2 border-dashed border-[#00BAAD]/30">
              <img :src="assetUrl(userInfo.avatar, DEFAULT_AVATAR)" class="rounded-full h-30 w-30 object-cover" alt="用户头像">
            </div>
            <div class="text-center"><h3 class="text-xl font-bold text-[#1E293B]">{{ userInfo.nick || userInfo.username }}</h3><p class="text-sm text-[#64748B] flex items-center justify-center gap-1 mt-1"><PhFingerprint :size="14"/>ID: {{ userInfo.id }}</p></div>
            <el-tag :type="userInfo.power === 1 ? 'danger' : 'success'">{{ userInfo.power === 1 ? '管理员' : '用户' }}</el-tag>
          </div>
          <div class="flex-1 flex flex-col gap-5">
            <div class="grid grid-cols-2 gap-4">
              <div class="info-item"><label><PhUser :size="16"/>账号名称</label><div class="value">{{ userInfo.username }}</div></div>
              <div class="info-item"><label><PhEnvelope :size="16"/>电子邮箱</label><div class="value">{{ userInfo.email || '未绑定邮箱' }}</div></div>
              <div class="info-item"><label><PhClock :size="16"/>注册时间</label><div class="value">{{ formatDate(userInfo.created) }}</div></div>
              <div class="info-item"><label><PhClock :size="16"/>最后更新</label><div class="value">{{ formatDate(userInfo.updated) }}</div></div>
            </div>
            <div class="info-item"><label><PhArticle :size="16"/>个人简介</label><div class="value bio-box">{{ userInfo.userdesc || '暂无简介' }}</div></div>
          </div>
        </template>
      </div>
    </el-dialog>

    <el-dialog v-model="showTitleDialog" title="称号管理" width="500px" align-center class="pm-manage-dialog" :close-on-click-modal="!titleSaving" :close-on-press-escape="!titleSaving">
      <div v-loading="detailLoading" class="flex flex-col gap-4 py-2">
        <div class="text-sm text-[#64748B]">正在管理 <strong class="text-[#00BAAD]">{{ titleTargetUser?.username }}</strong> 的称号</div>
        <div class="flex gap-2"><el-select v-model="selectedTitleId" placeholder="请选择或搜索称号" filterable class="flex-1"><el-option v-for="item in titleList" :key="item.id" :label="item.name" :value="item.id"/></el-select><pm-button text="新增" :icon="Plus" @click="handleAddTitle"/></div>
        <el-table :data="titleData" class="w-full pm-title-table"><el-table-column prop="id" label="ID" width="80" align="center"/><el-table-column prop="name" label="称号名称"/><el-table-column label="操作" width="80" align="center"><template #default="{ $index }"><el-button :icon="Delete" circle plain type="danger" size="small" @click="handleDeleteTitle($index)"/></template></el-table-column><template #empty><el-empty description="暂无称号" :image-size="60"/></template></el-table>
      </div>
      <template #footer><div class="flex justify-end gap-3 pb-2"><pm-button text="取消" color="white" :disabled="titleSaving" @click="showTitleDialog = false"/><pm-button text="保存" :loading="titleSaving" @click="saveTitles"/></div></template>
    </el-dialog>

    <el-dialog v-model="showResetDialog" title="重置用户密码" width="460px" align-center class="pm-manage-dialog" :close-on-click-modal="!resetSaving" :close-on-press-escape="!resetSaving">
      <el-form label-position="top">
        <p class="rounded-xl bg-[#F8FAFC] px-4 py-3 mb-4 text-sm text-[#64748B]">正在为 <strong class="text-[#00BAAD]">{{ resetTargetUser?.username }}</strong> 设置新密码</p>
        <el-form-item><template #label><div class="form-label"><PhLock :size="18"/>新密码</div></template><el-input v-model="resetPassword" type="password" show-password placeholder="请输入 6-32 位新密码" class="pm-manage-input"/></el-form-item>
        <el-form-item><template #label><div class="form-label"><PhLock :size="18"/>确认新密码</div></template><el-input v-model="resetConfirmPassword" type="password" show-password placeholder="请再次输入新密码" class="pm-manage-input"/></el-form-item>
      </el-form>
      <template #footer><div class="flex justify-end gap-3 pb-2"><pm-button text="取消" color="white" :disabled="resetSaving" @click="showResetDialog = false"/><pm-button text="确认重置" :loading="resetSaving" @click="saveResetPassword"/></div></template>
    </el-dialog>

    <el-dialog v-model="showPowerDialog" title="设置用户权限" width="480px" align-center class="pm-manage-dialog" :close-on-click-modal="!powerSaving" :close-on-press-escape="!powerSaving">
      <div class="flex flex-col gap-6 py-2"><div class="text-center text-[#64748B]">正在设置 <strong class="text-[#00BAAD]">{{ powerTargetUser?.username }}</strong> 的权限</div><div class="flex gap-4"><div class="power-card" :class="currentPower === 0 ? 'active-user' : ''" @click="currentPower = 0"><PhUser :size="26"/><strong>普通用户</strong><span>基础权限</span></div><div class="power-card" :class="currentPower === 1 ? 'active-admin' : ''" @click="currentPower = 1"><PhShieldCheck :size="26"/><strong>管理员</strong><span>全站管理权限</span></div></div></div>
      <template #footer><div class="flex justify-end gap-3 pb-2"><pm-button text="取消" color="white" :disabled="powerSaving" @click="showPowerDialog = false"/><pm-button text="确认修改" :loading="powerSaving" @click="savePower"/></div></template>
    </el-dialog>

    <manage-confirm-dialog v-model="deleteVisible" title="删除用户" :message="`确定删除用户「${deleteTarget?.username || ''}」吗？\n删除后该账号将无法登录，公开资料会显示为已注销用户。`" confirm-text="确认删除" confirm-color="red" :loading="deleteLoading" @confirm="confirmDelete"/>
  </div>
</template>

<style scoped lang="scss">
:deep(.el-tag) { border-radius: 10px; }
:deep(.el-select__wrapper), :deep(.el-input__wrapper) { border-radius: 12px !important; }
:deep(.pm-manage-input) { width: 100%; .el-input__wrapper { height: 48px; background: #F8FAFC; box-shadow: inset 0 0 0 1.5px transparent; padding: 8px 14px; &:hover { box-shadow: inset 0 0 0 1.5px #E2E8F0; } &.is-focus { box-shadow: inset 0 0 0 1.5px #00BAAD !important; } } }
:deep(.pm-title-table) { --el-table-border-color: #F1F5F9; border: 1px solid #F1F5F9; border-radius: 16px; overflow: hidden; &::before, &::after, .el-table__inner-wrapper::before { display: none !important; } .el-table__header-wrapper th { background: #F8FAFC; color: #64748B; } }
:deep(.pm-manage-dialog) { border: 1px solid #E2E8F0; border-radius: 20px !important; padding: 10px; overflow: hidden; .el-dialog__header { margin-right: 0; padding: 20px 24px 10px; .el-dialog__title { font-size: 18px; font-weight: 600; color: #1E293B; } } .el-dialog__body { padding: 10px 24px; } .el-dialog__footer { padding: 10px 24px 20px; } .el-dialog__headerbtn { top: 20px; right: 20px; &:hover .el-dialog__close { color: #00BAAD; } } }
.form-label { display: flex; align-items: center; gap: 6px; color: #64748B; font-weight: 500; }
.info-item { display: flex; flex-direction: column; gap: 8px; label { display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 600; color: #64748B; } .value { background: #F8FAFC; padding: 12px 16px; border-radius: 12px; color: #334155; border: 1px solid #F1F5F9; user-select: text; } .bio-box { min-height: 80px; white-space: pre-wrap; } }
.power-card { flex: 1; padding: 18px; border: 2px solid #F1F5F9; border-radius: 16px; background: #F8FAFC; cursor: pointer; display: flex; flex-direction: column; align-items: center; gap: 8px; color: #64748B; transition: .2s; span { font-size: 12px; color: #94A3B8; } &.active-user { border-color: #00BAAD; background: rgba(0,186,173,.05); color: #00A396; } &.active-admin { border-color: #EF4444; background: #FEF2F2; color: #EF4444; } }

@media (max-width: 767px) {
  .user-detail-container {
    flex-direction: column;
    gap: 20px;

    > .flex-1 {
      min-width: 0;
    }

    .grid-cols-2 {
      grid-template-columns: minmax(0, 1fr);
    }
  }

  .power-card {
    padding: 14px 8px;
  }

  :deep(.pm-title-table) {
    min-width: 0;
  }
}
</style>
