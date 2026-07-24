package model

import "time"

type Tag struct {
	ID      int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string    `json:"name" gorm:"size:255"`
	Created time.Time `json:"created" gorm:"autoCreateTime"`
	Updated time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (Tag) TableName() string {
	return "pm_tags"
}
