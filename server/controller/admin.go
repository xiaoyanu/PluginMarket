package controller

import (
	"errors"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pluginmarket-server/config"
	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ===== 插件审核 =====

// GetPendingPlugins 获取待审核插件列表
// GET /api/admin/plugin/pending
func GetPendingPlugins(c *gin.Context) {
	keywords := c.Query("keywords")
	pluginType := getOptionalInt(c, "type")
	page, pageSize := getPageParams(c)

	list, total, err := repository.GetPendingPlugins(keywords, pluginType, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// ApprovePlugin 审核通过
// PUT /api/admin/plugin/:id/approve
func ApprovePlugin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "操作失败")
		return
	}

	if err := repository.UpdatePluginFields(id, map[string]interface{}{
		"status":     0,
		"reject_msg": "",
	}); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "操作失败")
		return
	}

	sendPluginReviewNotification(pluginApprovedEmail, plugin, "")
	utils.OKMsg(c, "审核通过")
}

// RejectPlugin 审核拒绝
// PUT /api/admin/plugin/:id/reject
func RejectPlugin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入拒绝原因")
		return
	}
	reason, err := normalizeRejectReason(req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "操作失败")
		return
	}

	if err := repository.UpdatePluginFields(id, map[string]interface{}{
		"status":     1,
		"reject_msg": reason,
	}); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "操作失败")
		return
	}

	sendPluginReviewNotification(pluginRejectedEmail, plugin, reason)
	utils.OKMsg(c, "已拒绝")
}

func normalizeRejectReason(reason string) (string, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return "", errors.New("请输入拒绝原因")
	}
	if len([]rune(reason)) > 255 {
		return "", errors.New("拒绝原因不能超过255个字符")
	}
	return reason, nil
}

type pluginReviewEmailKind string

const (
	pluginPendingReviewEmail pluginReviewEmailKind = "pending"
	pluginApprovedEmail      pluginReviewEmailKind = "approved"
	pluginRejectedEmail      pluginReviewEmailKind = "rejected"
)

func buildPluginReviewEmail(kind pluginReviewEmailKind, pluginID int, pluginName, reason, baseURL string) (string, string, map[string]string) {
	values := map[string]string{"plugin_name": pluginName}
	switch kind {
	case pluginPendingReviewEmail:
		return "pendingPluginReviewTitle", "pendingPluginReviewBody", values
	case pluginApprovedEmail:
		return "pluginApprovedTitle", "pluginApprovedBody", values
	case pluginRejectedEmail:
		values["content"] = reason
		values["link"] = strings.TrimRight(baseURL, "/") + "/plugin/" + strconv.Itoa(pluginID)
		return "pluginRejectedTitle", "pluginRejectedBody", values
	default:
		return "", "", values
	}
}

func boundEmailAddresses(users []model.User) []string {
	emails := make([]string, 0, len(users))
	for _, user := range users {
		if email := strings.TrimSpace(user.Email); email != "" {
			emails = append(emails, email)
		}
	}
	return emails
}

func sendPendingPluginReviewNotification(plugin *model.Plugin) {
	admins, err := repository.GetAdminUsers()
	if err != nil {
		return
	}
	titleKey, bodyKey, values := buildPluginReviewEmail(pluginPendingReviewEmail, plugin.ID, plugin.Name, "", webURL())
	values["site_title"] = notificationSiteTitle()
	for _, email := range boundEmailAddresses(admins) {
		_ = utils.SendTemplateEmailAsync(email, titleKey, bodyKey, values)
	}
}

func sendPluginReviewNotification(kind pluginReviewEmailKind, plugin *model.Plugin, reason string) {
	notificationKind := workflowPluginApproved
	if kind == pluginRejectedEmail {
		notificationKind = workflowPluginRejected
	}
	createWorkflowNotification(notificationKind, plugin.User, 0, int64(plugin.ID), plugin.Name, "", reason, "/plugin/"+strconv.Itoa(plugin.ID))

	author, err := repository.GetUserByID(plugin.User)
	if err != nil || strings.TrimSpace(author.Email) == "" {
		return
	}
	titleKey, bodyKey, values := buildPluginReviewEmail(kind, plugin.ID, plugin.Name, reason, webURL())
	values["site_title"] = notificationSiteTitle()
	_ = utils.SendTemplateEmailAsync(strings.TrimSpace(author.Email), titleKey, bodyKey, values)
}

func notificationSiteTitle() string {
	if title := strings.TrimSpace(repository.GetSetting("siteTitle")); title != "" {
		return title
	}
	return "PluginMarket"
}

