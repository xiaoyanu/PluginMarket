package repository

import (
	"math"

	"pluginmarket-server/database"
	"pluginmarket-server/model"

	"gorm.io/gorm"
)

type CommentTarget struct {
	CommentID     int  `json:"commentId"`
	RootCommentID int  `json:"rootCommentId"`
	IsReply       bool `json:"isReply"`
	Page          int  `json:"page"`
}

func CreateComment(comment *model.Comment) error {
	return database.DB.Create(comment).Error
}

func GetCommentByID(id int) (*model.Comment, error) {
	var comment model.Comment
	err := database.DB.First(&comment, id).Error
	return &comment, err
}

func GetCommentTarget(pluginID, commentID, pageSize int) (*CommentTarget, error) {
	var target model.Comment
	if err := database.DB.Where("id = ? AND plugin = ?", commentID, pluginID).First(&target).Error; err != nil {
		return nil, err
	}
	rootID := target.ID
	if target.Root > 0 {
		rootID = target.Root
	}
	var root model.Comment
	if err := database.DB.Where("id = ? AND plugin = ? AND root = 0 AND parent = 0", rootID, pluginID).First(&root).Error; err != nil {
		return nil, err
	}
	var newer int64
	if err := database.DB.Model(&model.Comment{}).
		Where("plugin = ? AND root = 0 AND parent = 0 AND created > ?", pluginID, root.Created).
		Count(&newer).Error; err != nil {
		return nil, err
	}
	return &CommentTarget{
		CommentID: commentID, RootCommentID: rootID, IsReply: target.Root > 0,
		Page: int(math.Floor(float64(newer)/float64(pageSize))) + 1,
	}, nil
}

func DeleteComment(id int) error {
	// 删除该评论及其所有层级子回复，避免留下孤儿评论
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var ids []int
		queue := []int{id}
		seen := map[int]struct{}{}

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			if _, ok := seen[current]; ok {
				continue
			}
			seen[current] = struct{}{}
			ids = append(ids, current)

			var children []int
			if err := tx.Model(&model.Comment{}).
				Where("parent = ? OR root = ?", current, current).
				Pluck("id", &children).Error; err != nil {
				return err
			}
			queue = append(queue, children...)
		}

		return tx.Where("id IN ?", ids).Delete(&model.Comment{}).Error
	})
}

func GetCommentsByPlugin(pluginID, page, pageSize int) ([]model.CommentListItem, int64, error) {
	// 只查顶级评论（root = 0 且 parent = 0）
	var total int64
	database.DB.Model(&model.Comment{}).Where("plugin = ? AND root = 0 AND parent = 0", pluginID).Count(&total)

	var comments []model.Comment
	err := database.DB.Where("plugin = ? AND root = 0 AND parent = 0", pluginID).
		Order("created DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.CommentListItem, len(comments))
	for i, c := range comments {
		result[i] = model.CommentListItem{
			ID:      c.ID,
			Content: c.Content,
			Root:    c.Root,
			Parent:  c.Parent,
			Created: c.Created,
		}

		// 作者信息
		author := getUserCommentAuthor(c.User)
		result[i].Author = author

		// 获取子回复
		var replies []model.Comment
		database.DB.Where("root = ?", c.ID).Order("created ASC").Find(&replies)

		result[i].ReplyCount = len(replies)
		result[i].Replies = make([]model.CommentReply, len(replies))
		for j, r := range replies {
			rAuthor := getUserCommentAuthor(r.User)
			reply := model.CommentReply{
				ID:        r.ID,
				Content:   r.Content,
				Root:      r.Root,
				Parent:    r.Parent,
				ReplyUser: r.ReplyUser,
				Created:   r.Created,
				Author:    rAuthor,
			}
			if r.ReplyUser > 0 {
				var replyToUser model.User
				database.DB.First(&replyToUser, r.ReplyUser)
				replyToUser.AnonymizeIfDeleted()
				reply.ReplyTo = &model.ReplyTo{ID: replyToUser.ID, Nick: replyToUser.Nick}
			}
			result[i].Replies[j] = reply
		}
	}

	return result, total, nil
}

func getUserCommentAuthor(userID int) model.CommentAuthor {
	var user model.User
	database.DB.First(&user, userID)
	user.AnonymizeIfDeleted()
	var titles []model.Title
	database.DB.Joins("JOIN pm_user_titles ON pm_user_titles.title_id = pm_titles.id").
		Where("pm_user_titles.user_id = ?", userID).
		Find(&titles)
	return model.CommentAuthor{
		ID:     user.ID,
		Nick:   user.Nick,
		Avatar: user.Avatar,
		Titles: titles,
	}
}
