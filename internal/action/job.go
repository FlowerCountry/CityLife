package action

import (
	"fmt"

	"citylife/internal/game"
)

// WorkAction 工作行动
type WorkAction struct {
	id      string
	name    string
	seconds int
	income  int
	consume map[string]int
}

func (a *WorkAction) ID() string { return a.id }

func (a *WorkAction) Name() string { return a.name }

func (a *WorkAction) DurationSeconds() int { return a.seconds }

func (a *WorkAction) Income() int { return a.income }

func (a *WorkAction) Consume() map[string]int {
	out := make(map[string]int, len(a.consume))
	for k, v := range a.consume {
		out[k] = v
	}
	return out
}

func (a *WorkAction) Info() string {
	hours := float64(a.seconds) / 3600.0
	return fmt.Sprintf("%s（%.1f小时）+%d元", a.name, hours, a.income)
}

func (a *WorkAction) Execute(state *game.State) *Result {
	if a.income <= 0 || a.income%100 != 0 {
		return &Result{Message: "工作配置错误：收入必须为100的倍数", Success: false}
	}

	// 时间流逝
	state.World.UpdateTime(a.seconds)

	// 额外体力/心情消耗（工作强度）
	for k, v := range a.consume {
		if v > 0 {
			state.User.ConsumeNutrition(k, v)
		}
	}

	// 发工资（直接发100元面额）
	state.World.AddToWallet(0, a.income/100)

	return &Result{
		Message:     fmt.Sprintf("你完成了%s，获得%d元。", a.name, a.income),
		Success:     true,
		TimeElapsed: a.seconds,
	}
}

func (a *WorkAction) Category() EventCategory { return CategoryPrimary }

func GetJobActions() []Action {
	return []Action{
		&WorkAction{
			id:      "work_delivery",
			name:    "送外卖",
			seconds: 2 * 60 * 60,
			income:  300,
			consume: map[string]int{"饱腹感": 10, "饥渴": 12, "幸福感": 4},
		},
		&WorkAction{
			id:      "work_parttime",
			name:    "便利店兼职",
			seconds: 4 * 60 * 60,
			income:  600,
			consume: map[string]int{"饱腹感": 16, "饥渴": 18, "幸福感": 6},
		},
		&WorkAction{
			id:      "work_construction",
			name:    "工地搬砖",
			seconds: 6 * 60 * 60,
			income:  1000,
			consume: map[string]int{"饱腹感": 24, "饥渴": 26, "幸福感": 10, "精神振奋": 8},
		},
	}
}
