package handler

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"
	"citylife/internal/food"

	"github.com/gin-gonic/gin"
)

// GetCommodities 获取商品列表
func GetCommodities(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	commodities := make([]response.CommodityInfo, len(food.AllCommodities))
	for i, f := range food.AllCommodities {
		var effects []string
		for _, h := range f.GetEffectiveHealth() {
			if h.Amount > 0 {
				effects = append(effects, h.Name+"+"+string(rune('0'+h.Amount/10))+string(rune('0'+h.Amount%10)))
			}
		}

		commodities[i] = response.CommodityInfo{
			Name:           f.Name,
			Price:          f.Price,
			Type:           food.FoodTypeNames[f.Type],
			Freshness:      f.Freshness,
			FreshnessLabel: f.GetFreshnessLabel(),
			Effects:        effects,
		}
	}

	response.Success(c, commodities)
}

// BuyCommodity 购买商品
func BuyCommodity(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	// 检查游戏是否结束
	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	commodityName := c.Param("commodity")

	// 查找商品
	var targetFood *food.Food
	for _, f := range food.AllCommodities {
		if f.Name == commodityName {
			targetFood = f
			break
		}
	}

	if targetFood == nil {
		response.BadRequest(c, "商品不存在")
		return
	}

	// 创建并执行行动
	act := &action.CommodityAction{Food: targetFood}
	result := executor.Execute(act, state)

	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	response.Success(c, resp)
}
