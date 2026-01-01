package action

import (
	"fmt"

	"citylife/internal/game"
)

// PayAndApplyAction 通用：花钱+时间+营养效果
type PayAndApplyAction struct {
	id      string
	name    string
	price   int
	seconds int
	add     map[string]int
	consume map[string]int
}

func (a *PayAndApplyAction) ID() string { return a.id }

func (a *PayAndApplyAction) Name() string { return a.name }

func (a *PayAndApplyAction) Price() int { return a.price }

func (a *PayAndApplyAction) DurationSeconds() int { return a.seconds }

func (a *PayAndApplyAction) Add() map[string]int {
	out := make(map[string]int, len(a.add))
	for k, v := range a.add {
		out[k] = v
	}
	return out
}

func (a *PayAndApplyAction) Consume() map[string]int {
	out := make(map[string]int, len(a.consume))
	for k, v := range a.consume {
		out[k] = v
	}
	return out
}

func (a *PayAndApplyAction) Info() string {
	if a.price <= 0 {
		return a.name
	}
	return fmt.Sprintf("%s（%d元）", a.name, a.price)
}

func (a *PayAndApplyAction) Execute(state *game.State) *Result {
	if a.price > 0 {
		if err := state.World.SpendMoney(a.price); err != nil {
			return &Result{
				Message: fmt.Sprintf("现金不足，需要 %d 元，当前只有 %d 元", a.price, state.World.GetWalletTotal()),
				Success: false,
			}
		}
	}

	if a.seconds > 0 {
		state.World.UpdateTime(a.seconds)
	}

	for k, v := range a.consume {
		if v > 0 {
			state.User.ConsumeNutrition(k, v)
		}
	}
	for k, v := range a.add {
		if v > 0 {
			state.User.AddNutrition(k, v)
		}
	}

	return &Result{
		Message:     fmt.Sprintf("你进行了：%s。", a.name),
		Success:     true,
		TimeElapsed: a.seconds,
	}
}

func (a *PayAndApplyAction) Category() EventCategory { return CategoryPrimary }

func GetRestaurantActions() []Action {
	return []Action{
		&PayAndApplyAction{
			id:      "eat_fastfood",
			name:    "吃快餐",
			price:   50,
			seconds: 30 * 60,
			add:     map[string]int{"饱腹感": 22, "饥渴": 12, "幸福感": 6},
		},
		&PayAndApplyAction{
			id:      "eat_setmeal",
			name:    "吃套餐",
			price:   120,
			seconds: 45 * 60,
			add:     map[string]int{"饱腹感": 38, "饥渴": 18, "幸福感": 12},
		},
		&PayAndApplyAction{
			id:      "eat_feast",
			name:    "吃一顿好的",
			price:   300,
			seconds: 60 * 60,
			add:     map[string]int{"饱腹感": 55, "饥渴": 22, "幸福感": 22, "精神振奋": 10},
		},
	}
}

func GetParkActions() []Action {
	return []Action{
		&PayAndApplyAction{
			id:      "walk_park",
			name:    "公园散步",
			price:   0,
			seconds: 30 * 60,
			add:     map[string]int{"幸福感": 10},
			consume: map[string]int{"饱腹感": 6, "饥渴": 7},
		},
		&PayAndApplyAction{
			id:      "jog_park",
			name:    "公园慢跑",
			price:   0,
			seconds: 60 * 60,
			add:     map[string]int{"幸福感": 16, "精神振奋": 6},
			consume: map[string]int{"饱腹感": 10, "饥渴": 12},
		},
	}
}

func GetHotelActions() []Action {
	return []Action{
		&PayAndApplyAction{
			id:      "hotel_sleep",
			name:    "旅馆住店睡觉（8小时）",
			price:   500,
			seconds: 8 * 60 * 60,
			add:     map[string]int{"精神振奋": 35, "幸福感": 14},
			consume: map[string]int{"饱腹感": 20, "饥渴": 24},
		},
	}
}
