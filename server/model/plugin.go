package model

import "time"

type Plugin struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255"`
	DescText  string    `json:"desc_text" gorm:"size:255"`
	Type      int       `json:"type"`
	User      int       `json:"user"`
	URL       string    `json:"url" gorm:"size:255"`
	URLCode   string    `json:"url_code" gorm:"size:255"`
	Icon      string    `json:"icon" gorm:"size:255"`
	Status    int       `json:"status" gorm:"default:0"`
	Views     int       `json:"views" gorm:"default:0"`
	Downloads int       `json:"downloads" gorm:"default:0"`
	Content   string    `json:"content" gorm:"type:longtext"`
	RejectMsg string    `json:"reject_msg" gorm:"size:255"`
	Created   time.Time `json:"created" gorm:"autoCreateTime"`
	Updated   time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (Plugin) TableName() string {
	return "pm_plugins"
}

// PluginListItem 插件列表项
type PluginListItem struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	DescText  string  `json:"desc_text"`
	Icon      string  `json:"icon"`
	Type      int     `json:"type"`
	Views     int     `json:"views"`
	Downloads int     `json:"downloads"`
	Star      int     `json:"star"`
	HotScore  float64 `json:"hot_score"`
	FrameIcon string  `json:"frame_icon"`
}

// PluginMyItem 我的插件列表项
type PluginMyItem struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	Created   time.Time `json:"created"`
	Status    int       `json:"status"`
	RejectMsg string    `json:"reject_msg"`
}

// PluginDetail 插件详情
type PluginDetail struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	DescText   string       `json:"desc_text"`
	Content    string       `json:"content"`
	Icon       string       `json:"icon"`
	Type       int          `json:"type"`
	Status     int          `json:"status"`
	Views      int          `json:"views"`
	Downloads  int          `json:"downloads"`
	Star       int          `json:"star"`
	URL        string       `json:"url"`
	URLCode    string       `json:"url_code"`
	Created    time.Time    `json:"created"`
	Updated    time.Time    `json:"updated"`
	Author     PluginAuthor `json:"author"`
	Frameworks []Frame      `json:"frameworks" gorm:"many2many:pm_plugin_frames;foreignKey:ID;joinForeignKey:plugin_id;References:ID;joinReferences:frame_id"`
	Tags       []Tag        `json:"tags" gorm:"many2many:pm_plugin_tags;foreignKey:ID;joinForeignKey:plugin_id;References:ID;joinReferences:tag_id"`
	IsStarred  bool         `json:"is_starred" gorm:"-"`
}

type PluginAuthor struct {
	ID       int     `json:"id"`
	Nick     string  `json:"nick"`
	Avatar   string  `json:"avatar"`
	Userdesc string  `json:"userdesc"`
	Titles   []Title `json:"titles" gorm:"many2many:pm_user_titles;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:title_id"`
}

// PluginFrame 插件-框架关联
type PluginFrame struct {
	PluginID int `json:"plugin_id" gorm:"primaryKey"`
	FrameID  int `json:"frame_id" gorm:"primaryKey"`
}

func (PluginFrame) TableName() string {
	return "pm_plugin_frames"
}

// PluginTag 插件-标签关联
type PluginTag struct {
	PluginID int `json:"plugin_id" gorm:"primaryKey"`
	TagID    int `json:"tag_id" gorm:"primaryKey"`
}

func (PluginTag) TableName() string {
	return "pm_plugin_tags"
}
