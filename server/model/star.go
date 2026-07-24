package model

// UserStar 用户收藏
type UserStar struct {
	UserID   int `json:"user_id" gorm:"primaryKey"`
	PluginID int `json:"plugin_id" gorm:"primaryKey"`
}

func (UserStar) TableName() string {
	return "pm_user_stars"
}
