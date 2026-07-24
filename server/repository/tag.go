package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

func GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Order("id ASC").Find(&tags).Error
	return tags, err
}

func GetTagList(keywords string, page, pageSize int) ([]model.Tag, int64, error) {
	query := database.DB.Model(&model.Tag{})
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	var total int64
	query.Count(&total)
	var tags []model.Tag
	err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tags).Error
	return tags, total, err
}

func GetTagByID(id int) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.First(&tag, id).Error
	return &tag, err
}

func CreateTag(tag *model.Tag) error {
	return database.DB.Create(tag).Error
}

func UpdateTag(tag *model.Tag) error {
	return database.DB.Save(tag).Error
}

func DeleteTag(id int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tag_id = ?", id).Delete(&model.PluginTag{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.Tag{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
