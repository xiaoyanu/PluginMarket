package middleware

import (
	"errors"
	"strings"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

var errInactiveUser = errors.New("user is deleted")

type authenticatedIdentity struct {
	UserID int
	Power  int
}

type userLookup func(int) (*model.User, error)

// resolveDatabaseIdentity 只使用 JWT 中的用户 ID 定位账号，权限始终以数据库当前值为准。
func resolveDatabaseIdentity(userID int, lookup userLookup) (authenticatedIdentity, error) {
	user, err := lookup(userID)
	if err != nil {
		return authenticatedIdentity{}, err
	}
	if user == nil || user.IsDelete {
		return authenticatedIdentity{}, errInactiveUser
	}
	return authenticatedIdentity{UserID: user.ID, Power: user.Power}, nil
}

func setAuthenticatedIdentity(c *gin.Context, identity authenticatedIdentity) {
	c.Set("userId", identity.UserID)
	c.Set("power", identity.Power)
}

// AuthRequired JWT 认证中间件（必须登录）。JWT 只验证身份，实时权限来自数据库。
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

		identity, err := resolveDatabaseIdentity(claims.UserID, repository.GetUserByID)
		if err != nil {
			utils.Unauthorized(c, "用户状态无效，请重新登录")
			c.Abort()
			return
		}

		setAuthenticatedIdentity(c, identity)
		c.Next()
	}
}

// OptionalAuth 可选认证。有效 JWT 仍需通过数据库用户状态校验，否则按匿名访问。
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			if claims, err := utils.ParseToken(tokenStr); err == nil {
				if identity, lookupErr := resolveDatabaseIdentity(claims.UserID, repository.GetUserByID); lookupErr == nil {
					setAuthenticatedIdentity(c, identity)
				}
			}
		}
		c.Next()
	}
}
