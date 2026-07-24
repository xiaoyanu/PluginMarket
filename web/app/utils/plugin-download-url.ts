export const isValidPluginDownloadUrl = (value: string) => {
  const normalized = value.trim()
  if (!/^https?:\/\/[^/\s]+(?:[/?#]|$)/i.test(normalized)) return false

  try {
    const url = new URL(normalized)
    return (url.protocol === 'http:' || url.protocol === 'https:') && Boolean(url.hostname)
  } catch {
    return false
  }
}
