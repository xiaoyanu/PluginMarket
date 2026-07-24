package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

func GetAllFrames() ([]model.Frame, error) {
	var frames []model.Frame
	err := database.DB.Order("id ASC").Find(&frames).Error
	return frames, err
}

func GetFrameList(keywords string, page, pageSize int) ([]model.Frame, int64, error) {
	query := database.DB.Model(&model.Frame{})
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	var total int64
	query.Count(&total)
	var frames []model.Frame
	err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&frames).Error
	return frames, total, err
}

func GetFrameByID(id int) (*model.Frame, error) {
	var frame model.Frame
	err := database.DB.First(&frame, id).Error
	return &frame, err
}

func CreateFrame(frame *model.Frame) error {
	return database.DB.Create(frame).Error
}

func UpdateFrame(frame *model.Frame) error {
	return database.DB.Save(frame).Error
}

func DeleteFrame(id int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("frame_id = ?", id).Delete(&model.PluginFrame{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.Frame{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
