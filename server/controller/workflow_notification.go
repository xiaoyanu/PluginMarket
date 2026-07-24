package controller

import (
	"strings"
	"time"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
)

type workflowNotificationKind string

const (
	workflowNewComment     workflowNotificationKind = "new_comment"
	workflowReplyComment   workflowNotificationKind = "comment_reply"
	workflowPluginApproved workflowNotificationKind = "plugin_approved"
	workflowPluginRejected workflowNotificationKind = "plugin_rejected"
)

type workflowNotificationTemplateKeys struct {
	title   string
	content string
}

var workflowNotificationTemplates = map[workflowNotificationKind]workflowNotificationTemplateKeys{
	workflowNewComment:     {title: "notificationNewCommentTitle", content: "notificationNewCommentContent"},
	workflowReplyComment:   {title: "notificationReplyCommentTitle", content: "notificationReplyCommentContent"},
	workflowPluginApproved: {title: "notificationPluginApprovedTitle", content: "notificationPluginApprovedContent"},
	workflowPluginRejected: {title: "notificationPluginRejectedTitle", content: "notificationPluginRejectedContent"},
}

func renderWorkflowNotificationTemplate(template string, values map[string]string) string {
	result := template
	for key, value := range values {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

func buildWorkflowNotification(kind workflowNotificationKind, settings map[string]string, receiverID, senderID int, pluginID int64, pluginName, userName, content, actionURL string) (*model.Notification, bool) {
	keys, exists := workflowNotificationTemplates[kind]
	if !exists {
		return nil, false
	}
	titleTemplate := strings.TrimSpace(settings[keys.title])
	contentTemplate := strings.TrimSpace(settings[keys.content])
	if titleTemplate == "" || contentTemplate == "" {
		return nil, false
	}
	values := map[string]string{
		"user_name":   userName,
		"content":     content,
		"plugin_name": pluginName,
	}
	receiver := receiverID
	var sender *int
	if senderID > 0 {
		sender = &senderID
	}
	now := time.Now()
	return &model.Notification{
		Type:         string(kind),
		AudienceType: "user",
		ReceiverID:   &receiver,
		SenderID:     sender,
		Title:        renderWorkflowNotificationTemplate(titleTemplate, values),
		Content:      renderWorkflowNotificationTemplate(contentTemplate, values),
		IconName:     workflowNotificationIcon(kind),
		IconColor:    workflowNotificationColor(kind),
		SourceType:   "plugin",
		SourceID:     &pluginID,
		ActionURL:    actionURL,
		Status:       1,
		PublishedAt:  &now,
	}, true
}

func workflowNotificationIcon(kind workflowNotificationKind) string {
	switch kind {
	case workflowNewComment:
		return "PhChatText"
	case workflowReplyComment:
		return "PhChats"
	case workflowPluginApproved:
		return "PhSealCheck"
	case workflowPluginRejected:
		return "PhSealWarning"
	default:
		return model.NotificationDefaultIcon
	}
}

func workflowNotificationColor(kind workflowNotificationKind) string {
	switch kind {
	case workflowNewComment, workflowReplyComment:
		return "#71BBF4"
	case workflowPluginApproved:
		return "#27DD6D"
	case workflowPluginRejected:
		return "#FF3B0A"
	default:
		return model.NotificationDefaultColor
	}
}

func createWorkflowNotification(kind workflowNotificationKind, receiverID, senderID int, pluginID int64, pluginName, userName, content, actionURL string) {
	notification, ok := buildWorkflowNotification(kind, repository.GetAllSettings(), receiverID, senderID, pluginID, pluginName, userName, content, actionURL)
	if !ok {
		return
	}
	_ = repository.CreateNotification(notification)
}
