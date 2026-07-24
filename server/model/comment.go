package model

import "time"

type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Plugin    int       `json:"plugin" gorm:"index:idx_plugin"`
	User      int       `json:"user"`
	Content   string    `json:"content" gorm:"type:text"`
	Root      int       `json:"root" gorm:"index:idx_root;default:0"`
	Parent    int       `json:"parent" gorm:"default:0"`
	ReplyUser int       `json:"reply_user"`
	Created   time.Time `json:"created" gorm:"autoCreateTime"`
	Updated   time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (Comment) TableName() string {
	return "pm_comments"
}

// CommentListItem 评论列表项（含作者信息和回复）
type CommentListItem struct {
	ID         int            `json:"id"`
	Content    string         `json:"content"`
	Root       int            `json:"root"`
	Parent     int            `json:"parent"`
	Created    time.Time      `json:"created"`
	Author     CommentAuthor  `json:"author"`
	Replies    []CommentReply `json:"replies" gorm:"-"`
	ReplyCount int            `json:"replyCount" gorm:"-"`
}

type CommentAuthor struct {
	ID     int     `json:"id"`
	Nick   string  `json:"nick"`
	Avatar string  `json:"avatar"`
	Titles []Title `json:"titles" gorm:"many2many:pm_user_titles;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:title_id"`
}

type CommentReply struct {
	ID        int           `json:"id"`
	Content   string        `json:"content"`
	Root      int           `json:"root"`
	Parent    int           `json:"parent"`
	ReplyUser int           `json:"reply_user"`
	Created   time.Time     `json:"created"`
	Author    CommentAuthor `json:"author"`
	ReplyTo   *ReplyTo      `json:"replyTo,omitempty"`
}

type ReplyTo struct {
	ID   int    `json:"id"`
	Nick string `json:"nick"`
}
