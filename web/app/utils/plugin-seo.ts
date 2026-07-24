export type PluginSeoData = {
  name?: string
  desc_text?: string
  description?: string
}

export type SiteSeoData = {
  title: string
  description: string
}

export const getPluginSeoMeta = (
  plugin: PluginSeoData | null | undefined,
  site: SiteSeoData,
) => {
  const name = plugin?.name?.trim()
  const description = (plugin?.desc_text || plugin?.description)?.trim()

  return {
    title: name ? `${name} - ${site.title}` : site.title,
    description: description || site.description,
  }
}