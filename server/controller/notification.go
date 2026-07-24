package controller

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultNotificationColor = model.NotificationDefaultColor

var (
	notificationIconPattern  = regexp.MustCompile(`^Ph[A-Za-z0-9]+$`)
	notificationColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$`)
)

type NotificationVisibility struct {
	AudienceType       string
	ReceiverID         int
	IncludeFutureUsers bool
	PublishedAt        time.Time
}

type NotificationUserVisibility struct {
	ID        int
	Power     int
	CreatedAt time.Time
}

func notificationVisibleToUser(n NotificationVisibility, u NotificationUserVisibility) bool {
	switch n.AudienceType {
	case "user":
		return n.ReceiverID == u.ID
	case "all":
		return n.IncludeFutureUsers || !u.CreatedAt.After(n.PublishedAt)
	case "admin":
		return u.Power == 1 && (n.IncludeFutureUsers || !u.CreatedAt.After(n.PublishedAt))
	case "normal":
		return u.Power == 0 && (n.IncludeFutureUsers || !u.CreatedAt.After(n.PublishedAt))
	default:
		return false
	}
}

func validateNotificationIconName(name string) error {
	if len(name) > 64 || !notificationIconPattern.MatchString(name) {
		return errors.New("通知图标名称无效")
	}
	return nil
}

func normalizeNotificationColor(color string) string {
	color = strings.ToUpper(strings.TrimSpace(color))
	if !notificationColorPattern.MatchString(color) {
		return defaultNotificationColor
	}
	return color
}

type notificationRequest struct {
	Type               string          `json:"type" binding:"required"`
	AudienceType       string          `json:"audience_type" binding:"required"`
	ReceiverID         *int            `json:"receiver_id"`
	Title              string          `json:"title" binding:"required"`
	Content            string          `json:"content" binding:"required"`
	IconName           string          `json:"icon_name"`
	IconColor          string          `json:"icon_color"`
	SourceType         string          `json:"source_type"`
	SourceID           *int64          `json:"source_id"`
	ActionURL          string          `json:"action_url"`
	ExtraData          json.RawMessage `json:"extra_data"`
	IncludeFutureUsers bool            `json:"include_future_users"`
	Status             int             `json:"status"`
	PublishedAt        *time.Time      `json:"published_at"`
	ExpiresAt          *time.Time      `json:"expires_at"`
}

func (req *notificationRequest) normalizeAndValidate() error {
	req.Type = strings.TrimSpace(req.Type)
	req.AudienceType = strings.TrimSpace(req.AudienceType)
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.SourceType = strings.TrimSpace(req.SourceType)
	req.ActionURL = strings.TrimSpace(req.ActionURL)
	if req.IconName == "" {
		req.IconName = model.NotificationDefaultIcon
	}
	if req.IconColor == "" {
		req.IconColor = model.NotificationDefaultColor
	}
	if req.Type == "" || len(req.Type) > 32 {
		return errors.New("通知类型无效")
	}
	if req.Title == "" || len([]rune(req.Title)) > 255 {
		return errors.New("通知标题不能为空且不能超过255个字符")
	}
	if req.Content == "" {
		return errors.New("通知内容不能为空")
	}
	if err := validateNotificationIconName(req.IconName); err != nil {
		return err
	}
	if !notificationColorPattern.MatchString(strings.TrimSpace(req.IconColor)) {
		return errors.New("通知图标颜色无效")
	}
	req.IconColor = normalizeNotificationColor(req.IconColor)
	switch req.AudienceType {
	case "user":
		if req.ReceiverID == nil || *req.ReceiverID <= 0 {
			return errors.New("指定用户通知必须提供有效的接收用户ID")
		}
		req.IncludeFutureUsers = false
	case "all", "admin", "normal":
		req.ReceiverID = nil
	default:
		return errors.New("通知接收范围无效")
	}
	if req.Status < 0 || req.Status > 2 {
		return errors.New("通知状态无效")
	}
	if len(req.SourceType) > 32 || len(req.ActionURL) > 500 {
		return errors.New("通知来源或跳转地址过长")
	}
	if len(req.ExtraData) > 0 && !json.Valid(req.ExtraData) {
		return errors.New("扩展数据必须是有效JSON")
	}
	if req.Status == 1 && req.PublishedAt == nil {
		now := time.Now()
		req.PublishedAt = &now
	}
	if req.PublishedAt != nil && req.ExpiresAt != nil && !req.ExpiresAt.After(*req.PublishedAt) {
		return errors.New("过期时间必须晚于发布时间")
	}
	return nil
}

func notificationFromRequest(req notificationRequest, senderID *int) *model.Notification {
	return &model.Notification{
		Type: req.Type, AudienceType: req.AudienceType, ReceiverID: req.ReceiverID, SenderID: senderID,
		Title: req.Title, Content: req.Content, IconName: req.IconName, IconColor: req.IconColor,
		SourceType: req.SourceType, SourceID: req.SourceID, ActionURL: req.ActionURL,
		ExtraData: string(req.ExtraData), IncludeFutureUsers: req.IncludeFutureUsers,
		Status: req.Status, PublishedAt: req.PublishedAt, ExpiresAt: req.ExpiresAt,
	}
}

func AdminCreateNotification(c *gin.Context) {
	var req notificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if err := req.normalizeAndValidate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.ReceiverID != nil {
		if _, err := repository.GetUserByID(*req.ReceiverID); errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BadRequest(c, "接收用户不存在")
			return
		} else if err != nil {
			utils.ServerError(c, "查询接收用户失败")
			return
		}
	}
	senderID := c.GetInt("userId")
	notification := notificationFromRequest(req, &senderID)
	if err := repository.CreateNotification(notification); err != nil {
		utils.ServerError(c, "创建通知失败")
		return
	}
	utils.OKData(c, gin.H{"id": notification.ID})
}

func AdminGetNotificationList(c *gin.Context) {
	page, pageSize := getPageParams(c)
	list, total, err := repository.GetAdminNotificationList(strings.TrimSpace(c.Query("keywords")), strings.TrimSpace(c.Query("type")), getOptionalInt(c, "status"), page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询通知失败")
		return
	}
	utils.OKPage(c, list, total)
}

func AdminGetNotificationDetail(c *gin.Context) {
	id, ok := parseNotificationID(c)
	if !ok {
		return
	}
	notification, err := repository.GetNotificationByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "通知不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "查询通知失败")
		return
	}
	utils.OKData(c, notification)
}

func AdminUpdateNotification(c *gin.Context) {
	id, ok := parseNotificationID(c)
	if !ok {
		return
	}
	existing, err := repository.GetNotificationByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "通知不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "查询通知失败")
		return
	}
	var req notificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if err := req.normalizeAndValidate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.ReceiverID != nil {
		if _, err := repository.GetUserByID(*req.ReceiverID); errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BadRequest(c, "接收用户不存在")
			return
		} else if err != nil {
			utils.ServerError(c, "查询接收用户失败")
			return
		}
	}
	updated := notificationFromRequest(req, existing.SenderID)
	updated.ID, updated.CreatedAt = existing.ID, existing.CreatedAt
	if updated.PublishedAt == nil && existing.PublishedAt != nil {
		updated.PublishedAt = existing.PublishedAt
	}
	if err := repository.UpdateNotification(updated); err != nil {
		utils.ServerError(c, "更新通知失败")
		return
	}
	utils.OKMsg(c, "更新成功")
}

func AdminDeleteNotification(c *gin.Context) {
	id, ok := parseNotificationID(c)
	if !ok {
		return
	}
	if err := repository.DeleteNotification(id); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "通知不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "删除通知失败")
		return
	}
	utils.OKMsg(c, "删除成功")
}

func GetNotificationList(c *gin.Context) {
	page, pageSize := getPageParams(c)
	userID, power := c.GetInt("userId"), c.GetInt("power")
	list, total, err := repository.GetVisibleNotifications(userID, power, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询通知失败")
		return
	}
	unread, err := repository.GetUnreadNotificationCount(userID, power)
	if err != nil {
		utils.ServerError(c, "查询未读通知失败")
		return
	}
	utils.OKData(c, gin.H{"list": list, "total": total, "unread": unread})
}

func MarkNotificationRead(c *gin.Context) {
	id, ok := parseNotificationID(c)
	if !ok {
		return
	}
	if err := repository.MarkNotificationRead(c.GetInt("userId"), c.GetInt("power"), id); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "通知不存在或无权访问")
		return
	} else if err != nil {
		utils.ServerError(c, "标记已读失败")
		return
	}
	utils.OKMsg(c, "已标记为已读")
}

func MarkAllNotificationsRead(c *gin.Context) {
	if err := repository.MarkAllNotificationsRead(c.GetInt("userId"), c.GetInt("power")); err != nil {
		utils.ServerError(c, "全部已读失败")
		return
	}
	utils.OKMsg(c, "已全部标记为已读")
}

func HideNotification(c *gin.Context) {
	id, ok := parseNotificationID(c)
	if !ok {
		return
	}
	if err := repository.HideNotification(c.GetInt("userId"), c.GetInt("power"), id); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "通知不存在或无权访问")
		return
	} else if err != nil {
		utils.ServerError(c, "隐藏通知失败")
		return
	}
	utils.OKMsg(c, "通知已隐藏")
}

func parseNotificationID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(c, "通知ID无效")
		return 0, false
	}
	return id, true
}
