package middleware

import (
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// AdminRequired 管理员权限中间件（需先经过 AuthRequired）
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		powerRaw, exists := c.Get("power")
		if !exists {
			utils.Forbidden(c, "无权限")
			c.Abort()
			return
		}

		power, ok := powerRaw.(int)
		if !ok || power != 1 {
			utils.Forbidden(c, "无权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