// AdminDeletePlugin 管理员删除插件
// DELETE /api/admin/plugin/:id
func AdminDeletePlugin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	values, err := repository.GetPluginUploadReferenceValues(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	if err := repository.DeletePlugin(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	if err := deleteUnreferencedPluginUploads(config.C.Uploads.Path, extractUploadReferences(values), loadCurrentUploadReferences); err != nil {
		log.Printf("插件 %d 已删除，关联图片清理失败: %v", id, err)
	}

	utils.OKMsg(c, "删除成功")
}

func validateAdminPluginStatus(status int) error {
	if status != 0 && status != 2 {
		return errors.New("插件状态无效")
	}
	return nil
}

func adminPluginStatusNotification(previousStatus, nextStatus int) (pluginReviewEmailKind, bool) {
	if previousStatus != 0 && nextStatus == 0 {
		return pluginApprovedEmail, true
	}
	return "", false
}

// AdminUpdatePluginStatus 管理员将插件设为审核通过或待审核。
func AdminUpdatePluginStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		utils.BadRequest(c, "参数错误")
		return
	}
	var req struct {
		Status *int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if err := validateAdminPluginStatus(*req.Status); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	plugin, err := repository.GetPluginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	fields := map[string]interface{}{"status": *req.Status}
	if *req.Status == 0 {
		fields["reject_msg"] = ""
	}
	if err := repository.UpdatePluginFields(id, fields); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "插件不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	if kind, ok := adminPluginStatusNotification(plugin.Status, *req.Status); ok {
		sendPluginReviewNotification(kind, plugin, "")
	}
	utils.OKMsg(c, "插件状态已更新")
}

// ===== 用户管理 =====

// AdminGetUserList 获取用户列表
// GET /api/admin/user/list
func getOptionalBool(c *gin.Context, key string) *bool {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &parsed
}

func AdminGetUserList(c *gin.Context) {
	keywords := strings.TrimSpace(c.Query("keywords"))
	power := getOptionalInt(c, "power")
	isDelete := getOptionalBool(c, "isDelete")
	page, pageSize := getPageParams(c)

	list, total, err := repository.GetUserList(keywords, power, isDelete, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// AdminGetUserDetail 获取用户详情与称号
// GET /api/admin/user/:id
func AdminGetUserDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		utils.BadRequest(c, "参数错误")
		return
	}
	user, err := repository.GetUserByIDWithTitles(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "用户不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKData(c, user)
}

// AdminGetUserPlugins 获取指定用户的全部插件（包括未过审插件）。
func AdminGetUserPlugins(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		utils.BadRequest(c, "参数错误")
		return
	}
	if _, err := repository.GetUserByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "用户不存在")
		return
	} else if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	page, pageSize := getPageParams(c)
	list, total, err := repository.GetAdminPluginsByUser(id, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

type adminUserAction string

const (
	adminUserResetPassword adminUserAction = "reset-password"
	adminUserChangeRole    adminUserAction = "change-role"
	adminUserChangeTitles  adminUserAction = "change-titles"
	adminUserDelete        adminUserAction = "delete"
)

func validateAdminUserPower(power int) error {
	if power != 0 && power != 1 {
		return errors.New("用户权限无效")
	}
	return nil
}

func validateAdminUserAction(action adminUserAction, operatorID, targetID int) error {
	if operatorID == targetID && (action == adminUserChangeRole || action == adminUserDelete) {
		return errors.New("不能对当前登录账号执行此操作")
	}
	return nil
}

// AdminResetPassword 管理员重置用户密码
// PUT /api/admin/user/:id/reset-password
func AdminResetPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if message := validatePassword(req.NewPassword); message != "" {
		utils.BadRequest(c, message)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	if err := repository.UpdateUserFields(id, map[string]interface{}{
		"password": string(hashed),
	}); err != nil {
		utils.ServerError(c, "重置失败")
		return
	}

	utils.OKMsg(c, "密码重置成功")
}

// AdminUpdateRole 修改用户角色
// PUT /api/admin/user/:id/role
func AdminUpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if err := validateAdminUserAction(adminUserChangeRole, c.GetInt("userId"), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var req struct {
		Power *int `json:"power" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if err := validateAdminUserPower(*req.Power); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := repository.UpdateUserFields(id, map[string]interface{}{
		"power": *req.Power,
	}); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKMsg(c, "角色更新成功")
}

// AdminUpdateTitles 管理用户称号
// PUT /api/admin/user/:id/titles
func AdminUpdateTitles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		TitleIds []int `json:"titleIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := repository.SetUserTitles(id, req.TitleIds); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKMsg(c, "称号更新成功")
}

// AdminDeleteUser 删除用户
// DELETE /api/admin/user/:id
func AdminDeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if err := validateAdminUserAction(adminUserDelete, c.GetInt("userId"), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := repository.DeleteUser(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "用户已删除")
}

// ===== 全局设置 =====

// GetSettings 获取全局设置
// GET /api/admin/setting
func GetSettings(c *gin.Context) {
	settings := repository.GetAllSettings()
	if settings["emailPass"] != "" {
		settings["emailPass"] = utils.SMTPPasswordMask
	}
	utils.OKData(c, settings)
}

// UpdateSettings 更新全局设置
// PUT /api/admin/setting
func UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if req["emailPass"] == utils.SMTPPasswordMask {
		delete(req, "emailPass")
	}

	if err := repository.UpdateSettings(req); err != nil {
		utils.ServerError(c, "保存失败")
		return
	}

	utils.OKMsg(c, "设置已保存")
}

// UploadLogo 上传站点 Logo
// POST /api/admin/setting/logo
func UploadLogo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}

	savePath, err := doSaveFile(c, file, "setting")
	if err != nil {
		utils.BadRequest(c, uploadErrorMessage(err))
		return
	}

	repository.UpdateSetting("siteLogo", savePath)
	utils.OKData(c, gin.H{"logo": savePath})
}

