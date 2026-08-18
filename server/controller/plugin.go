package controller

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"pluginmarket-server/model"
	"pluginmarket-server/repository"
	"pluginmarket-server/utils"

	"github.com/gin-gonic/gin"
)

// GetPluginListHot 热门插件列表
// GET /api/plugin/list/hot
func GetPluginListHot(c *gin.Context) {
	page, pageSize := getPageParams(c)
	frameID := getOptionalInt(c, "frameId")
	tagID := getOptionalInt(c, "tagId")
	pluginType := getOptionalInt(c, "type")

	list, total, err := repository.GetPluginListHot(frameID, tagID, pluginType, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// GetPluginListLatest 最新插件列表
// GET /api/plugin/list/latest
func GetPluginListLatest(c *gin.Context) {
	page, pageSize := getPageParams(c)
	frameID := getOptionalInt(c, "frameId")
	tagID := getOptionalInt(c, "tagId")
	pluginType := getOptionalInt(c, "type")

	list, total, err := repository.GetPluginListLatest(frameID, tagID, pluginType, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// GetPluginDetail 插件详情
// GET /api/plugin/:id
func GetPluginDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if err != nil {
		utils.NotFound(c, "插件不存在")
		return
	}

	if !canAccessPluginDetail(plugin.Status, plugin.User, c.GetInt("userId"), c.GetInt("power")) {
		utils.Forbidden(c, "该插件尚未通过审核，暂时无法访问")
		return
	}

	detail, err := repository.GetPluginDetail(id)
	if err != nil {
		utils.ServerError(c, "查询插件详情失败")
		return
	}

	// 增加浏览量
	repository.IncrementPluginViews(id)

	// 检查当前用户是否已收藏（如果已登录）
	userID, _ := c.Get("userId")
	if uid, ok := userID.(int); ok && uid > 0 {
		detail.IsStarred = repository.IsPluginStarred(uid, id)
	}

	utils.OKData(c, detail)
}

func canAccessPluginDetail(status, ownerID, currentUserID, power int) bool {
	return status == 0 || currentUserID == ownerID || power == 1
}

// SearchPlugins 搜索插件
// GET /api/plugin/search
func SearchPlugins(c *gin.Context) {
	keywords := c.Query("keywords")
	if keywords == "" {
		utils.BadRequest(c, "请输入搜索关键词")
		return
	}
	page, pageSize := getPageParams(c)
	frameID := getOptionalInt(c, "frameId")
	tagID := getOptionalInt(c, "tagId")
	pluginType := getOptionalInt(c, "type")

	list, total, err := repository.SearchPlugins(keywords, frameID, tagID, pluginType, page, pageSize)
	if err != nil {
		utils.ServerError(c, "搜索失败")
		return
	}
	utils.OKPage(c, list, total)
}

// GetPluginsByTag 按标签获取插件
// GET /api/plugin/list/by-tag
func GetPluginsByTag(c *gin.Context) {
	tagIDStr := c.Query("tagId")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	page, pageSize := getPageParams(c)

	list, total, err := repository.GetPluginsByTag(tagID, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// GetPluginsByUser 获取用户发布的插件
// GET /api/plugin/list/by-user
func GetPluginsByUser(c *gin.Context) {
	params, err := getPluginsByUserParams(c)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	list, total, err := repository.GetPluginsByUser(
		params.UserID,
		params.FrameID,
		params.TagID,
		params.PluginType,
		params.Page,
		params.PageSize,
	)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

type pluginsByUserParams struct {
	UserID     int
	FrameID    *int
	TagID      *int
	PluginType *int
	Page       int
	PageSize   int
}

func getPluginsByUserParams(c *gin.Context) (pluginsByUserParams, error) {
	userID, err := strconv.Atoi(c.Query("userId"))
	if err != nil || userID < 1 {
		return pluginsByUserParams{}, fmt.Errorf("invalid user id")
	}
	page, pageSize := getPageParams(c)
	return pluginsByUserParams{
		UserID:     userID,
		FrameID:    getOptionalInt(c, "frameId"),
		TagID:      getOptionalInt(c, "tagId"),
		PluginType: getOptionalInt(c, "type"),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetMyPlugins 我的插件列表
// GET /api/plugin/my
func GetMyPlugins(c *gin.Context) {
	userID := c.GetInt("userId")
	page, pageSize := getPageParams(c)
	keywords := c.Query("keywords")
	pluginType := getOptionalInt(c, "type")
	status := getOptionalInt(c, "status")

	list, total, err := repository.GetMyPlugins(userID, keywords, pluginType, status, page, pageSize)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.OKPage(c, list, total)
}

// PublishPlugin 发布插件
// POST /api/plugin/publish
func PublishPlugin(c *gin.Context) {
	userID := c.GetInt("userId")

	name := c.PostForm("name")
	descText := c.PostForm("desc_text")
	typeStr := c.PostForm("type")
	content := c.PostForm("content")
	url := strings.TrimSpace(c.PostForm("url"))
	urlCode := c.PostForm("url_code")
	frameworkIdsStr := getFormArray(c, "frameworkIds")
	tagIdsStr := getFormArray(c, "tagIds")

	if name == "" || descText == "" || typeStr == "" || content == "" || url == "" {
		utils.BadRequest(c, "缺少必填字段")
		return
	}
	if !isValidPluginDownloadURL(url) {
		utils.BadRequest(c, "下载链接必须是有效的 HTTP 或 HTTPS 地址")
		return
	}

	pluginType, _ := strconv.Atoi(typeStr)
	frameworkIDs, err := parseRelationIDs(frameworkIdsStr)
	if err != nil || len(frameworkIDs) == 0 {
		utils.BadRequest(c, "适用框架参数无效")
		return
	}
	tagIDs, err := parseRelationIDs(tagIdsStr)
	if err != nil {
		utils.BadRequest(c, "标签参数无效")
		return
	}
	if err := repository.ValidatePluginRelations(frameworkIDs, tagIDs); err != nil {
		if repository.IsInvalidPluginRelation(err) {
			utils.BadRequest(c, err.Error())
		} else {
			utils.ServerError(c, "关联数据校验失败")
		}
		return
	}

	// 处理图标上传
	iconPath := ""
	file, err := c.FormFile("icon")
	if err == nil {
		iconPath, err = doSaveFile(c, file, "plugin/icon")
		if err != nil {
			utils.BadRequest(c, uploadErrorMessage(err))
			return
		}
	}

	status := pluginPublishStatus(
		c.GetInt("power"),
		settingEnabled(repository.GetSetting("skipAudit"), false),
	)

	plugin := &model.Plugin{
		Name:     name,
		DescText: descText,
		Type:     pluginType,
		User:     userID,
		URL:      url,
		URLCode:  urlCode,
		Icon:     iconPath,
		Status:   status,
		Content:  content,
	}

	if err := repository.CreatePluginWithRelations(plugin, frameworkIDs, tagIDs); err != nil {
		if repository.IsInvalidPluginRelation(err) {
			utils.BadRequest(c, err.Error())
		} else {
			utils.ServerError(c, "发布失败")
		}
		return
	}
	if plugin.Status == 2 {
		sendPendingPluginReviewNotification(plugin)
	}

	utils.OKData(c, gin.H{"id": plugin.ID})
}

func isValidPluginDownloadURL(value string) bool {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(value))
	if err != nil || parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func pluginPublishStatus(power int, skipAudit bool) int {
	if power == 1 || skipAudit {
		return 0
	}
	return 2
}

func setStringFieldFromForm(fields map[string]interface{}, c *gin.Context, key string) {
	if value, ok := c.GetPostForm(key); ok {
		fields[key] = value
	}
}

// EditPlugin 编辑插件
// PUT /api/plugin/:id
func EditPlugin(c *gin.Context) {
	userID := c.GetInt("userId")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if err != nil {
		utils.NotFound(c, "插件不存在")
		return
	}

	if plugin.User != userID {
		utils.Forbidden(c, "无权限")
		return
	}

	fields := map[string]interface{}{}
	var frameworkIDs, tagIDs []int
	var updateFrameworks, updateTags bool
	if c.PostForm("frameworkIds") != "" || len(c.PostFormArray("frameworkIds")) > 0 {
		updateFrameworks = true
		frameworkIDs, err = parseRelationIDs(getFormArray(c, "frameworkIds"))
		if err != nil || len(frameworkIDs) == 0 {
			utils.BadRequest(c, "适用框架参数无效")
			return
		}
	}
	if c.PostForm("tagIds") != "" || len(c.PostFormArray("tagIds")) > 0 {
		updateTags = true
		tagIDs, err = parseRelationIDs(getFormArray(c, "tagIds"))
		if err != nil {
			utils.BadRequest(c, "标签参数无效")
			return
		}
	}
	if err := repository.ValidatePluginRelations(optionalRelationIDs(updateFrameworks, frameworkIDs), optionalRelationIDs(updateTags, tagIDs)); err != nil {
		if repository.IsInvalidPluginRelation(err) {
			utils.BadRequest(c, err.Error())
		} else {
			utils.ServerError(c, "关联数据校验失败")
		}
		return
	}
	if name := c.PostForm("name"); name != "" {
		fields["name"] = name
	}
	if desc := c.PostForm("desc_text"); desc != "" {
		fields["desc_text"] = desc
	}
	if typeStr := c.PostForm("type"); typeStr != "" {
		fields["type"], _ = strconv.Atoi(typeStr)
	}
	if content := c.PostForm("content"); content != "" {
		fields["content"] = content
	}
	if downloadURL := c.PostForm("url"); downloadURL != "" {
		downloadURL = strings.TrimSpace(downloadURL)
		if !isValidPluginDownloadURL(downloadURL) {
			utils.BadRequest(c, "下载链接必须是有效的 HTTP 或 HTTPS 地址")
			return
		}
		fields["url"] = downloadURL
	}
	setStringFieldFromForm(fields, c, "url_code")

	// 图标
	file, err := c.FormFile("icon")
	if err == nil {
		iconPath, err := doSaveFile(c, file, "plugin/icon")
		if err != nil {
			utils.BadRequest(c, uploadErrorMessage(err))
			return
		}
		fields["icon"] = iconPath
	}

	if len(fields) > 0 || updateFrameworks || updateTags {
		if c.GetInt("power") == 1 {
			fields["status"] = 0 // 管理员编辑后直接通过
		} else {
			fields["status"] = 2 // 普通用户编辑后重新审核
		}
		if err := repository.UpdatePluginWithRelations(id, fields, updateFrameworks, frameworkIDs, updateTags, tagIDs); err != nil {
			if repository.IsInvalidPluginRelation(err) {
				utils.BadRequest(c, err.Error())
			} else {
				utils.ServerError(c, "更新失败")
			}
			return
		}
		if fields["status"] == 2 {
			if updatedPlugin, err := repository.GetPluginByID(id); err == nil {
				sendPendingPluginReviewNotification(updatedPlugin)
			}
		}
	}

	utils.OKMsg(c, "更新成功")
}

// DeletePlugin 删除插件
// DELETE /api/plugin/:id
func DeletePlugin(c *gin.Context) {
	userID := c.GetInt("userId")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if err != nil {
		utils.NotFound(c, "插件不存在")
		return
	}

	if plugin.User != userID {
		utils.Forbidden(c, "无权限")
		return
	}

	if err := repository.DeletePlugin(id); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.OKMsg(c, "删除成功")
}

// DownloadPlugin 下载插件
// GET /api/plugin/:id/download
func DownloadPlugin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	plugin, err := repository.GetPluginByID(id)
	if err != nil {
		utils.NotFound(c, "插件不存在")
		return
	}

	if plugin.URL == "" {
		utils.NotFound(c, "该插件未配置下载地址")
		return
	}

	if err := repository.IncrementPluginDownloads(id); err != nil {
		utils.ServerError(c, "下载次数记录失败")
		return
	}

	utils.OKData(c, gin.H{
		"url":      plugin.URL,
		"url_code": plugin.URLCode,
	})
}

// 工具函数
func getPageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

func getOptionalInt(c *gin.Context, key string) *int {
	val := c.Query(key)
	if val == "" {
		return nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}
	return &n
}

func parseIds(strs []string) []int {
	ids := make([]int, 0, len(strs))
	for _, s := range strs {
		if n, err := strconv.Atoi(s); err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}

func parseRelationIDs(values []string) ([]int, error) {
	ids := make([]int, 0, len(values))
	seen := make(map[int]struct{}, len(values))
	for _, value := range values {
		id, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("关联ID无效")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func optionalRelationIDs(update bool, ids []int) []int {
	if !update {
		return nil
	}
	return ids
}

// getFormArray 从表单中获取数组字段，支持 JSON 数组格式和重复字段格式
func getFormArray(c *gin.Context, key string) []string {
	// 先尝试重复字段格式，若只有一个值且它本身是 JSON 数组，再继续按 JSON 解析
	vals := c.PostFormArray(key)
	if len(vals) > 1 {
		return vals
	}

	raw := ""
	if len(vals) == 1 {
		raw = vals[0]
	} else {
		raw = c.PostForm(key)
	}
	if raw == "" {
		return nil
	}
	// 尝试解析为 int 数组
	var intArr []int
	if err := json.Unmarshal([]byte(raw), &intArr); err == nil {
		strs := make([]string, len(intArr))
		for i, v := range intArr {
			strs[i] = strconv.Itoa(v)
		}
		return strs
	}
	// 尝试解析为 string 数组
	var strArr []string
	if err := json.Unmarshal([]byte(raw), &strArr); err == nil {
		return strArr
	}
	// 逗号分隔格式 "1,2,3"
	if strings.Contains(raw, ",") {
		return strings.Split(raw, ",")
	}
	return []string{raw}
}
