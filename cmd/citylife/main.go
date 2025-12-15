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
	log.Println("API文档:")
	log.Println("  POST   /api/v1/sessions                    创建游戏会话")
	log.Println("  GET    /api/v1/sessions/:id/state          获取游戏状态")
	log.Println("  GET    /api/v1/sessions/:id/actions        获取可用行动")
	log.Println("  POST   /api/v1/sessions/:id/actions        执行行动")
	log.Println("  POST   /api/v1/sessions/:id/navigate/:loc  导航到位置")
	log.Println("  POST   /api/v1/sessions/:id/bank/*         银行操作")
	log.Println("  GET    /api/v1/sessions/:id/shop/*         购物")
	log.Println("  POST   /api/v1/sessions/:id/hospital/*     医疗")
	log.Println("  GET    /api/v1/sessions/:id/saves          存档")

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
