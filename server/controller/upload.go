package controller

import (
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// UploadImage 通用图片上传
// POST /api/upload/image
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}

	savePath, err := doSaveFile(c, file, "image")
	if err != nil {
		utils.BadRequest(c, imageUploadErrorMessage(err))
		return
	}

	utils.OKData(c, gin.H{"url": savePath})
}
