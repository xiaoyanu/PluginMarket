package model

import "time"

type Title struct {
	ID      int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string    `json:"name" gorm:"size:255"`
	Icon    string    `json:"icon" gorm:"size:255"`
	Created time.Time `json:"created" gorm:"autoCreateTime"`
	Updated time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (Title) TableName() string {
	return "pm_titles"
}

// UserTitle 用户-称号关联
type UserTitle struct {
	UserID  int `json:"user_id" gorm:"primaryKey"`
	TitleID int `json:"title_id" gorm:"primaryKey"`
}

func (UserTitle) TableName() string {
	return "pm_user_titles"
}
