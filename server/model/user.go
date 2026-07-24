package model

import "time"

const (
	DeletedUserNick        = "已注销用户"
	DeletedUserDescription = "此用户已注销账号"
)

type User struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string    `json:"username" gorm:"size:255;not null"`
	Password string    `json:"-" gorm:"size:255;not null"`
	Nick     string    `json:"nick" gorm:"size:255"`
	Email    string    `json:"email" gorm:"size:191;default:null"`
	Userdesc string    `json:"userdesc" gorm:"size:255"`
	Avatar   string    `json:"avatar" gorm:"size:255"`
	Power    int       `json:"power"`
	IsDelete bool      `json:"is_delete" gorm:"not null;default:false;index"`
	Created  time.Time `json:"created" gorm:"autoCreateTime"`
	Updated  time.Time `json:"updated" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "pm_users"
}

// AnonymizeIfDeleted 统一隐藏已注销用户的公开资料。
func (u *User) AnonymizeIfDeleted() {
	if !u.IsDelete {
		return
	}
	u.Nick = DeletedUserNick
	u.Avatar = ""
	u.Userdesc = DeletedUserDescription
}

// UserPublic 用户公开信息（不含 username、email 等隐私字段）
type UserPublic struct {
	ID       int       `json:"id"`
	Nick     string    `json:"nick"`
	Avatar   string    `json:"avatar"`
	Userdesc string    `json:"userdesc"`
	Power    int       `json:"power"`
	IsDelete bool      `json:"is_delete"`
	Titles   []Title   `json:"titles" gorm:"many2many:pm_user_titles;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:title_id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Public 返回只包含公开字段的用户信息。
func (u UserWithTitle) Public() UserPublic {
	return UserPublic{
		ID:       u.ID,
		Nick:     u.Nick,
		Avatar:   u.Avatar,
		Userdesc: u.Userdesc,
		Power:    u.Power,
		IsDelete: u.IsDelete,
		Titles:   u.Titles,
		Created:  u.Created,
		Updated:  u.Updated,
	}
}

// UserWithTitle 用户信息（含称号）
type UserWithTitle struct {
	User
	Titles []Title `json:"titles" gorm:"many2many:pm_user_titles;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:title_id"`
}
