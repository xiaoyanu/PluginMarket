package utils

import (
	"fmt"
	"html"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"pluginmarket-server/repository"

	"gopkg.in/gomail.v2"
)

const SMTPPasswordMask = "**********"

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string
	FromName string
}

var allowedEmailVariables = map[string]struct{}{
	"site_title":      {},
	"plugin_name":     {},
	"comment_user":    {},
	"comment_content": {},
	"reply_user":      {},
	"reply_content":   {},
	"user_name":       {},
	"content":         {},
	"link":            {},
	"time":            {},
	"user":            {},
}

func BuildSMTPConfig(settings map[string]string) (SMTPConfig, error) {
	config := SMTPConfig{
		Host:     strings.TrimSpace(settings["emailHost"]),
		User:     strings.TrimSpace(settings["emailUser"]),
		Pass:     settings["emailPass"],
		From:     strings.TrimSpace(settings["emailFrom"]),
		FromName: strings.TrimSpace(settings["emailFromName"]),
	}
	port, err := strconv.Atoi(strings.TrimSpace(settings["emailPort"]))
	if err != nil || port < 1 || port > 65535 {
		return SMTPConfig{}, fmt.Errorf("SMTP端口无效")
	}
	config.Port = port
	if config.Host == "" || config.User == "" || config.Pass == "" {
		return SMTPConfig{}, fmt.Errorf("SMTP服务器、账号和授权码不能为空")
	}
	if config.From == "" {
		config.From = config.User
	}
	address, err := mail.ParseAddress(config.From)
	if err != nil {
		return SMTPConfig{}, fmt.Errorf("发件邮箱格式不正确")
	}
	config.From = address.Address
	if config.FromName == "" {
		config.FromName = "PluginMarket"
	}
	return config, nil
}

func RenderEmailTemplate(template string, values map[string]string) string {
	result := template
	for key := range allowedEmailVariables {
		value, exists := values[key]
		if !exists {
			continue
		}
		result = strings.ReplaceAll(result, "{{"+key+"}}", html.EscapeString(value))
	}
	return result
}

func SendEmail(to, subject, body string) error {
	settings := repository.GetAllSettings()
	config, err := BuildSMTPConfig(settings)
	if err != nil {
		return err
	}
	return SendEmailWithConfig(config, to, subject, body)
}

func SendEmailAsync(to, subject, body string) error {
	settings := repository.GetAllSettings()
	config, err := BuildSMTPConfig(settings)
	if err != nil {
		return err
	}
	if _, err := mail.ParseAddress(strings.TrimSpace(to)); err != nil {
		return fmt.Errorf("收件邮箱格式不正确")
	}
	go func() {
		if err := SendEmailWithConfig(config, to, subject, body); err != nil {
			fmt.Printf("发送邮件到 %s 失败: %v\n", to, err)
		}
	}()
	return nil
}

func SendEmailWithConfig(config SMTPConfig, to, subject, body string) error {
	address, err := mail.ParseAddress(strings.TrimSpace(to))
	if err != nil {
		return fmt.Errorf("收件邮箱格式不正确")
	}

	message := gomail.NewMessage()
	defer message.Reset()
	message.SetHeader("From", message.FormatAddress(config.From, config.FromName))
	message.SetHeader("To", address.Address)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Pass)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("SMTP连接或发送失败: %w", err)
	}
	return nil
}

func SendTemplateEmail(to, titleKey, bodyKey string, values map[string]string) error {
	settings := repository.GetAllSettings()
	title := strings.TrimSpace(settings[titleKey])
	body := strings.TrimSpace(settings[bodyKey])
	if title == "" || body == "" {
		return fmt.Errorf("邮件模板未配置")
	}
	return SendEmail(to, RenderEmailTemplate(title, values), RenderEmailTemplate(body, values))
}

func SendTemplateEmailAsync(to, titleKey, bodyKey string, values map[string]string) error {
	settings := repository.GetAllSettings()
	title := strings.TrimSpace(settings[titleKey])
	body := strings.TrimSpace(settings[bodyKey])
	if title == "" || body == "" {
		return nil
	}
	return SendEmailAsync(to, RenderEmailTemplate(title, values), RenderEmailTemplate(body, values))
}

func SendResetPasswordEmail(to, userName, resetURL string) error {
	return SendTemplateEmail(to, "resetPasswordTitle", "resetPasswordBody", map[string]string{
		"site_title": siteTitle(),
		"user_name":  userName,
		"link":       resetURL,
	})
}

func SendEmailVerifyEmail(to, userName, verifyURL string) error {
	return SendTemplateEmail(to, "emailVerifyTitle", "emailVerifyBody", map[string]string{
		"site_title": siteTitle(),
		"user_name":  userName,
		"link":       verifyURL,
	})
}

func NewCommentEmailValues(siteTitle, pluginName, commentUser, commentContent, link string, commentTime time.Time) map[string]string {
	return map[string]string{
		"site_title":      siteTitle,
		"plugin_name":     pluginName,
		"comment_user":    commentUser,
		"comment_content": commentContent,
		"link":            link,
		"time":            commentTime.Format("2006-01-02 15:04:05"),
	}
}

func SendNewCommentEmail(to, pluginName, commentUser, commentContent, link string, commentTime time.Time) error {
	return SendTemplateEmailAsync(to, "newCommentTitle", "newCommentBody", NewCommentEmailValues(
		siteTitle(), pluginName, commentUser, commentContent, link, commentTime,
	))
}

func ReplyCommentEmailValues(siteTitle, pluginName, replyUser, replyContent, repliedUser, repliedContent, link string, replyTime time.Time) map[string]string {
	return map[string]string{
		"site_title":    siteTitle,
		"plugin_name":   pluginName,
		"reply_user":    replyUser,
		"reply_content": replyContent,
		"user":          repliedUser,
		"content":       repliedContent,
		"link":          link,
		"time":          replyTime.Format("2006-01-02 15:04:05"),
	}
}

func SendReplyCommentEmail(to, pluginName, replyUser, replyContent, repliedUser, repliedContent, link string, replyTime time.Time) error {
	return SendTemplateEmailAsync(to, "replyCommentTitle", "replyCommentBody", ReplyCommentEmailValues(
		siteTitle(), pluginName, replyUser, replyContent, repliedUser, repliedContent, link, replyTime,
	))
}

func siteTitle() string {
	if title := repository.GetSetting("siteTitle"); title != "" {
		return title
	}
	return "PluginMarket"
}
