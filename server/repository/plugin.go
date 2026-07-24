package repository

import (
	"errors"
	"fmt"

	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

var errInvalidPluginRelation = errors.New("插件关联数据无效")

func IsInvalidPluginRelation(err error) bool {
	return errors.Is(err, errInvalidPluginRelation)
}

func ValidatePluginRelations(frameIDs, tagIDs []int) error {
	return validatePluginRelations(database.DB, frameIDs, tagIDs)
}

func GetPluginByID(id int) (*model.Plugin, error) {
	var plugin model.Plugin
	err := database.DB.First(&plugin, id).Error
	return &plugin, err
}

func GetPluginDetail(id int) (*model.PluginDetail, error) {
	var plugin model.Plugin
	if err := database.DB.First(&plugin, id).Error; err != nil {
		return nil, err
	}

	detail := &model.PluginDetail{
		ID:        plugin.ID,
		Name:      plugin.Name,
		DescText:  plugin.DescText,
		Content:   plugin.Content,
		Icon:      plugin.Icon,
		Type:      plugin.Type,
		Status:    plugin.Status,
		Views:     plugin.Views,
		Downloads: plugin.Downloads,
		URL:       plugin.URL,
		URLCode:   plugin.URLCode,
		Created:   plugin.Created,
		Updated:   plugin.Updated,
	}

	// 获取作者信息
	var author model.User
	if err := database.DB.First(&author, plugin.User).Error; err == nil {
		author.AnonymizeIfDeleted()
		var authorTitleIDs []int
		database.DB.Model(&model.UserTitle{}).
			Where("user_id = ?", author.ID).
			Pluck("title_id", &authorTitleIDs)

		var authorTitles []model.Title
		if len(authorTitleIDs) > 0 {
			database.DB.Where("id IN ?", authorTitleIDs).
				Order("id ASC").
				Find(&authorTitles)
		}

		detail.Author = model.PluginAuthor{
			ID:       author.ID,
			Nick:     author.Nick,
			Avatar:   author.Avatar,
			Userdesc: author.Userdesc,
			Titles:   authorTitles,
		}
	}

	// 获取框架
	var frameIDs []int
	database.DB.Model(&model.PluginFrame{}).
		Where("plugin_id = ?", id).
		Pluck("frame_id", &frameIDs)
	if len(frameIDs) > 0 {
		var frameworks []model.Frame
		database.DB.Where("id IN ?", frameIDs).
			Order("id ASC").
			Find(&frameworks)
		detail.Frameworks = frameworks
	} else {
		detail.Frameworks = []model.Frame{}
	}

	// 获取标签
	var tagIDs []int
	database.DB.Model(&model.PluginTag{}).
		Where("plugin_id = ?", id).
		Pluck("tag_id", &tagIDs)
	if len(tagIDs) > 0 {
		var tags []model.Tag
		database.DB.Where("id IN ?", tagIDs).
			Order("id ASC").
			Find(&tags)
		detail.Tags = tags
	} else {
		detail.Tags = []model.Tag{}
	}

	// 获取收藏数
	var starCount int64
	database.DB.Model(&model.UserStar{}).Where("plugin_id = ?", id).Count(&starCount)
	detail.Star = int(starCount)

	return detail, nil
}

func CreatePluginWithRelations(plugin *model.Plugin, frameIDs, tagIDs []int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := validatePluginRelations(tx, frameIDs, tagIDs); err != nil {
			return err
		}
		if err := tx.Create(plugin).Error; err != nil {
			return err
		}
		if err := replacePluginFrames(tx, plugin.ID, frameIDs); err != nil {
			return err
		}
		return replacePluginTags(tx, plugin.ID, tagIDs)
	})
}

func UpdatePluginWithRelations(pluginID int, fields map[string]interface{}, updateFrames bool, frameIDs []int, updateTags bool, tagIDs []int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := validatePluginRelations(tx, optionalIDs(updateFrames, frameIDs), optionalIDs(updateTags, tagIDs)); err != nil {
			return err
		}
		if len(fields) > 0 {
			result := tx.Model(&model.Plugin{}).Where("id = ?", pluginID).Updates(fields)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		if updateFrames {
			if err := replacePluginFrames(tx, pluginID, frameIDs); err != nil {
				return err
			}
		}
		if updateTags {
			return replacePluginTags(tx, pluginID, tagIDs)
		}
		return nil
	})
}

func optionalIDs(update bool, ids []int) []int {
	if !update {
		return nil
	}
	return ids
}

