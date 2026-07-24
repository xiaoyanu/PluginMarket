package controller

import (
	"encoding/json"
	"image/color"
	"strings"

	"pluginmarket-server/config"
	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

var captchaStore = base64Captcha.DefaultMemStore

func deletedUserLoginMessage(user *model.User) string {
	if user != nil && user.IsDelete {
		return "此用户已被删除，无法登录。"
	}
	return ""
}

// GetCaptcha 获取验证码
// GET /api/auth/captcha
func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverString(
		80,
		240,
		2,
		2,
		4,
		"1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		&color.RGBA{R: 255, G: 255, B: 255, A: 255},
		nil,
		nil,
	)
	cp := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		utils.ServerError(c, "验证码生成失败")
		return
	}
	utils.OKData(c, gin.H{
		"captchaId": id,
		"image":     b64s,
	})
}

// Login 用户登录
// POST /api/auth/login
func Login(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
		CaptchaID string `json:"captchaId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		utils.BadRequest(c, "验证码错误")
		return
	}

	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}
	if message := deletedUserLoginMessage(user); message != "" {
		utils.BadRequest(c, message)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Power)
	if err != nil {
		utils.ServerError(c, "Token生成失败")
		return
	}

	utils.OKData(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"power": user.Power,
		},
	})
}

// Register 用户注册
// POST /api/auth/register
func Register(c *gin.Context) {
	if !settingEnabled(repository.GetSetting("allowRegister"), true) {
		utils.Forbidden(c, "网站暂未开放注册")
		return
	}

	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
		CaptchaID string `json:"captchaId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if message := validateRegistrationCredentials(req.Username, req.Password); message != "" {
		utils.BadRequest(c, message)
		return
	}

	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		utils.BadRequest(c, "验证码错误")
		return
	}

	if _, err := repository.GetUserByUsername(req.Username); err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashed),
		Nick:     req.Username,
	}
	if err := repository.CreateUser(user); err != nil {
		utils.ServerError(c, "注册失败")
		return
	}

	utils.OKMsg(c, "注册成功")
}

// ForgotPassword 忘记密码（发送重置邮件）
// POST /api/auth/forgot-password
func ForgotPassword(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
		CaptchaID string `json:"captchaId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		utils.BadRequest(c, "验证码错误")
		return
	}

	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		utils.BadRequest(c, "用户不存在")
		return
	}

	if user.Email == "" {
		utils.BadRequest(c, "该用户未绑定邮箱")
		return
	}

	token, err := utils.GenerateShortToken(user.ID, user.Email, "reset_password", user.Updated.UnixMilli())
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	resetURL := webURL() + "/user/reset?token=" + token
	if err := utils.SendResetPasswordEmail(user.Email, user.Nick, resetURL); err != nil {
		utils.ServerError(c, "邮件发送失败："+err.Error())
		return
	}

	maskedEmail := maskEmail(user.Email)
	utils.OKMsg(c, "重置邮件已发送，请查收邮箱:"+maskedEmail)
}

// ResetPassword 重置密码（通过邮件链接）
// POST /api/auth/reset-password
func ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	if message := validatePassword(req.Password); message != "" {
		utils.BadRequest(c, message)
		return
	}

	claims, err := utils.ParseShortToken(req.Token)
	if err != nil {
		utils.BadRequest(c, "链接无效或已过期")
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "reset_password" {
		utils.BadRequest(c, "链接类型错误")
		return
	}

	userIDf, ok := claims["userId"].(float64)
	if !ok || userIDf <= 0 || userIDf != float64(int(userIDf)) {
		utils.BadRequest(c, "链接无效")
		return
	}
	userID := int(userIDf)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		utils.BadRequest(c, "链接无效")
		return
	}
	version, versionOK := claimInt64(claims["version"])
	if !versionOK || version != user.Updated.UnixMilli() {
		utils.BadRequest(c, "链接已失效")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "服务器错误")
		return
	}

	if err := repository.UpdateUserFields(userID, map[string]interface{}{
		"password": string(hashed),
	}); err != nil {
		utils.ServerError(c, "密码重置失败")
		return
	}

	utils.OKMsg(c, "密码重置成功")
}

func claimInt64(value interface{}) (int64, bool) {
	switch number := value.(type) {
	case float64:
		return int64(number), number == float64(int64(number))
	case json.Number:
		parsed, err := number.Int64()
		return parsed, err == nil
	default:
		return 0, false
	}
}

func webURL() string {
	if value := strings.TrimRight(config.C.Server.WebURL, "/"); value != "" {
		return value
	}
	return "http://localhost:3000"
}

func maskEmail(email string) string {
	at := -1
	for i, ch := range email {
		if ch == '@' {
			at = i
			break
		}
	}
	if at <= 1 {
		return email
	}
	return string(email[0]) + "***" + email[at:]
}
