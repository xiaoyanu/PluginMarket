// 后端地址
const API_HOST = 'https://pm.zxz.ee/service'

export const siteConfig = {
  siteLogo: '/images/logo.png',
  siteFavicon: '/images/favicon.ico',
  siteTitle: 'PluginMarket 插件市场',
  siteKeywords: 'Bee框架,插件市场,插件资源,下载插件,plugin,pluginmarket,插件,QQ机器人',
  siteDescription: '一个专注于插件分享的平台，提供丰富的插件资源供用户下载。',
  allowRegister: true,
  allowComment: true,
  skipAudit: false,
  allowUpload: true,
  maxImageSize: 1.5,
  allowedImageExtensions: 'jpg,png,gif,webp',
  defaultAvatar: '/images/default-avatar.png',
  defaultPluginIcon: '/images/plugin-icon.png',
  notificationTemplates: {
    newComment: {
      title: '您的插件 {{plugin_name}} 收到了新评论',
      content: '{{user_name}} ：{{content}}',
    },
    replyComment: {
      title: '{{user_name}} 回复了您在 {{plugin_name}} 的评论',
      content: '{{content}}',
    },
    pluginApproved: {
      title: '插件审核通过',
      content: '您的插件《{{plugin_name}}》已审核通过并正式发布。',
    },
    pluginRejected: {
      title: '插件审核未通过',
      content: '您的插件《{{plugin_name}}》未通过审核，拒绝理由：{{content}}',
    },
  },
  emailTemplates: {
    newComment: {
      title: '[ {{site_title}} ] 您的插件收到了新评论',
      body: "<div style=\"width: 510px; height: auto; margin: 40px auto; padding: 10px 30px; border-radius: 5px; border: 1px solid #ffb0b0; box-shadow: 0px 0px 20px #888888; user-select: none;\"><a href=\"https://pm.zxz.ee\" target=\"_blank\" style=\"text-decoration: none\"><div style=\"height: 60px; width: 320px; line-height: 60px; text-align: center; font-size: 24px; background-color: #ff7272; border-radius: 20px; color: #fff; margin: 10px auto; margin-bottom: 60px;\">🍰{{site_title}}</div></a><div style=\"width: auto; height: auto; margin: 0 10px; color: rgb(0, 0, 0); text-align: center; border-radius: 10px;\"><span style=\"font-size: 20px; font-weight: bold\">《{{plugin_name}}》有了新的评论</span><hr style=\"width: 80%; border: 0; height: 1px; background-color: #f1f1f1\"/></div><div style=\"text-align: center; font-size: 12px; color: #969696; padding: 1rem 0\">{{time}}</div><div style=\"display: flex; flex-direction: column\"><div style=\"line-height: 2rem; font-size: 13px; font-weight: bold\">{{comment_user}}</div><div style=\"display: flex; margin-bottom: 5px\"><div style=\"width: auto; flex: none\"><img id=\"cravatar\" src=\"data:image/png;base64,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\" style=\"width: 35px; height: 35px; border-radius: 50%\"/></div><div style=\"position: relative; margin-left: 20px\"><span style=\"width: 0; height: 0; border-top: 8px solid transparent; border-bottom: 8px solid transparent; border-right: 8px solid; border-right-color: #373737; left: -8px; right: auto; top: 12px; position: absolute;\"></span><div style=\"background-color: #373737; padding: 10px; border-radius: 9px; margin-bottom: 3px; font-size: 13px; color: #ffffff; word-break: break-all;\">{{comment_content}}</div></div></div></div><div style=\"text-align: center; margin-top: 60px; font-size: 12px; color: #888\"><a href=\"{{link}}\" target=\"_blank\">点击查看评论</a><br><br>由系统发送，请勿回复。</div></div>",
    },
    replyComment: {
      title: '[ {{site_title}} ] 您的评论收到了回复',
      body: "<div style=\"width:510px;height:auto;margin:40px auto;padding:10px;padding-left:30px;padding-right:30px;border-radius:5px;border:1px solid #ffb0b0;box-shadow:0px 0px 20px #888888;user-select:none;\"><a href=\"https://pm.zxz.ee\" target=\"_blank\" style=\"text-decoration:none;\"><div style=\"height:60px;width:320px;line-height:60px;text-align:center;font-size:24px;background-color:#ff7272;border-radius:20px;color:#fff;margin:10px auto;margin-bottom:60px;\">🍰{{site_title}}</div></a><div style=\"width:auto;height:auto;margin:0 10px;color:rgb(0,0,0);text-align:center;border-radius:10px;\"><span style=\"font-size:20px;font-weight:bold;\">《{{plugin_name}}》的评论有了新的回复</span><hr style=\"width:80%;border:0;height:1px;background-color:#f1f1f1;\"></div><div style=\"text-align:center;font-size:12px;color:#969696;padding:1rem 0;\">{{time}}</div><div style=\"display:flex;flex-direction:column;align-items:flex-end;\"><div style=\"line-height:2rem;font-size:13px;font-weight:bold;\">{{user}}</div><div style=\"display:flex;margin-bottom:5px;\"><div style=\"position:relative;margin-right:20px;\"><span style=\"width:0;height:0;border-top:8px solid transparent;border-bottom:8px solid transparent;border-left:8px solid;border-left-color:#373737;left:auto;right:-8px;top:12px;position:absolute;\"></span><div style=\"background-color:#373737;padding:10px;border-radius:9px;margin-bottom:3px;font-size:13px;color:#ffffff;word-break:break-all;\">{{content}}</div></div><div style=\"width:auto;flex:none;\"><img id=\"cravatar\" src=\"data:image/png;base64,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\" style=\"width:35px;height:35px;border-radius:50%;\"></div></div></div><div style=\"display:flex;flex-direction:column;\"><div style=\"line-height:2rem;font-size:13px;font-weight:bold;\">{{reply_user}}</div><div style=\"display:flex;margin-bottom:5px;\"><div style=\"width:auto;flex:none;\"><img id=\"cravatar\" src=\"data:image/png;base64,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\" style=\"width:35px;height:35px;border-radius:50%;\"></div><div style=\"position:relative;margin-left:20px;\"><span style=\"width:0;height:0;border-top:8px solid transparent;border-bottom:8px solid transparent;border-right:8px solid;border-right-color:#373737;left:-8px;right:auto;top:12px;position:absolute;\"></span><div style=\"background-color:#373737;padding:10px;border-radius:9px;margin-bottom:3px;font-size:13px;color:#ffffff;word-break:break-all;\">{{reply_content}}</div></div></div></div><div style=\"text-align: center; margin-top: 60px; font-size: 12px; color: #888\"><a href=\"{{link}}\" target=\"_blank\">点击查看评论</a><br><br>由系统发送，请勿回复。</div></div>",
    },
    pendingPluginReview: {
      title: '[ {{site_title}} ] 有插件待审核',
      body: '<h2>有插件待审核</h2><p>用户发布或修改了插件《{{plugin_name}}》，请及时前往管理后台审核。</p>',
    },
    pluginApproved: {
      title: '[ {{site_title}} ] 您的插件已审核通过',
      body: '<h2>插件审核通过</h2><p>您的插件《{{plugin_name}}》已审核通过并正式发布。</p>',
    },
    pluginRejected: {
      title: '[ {{site_title}} ] 您的插件审核未通过',
      body: '<h2>插件审核未通过</h2><p>您的插件《{{plugin_name}}》未通过审核。</p><p><strong>拒绝理由：</strong>{{content}}</p><p><a href="{{link}}">查看并修改插件</a></p>',
    },
    emailVerify: {
      title: '[ {{site_title}} ] 请验证您的邮箱',
      body: "<div style=\"width: 550px; height: auto; margin: 40px auto; padding: 10px; border-radius: 5px; border: 1px solid #ffb0b0; box-shadow: 0px 0px 20px #888888; user-select: none;\"><a href=\"https://pm.zxz.ee\" target=\"_blank\" style=\"text-decoration: none\"><div style=\"height: 60px; width: 320px; line-height: 60px; text-align: center; font-size: 24px; background-color: #ff7272; border-radius: 20px; color: #fff; margin: 10px auto; margin-bottom: 60px;\">🍰{{site_title}}</div></a><div style=\"width: auto; height: auto; margin: 0 10px; color: rgb(0, 0, 0); text-align: center; border-radius: 10px;\"><a style=\"text-decoration: none; color: #000\" href=\"{{link}}\"><div style=\"background-color: #ecf5ff; width: 260px; height: 40px; line-height: 40px; font-size: 20px; font-weight: bold; margin: 0 auto; border-radius: 10px; border: 1px solid #409eff; color: #409eff;\">点击验证邮箱，并完成绑定</div></a></div><div style=\"text-align: center; margin-top: 60px; font-size: 12px; color: #888\">链接将在10分钟后失效，如非本人操作请忽略此邮件</div></div>",
    },
    resetPassword: {
      title: '[ {{site_title}} ] 重置您的账户密码',
      body: "<div style=\"width: 550px; height: auto; margin: 40px auto; padding: 10px; border-radius: 5px; border: 1px solid #ffb0b0; box-shadow: 0px 0px 20px #888888; user-select: none;\"><a href=\"https://pm.zxz.ee\" target=\"_blank\" style=\"text-decoration: none\"><div style=\"height: 60px; width: 320px; line-height: 60px; text-align: center; font-size: 24px; background-color: #ff7272; border-radius: 20px; color: #fff; margin: 10px auto; margin-bottom: 60px;\">🍰{{site_title}}</div></a><div style=\"width: auto; height: auto; margin: 0 10px; color: rgb(0, 0, 0); text-align: center; border-radius: 10px;\"><a style=\"text-decoration: none; color: #000\" href=\"{{link}}\"><div style=\"background-color: #ecf5ff; width: 260px; height: 40px; line-height: 40px; font-size: 20px; font-weight: bold; margin: 0 auto; border-radius: 10px; border: 1px solid #409eff; color: #409eff;\">点击重置密码</div></a></div><div style=\"text-align: center; margin-top: 60px; font-size: 12px; color: #888\">链接将在10分钟后失效，如非本人操作请忽略此邮件</div></div>",
    },
  },
  apiBase: `${API_HOST}/api`,
  assetBase: API_HOST,
  defaultUserProfile: '这个人很懒……什么也没留下。',
} as const

