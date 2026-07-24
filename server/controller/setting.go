package controller

import (
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

var publicSiteSettingKeys = [...]string{
	"siteLogo",
	"siteTitle",
	"siteKeywords",
	"siteDesc",
}

func buildPublicSiteSettings(settings map[string]string) map[string]string {
	result := make(map[string]string, len(publicSiteSettingKeys))
	for _, key := range publicSiteSettingKeys {
		result[key] = settings[key]
	}
	return result
}

// GetPublicSiteSettings 获取无需登录即可使用的网站基础设置。
// GET /api/setting/public
func GetPublicSiteSettings(c *gin.Context) {
	utils.OKData(c, buildPublicSiteSettings(repository.GetAllSettings()))
}
