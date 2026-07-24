package database

import (
	"fmt"
	"time"

	"pluginmarket-server/model"
)

const boundEmailUniqueIndex = "uq_pm_users_bound_email"

// BackfillTagTimestamps 为升级前已有的标签补齐新增的日期字段。
func BackfillTagTimestamps(now time.Time) error {
	return DB.Model(&model.Tag{}).
		Where("created = ? OR created IS NULL", time.Time{}).
		Updates(map[string]interface{}{"created": now, "updated": now}).Error
}

// EnsureUniqueBoundEmails 为非空邮箱建立唯一索引。
// MySQL 允许 UNIQUE 索引中存在多个 NULL，因此未绑定邮箱统一从空字符串迁移为 NULL。
func EnsureUniqueBoundEmails() error {
	var duplicateCount int64
	if err := DB.Raw(`
		SELECT COUNT(*)
		FROM (
			SELECT LOWER(TRIM(email)) AS normalized_email
			FROM pm_users
			WHERE TRIM(COALESCE(email, '')) <> ''
			GROUP BY LOWER(TRIM(email))
			HAVING COUNT(*) > 1
		) duplicates
	`).Scan(&duplicateCount).Error; err != nil {
		return err
	}
	if duplicateCount > 0 {
		return fmt.Errorf("存在 %d 组重复绑定邮箱，请先处理后再创建唯一索引", duplicateCount)
	}

	if err := DB.Exec("ALTER TABLE pm_users MODIFY COLUMN email VARCHAR(191) NULL DEFAULT NULL").Error; err != nil {
		return err
	}
	if err := DB.Exec("UPDATE pm_users SET email = NULL WHERE TRIM(COALESCE(email, '')) = ''").Error; err != nil {
		return err
	}
	if err := DB.Exec("UPDATE pm_users SET email = LOWER(TRIM(email)) WHERE email IS NOT NULL").Error; err != nil {
		return err
	}
	if DB.Migrator().HasIndex(&model.User{}, boundEmailUniqueIndex) {
		return nil
	}
	return DB.Exec("CREATE UNIQUE INDEX `" + boundEmailUniqueIndex + "` ON `pm_users` (`email`)").Error
}
