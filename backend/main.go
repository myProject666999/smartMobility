package main

import (
	"fmt"
	"log"

	"smartMobility/config"
	"smartMobility/routes"
	"smartMobility/utils"
	"smartMobility/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	if err := utils.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	go websocket.WSHub.Run()
	log.Println("WebSocket服务已启动")

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("服务器已启动, 监听端口: %d", config.AppConfig.Server.Port)
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
