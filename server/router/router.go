package router

import (
	"pluginmarket-server/controller"
	"pluginmarket-server/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		// 网站公开基础设置，不返回邮件等私有配置。
		api.GET("/setting/public", controller.GetPublicSiteSettings)

		// ===== 认证模块 =====
		auth := api.Group("/auth")
		{
			auth.GET("/captcha", controller.GetCaptcha)
			auth.POST("/login", controller.Login)
			auth.POST("/register", controller.Register)
			auth.POST("/forgot-password", controller.ForgotPassword)
			auth.POST("/reset-password", controller.ResetPassword)
		}

		// ===== 用户模块 =====
		user := api.Group("/user")
		{
			// 需要登录
			userAuth := user.Group("", middleware.AuthRequired())
			{
				userAuth.GET("/info", controller.GetUserInfo)
				userAuth.PUT("/info", controller.UpdateUserInfo)
				userAuth.PUT("/password", controller.ChangePassword)
				userAuth.POST("/avatar", controller.UploadAvatar)
				userAuth.POST("/email/send-verify", controller.SendEmailVerify)
			}
			// 公开接口
			user.GET("/email/verify", controller.EmailVerify)
			user.GET("/:id", controller.GetUserPublic)
		}

		// ===== 插件模块 =====
		plugin := api.Group("/plugin")
		{
			// 公开接口
			plugin.GET("/list/hot", controller.GetPluginListHot)
			plugin.GET("/list/latest", controller.GetPluginListLatest)
			plugin.GET("/search", controller.SearchPlugins)
			plugin.GET("/list/by-tag", controller.GetPluginsByTag)
			plugin.GET("/list/by-user", controller.GetPluginsByUser)
			// 插件详情（可选登录判断收藏状态）
			plugin.GET("/:id", middleware.OptionalAuth(), controller.GetPluginDetail)
			plugin.GET("/:id/download", controller.DownloadPlugin)

			// 需要登录
			pluginAuth := plugin.Group("", middleware.AuthRequired())
			{
				pluginAuth.GET("/my", controller.GetMyPlugins)
				pluginAuth.POST("/publish", controller.PublishPlugin)
				pluginAuth.PUT("/:id", controller.EditPlugin)
				pluginAuth.DELETE("/:id", controller.DeletePlugin)
			}
		}

		// ===== 评论模块 =====
		comment := api.Group("/comment")
		{
			comment.GET("/list", controller.GetCommentList)
			commentAuth := comment.Group("", middleware.AuthRequired())
			{
				commentAuth.POST("", controller.CreateComment)
				commentAuth.DELETE("/:id", controller.DeleteComment)
			}
		}

		// ===== 收藏模块 =====
		star := api.Group("/star", middleware.AuthRequired())
		{
			star.POST("/toggle", controller.ToggleStar)
			star.GET("/list", controller.GetStarList)
		}

		// ===== 通知模块 =====
		notification := api.Group("/notification", middleware.AuthRequired())
		{
			notification.GET("/list", controller.GetNotificationList)
			notification.GET("/unread-count", controller.GetUnreadNotificationCount)
			notification.PUT("/read-all", controller.MarkAllNotificationsRead)
			notification.PUT("/:id/read", controller.MarkNotificationRead)
			notification.DELETE("/:id", controller.HideNotification)
		}

		// ===== 标签模块 =====
		tag := api.Group("/tag")
		{
			tag.GET("/list", controller.GetTagList)
			tagAdmin := tag.Group("", middleware.AuthRequired(), middleware.AdminRequired())
			{
				tagAdmin.POST("", controller.CreateTag)
				tagAdmin.PUT("/:id", controller.UpdateTag)
				tagAdmin.DELETE("/:id", controller.DeleteTag)
			}
		}

		// ===== 框架模块 =====
		frame := api.Group("/frame")
		{
			frame.GET("/list", controller.GetFrameList)
			frameAdmin := frame.Group("", middleware.AuthRequired(), middleware.AdminRequired())
			{
				frameAdmin.POST("", controller.CreateFrame)
				frameAdmin.PUT("/:id", controller.UpdateFrame)
				frameAdmin.DELETE("/:id", controller.DeleteFrame)
			}
		}

		// ===== 称号模块 =====
		title := api.Group("/title")
		{
			title.GET("/list", controller.GetTitleList)
			titleAdmin := title.Group("", middleware.AuthRequired(), middleware.AdminRequired())
			{
				titleAdmin.POST("", controller.CreateTitle)
				titleAdmin.PUT("/:id", controller.UpdateTitle)
				titleAdmin.DELETE("/:id", controller.DeleteTitle)
			}
		}

		// ===== 管理员模块 =====
		admin := api.Group("/admin", middleware.AuthRequired(), middleware.AdminRequired())
		{
			// 插件审核
			admin.GET("/plugin/pending", controller.GetPendingPlugins)
			admin.PUT("/plugin/:id/approve", controller.ApprovePlugin)
			admin.PUT("/plugin/:id/reject", controller.RejectPlugin)
			admin.PUT("/plugin/:id/status", controller.AdminUpdatePluginStatus)
			admin.DELETE("/plugin/:id", controller.AdminDeletePlugin)

			// 用户管理
			admin.GET("/user/list", controller.AdminGetUserList)
			admin.GET("/user/:id", controller.AdminGetUserDetail)
			admin.GET("/user/:id/plugins", controller.AdminGetUserPlugins)
			admin.PUT("/user/:id/reset-password", controller.AdminResetPassword)
			admin.PUT("/user/:id/role", controller.AdminUpdateRole)
			admin.PUT("/user/:id/titles", controller.AdminUpdateTitles)
			admin.DELETE("/user/:id", controller.AdminDeleteUser)

			// 通知管理
			admin.GET("/notification/list", controller.AdminGetNotificationList)
			admin.GET("/notification/:id", controller.AdminGetNotificationDetail)
			admin.POST("/notification", controller.AdminCreateNotification)
			admin.PUT("/notification/:id", controller.AdminUpdateNotification)
			admin.DELETE("/notification/:id", controller.AdminDeleteNotification)

			// 全局设置
			admin.GET("/setting", controller.GetSettings)
			admin.PUT("/setting", controller.UpdateSettings)
			admin.POST("/setting/logo", controller.UploadLogo)
			admin.DELETE("/setting/logo", controller.ClearLogo)
			admin.POST("/setting/test-email", controller.TestEmail)
			admin.POST("/setting/clean-images", controller.CleanImages)
		}

		// ===== 文件上传 =====
		upload := api.Group("/upload", middleware.AuthRequired())
		{
			upload.POST("/image", controller.UploadImage)
		}
	}

	return r
}
