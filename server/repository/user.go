package repository

import (
	"errors"

	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrEmailAlreadyBound = errors.New("邮箱已被其他用户绑定")

func GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	user.AnonymizeIfDeleted()
	return &user, err
}

func GetUserByIDWithTitles(id int) (*model.UserWithTitle, error) {
	var user model.UserWithTitle
	err := database.DB.Preload("Titles").First(&user, id).Error
	user.User.AnonymizeIfDeleted()
	return &user, err
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func UpdateUserFields(id int, fields map[string]interface{}) error {
	result := database.DB.Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// BindUserEmail 依赖数据库唯一索引原子绑定邮箱，避免并发重复绑定。
func BindUserEmail(id int, email string) error {
	return normalizeEmailBindingError(UpdateUserFields(id, map[string]interface{}{"email": email}))
}

func normalizeEmailBindingError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrEmailAlreadyBound
	}
	return err
}

func DeleteUser(id int) error {
	return UpdateUserFields(id, map[string]interface{}{
		"is_delete": true,
		"nick":      model.DeletedUserNick,
		"avatar":    "",
		"userdesc":  model.DeletedUserDescription,
		"power":     0,
	})
}

func GetUserList(keywords string, power *int, isDelete *bool, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.DB.Model(&model.User{})
	if keywords != "" {
		query = query.Where("username LIKE ? OR nick LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if power != nil {
		query = query.Where("power = ?", *power)
	}
	if isDelete != nil {
		query = query.Where("is_delete = ?", *isDelete)
	}

	query.Count(&total)
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	for i := range users {
		users[i].AnonymizeIfDeleted()
	}
	return users, total, err
}

func GetAdminUsers() ([]model.User, error) {
	var users []model.User
	err := database.DB.Where("power = ? AND is_delete = ?", 1, false).Find(&users).Error
	return users, err
}

// GetUserTitles 获取用户的称号ID列表
func GetUserTitles(userID int) []int {
	var titles []model.UserTitle
	database.DB.Where("user_id = ?", userID).Find(&titles)
	ids := make([]int, len(titles))
	for i, t := range titles {
		ids[i] = t.TitleID
	}
	return ids
}

// SetUserTitles 全量替换用户称号
func SetUserTitles(userID int, titleIDs []int) error {
	// 先删除旧的
	if err := database.DB.Where("user_id = ?", userID).Delete(&model.UserTitle{}).Error; err != nil {
		return err
	}
	// 插入新的
	for _, tid := range titleIDs {
		ut := model.UserTitle{UserID: userID, TitleID: tid}
		if err := database.DB.Create(&ut).Error; err != nil {
			return err
		}
	}
	return nil
}
