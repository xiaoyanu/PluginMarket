package model

import "time"

type Frame struct {
	ID      int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string    `json:"name" gorm:"size:255"`
	Icon    string    `json:"icon" gorm:"size:255"`
	Created time.Time `json:"created" gorm:"autoCreateTime"`
	Updated time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (Frame) TableName() string {
	return "pm_frames"
}
