export type ApiResponse<T> = {
  code: number
  msg: string
  data: T
}

export type CommentAuthor = {
  id: number
  nick?: string
  avatar?: string
  titles?: any[]
}

export type CommentReply = {
  id: number
  content: string
  root: number
  parent: number
  reply_user: number
  created: string
  author: CommentAuthor
  replyTo?: { id: number; nick?: string }
}

export type CommentItem = {
  id: number
  content: string
  root: number
  parent: number
  created: string
  author: CommentAuthor
  replies: CommentReply[]
  replyCount: number
}

export type CommentListData = {
  list: CommentItem[]
  total: number
  target?: {
    commentId: number
    rootCommentId: number
    isReply: boolean
    page: number
  } | null
}

export type PluginAuthor = {
  id: number
  nick?: string
  avatar?: string
  userdesc?: string
  titles?: any[]
}

export type PluginDetail = {
  id: number
  name?: string
  desc_text?: string
  content?: string
  icon?: string
  type?: number
  views?: number
  downloads?: number
  star?: number
  url?: string
  url_code?: string
  created?: string
  updated?: string
  author?: PluginAuthor
  frameworks?: Array<{ id: number; name: string; icon?: string }>
  tags?: Array<{ id: number; name: string }>
  is_starred?: boolean
}

export const getPluginDetail = (id: number | string) => {
  return useApiFetch<PluginDetail>(`/plugin/${id}`, { method: 'GET' })
}

export type DownloadInfo = {
  url: string
  url_code?: string
}

export const getPluginDownloadInfo = (id: number | string) => {
  return useApiFetch<DownloadInfo>(`/plugin/${id}/download`, { method: 'GET' })
}

export const getCommentsByPlugin = (pluginId: number | string, page: number = 1, pageSize: number = 20, targetCommentId?: number | null) => {
  return useApiFetch<CommentListData>('/comment/list', {
    method: 'GET',
    query: { pluginId, page, pageSize, ...(targetCommentId ? { targetCommentId } : {}) },
  })
}

export const createComment = (body: { pluginId: number; content: string; parentId?: number }) => {
  return useApiFetch('/comment', {
    method: 'POST',
    body,
  })
}

export const deleteComment = (id: number) => {
  return useApiFetch(`/comment/${id}`, {
    method: 'DELETE',
  })
}
