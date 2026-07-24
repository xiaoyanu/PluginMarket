package repository

import (
	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

func CreateNotification(notification *model.Notification) error {
	return database.DB.Create(notification).Error
}

func GetNotificationByID(id int64) (*model.Notification, error) {
	var notification model.Notification
	err := database.DB.First(&notification, id).Error
	return &notification, err
}

func UpdateNotification(notification *model.Notification) error {
	return database.DB.Save(notification).Error
}

func DeleteNotification(id int64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("notification_id = ?", id).Delete(&model.UserNotificationState{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.Notification{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func GetAdminNotificationList(keywords, notificationType string, status *int, page, pageSize int) ([]model.Notification, int64, error) {
	query := database.DB.Model(&model.Notification{})
	if keywords != "" {
		like := "%" + keywords + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Notification
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func notificationAudienceSQL(userID int, power int) (string, []interface{}) {
	return `(
		(n.audience_type = 'user' AND n.receiver_id = ?)
		OR (n.audience_type = 'all' AND (n.include_future_users = 1 OR u.created <= n.published_at))
		OR (n.audience_type = 'normal' AND ? = 0 AND (n.include_future_users = 1 OR u.created <= n.published_at))
		OR (n.audience_type = 'admin' AND ? = 1 AND (n.include_future_users = 1 OR u.created <= n.published_at))
	)`, []interface{}{userID, power, power}
}

func visibleNotificationQuery(userID, power int) *gorm.DB {
	audienceSQL, audienceArgs := notificationAudienceSQL(userID, power)
	return database.DB.Table("pm_notifications AS n").
		Select(`n.*, CASE WHEN s.is_read = 1 THEN 1 ELSE 0 END AS is_read`).
		Joins("LEFT JOIN pm_user_notification_states AS s ON s.notification_id = n.id AND s.user_id = ?", userID).
		Joins("INNER JOIN pm_users AS u ON u.id = ?", userID).
		Where("n.status = ? AND n.published_at IS NOT NULL AND n.published_at <= NOW()", 1).
		Where("(n.expires_at IS NULL OR n.expires_at > NOW())").
		Where(audienceSQL, audienceArgs...).
		Where("s.is_hidden IS NULL OR s.is_hidden = 0")
}

func GetVisibleNotifications(userID, power, page, pageSize int) ([]model.NotificationListItem, int64, error) {
	query := visibleNotificationQuery(userID, power)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.NotificationListItem
	err := query.Order("CASE WHEN s.is_read = 1 THEN 1 ELSE 0 END ASC").
		Order("n.published_at DESC, n.id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func GetUnreadNotificationCount(userID, power int) (int64, error) {
	query := visibleNotificationQuery(userID, power).Where("s.is_read IS NULL OR s.is_read = 0")
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func MarkNotificationRead(userID, power int, notificationID int64) error {
	var visibleID int64
	result := visibleNotificationQuery(userID, power).Where("n.id = ?", notificationID).Select("n.id").Scan(&visibleID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	state := &model.UserNotificationState{UserID: userID, NotificationID: notificationID, IsRead: true}
	return database.DB.Where("user_id = ? AND notification_id = ?", userID, notificationID).
		Assign(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")}).FirstOrCreate(state).Error
}

func MarkAllNotificationsRead(userID, power int) error {
	audienceSQL, audienceArgs := notificationAudienceSQL(userID, power)
	return database.DB.Exec(`
		INSERT INTO pm_user_notification_states (user_id, notification_id, is_read, read_at, created_at, updated_at)
		SELECT ?, n.id, 1, NOW(), NOW(), NOW()
		FROM pm_notifications AS n
		INNER JOIN pm_users AS u ON u.id = ?
		WHERE n.status = 1 AND n.published_at IS NOT NULL AND n.published_at <= NOW()
		  AND (n.expires_at IS NULL OR n.expires_at > NOW())
		  AND `+audienceSQL+`
		  AND (NOT EXISTS (
			SELECT 1 FROM pm_user_notification_states AS existing
			WHERE existing.user_id = ? AND existing.notification_id = n.id AND existing.is_hidden = 1
		  ))
		ON DUPLICATE KEY UPDATE is_read = 1, read_at = NOW(), updated_at = NOW()
	`, append([]interface{}{userID, userID}, append(audienceArgs, userID)...)...).Error
}

func HideNotification(userID, power int, notificationID int64) error {
	var visibleID int64
	result := visibleNotificationQuery(userID, power).Where("n.id = ?", notificationID).Select("n.id").Scan(&visibleID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	state := &model.UserNotificationState{UserID: userID, NotificationID: notificationID, IsHidden: true}
	return database.DB.Where("user_id = ? AND notification_id = ?", userID, notificationID).
		Assign(map[string]interface{}{"is_hidden": true}).FirstOrCreate(state).Error
}
