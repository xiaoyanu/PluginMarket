package repository

import (
	"strings"
	"time"

	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const pluginAccessWindow = 10 * time.Minute

type pluginAccessType string

const (
	pluginAccessView     pluginAccessType = "view"
	pluginAccessDownload pluginAccessType = "download"
)

func shouldCountPluginAccess(lastAt time.Time, lastIP, currentIP string, now time.Time, window time.Duration) bool {
	if lastAt.IsZero() {
		return true
	}
	if strings.TrimSpace(lastIP) != strings.TrimSpace(currentIP) {
		return true
	}
	return now.Sub(lastAt) > window
}

func RecordPluginViewAccess(pluginID int, ip string, enabled bool) error {
	if !enabled {
		return IncrementPluginViews(pluginID)
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		return recordPluginAccessAt(tx, pluginID, ip, pluginAccessView, time.Now())
	})
}

func RecordPluginDownloadAccess(pluginID int, ip string, enabled bool) error {
	if !enabled {
		return IncrementPluginDownloads(pluginID)
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		return recordPluginAccessAt(tx, pluginID, ip, pluginAccessDownload, time.Now())
	})
}

func recordPluginAccessAt(tx *gorm.DB, pluginID int, ip string, accessType pluginAccessType, now time.Time) error {
	if err := ensurePluginAccessState(tx, pluginID); err != nil {
		return err
	}

	var state model.PluginAccessState
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&state, "plugin_id = ?", pluginID).Error; err != nil {
		return err
	}

	ip = strings.TrimSpace(ip)
	switch accessType {
	case pluginAccessView:
		if !shouldCountPluginAccess(state.LastViewAt, state.LastViewIP, ip, now, pluginAccessWindow) {
			return nil
		}
		if err := incrementPluginCounter(tx, pluginID, "views"); err != nil {
			return err
		}
		return tx.Model(&model.PluginAccessState{}).Where("plugin_id = ?", pluginID).Updates(map[string]interface{}{
			"last_view_at": now,
			"last_view_ip": ip,
		}).Error
	case pluginAccessDownload:
		if !shouldCountPluginAccess(state.LastDownloadAt, state.LastDownloadIP, ip, now, pluginAccessWindow) {
			return nil
		}
		if err := incrementPluginCounter(tx, pluginID, "downloads"); err != nil {
			return err
		}
		return tx.Model(&model.PluginAccessState{}).Where("plugin_id = ?", pluginID).Updates(map[string]interface{}{
			"last_download_at": now,
			"last_download_ip": ip,
		}).Error
	default:
		return gorm.ErrInvalidData
	}
}

func ensurePluginAccessState(tx *gorm.DB, pluginID int) error {
	state := &model.PluginAccessState{PluginID: pluginID}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}},
		DoNothing: true,
	}).Create(state).Error
}

func incrementPluginCounter(tx *gorm.DB, pluginID int, column string) error {
	result := tx.Model(&model.Plugin{}).Where("id = ?", pluginID).UpdateColumn(column, gorm.Expr(column+" + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
