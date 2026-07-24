package model

import "time"

const (
	NotificationDefaultIcon  = "PhBellRinging"
	NotificationDefaultColor = "#64748B"
)

type Notification struct {
	ID                 int64      `json:"id" gorm:"type:bigint;primaryKey;autoIncrement"`
	Type               string     `json:"type" gorm:"size:32;not null"`
	AudienceType       string     `json:"audience_type" gorm:"size:20;not null;default:user;index:idx_audience_status_publish,priority:1"`
	ReceiverID         *int       `json:"receiver_id" gorm:"type:int;index:idx_receiver_status_publish,priority:1"`
	SenderID           *int       `json:"sender_id" gorm:"type:int"`
	Title              string     `json:"title" gorm:"size:255;not null"`
	Content            string     `json:"content" gorm:"type:text;not null"`
	IconName           string     `json:"icon_name" gorm:"size:64;not null;default:PhBellRinging"`
	IconColor          string     `json:"icon_color" gorm:"size:16;not null;default:#64748B"`
	SourceType         string     `json:"source_type" gorm:"size:32;index:idx_source,priority:1"`
	SourceID           *int64     `json:"source_id" gorm:"type:bigint;index:idx_source,priority:2"`
	ActionURL          string     `json:"action_url" gorm:"size:500"`
	ExtraData          string     `json:"extra_data" gorm:"type:text"`
	IncludeFutureUsers bool       `json:"include_future_users" gorm:"not null;default:false"`
	Status             int        `json:"status" gorm:"type:tinyint;not null;default:2;index:idx_audience_status_publish,priority:2;index:idx_receiver_status_publish,priority:2"`
	PublishedAt        *time.Time `json:"published_at" gorm:"index:idx_audience_status_publish,priority:3;index:idx_receiver_status_publish,priority:3"`
	ExpiresAt          *time.Time `json:"expires_at" gorm:"index:idx_expires_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Notification) TableName() string { return "pm_notifications" }

type UserNotificationState struct {
	UserID         int        `json:"user_id" gorm:"primaryKey"`
	NotificationID int64      `json:"notification_id" gorm:"primaryKey"`
	IsRead         bool       `json:"is_read" gorm:"not null;default:false"`
	IsHidden       bool       `json:"is_hidden" gorm:"not null;default:false"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (UserNotificationState) TableName() string { return "pm_user_notification_states" }

type NotificationListItem struct {
	Notification
	IsRead bool `json:"is_read" gorm:"column:is_read"`
}
