// Package action 包含购物相关行动
package action

import (
	"fmt"
	"strings"

	"citylife/internal/food"
	"citylife/internal/game"
	"citylife/internal/payment"
)

// CommodityAction 购买商品行动
type CommodityAction struct {
	Food *food.Food
}

func (a *CommodityAction) ID() string {
	return "buy_" + a.Food.Name
}

func (a *CommodityAction) Info() string {
	return a.Food.GetInfo()
}

func (a *CommodityAction) Execute(state *game.State) *Result {
	price := a.Food.Price

	// 检查钱包余额
	if state.World.GetWalletTotal() < price {
		return &Result{
			Message: fmt.Sprintf("购买失败！现金不足，需要 %d 元，当前只有 %d 元", price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	// 使用Cashier进行支付
	cashier := payment.NewCashier()
	result := cashier.Pay(price, state.World.Wallet)

	if !result.Success {
		return &Result{
			Message: state.Locale.Get("purchase.error", "insufficient"),
			Success: false,
		}
	}

	// 更新钱包
	for i := range state.World.Wallet {
		state.World.Wallet[i] -= result.Used[i]
		state.World.Wallet[i] += result.Change[i]
	}

	// 应用营养效果
	effectiveHealth := a.Food.GetEffectiveHealth()
	var effects []string
	for _, h := range effectiveHealth {
		// 处理"饥饿"作为"饱腹感"的别名
		name := h.Name
		if name == "饥饿" {
			name = "饱腹感"
		}
		state.User.AddNutrition(name, h.Amount)
		if h.Amount > 0 {
			effects = append(effects, fmt.Sprintf("%s+%d", name, h.Amount))
		} else {
			effects = append(effects, fmt.Sprintf("%s%d", name, h.Amount))
		}
	}

	// 检查食物中毒风险（基于新鲜度）
	poisonRisk := a.Food.GetFoodPoisoningRisk()
	poisonWarning := ""
	if poisonRisk > 0 {
		state.Disease.TriggerFoodPoisoning(poisonRisk)
		if state.Disease.HasDisease("food_poisoning") {
			poisonWarning = "\n⚠️ 你感到肚子不太舒服..."
		}
	}

	// 更新时间（购物消耗1分钟）
	state.World.UpdateTime(60)

	// 构建消息
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("购买成功！%s（%d元）\n", a.Food.Name, price))

	// 显示支付详情
	if result.Paid > price {
		msg.WriteString(fmt.Sprintf("实付: %d元，找零: %d元\n", result.Paid, result.Paid-price))
	}

	msg.WriteString("营养效果: " + strings.Join(effects, ", "))
	msg.WriteString(poisonWarning)

	return &Result{
		Message:     msg.String(),
		Success:     true,
		TimeElapsed: 60,
	}
}

func (a *CommodityAction) Category() EventCategory {
	return CategoryPrimary
}

// GetCommodityActions 获取所有商品行动
func GetCommodityActions() []Action {
	actions := make([]Action, len(food.AllCommodities))
	for i, f := range food.AllCommodities {
		actions[i] = &CommodityAction{Food: f}
	}
	return actions
}
