package main

import (
	"flag"
	"log"
	"time"

	"citylife/internal/api"
	"citylife/internal/path"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

func main() {
	// 命令行参数
	port := flag.String("port", "8080", "服务器端口")
	debug := flag.Bool("debug", false, "是否开启调试模式")
	flag.Parse()

	// 初始化路径管理器
	if err := initPathManager(); err != nil {
		log.Printf("路径管理器初始化警告: %v\n", err)
	}

	// 设置Gin模式
	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Session管理器（30分钟超时）
	sessionMgr := session.NewManager(30 * time.Minute)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r, sessionMgr)

	// 启动服务器
	log.Printf("CityLife API Server starting on :%s\n", *port)
	log.Println("Base path: /api/v2")
	if *debug {
		log.Println("Registered routes:")
		for _, rt := range r.Routes() {
			log.Printf("  %-6s %s\n", rt.Method, rt.Path)
		}
	}

	if err := r.Run(":" + *port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// initPathManager 初始化路径管理器
func initPathManager() error {
	pm := path.GetInstance()

	// 从可执行文件路径自动初始化
	if err := pm.InitializeFromExecutable(); err != nil {
		return err
	}

	// 确保所有目录存在
	if err := pm.EnsureDirectories(); err != nil {
		return err
	}

	return nil
}
