<script setup lang="ts">
import {PhPaperPlaneTilt} from "@phosphor-icons/vue";

type ReplyTarget = {
  id: number
  author?: {
    nick?: string
  }
}

const comment = defineModel<string>({ default: '' })
const props = defineProps<{
  disabled?: boolean
  replyTarget?: ReplyTarget | null
  loading?: boolean
}>()

const emit = defineEmits<{
  (event: 'submit'): void
  (event: 'cancel-reply'): void
}>()

const placeholder = computed(() => props.disabled ? '请登录后评论' : '快来这里留下你的评论吧！')

const handleSubmit = () => {
  if (props.disabled) return
  emit('submit')
}
</script>

<template>
  <div class="flex flex-col">
    <div>
      <div v-if="replyTarget" class="mb-2 flex items-center justify-between rounded-xl bg-[#F1F5F9] px-3 py-2 text-sm text-[#475569]">
        <span>正在回复：{{ replyTarget.author?.nick || '匿名用户' }}</span>
        <span class="cursor-pointer text-[#006affcc] hover:text-[#006aff]" @click="emit('cancel-reply')">取消回复</span>
      </div>
      <el-input
          v-model="comment"
          :disabled="props.disabled"
          :placeholder="placeholder"
          type="textarea"
          class="pm-input"
      />
      <div class="flex mt-2.5 justify-end">
        <pm-button :text="props.loading ? '发送中' : '发送评论'" :disabled="props.disabled || props.loading" @click="handleSubmit">
          <template #icon>
            <PhPaperPlaneTilt/>
          </template>
        </pm-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
:deep(.pm-input) {
  .el-textarea__inner {
    background-color: #F1F5F9;
    box-shadow: none;
    border: 1.5px solid transparent;
    border-radius: 12px;
    padding: 10px;
    transition: background-color 0.2s, border-color 0.2s, box-shadow 0.2s;
    font-size: 14px;
    color: #1e293b;
    min-height: 100px !important;

    &::placeholder {
      color: #94a3b8;
    }

    &:focus {
      background-color: #fff;
      border-color: #60a5fa !important;
      box-shadow: 0 0 0 4px rgba(96, 165, 250, 0.1) !important;
      outline: none;
    }

    &:hover:not(:focus) {
      border-color: #cbd5e1;
    }
  }
}
</style>