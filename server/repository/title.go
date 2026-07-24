package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

func GetAllTitles() ([]model.Title, error) {
	var titles []model.Title
	err := database.DB.Order("id ASC").Find(&titles).Error
	return titles, err
}

func GetTitleList(keywords string, page, pageSize int) ([]model.Title, int64, error) {
	query := database.DB.Model(&model.Title{})
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	var total int64
	query.Count(&total)
	var titles []model.Title
	err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&titles).Error
	return titles, total, err
}

func GetTitleByID(id int) (*model.Title, error) {
	var title model.Title
	err := database.DB.First(&title, id).Error
	return &title, err
}

func CreateTitle(title *model.Title) error {
	return database.DB.Create(title).Error
}

func UpdateTitle(title *model.Title) error {
	return database.DB.Save(title).Error
}

func DeleteTitle(id int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("title_id = ?", id).Delete(&model.UserTitle{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.Title{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
