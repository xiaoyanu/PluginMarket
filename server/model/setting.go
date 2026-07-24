package model

// Setting 全局设置（键值对）
type Setting struct {
	Key   string `json:"s_key" gorm:"primaryKey;size:100;column:s_key"`
	Value string `json:"s_value" gorm:"type:longtext;column:s_value"`
}

func (Setting) TableName() string {
	return "pm_settings"
}
