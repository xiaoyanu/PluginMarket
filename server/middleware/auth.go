package middleware

import (
	"strings"

	"pluginmarket-server/database"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT 认证中间件（必须登录）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}
		if database.DB != nil {
			user, err := repository.GetUserByID(claims.UserID)
			if err != nil || user.IsDelete {
				utils.Unauthorized(c, "此用户已被删除，无法登录。")
				c.Abort()
				return
			}
		}
		c.Set("userId", claims.UserID)
		c.Set("power", claims.Power)
		c.Next()
	}
}

// OptionalAuth 可选认证（不强制登录，有 token 则解析）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := utils.ParseToken(tokenStr)
			if err == nil {
				if database.DB == nil {
					c.Set("userId", claims.UserID)
					c.Set("power", claims.Power)
				} else if user, userErr := repository.GetUserByID(claims.UserID); userErr == nil && !user.IsDelete {
					c.Set("userId", claims.UserID)
					c.Set("power", claims.Power)
				}
			}
		}
		c.Next()
	}
}
