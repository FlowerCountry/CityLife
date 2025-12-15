// Package api 提供API路由配置
package api

import (
	"citylife/internal/api/handler"
	"citylife/internal/api/middleware"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine, sm *session.Manager) {
	// 全局中间件
	r.Use(middleware.CORS())

	// API v1
	v1 := r.Group("/api/v1")

	// Session管理（不需要验证）
	v1.POST("/sessions", handler.CreateSession(sm))

	// 需要Session验证的路由
	sessionGroup := v1.Group("/sessions/:session_id")
	sessionGroup.Use(middleware.RequireSession(sm))
	{
		// Session基础
		sessionGroup.GET("", handler.GetSession(sm))
		sessionGroup.DELETE("", handler.DeleteSession(sm))

		// 游戏状态
		sessionGroup.GET("/state", handler.GetState)
		sessionGroup.GET("/status", handler.GetStatus)

		// 行动系统
		sessionGroup.GET("/actions", handler.GetActions)
		sessionGroup.POST("/actions", handler.ExecuteAction)

		// 导航
		sessionGroup.POST("/navigate/:location_id", handler.Navigate)

		// 银行操作
		bank := sessionGroup.Group("/bank")
		{
			bank.POST("/deposit", handler.Deposit)
			bank.POST("/withdraw", handler.Withdraw)
			bank.POST("/deposit-all", handler.DepositAll)
			bank.POST("/withdraw-all", handler.WithdrawAll)
		}

		// 购物
		shop := sessionGroup.Group("/shop")
		{
			shop.GET("/commodities", handler.GetCommodities)
			shop.POST("/buy/:commodity", handler.BuyCommodity)
		}

		// 医疗
		hospital := sessionGroup.Group("/hospital")
		{
			hospital.POST("/doctor", handler.SeeDoctor)
			hospital.GET("/checkups", handler.GetCheckups)
			hospital.POST("/checkup/:id", handler.DoCheckup)
			hospital.GET("/medicines", handler.GetMedicines)
			hospital.POST("/medicine/:id", handler.BuyMedicine)
		}

		// 存档
		saves := sessionGroup.Group("/saves")
		{
			saves.GET("", handler.GetSaveSlots)
			saves.POST("/:slot", handler.SaveGame)
			saves.POST("/:slot/load", handler.LoadGame)
		}

		// 查看
		sessionGroup.GET("/wallet", handler.GetWallet)
		sessionGroup.GET("/bank", handler.GetBankBalance)
		sessionGroup.GET("/health", handler.GetHealth)
		sessionGroup.GET("/diseases", handler.GetDiseases)
		sessionGroup.GET("/map", handler.GetMap)
	}
}
