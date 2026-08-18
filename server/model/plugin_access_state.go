package model

import "time"

// PluginAccessState 插件浏览/下载去重状态
type PluginAccessState struct {
	PluginID       int       `json:"plugin_id" gorm:"primaryKey;column:plugin_id"`
	LastViewAt     time.Time `json:"last_view_at" gorm:"column:last_view_at"`
	LastViewIP     string    `json:"last_view_ip" gorm:"size:45;column:last_view_ip"`
	LastDownloadAt time.Time `json:"last_download_at" gorm:"column:last_download_at"`
	LastDownloadIP string    `json:"last_download_ip" gorm:"size:45;column:last_download_ip"`
}

func (PluginAccessState) TableName() string {
	return "pm_plugin_access_state"
}
