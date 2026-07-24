<script setup lang="ts">
import {PhCaretRight, PhCaretUp, PhCaretDown, PhTrash} from "@phosphor-icons/vue";
import dayjs from 'dayjs'
import type { CommentItem, CommentReply, CommentListData } from '~/composables/api/plugin'
import { DEFAULT_AVATAR } from '~/config'
import { targetCommentElementId } from '~/utils/comment-target'

const props = defineProps<{
  comments: CommentItem[]
  total: number
  loading?: boolean
  pluginAuthorId: number
  commentTarget?: NonNullable<CommentListData['target']> | null
}>()

const emit = defineEmits<{
  (event: 'reply', comment: CommentItem | CommentReply): void
  (event: 'delete', comment: CommentItem | CommentReply): void
}>()

const userStore = useUserStore()
const assetUrl = useAssetUrl()
const expandedReplies = ref<Record<number, boolean>>({})

const isAuthor = (userId?: number) => !!userId && userId === props.pluginAuthorId
const formatTime = (value?: string) => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '-'
const canShowAuthorTag = (userId?: number) => isAuthor(userId)
const canDelete = (item: CommentItem | CommentReply) => {
  const currentUser = userStore.userInfo
  if (!userStore.isLogin || !currentUser?.id) return false
  if (currentUser.power === 1) return true
  if (currentUser.id === props.pluginAuthorId) return true
  return currentUser.id === item.author?.id
}
const isExpanded = (id: number) => expandedReplies.value[id] ?? false
const toggleReplies = (id: number) => {
  expandedReplies.value[id] = !isExpanded(id)
}

watch(() => props.commentTarget, (target) => {
  if (target?.isReply) expandedReplies.value[target.rootCommentId] = true
}, { immediate: true })
</script>

<template>
  <div class="flex flex-col gap-5">
    <div class="text-[#1E293B] font-bold text-[16px]">
      {{ total }} 条评论
    </div>
    <div v-if="loading" class="py-6 text-center text-[#64748B]">评论加载中...</div>
    <div v-else-if="comments.length === 0" class="py-6 text-center text-[#64748B]">还没有评论，来抢沙发吧</div>
    <div v-else class="flex flex-col gap-5 text-[15px]">
      <div
        v-for="item in comments"
        :id="targetCommentElementId(item.id)"
        :key="item.id"
        class="comment-entry flex min-w-0 gap-2"
        :class="{ 'comment-target': commentTarget?.commentId === item.id }"
      >
        <div class="shrink-0">
          <nuxt-link :to="`/home/${item.author?.id || 1}`" target="_blank">
            <img
              class="comment-avatar cursor-pointer w-12.5 h-12.5 rounded-[10px]"
              :src="assetUrl(item.author?.avatar, DEFAULT_AVATAR)"
              alt="avatar"
              draggable="false"
            />
          </nuxt-link>
        </div>
        <div class="flex min-w-0 flex-1 flex-col gap-1.5">
          <div class="flex gap-2 items-center flex-wrap">
            <div class="flex items-center gap-1">
              <span class="text-[#1E293B] select-text">{{ item.author?.nick || '匿名用户' }}</span>
              <span v-if="canShowAuthorTag(item.author?.id)" class="author">作者</span>
            </div>
            <span class="text-[14px] text-[#64748ba5]">{{ formatTime(item.created) }}</span>
            <span v-if="userStore.isLogin" class="text-[#006affcc] cursor-pointer hover:text-[#006aff]" @click="emit('reply', item)">回复</span>
            <span v-if="canDelete(item)" class="text-red-400 cursor-pointer hover:text-red-500" @click="emit('delete', item)">
              <el-tooltip content="删除评论" effect="light" placement="right">
                <PhTrash weight="bold"/>
              </el-tooltip>
            </span>
          </div>
          <div class="comment-content text-[#64748B] select-text whitespace-pre-wrap">
            {{ item.content }}
          </div>

          <div v-if="item.replies?.length" class="mt-2">
            <div v-if="isExpanded(item.id)" class="flex flex-1 flex-col gap-3 pl-3 border-l border-[#E2E8F0]">
              <div
                v-for="reply in item.replies"
                :id="targetCommentElementId(reply.id)"
                :key="reply.id"
                class="comment-entry flex min-w-0 gap-2"
                :class="{ 'comment-target': commentTarget?.commentId === reply.id }"
              >
                <div class="shrink-0">
                  <nuxt-link :to="`/home/${reply.author?.id || 1}`" target="_blank">
                    <img
                      class="comment-avatar cursor-pointer w-8 h-8 rounded-lg"
                      :src="assetUrl(reply.author?.avatar, DEFAULT_AVATAR)"
                      alt="avatar"
                      draggable="false"
                    />
                  </nuxt-link>
                </div>
                <div class="flex min-w-0 flex-1 flex-col gap-1.5">
                  <div class="flex gap-2 items-center flex-wrap">
                    <div class="flex items-center gap-1 flex-wrap">
                      <span class="text-[#1E293B] select-text">{{ reply.author?.nick || '匿名用户' }}</span>
                      <span v-if="canShowAuthorTag(reply.author?.id)" class="author">作者</span>
                      <span v-if="reply.replyTo?.nick"><PhCaretRight weight="fill" color="#A0A0A0"/></span>
                      <span v-if="reply.replyTo?.nick" class="text-[#1E293B] select-text">{{ reply.replyTo.nick }}</span>
                      <span v-if="reply.replyTo?.nick && isAuthor(reply.replyTo?.id)" class="author">作者</span>
                    </div>
                    <span class="text-[14px] text-[#64748ba5]">{{ formatTime(reply.created) }}</span>
                    <span v-if="userStore.isLogin" class="text-[#006affcc] cursor-pointer hover:text-[#006aff]" @click="emit('reply', reply)">回复</span>
                    <span v-if="canDelete(reply)" class="text-red-400 cursor-pointer hover:text-red-500" @click="emit('delete', reply)">
                      <el-tooltip content="删除评论" effect="light" placement="right">
                        <PhTrash weight="bold"/>
                      </el-tooltip>
                    </span>
                  </div>
                  <div class="comment-content text-[#64748B] select-text whitespace-pre-wrap">
                    {{ reply.content }}
                  </div>
                </div>
              </div>
            </div>
            <div class="cursor-pointer flex gap-1 items-center text-[13px] text-[#4D5663] font-bold" @click="toggleReplies(item.id)">
              {{ isExpanded(item.id) ? '收起' : `展开 ${item.replies.length} 条回复` }}
              <PhCaretUp v-if="isExpanded(item.id)" weight="bold"/>
              <PhCaretDown v-else weight="bold"/>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.comment-avatar {
  background-color: #fff;
  border: 1px solid #FFF;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.comment-entry {
  border-radius: 10px;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

.comment-content {
  overflow-wrap: anywhere;
}

@media (max-width: 767px) {
  .comment-avatar {
    width: 40px !important;
    height: 40px !important;
  }

  .comment-entry .comment-entry .comment-avatar {
    width: 30px !important;
    height: 30px !important;
  }
}

.comment-target {
  background-color: #EFF6FF;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.25);
}

.author {
  font-size: 12px;
  color: #fd3054;
  padding: 2px 4px;
  border-radius: 5px;
  background-color: #fcdfe7;
}
</style>
