package controller

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"path/filepath"
	"strconv"
	"strings"

	"pluginmarket-server/config"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUserInfo 获取当前登录用户信息
// GET /api/user/info
func GetUserInfo(c *gin.Context) {
	userID := c.GetInt("userId")

	user, err := repository.GetUserByIDWithTitles(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.OKData(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nick":     user.Nick,
		"avatar":   user.Avatar,
		"email":    user.Email,
		"userdesc": user.Userdesc,
		"power":    user.Power,
		"titles":   user.Titles,
		"created":  user.Created,
		"updated":  user.Updated,
	})
}

// GetUserPublic 获取指定用户公开信息
// GET /api/user/:id
func GetUserPublic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := repository.GetUserByIDWithTitles(id)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.OKData(c, user.Public())
}

// UpdateUserInfo 更新个人信息
// PUT /api/user/info
func UpdateUserInfo(c *gin.Context) {
	userID := c.GetInt("userId")

	var req struct {
		Nick     string `json:"nick"`
		Userdesc string `json:"userdesc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	fields := map[string]interface{}{}
	if req.Nick != "" {
		fields["nick"] = req.Nick
	}
	if req.Userdesc != "" {
		fields["userdesc"] = req.Userdesc
	}

	if len(fields) > 0 {
		if err := repository.UpdateUserFields(userID, fields); err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	utils.OKMsg(c, "更新成功")
}

// ChangePassword 修改密码
// PUT /api/user/password
func ChangePassword(c *gin.Context) {
	userID := c.GetInt("userId")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
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

	user, err := repository.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.BadRequest(c, "旧密码错误")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	if err := repository.UpdateUserFields(userID, map[string]interface{}{
		"password": string(hashed),
	}); err != nil {
		utils.ServerError(c, "密码修改失败")
		return
	}

	utils.OKMsg(c, "密码修改成功")
}

// UploadAvatar 上传头像
// POST /api/user/avatar
func UploadAvatar(c *gin.Context) {
	userID := c.GetInt("userId")

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}

	savePath, err := doSaveFile(c, file, "avatar")
	if err != nil {
		utils.BadRequest(c, imageUploadErrorMessage(err))
		return
	}

	if err := repository.UpdateUserFields(userID, map[string]interface{}{
		"avatar": savePath,
	}); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKData(c, gin.H{"avatar": savePath})
}

// SendEmailVerify 发送邮箱验证链接
// POST /api/user/email/send-verify
func SendEmailVerify(c *gin.Context) {
	userID := c.GetInt("userId")

	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	email := strings.TrimSpace(req.Email)
	address, err := mail.ParseAddress(email)
	if err != nil {
		utils.BadRequest(c, "邮箱格式不正确")
		return
	}
	email = strings.ToLower(address.Address)
	if existing, err := repository.GetUserByEmail(email); err == nil && existing.ID != userID {
		utils.BadRequest(c, "该邮箱已被其他用户绑定")
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerError(c, "邮箱检查失败")
		return
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	token, err := utils.GenerateShortToken(userID, email, "email_bind", user.Updated.UnixMilli())
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	verifyURL := webURL() + "/user/verify-email?token=" + token
	if err := utils.SendEmailVerifyEmail(email, user.Nick, verifyURL); err != nil {
		utils.ServerError(c, "邮件发送失败："+err.Error())
		return
	}

	utils.OKMsg(c, "验证邮件已发送，请查收邮箱")
}

// EmailVerify 邮箱验证回调
// GET /api/user/email/verify
func EmailVerify(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		utils.BadRequest(c, "缺少token参数")
		return
	}

	claims, err := utils.ParseShortToken(tokenStr)
	if err != nil {
		utils.BadRequest(c, "链接无效或已过期")
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "email_bind" {
		utils.BadRequest(c, "链接类型错误")
		return
	}

	userIDf, userOK := claims["userId"].(float64)
	email, emailOK := claims["email"].(string)
	address, emailErr := mail.ParseAddress(strings.TrimSpace(email))
	if !userOK || userIDf <= 0 || userIDf != float64(int(userIDf)) || !emailOK || emailErr != nil {
		utils.BadRequest(c, "链接无效")
		return
	}
	user, err := repository.GetUserByID(int(userIDf))
	if err != nil {
		utils.BadRequest(c, "链接无效")
		return
	}
	version, versionOK := claimInt64(claims["version"])
	if !versionOK || version != user.Updated.UnixMilli() {
		utils.BadRequest(c, "链接已失效")
		return
	}

	if err := repository.BindUserEmail(int(userIDf), strings.ToLower(address.Address)); errors.Is(err, repository.ErrEmailAlreadyBound) {
		utils.BadRequest(c, "该邮箱已被其他用户绑定")
		return
	} else if err != nil {
		utils.ServerError(c, "邮箱绑定失败")
		return
	}

	utils.OKMsg(c, "邮箱绑定成功")
}

// doSaveFile 保存上传文件，返回相对路径
func doSaveFile(c *gin.Context, file *multipart.FileHeader, subDirKey string) (string, error) {
	if !settingEnabled(repository.GetSetting("allowUpload"), true) {
		return "", fmt.Errorf("图片上传已关闭")
	}

	// 获取具体的子目录
	var subDir string
	switch subDirKey {
	case "avatar":
		subDir = config.C.Uploads.Avatars
	case "image", "plugin/icon":
		subDir = config.C.Uploads.Images
	case "frame":
		subDir = config.C.Uploads.Frames
	case "title":
		subDir = config.C.Uploads.Titles
	default:
		subDir = subDirKey
	}

	// 获取文件大小限制 (优先从数据库获取)
	maxSizeMB := config.C.Uploads.MaxSize
	if dbSizeStr := repository.GetSetting("maxFileSize"); dbSizeStr != "" {
		if s, err := strconv.ParseFloat(dbSizeStr, 64); err == nil && s > 0 {
			maxSizeMB = s
		}
	}
	if maxSizeMB <= 0 {
		maxSizeMB = 10 // 最终保底
	}

	// 校验真实 MIME 类型
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:n])
	if err := validateImageUpload(file.Size, file.Filename, contentType, maxSizeMB, repository.GetSetting("allowedExtensions")); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	return saveDeduplicatedUpload(config.C.Uploads.Path, subDir, ext, src, file.Size)
}

func validateImageUpload(size int64, filename, contentType string, maxSizeMB float64, allowedExtensions string) error {
	if float64(size) > maxSizeMB*1024*1024 {
		return fmt.Errorf("文件大小超过限制，最大允许 %gMB", maxSizeMB)
	}
	if !imageExtensionAllowed(filename, allowedExtensions) {
		return fmt.Errorf("不支持的图片格式")
	}
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("上传文件必须是图片")
	}
	return nil
}

func settingEnabled(value string, fallback bool) bool {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	return fallback
}

func imageExtensionAllowed(filename, configured string) bool {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if ext == "" {
		return false
	}
	if configured == "" {
		configured = "jpg,png,gif,webp"
	}
	for _, allowed := range strings.Split(configured, ",") {
		if strings.TrimSpace(strings.ToLower(allowed)) == ext {
			return true
		}
	}
	return false
}
