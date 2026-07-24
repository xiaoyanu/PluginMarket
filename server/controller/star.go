package controller

import (
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// ToggleStar 收藏/取消收藏
// POST /api/star/toggle
func ToggleStar(c *gin.Context) {
	userID := c.GetInt("userId")

	var req struct {
		PluginID int `json:"pluginId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	isStarred, err := repository.ToggleStar(userID, req.PluginID)
	if err != nil {
		utils.ServerError(c, "操作失败")
		return
	}

	msg := "收藏成功"
	if !isStarred {
		msg = "已取消收藏"
	}

	utils.OKDataMsg(c, gin.H{"is_starred": isStarred}, msg)
}

// GetStarList 我的收藏列表
// GET /api/star/list
func GetStarList(c *gin.Context) {
	userID := c.GetInt("userId")
	page, pageSize := getPageParams(c)

	list, total, err := repository.GetStarList(userID, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}
