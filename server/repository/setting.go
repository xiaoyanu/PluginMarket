package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"
)

// GetSetting 获取单个设置值
func GetSetting(key string) string {
	var s model.Setting
	if err := database.DB.Where("s_key = ?", key).First(&s).Error; err != nil {
		return ""
	}
	return s.Value
}

// GetAllSettings 获取所有设置
func GetAllSettings() map[string]string {
	var settings []model.Setting
	database.DB.Find(&settings)
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}

// UpdateSetting 更新单个设置
func UpdateSetting(key, value string) error {
	var setting model.Setting
	if err := database.DB.Where("s_key = ?", key).FirstOrCreate(&setting, model.Setting{Key: key}).Error; err != nil {
		return err
	}
	// 使用 Update 而不是结构体 Assign，确保空字符串也能写入数据库。
	return database.DB.Model(&setting).Update("s_value", value).Error
}

// UpdateSettings 批量更新设置
func UpdateSettings(settings map[string]string) error {
	for key, val := range settings {
		if err := UpdateSetting(key, val); err != nil {
			return err
		}
	}
	return nil
}

// GetUploadReferenceValues 查询所有可能引用 uploads 文件的数据库字段。
// 除头像、各类图标和站点设置外，还会扫描插件正文与评论正文中的 Markdown/HTML 图片路径。
func GetUploadReferenceValues() ([]string, error) {
	values := make([]string, 0)
	queries := []struct {
		model  interface{}
		column string
	}{
		{&model.User{}, "avatar"},
		{&model.Plugin{}, "icon"},
		{&model.Plugin{}, "content"},
		{&model.Comment{}, "content"},
		{&model.Frame{}, "icon"},
		{&model.Title{}, "icon"},
		{&model.Setting{}, "s_value"},
	}

	for _, query := range queries {
		var fieldValues []string
		if err := database.DB.Model(query.model).
			Where(query.column+" LIKE ?", "%/uploads/%").
			Pluck(query.column, &fieldValues).Error; err != nil {
			return nil, err
		}
		values = append(values, fieldValues...)
	}
	return values, nil
}