func validatePluginRelations(tx *gorm.DB, frameIDs, tagIDs []int) error {
	if len(frameIDs) > 0 {
		var count int64
		if err := tx.Model(&model.Frame{}).Where("id IN ?", frameIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(frameIDs)) {
			return fmt.Errorf("%w: 适用框架不存在或已被删除", errInvalidPluginRelation)
		}
	}
	if len(tagIDs) > 0 {
		var count int64
		if err := tx.Model(&model.Tag{}).Where("id IN ?", tagIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(tagIDs)) {
			return fmt.Errorf("%w: 标签不存在或已被删除", errInvalidPluginRelation)
		}
	}
	return nil
}

func replacePluginFrames(tx *gorm.DB, pluginID int, frameIDs []int) error {
	if err := tx.Where("plugin_id = ?", pluginID).Delete(&model.PluginFrame{}).Error; err != nil {
		return err
	}
	for _, frameID := range frameIDs {
		if err := tx.Create(&model.PluginFrame{PluginID: pluginID, FrameID: frameID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func replacePluginTags(tx *gorm.DB, pluginID int, tagIDs []int) error {
	if err := tx.Where("plugin_id = ?", pluginID).Delete(&model.PluginTag{}).Error; err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := tx.Create(&model.PluginTag{PluginID: pluginID, TagID: tagID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func UpdatePlugin(plugin *model.Plugin) error {
	return database.DB.Save(plugin).Error
}

func UpdatePluginFields(id int, fields map[string]interface{}) error {
	result := database.DB.Model(&model.Plugin{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeletePlugin(id int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plugin_id = ?", id).Delete(&model.PluginFrame{}).Error; err != nil {
			return err
		}
		if err := tx.Where("plugin_id = ?", id).Delete(&model.PluginTag{}).Error; err != nil {
			return err
		}
		if err := tx.Where("plugin = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("plugin_id = ?", id).Delete(&model.UserStar{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Plugin{}, id).Error
	})
}

func IncrementPluginViews(id int) error {
	return database.DB.Model(&model.Plugin{}).Where("id = ?", id).
		UpdateColumn("views", database.DB.Raw("`views` + 1")).Error
}

func IncrementPluginDownloads(id int) error {
	return database.DB.Model(&model.Plugin{}).Where("id = ?", id).
		UpdateColumn("downloads", database.DB.Raw("`downloads` + 1")).Error
}

func fillPluginListExtra(plugins []model.PluginListItem) {
	if len(plugins) == 0 {
		return
	}

	var pluginIDs []int
	for _, p := range plugins {
		pluginIDs = append(pluginIDs, p.ID)
	}

	// 批量查询收藏数
	var starCounts []struct {
		PluginID int
		Count    int
	}
	database.DB.Model(&model.UserStar{}).
		Select("plugin_id, count(*) as count").
		Where("plugin_id IN ?", pluginIDs).
		Group("plugin_id").
		Scan(&starCounts)

	starMap := make(map[int]int)
	for _, sc := range starCounts {
		starMap[sc.PluginID] = sc.Count
	}

	// 批量查询首个框架图标
	var frames []struct {
		PluginID int
		Icon     string
	}
	database.DB.Table("pm_plugin_frames").
		Select("pm_plugin_frames.plugin_id, pm_frames.icon").
		Joins("JOIN pm_frames ON pm_frames.id = pm_plugin_frames.frame_id").
		Where("pm_plugin_frames.plugin_id IN ?", pluginIDs).
		Scan(&frames)

	frameMap := make(map[int]string)
	for _, f := range frames {
		if _, ok := frameMap[f.PluginID]; !ok {
			frameMap[f.PluginID] = f.Icon
		}
	}

	// 填充数据
	for i := range plugins {
		plugins[i].Star = starMap[plugins[i].ID]
		plugins[i].FrameIcon = frameMap[plugins[i].ID]
	}
}

// 插件列表查询（通用）
func getPluginListQuery(frameID, tagID *int, pluginType *int) *gorm.DB {
	query := database.DB.Model(&model.Plugin{}).Where("status = 0")
	if frameID != nil {
		query = query.Joins("JOIN pm_plugin_frames ON pm_plugin_frames.plugin_id = pm_plugins.id").
			Where("pm_plugin_frames.frame_id = ?", *frameID)
	}
	if tagID != nil {
		query = query.Joins("JOIN pm_plugin_tags ON pm_plugin_tags.plugin_id = pm_plugins.id").
			Where("pm_plugin_tags.tag_id = ?", *tagID)
	}
	if pluginType != nil {
		query = query.Where("type = ?", *pluginType)
	}
	return query
}

func GetPluginListHot(frameID, tagID, pluginType *int, page, pageSize int) ([]model.PluginListItem, int64, error) {
	query := getPluginListQuery(frameID, tagID, pluginType)

	var total int64
	query.Session(&gorm.Session{}).Distinct("pm_plugins.id").Count(&total)

	var plugins []model.PluginListItem
	err := query.Select(`
		pm_plugins.id,
		pm_plugins.name,
		pm_plugins.desc_text,
		pm_plugins.icon,
		pm_plugins.type,
		pm_plugins.views,
		pm_plugins.downloads,
		COUNT(pm_user_stars.plugin_id) AS star,
		ROUND((pm_plugins.downloads + pm_plugins.views + COUNT(pm_user_stars.plugin_id)) / 3, 2) AS hot_score
	`).
		Joins("LEFT JOIN pm_user_stars ON pm_user_stars.plugin_id = pm_plugins.id").
		Group("pm_plugins.id").
		Order("hot_score DESC, pm_plugins.downloads DESC, pm_plugins.views DESC, pm_plugins.id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&plugins).Error

	fillPluginListExtra(plugins)

	return plugins, total, err
}

func GetPluginListLatest(frameID, tagID, pluginType *int, page, pageSize int) ([]model.PluginListItem, int64, error) {
	return getPluginListSorted("created DESC", frameID, tagID, pluginType, page, pageSize)
}

func getPluginListSorted(order string, frameID, tagID, pluginType *int, page, pageSize int) ([]model.PluginListItem, int64, error) {
	query := getPluginListQuery(frameID, tagID, pluginType)

	var total int64
	query.Count(&total)

	var plugins []model.PluginListItem
	err := query.Select("pm_plugins.id, pm_plugins.name, pm_plugins.desc_text, pm_plugins.icon, pm_plugins.type, pm_plugins.views, pm_plugins.downloads").
		Order("pm_plugins." + order).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&plugins).Error

	// 填充收藏数和首个框架图标
	fillPluginListExtra(plugins)

	return plugins, total, err
}

func SearchPlugins(keywords string, frameID, tagID, pluginType *int, page, pageSize int) ([]model.PluginListItem, int64, error) {
	query := getPluginListQuery(frameID, tagID, pluginType).
		Where("pm_plugins.name LIKE ? OR pm_plugins.desc_text LIKE ?", "%"+keywords+"%", "%"+keywords+"%")

	var total int64
	query.Count(&total)

	var plugins []model.PluginListItem
	err := query.Select("pm_plugins.id, pm_plugins.name, pm_plugins.desc_text, pm_plugins.icon, pm_plugins.type, pm_plugins.views, pm_plugins.downloads").
		Order("pm_plugins.views DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&plugins).Error

	fillPluginListExtra(plugins)

	return plugins, total, err
}

func GetPluginsByTag(tagID, page, pageSize int) ([]model.PluginListItem, int64, error) {
	return getPluginListSorted("created DESC", nil, &tagID, nil, page, pageSize)
}

func GetPluginsByUser(userID int, frameID, tagID, pluginType *int, page, pageSize int) ([]model.PluginListItem, int64, error) {
	query := getPluginListQuery(frameID, tagID, pluginType).Where("pm_plugins.user = ?", userID)

	var total int64
	query.Count(&total)

	var plugins []model.PluginListItem
	err := query.Select("pm_plugins.id, pm_plugins.name, pm_plugins.desc_text, pm_plugins.icon, pm_plugins.type, pm_plugins.views, pm_plugins.downloads").
		Order("pm_plugins.created DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&plugins).Error

	fillPluginListExtra(plugins)

	return plugins, total, err
}

func GetMyPlugins(userID int, keywords string, pluginType *int, status *int, page, pageSize int) ([]model.PluginMyItem, int64, error) {
	query := database.DB.Model(&model.Plugin{}).Where("user = ?", userID)
	if keywords != "" {
		query = query.Where("name LIKE ? OR desc_text LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if pluginType != nil {
		query = query.Where("type = ?", *pluginType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var total int64
	query.Count(&total)

	var plugins []model.PluginMyItem
	err := query.Select("id, name, type, created, status, reject_msg").
		Order("created DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&plugins).Error
	return plugins, total, err
}

// GetPendingPlugins 获取待审核插件
func GetPendingPlugins(keywords string, pluginType *int, page, pageSize int) ([]model.Plugin, int64, error) {
	query := database.DB.Model(&model.Plugin{}).Where("status = 2")
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	if pluginType != nil {
		query = query.Where("type = ?", *pluginType)
	}

	var total int64
	query.Count(&total)

	var plugins []model.Plugin
	err := query.Order("updated DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&plugins).Error
	return plugins, total, err
}

// IsPluginStarred 检查用户是否已收藏插件
func IsPluginStarred(userID, pluginID int) bool {
	var count int64
	database.DB.Model(&model.UserStar{}).Where("user_id = ? AND plugin_id = ?", userID, pluginID).Count(&count)
	return count > 0
}
