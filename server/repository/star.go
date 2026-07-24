package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"
)

func ToggleStar(userID, pluginID int) (bool, error) {
	var star model.UserStar
	err := database.DB.Where("user_id = ? AND plugin_id = ?", userID, pluginID).First(&star).Error
	if err == nil {
		// 已收藏，取消
		database.DB.Delete(&star)
		return false, nil
	}
	// 未收藏，添加
	err = database.DB.Create(&model.UserStar{UserID: userID, PluginID: pluginID}).Error
	return true, err
}

func GetStarList(userID, page, pageSize int) ([]model.PluginListItem, int64, error) {
	query := database.DB.Model(&model.UserStar{}).Where("user_id = ?", userID)
	var total int64
	query.Count(&total)

	var pluginIDs []int
	database.DB.Model(&model.UserStar{}).Where("user_id = ?", userID).
		Order("plugin_id DESC").
		Offset((page-1)*pageSize).Limit(pageSize).
		Pluck("plugin_id", &pluginIDs)

	var plugins []model.PluginListItem
	if len(pluginIDs) > 0 {
		database.DB.Model(&model.Plugin{}).Where("id IN ?", pluginIDs).
			Select("id, name, desc_text, icon, type, views, downloads").
			Find(&plugins)

		for i := range plugins {
			var starCount int64
			database.DB.Model(&model.UserStar{}).Where("plugin_id = ?", plugins[i].ID).Count(&starCount)
			plugins[i].Star = int(starCount)

			var frame model.Frame
			database.DB.Joins("JOIN pm_plugin_frames ON pm_plugin_frames.frame_id = pm_frames.id").
				Where("pm_plugin_frames.plugin_id = ?", plugins[i].ID).
				First(&frame)
			plugins[i].FrameIcon = frame.Icon
		}
	}

	return plugins, total, nil
}
