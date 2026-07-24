export const parseTargetCommentId = (value: string | string[] | undefined): number | null => {
  const raw = Array.isArray(value) ? value[0] : value
  if (!raw || !/^\d+$/.test(raw)) return null
  const id = Number(raw)
  return Number.isSafeInteger(id) && id > 0 ? id : null
}

export const targetCommentElementId = (commentId: number) => `comment-${commentId}`