// ClearLogo 清除自定义站点 Logo，前端将回退到 config.ts 默认 Logo。
// DELETE /api/admin/setting/logo
func ClearLogo(c *gin.Context) {
	settings := map[string]string{"siteLogo": repository.GetSetting("siteLogo")}
	clearSiteLogoValue(settings)
	if err := repository.UpdateSetting("siteLogo", settings["siteLogo"]); err != nil {
		utils.ServerError(c, "清除 Logo 失败")
		return
	}
	utils.OKMsg(c, "Logo 已恢复默认")
}

func clearSiteLogoValue(settings map[string]string) {
	settings["siteLogo"] = ""
}

// TestEmail 使用数据库中的SMTP配置发送固定内容的测试邮件。
// POST /api/admin/setting/test-email
func TestEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入收件邮箱")
		return
	}
	address, err := mail.ParseAddress(strings.TrimSpace(req.Email))
	if err != nil {
		utils.BadRequest(c, "收件邮箱格式不正确")
		return
	}

	title := "【" + repository.GetSetting("siteTitle") + "】SMTP测试邮件"
	body := "<h2>SMTP配置测试成功</h2><p>如果您收到这封邮件，说明当前SMTP服务器配置可正常连接并发送邮件。</p>"
	if err := utils.SendEmail(address.Address, title, body); err != nil {
		log.Printf("SMTP测试邮件发送失败: %v", err)
		utils.ServerError(c, "测试邮件发送失败，请检查SMTP配置")
		return
	}
	utils.OKMsg(c, "测试邮件已发送")
}

func uploadErrorMessage(err error) string {
	return imageUploadErrorMessage(err)
}

func imageUploadErrorMessage(err error) string {
	if err == nil {
		return "文件保存失败"
	}
	for _, storageError := range []struct {
		target  error
		message string
	}{
		{syscall.EACCES, "上传目录没有写入权限"},
		{syscall.EPERM, "上传目录没有写入权限"},
		{syscall.ENOSPC, "服务器存储空间不足"},
		{syscall.EDQUOT, "服务器存储配额已用尽"},
		{syscall.EROFS, "上传目录所在文件系统为只读"},
		{syscall.ENOENT, "上传目录不存在或配置错误"},
	} {
		if errors.Is(err, storageError.target) {
			return storageError.message
		}
	}
	message := err.Error()
	for _, prefix := range []string{
		"文件大小超过限制",
		"不支持的图片格式",
		"不支持的文件类型",
		"上传文件必须是图片",
		"图片上传已关闭",
	} {
		if strings.HasPrefix(message, prefix) {
			return message
		}
	}
	return "文件保存失败"
}

// CleanImages 查询数据库引用并清理确认未使用的 uploads 文件。
// POST /api/admin/setting/clean-images
func CleanImages(c *gin.Context) {
	result, err := deleteUnusedUploads(
		config.C.Uploads.Path,
		time.Now(),
		uploadCleanupGracePeriod,
		loadUsedUploadReferences,
	)
	if err != nil {
		utils.ServerError(c, "清理未使用文件失败")
		return
	}
	utils.OKData(c, result)
}

func loadUsedUploadReferences() (map[string]struct{}, error) {
	values, err := repository.GetUploadReferenceValues()
	if err != nil {
		return nil, err
	}
	return extractUploadReferences(values), nil
}
