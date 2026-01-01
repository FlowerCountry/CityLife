// Package action 定义游戏行动接口
package action

import (
	"citylife/internal/game"
)

// EventCategory 行动分类（用于菜单排序）
type EventCategory int

const (
	CategoryPrimary    EventCategory = iota // 主要行动（银行、交易、购买）
	CategoryInsight                         // 信息行动（查看状态、信息）
	CategoryNavigation                      // 导航行动（移动）
)

// CategoryName 返回分类名称
func (c EventCategory) String() string {
	switch c {
	case CategoryPrimary:
		return "primary"
	case CategoryInsight:
		return "insight"
	case CategoryNavigation:
		return "navigation"
	default:
		return "unknown"
	}
}

// Action 行动接口
type Action interface {
	// ID 返回行动唯一标识
	ID() string

	// Info 返回显示文本
	Info() string

	// Execute 执行行动，返回结果
	Execute(state *game.State) *Result

	// Category 返回行动分类
	Category() EventCategory
}

// Result 行动执行结果
type Result struct {
	Message     string // 显示消息
	Success     bool   // 是否成功
	TimeElapsed int    // 消耗的秒数
}

// SortActions 按分类排序行动
func SortActions(actions []Action) []Action {
	// 简单的稳定排序：按分类排序
	result := make([]Action, len(actions))
	copy(result, actions)

	// 冒泡排序（稳定）
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-1-i; j++ {
			if result[j].Category() > result[j+1].Category() {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	return result
}
