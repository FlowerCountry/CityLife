// Package api 提供API路由配置
package api

import (
	handlerv2 "citylife/internal/api/handler/v2"
	"citylife/internal/api/middleware"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine, sm *session.Manager) {
	// 全局中间件
	r.Use(middleware.CORS())

	// API v2
	v2 := r.Group("/api/v2")

	v2.POST("/sessions", handlerv2.CreateSession(sm))

	v2Session := v2.Group("/sessions/:session_id")
	v2Session.Use(middleware.RequireSession(sm))
	{
		v2Session.GET("", handlerv2.GetSession(sm))
		v2Session.DELETE("", handlerv2.DeleteSession(sm))

		v2Session.GET("/state", handlerv2.GetState)
		v2Session.GET("/status", handlerv2.GetStatus)

		v2Session.GET("/actions", handlerv2.GetActions)
		v2Session.POST("/actions/:action_id", handlerv2.ExecuteAction)

		v2Session.POST("/navigate/:location_id", handlerv2.Navigate)

		bank := v2Session.Group("/bank")
		{
			bank.POST("/deposit", handlerv2.Deposit)
			bank.POST("/withdraw", handlerv2.Withdraw)
			bank.POST("/deposit-all", handlerv2.DepositAll)
			bank.POST("/withdraw-all", handlerv2.WithdrawAll)
		}

		shop := v2Session.Group("/shop")
		{
			shop.GET("/commodities", handlerv2.GetCommodities)
			shop.POST("/buy/:commodity", handlerv2.BuyCommodity)
		}

		hospital := v2Session.Group("/hospital")
		{
			hospital.POST("/doctor", handlerv2.SeeDoctor)
			hospital.GET("/checkups", handlerv2.GetCheckups)
			hospital.POST("/checkup/:id", handlerv2.DoCheckup)
			hospital.GET("/medicines", handlerv2.GetMedicines)
			hospital.POST("/medicine/:id", handlerv2.BuyMedicine)
		}

		saves := v2Session.Group("/saves")
		{
			saves.GET("", handlerv2.GetSaveSlots)
			saves.POST("/:slot", handlerv2.SaveGame)
			saves.POST("/:slot/load", handlerv2.LoadGame)
		}

		housing := v2Session.Group("/housing")
		{
			housing.GET("", handlerv2.GetHousing)
			housing.GET("/offers", handlerv2.GetHousingOffers)
			housing.POST("/rent/:level", handlerv2.RentHousing)
			housing.POST("/buy/:level", handlerv2.BuyHousing)
			housing.POST("/renew", handlerv2.RenewRent)
			housing.POST("/cancel", handlerv2.CancelRent)
			housing.POST("/go-home", handlerv2.GoHome)
			housing.POST("/leave-home", handlerv2.LeaveHome)
			housing.POST("/sleep", handlerv2.SleepAtHome)
		}

		jobs := v2Session.Group("/jobs")
		{
			jobs.GET("", handlerv2.GetJobs)
			jobs.POST("/:id", handlerv2.DoJob)
		}

		restaurant := v2Session.Group("/restaurant")
		{
			restaurant.GET("/menu", handlerv2.GetRestaurantMenu)
			restaurant.POST("/:id", handlerv2.DoRestaurant)
		}
		park := v2Session.Group("/park")
		{
			park.GET("/activities", handlerv2.GetParkActivities)
			park.POST("/:id", handlerv2.DoPark)
		}
		hotel := v2Session.Group("/hotel")
		{
			hotel.GET("/services", handlerv2.GetHotelServices)
			hotel.POST("/:id", handlerv2.DoHotel)
		}

		v2Session.GET("/wallet", handlerv2.GetWallet)
		v2Session.GET("/bank", handlerv2.GetBankBalance)
		v2Session.GET("/health", handlerv2.GetHealth)
		v2Session.GET("/diseases", handlerv2.GetDiseases)
		v2Session.GET("/map", handlerv2.GetMap)
	}
}
