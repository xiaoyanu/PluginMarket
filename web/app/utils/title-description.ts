import sanitizeHtml from 'sanitize-html'

const allowedTags = [
  'br', 'b', 'strong', 'i', 'em', 'u', 's', 'del', 'code',
  'small', 'sub', 'sup', 'span', 'p', 'ul', 'ol', 'li',
]

export const renderTitleDescription = (value?: string) => sanitizeHtml(
  (value || '').replace(/\r\n?/g, '\n').replace(/\n/g, '<br>'),
  {
    allowedTags,
    allowedAttributes: {},
  },
)
