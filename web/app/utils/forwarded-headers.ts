type IncomingHeaders = Record<string, string | string[] | undefined>

const FORWARDED_CLIENT_IP_HEADERS = ['x-forwarded-for', 'x-real-ip'] as const

const firstHeaderValue = (value: string | string[] | undefined) => {
  if (Array.isArray(value)) return value.find(item => item.trim())?.trim() || ''
  return typeof value === 'string' ? value.trim() : ''
}

export const appendForwardedClientIpHeaders = (headers: Headers, incomingHeaders: IncomingHeaders = {}) => {
  for (const name of FORWARDED_CLIENT_IP_HEADERS) {
    const value = firstHeaderValue(incomingHeaders[name])
    if (value && !headers.has(name)) headers.set(name, value)
  }
  return headers
}
