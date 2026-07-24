package controller

import (
	"strconv"
	"strings"
	"time"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// GetCommentList 获取插件评论列表
// GET /api/comment/list
func GetCommentList(c *gin.Context) {
	pluginIDStr := c.Query("pluginId")
	pluginID, err := strconv.Atoi(pluginIDStr)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	page, pageSize := getPageParams(c)
	var target *repository.CommentTarget
	if targetIDText := strings.TrimSpace(c.Query("targetCommentId")); targetIDText != "" {
		targetID, parseErr := strconv.Atoi(targetIDText)
		if parseErr != nil || targetID <= 0 {
			utils.BadRequest(c, "目标评论ID无效")
			return
		}
		target, err = repository.GetCommentTarget(pluginID, targetID, pageSize)
		if err != nil {
			utils.NotFound(c, "该评论已被删除、不存在或不属于当前插件")
			return
		}
		page = target.Page
	}

	list, total, err := repository.GetCommentsByPlugin(pluginID, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKData(c, gin.H{"list": list, "total": total, "target": target})
}

// CreateComment 发表评论
// POST /api/comment
func CreateComment(c *gin.Context) {
	if !settingEnabled(repository.GetSetting("allowComment"), true) {
		utils.Forbidden(c, "网站暂未开放评论")
		return
	}

	userID := c.GetInt("userId")

	var req struct {
		PluginID int    `json:"pluginId" binding:"required"`
		Content  string `json:"content" binding:"required"`
		ParentID int    `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	// 计算 root：如果有 parentId，找到该评论的 root
	root := 0
	replyUser := 0
	if req.ParentID > 0 {
		parent, err := repository.GetCommentByID(req.ParentID)
		if err != nil {
			utils.BadRequest(c, "回复的评论不存在")
			return
		}
		if parent.Plugin != req.PluginID {
			utils.BadRequest(c, "回复评论与插件不匹配")
			return
		}
		if parent.Root == 0 {
			root = parent.ID
		} else {
			root = parent.Root
		}
		replyUser = parent.User
	}

	comment := &model.Comment{
		Plugin:    req.PluginID,
		User:      userID,
		Content:   req.Content,
		Root:      root,
		Parent:    req.ParentID,
		ReplyUser: replyUser,
	}

	if err := repository.CreateComment(comment); err != nil {
		utils.ServerError(c, "评论失败")
		return
	}

	sendCommentNotification(req.PluginID, comment.ID, userID, replyUser, req.Content, req.ParentID, comment.Created)
	utils.OKData(c, gin.H{"id": comment.ID})
}

// DeleteComment 删除评论
// DELETE /api/comment/:id
func DeleteComment(c *gin.Context) {
	userID := c.GetInt("userId")
	power := c.GetInt("power")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	comment, err := repository.GetCommentByID(id)
	if err != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	// 权限检查：评论作者本人、插件作者本人或管理员
	if comment.User != userID {
		plugin, _ := repository.GetPluginByID(comment.Plugin)
		if plugin == nil || plugin.User != userID {
			if power != 1 {
				utils.Forbidden(c, "无权限")
				return
			}
		}
	}

	if err := repository.DeleteComment(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "删除成功")
}

func buildCommentLinks(baseURL string, pluginID, commentID int) (string, string) {
	relative := "/plugin/" + strconv.Itoa(pluginID) + "?commentId=" + strconv.Itoa(commentID) + "#comment-" + strconv.Itoa(commentID)
	return strings.TrimRight(baseURL, "/") + relative, relative
}

func sendCommentNotification(pluginID, commentID, senderID, replyUserID int, content string, parentID int, created time.Time) {
	plugin, err := repository.GetPluginByID(pluginID)
	if err != nil {
		return
	}
	sender, err := repository.GetUserByID(senderID)
	if err != nil {
		return
	}
	link, actionURL := buildCommentLinks(webURL(), pluginID, commentID)
	if parentID > 0 {
		if replyUserID == senderID {
			return
		}
		createWorkflowNotification(workflowReplyComment, replyUserID, senderID, int64(pluginID), plugin.Name, sender.Nick, content, actionURL)
		recipient, err := repository.GetUserByID(replyUserID)
		if err == nil && recipient.Email != "" {
			parent, parentErr := repository.GetCommentByID(parentID)
			if parentErr == nil {
				_ = utils.SendReplyCommentEmail(recipient.Email, plugin.Name, sender.Nick, content, recipient.Nick, parent.Content, link, created)
			}
		}
		return
	}
	if plugin.User == senderID {
		return
	}
	createWorkflowNotification(workflowNewComment, plugin.User, senderID, int64(pluginID), plugin.Name, sender.Nick, content, actionURL)
	recipient, err := repository.GetUserByID(plugin.User)
	if err == nil && recipient.Email != "" {
		_ = utils.SendNewCommentEmail(recipient.Email, plugin.Name, sender.Nick, content, link, created)
	}
}
