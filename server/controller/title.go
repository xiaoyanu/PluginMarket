package controller

import (
	"errors"
	"strconv"
	"strings"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

func validateTitleFields(name, description string) (string, string, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return "", "", errors.New("请输入称号名称")
	}
	if description == "" {
		return "", "", errors.New("请输入称号描述")
	}
	if len([]rune(description)) > 500 {
		return "", "", errors.New("称号描述不能超过500个字符")
	}
	return name, description, nil
}

// GetTitleList 获取所有称号
// GET /api/title/list
func GetTitleList(c *gin.Context) {
	page, pageSize := getPageParams(c)
	keywords := c.Query("keywords")
	titles, total, err := repository.GetTitleList(keywords, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, titles, total)
}

// CreateTitle 添加称号（管理员）
// POST /api/title
func CreateTitle(c *gin.Context) {
	name, description, err := validateTitleFields(c.PostForm("name"), c.PostForm("description"))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	file, err := c.FormFile("icon")
	if err != nil {
		utils.BadRequest(c, "请上传称号图标")
		return
	}

	iconPath, err := doSaveFile(c, file, "title")
	if err != nil {
		utils.ServerError(c, "图标上传失败")
		return
	}

	title := &model.Title{Name: name, Description: description, Icon: iconPath}
	if err := repository.CreateTitle(title); err != nil {
		utils.ServerError(c, "添加失败")
		return
	}

	utils.OKData(c, gin.H{"id": title.ID})
}

// UpdateTitle 编辑称号（管理员）
// PUT /api/title/:id
func UpdateTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	title, err := repository.GetTitleByID(id)
	if err != nil {
		utils.NotFound(c, "称号不存在")
		return
	}

	name, description, fieldErr := validateTitleFields(c.PostForm("name"), c.PostForm("description"))
	if fieldErr != nil {
		utils.BadRequest(c, fieldErr.Error())
		return
	}
	title.Name = name
	title.Description = description

	file, err := c.FormFile("icon")
	if err == nil {
		iconPath, err := doSaveFile(c, file, "title")
		if err == nil {
			title.Icon = iconPath
		}
	}

	if err := repository.UpdateTitle(title); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.OKMsg(c, "更新成功")
}

// DeleteTitle 删除称号（管理员）
// DELETE /api/title/:id
func DeleteTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := repository.DeleteTitle(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "删除成功")
}
