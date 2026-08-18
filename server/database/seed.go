package database

import (
	"errors"
	"log"

	"pluginmarket-server/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed 初始化种子数据
func Seed() {
	seedSettings()
	seedAdmin()
}

func seedSettings() {
	defaults := defaultSettings()

	for key, val := range defaults {
		var s model.Setting
		err := DB.Where("s_key = ?", key).First(&s).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("种子数据：查询设置 %s 失败: %v", key, err)
			continue
		}
		if err := DB.Create(&model.Setting{Key: key, Value: val}).Error; err != nil {
			log.Printf("种子数据：创建设置 %s 失败: %v", key, err)
		}
	}
	log.Println("种子数据：设置初始化完成")
}

func defaultSettings() map[string]string {
	return map[string]string{
		"siteTitle":           "插件市场 PluginMarket",
		"siteDesc":            "一个专注于插件分享的平台",
		"siteKeywords":        "插件,市场,Plugin,Market",
		"siteLogo":            "",
		"defaultAvatar":       "",
		"defaultPluginIcon":   "",
		"allowRegister":       "true",
		"allowComment":        "true",
		"skipAudit":           "false",
		"allowUpload":         "true",
		"antiBrushPluginData": "false",
		"maxFileSize":         "1.5",
		"allowedExtensions":   "jpg,png,gif,webp",
		"emailHost":           "",
		"emailPort":           "465",
		"emailUser":           "",
		"emailPass":           "",
		"emailFrom":           "",
		"emailFromName":       "PluginMarket",
	}
}

func seedAdmin() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.User{
		Username: "admin",
		Password: string(hashed),
		Nick:     "管理员",
		Email:    "admin@example.com",
		Power:    1,
	}
	DB.Create(&admin)
	log.Println("种子数据：默认管理员已创建 (admin / admin123)")
}
