package main

import (
	"fmt"
	"log"
	"time"

	"pluginmarket-server/config"
	"pluginmarket-server/database"
	"pluginmarket-server/model"
	"pluginmarket-server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	if err := database.Init(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构
	_ = database.DB.AutoMigrate(
		&model.User{},
		&model.Plugin{},
		&model.Comment{},
		&model.Tag{},
		&model.Frame{},
		&model.Title{},
		&model.UserTitle{},
		&model.UserStar{},
		&model.PluginFrame{},
		&model.PluginTag{},
		&model.Setting{},
		&model.Notification{},
		&model.UserNotificationState{},
	)
	if err := database.BackfillTagTimestamps(time.Now()); err != nil {
		log.Printf("补齐标签日期失败: %v", err)
	}
	if err := database.EnsureUniqueBoundEmails(); err != nil {
		log.Fatalf("创建用户邮箱唯一约束失败: %v", err)
	}

	// 初始化种子数据
	database.Seed()

	// 设置 Gin 模式
	gin.SetMode(config.C.Server.Mode)

	// 初始化路由
	r := router.Setup()

	// 启动服务
	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	log.Println("\n╔═╗┬  ┬ ┬┌─┐┬┌┐┌╔╦╗┌─┐┬─┐┬┌─┌─┐┌┬┐\n╠═╝│  │ ││ ┬││││║║║├─┤├┬┘├┴┐├┤  │ \n╩  ┴─┘└─┘└─┘┴┘└┘╩ ╩┴ ┴┴└─┴ ┴└─┘ ┴ ")
	log.Printf("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
