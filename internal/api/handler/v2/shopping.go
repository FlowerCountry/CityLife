package v2

import (
	"fmt"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/food"
	"citylife/internal/world"

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
			effects = append(effects, fmt.Sprintf("%s%+d", h.Name, h.Amount))
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

	ok(c, commodities)
}

// BuyCommodity 购买商品
func BuyCommodity(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	if !state.User.IsAlive() {
		failGameOver(c, state)
		return
	}

	if state.World.Where != world.LocationSupermarket && state.World.Where != world.LocationSupermarketInner {
		fail(c, response.ErrCodeInvalidLocation, "你需要在超市才能购买商品", actionData{State: buildState(state)})
		return
	}

	commodityName := c.Param("commodity")

	var targetFood *food.Food
	for _, f := range food.AllCommodities {
		if f.Name == commodityName {
			targetFood = f
			break
		}
	}

	if targetFood == nil {
		response.NotFound(c, response.ErrCodeInvalidRequest, "商品不存在")
		return
	}

	failCode := ""
	if state.World.GetWalletTotal() < targetFood.Price {
		failCode = response.ErrCodeInsufficientFunds
	}

	executeAndRespond(c, state, &action.CommodityAction{Food: targetFood}, failCode)
}
