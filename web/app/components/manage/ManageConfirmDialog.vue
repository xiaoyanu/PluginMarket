<script setup lang="ts">
type ConfirmColor = 'blue' | 'red' | 'orange' | 'green' | 'gray' | 'black' | 'white'

withDefaults(defineProps<{
  title: string
  message: string
  confirmText?: string
  confirmColor?: ConfirmColor
  loading?: boolean
}>(), {
  confirmText: '确认',
  confirmColor: 'red',
  loading: false,
})

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ confirm: [] }>()
</script>

<template>
  <el-dialog
      v-model="visible"
      :title="title"
      width="460px"
      align-center
      class="pm-manage-dialog"
      :close-on-click-modal="!loading"
      :close-on-press-escape="!loading"
      :show-close="!loading"
  >
    <div class="rounded-xl bg-[#F8FAFC] px-4 py-4 text-sm leading-6 text-[#475569] whitespace-pre-wrap break-words">
      {{ message }}
    </div>
    <template #footer>
      <div class="flex justify-end gap-3 pb-2">
        <pm-button text="取消" color="white" :disabled="loading" @click="visible = false"/>
        <pm-button :text="confirmText" :color="confirmColor" :loading="loading" @click="emit('confirm')"/>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
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
</style>
