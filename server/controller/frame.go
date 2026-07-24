package controller

import (
	"strconv"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// GetFrameList 获取所有框架
// GET /api/frame/list
func GetFrameList(c *gin.Context) {
	page, pageSize := getPageParams(c)
	keywords := c.Query("keywords")
	frames, total, err := repository.GetFrameList(keywords, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, frames, total)
}

// CreateFrame 添加框架（管理员）
// POST /api/frame
func CreateFrame(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequest(c, "请输入框架名称")
		return
	}

	file, err := c.FormFile("icon")
	if err != nil {
		utils.BadRequest(c, "请上传框架图标")
		return
	}

	iconPath, err := doSaveFile(c, file, "frame")
	if err != nil {
		utils.ServerError(c, "图标上传失败")
		return
	}

	frame := &model.Frame{Name: name, Icon: iconPath}
	if err := repository.CreateFrame(frame); err != nil {
		utils.ServerError(c, "添加失败")
		return
	}

	utils.OKData(c, gin.H{"id": frame.ID})
}

// UpdateFrame 编辑框架（管理员）
// PUT /api/frame/:id
func UpdateFrame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	frame, err := repository.GetFrameByID(id)
	if err != nil {
		utils.NotFound(c, "框架不存在")
		return
	}

	if name := c.PostForm("name"); name != "" {
		frame.Name = name
	}

	file, err := c.FormFile("icon")
	if err == nil {
		iconPath, err := doSaveFile(c, file, "frame")
		if err == nil {
			frame.Icon = iconPath
		}
	}

	if err := repository.UpdateFrame(frame); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKMsg(c, "更新成功")
}

// DeleteFrame 删除框架（管理员）
// DELETE /api/frame/:id
func DeleteFrame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DeleteFrame(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "删除成功")
}
