import {
  DEFAULT_SITE_DESCRIPTION,
  DEFAULT_SITE_KEYWORDS,
  DEFAULT_SITE_LOGO,
  DEFAULT_SITE_TITLE,
} from '~/config'

export interface SiteSettings {
  siteLogo: string
  siteTitle: string
  siteKeywords: string
  siteDescription: string
}

const defaultSiteSettings: SiteSettings = {
  siteLogo: DEFAULT_SITE_LOGO,
  siteTitle: DEFAULT_SITE_TITLE,
  siteKeywords: DEFAULT_SITE_KEYWORDS,
  siteDescription: DEFAULT_SITE_DESCRIPTION,
}

interface PublicSiteSettingsResponse {
  data?: {
    siteLogo?: string
    siteTitle?: string
    siteKeywords?: string
    siteDesc?: string
  }
}

export const useSiteSettings = () => {
  const settings = useState<SiteSettings>('site-settings', () => ({ ...defaultSiteSettings }))

  const apply = (data: PublicSiteSettingsResponse['data'] = {}) => {
    settings.value = {
      siteLogo: data?.siteLogo || defaultSiteSettings.siteLogo,
      siteTitle: data?.siteTitle || defaultSiteSettings.siteTitle,
      siteKeywords: data?.siteKeywords || defaultSiteSettings.siteKeywords,
      siteDescription: data?.siteDesc || defaultSiteSettings.siteDescription,
    }
  }

  const load = async () => {
    try {
      const response = await useApiFetch<PublicSiteSettingsResponse>('/setting/public', { suppressErrorMessage: true })
      apply(response?.data)
    } catch {
      // 公共接口不可用时继续使用 config.ts 中的默认设置。
    }
  }

  return { settings, apply, load, defaults: defaultSiteSettings }
}