export const DEFAULT_AVATAR = siteConfig.defaultAvatar
export const DEFAULT_SITE_LOGO = siteConfig.siteLogo
export const DEFAULT_SITE_FAVICON = siteConfig.siteFavicon
export const DEFAULT_SITE_TITLE = siteConfig.siteTitle
export const DEFAULT_SITE_KEYWORDS = siteConfig.siteKeywords
export const DEFAULT_SITE_DESCRIPTION = siteConfig.siteDescription
export const DEFAULT_ALLOW_REGISTER = siteConfig.allowRegister
export const DEFAULT_ALLOW_COMMENT = siteConfig.allowComment
export const DEFAULT_SKIP_AUDIT = siteConfig.skipAudit
export const DEFAULT_ALLOW_UPLOAD = siteConfig.allowUpload
export const DEFAULT_MAX_IMAGE_SIZE = siteConfig.maxImageSize
export const DEFAULT_ALLOWED_IMAGE_EXTENSIONS = siteConfig.allowedImageExtensions
export const DEFAULT_PLUGIN_ICON = siteConfig.defaultPluginIcon
export const DEFAULT_NOTIFICATION_TEMPLATES = siteConfig.notificationTemplates
export const DEFAULT_EMAIL_TEMPLATES = siteConfig.emailTemplates
export const API_BASE = siteConfig.apiBase
export const ASSET_BASE = siteConfig.assetBase
export const DEFAULT_USER_PROFILE = siteConfig.defaultUserProfile
