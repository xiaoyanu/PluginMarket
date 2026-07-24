package controller

import (
	"strconv"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// GetTagList 获取所有标签
// GET /api/tag/list
func GetTagList(c *gin.Context) {
	page, pageSize := getPageParams(c)
	keywords := c.Query("keywords")
	tags, total, err := repository.GetTagList(keywords, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, tags, total)
}

// CreateTag 添加标签（管理员）
// POST /api/tag
func CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	tag := &model.Tag{Name: req.Name}
	if err := repository.CreateTag(tag); err != nil {
		utils.ServerError(c, "添加失败")
		return
	}

	utils.OKData(c, gin.H{"id": tag.ID})
}

// UpdateTag 编辑标签（管理员）
// PUT /api/tag/:id
func UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	tag, err := repository.GetTagByID(id)
	if err != nil {
		utils.NotFound(c, "标签不存在")
		return
	}

	tag.Name = req.Name
	if err := repository.UpdateTag(tag); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKMsg(c, "更新成功")
}

// DeleteTag 删除标签（管理员）
// DELETE /api/tag/:id
func DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DeleteTag(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "删除成功")
}
