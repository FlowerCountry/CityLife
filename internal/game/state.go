// Package game 管理游戏状态
package game

import (
	"citylife/internal/disease"
	"citylife/internal/locale"
	"citylife/internal/nutrition"
	"citylife/internal/user"
	"citylife/internal/world"
)

// State 游戏状态
type State struct {
	World     *world.World
	User      *user.User
	Nutrition *nutrition.Manager
	Disease   *disease.Manager
	Locale    *locale.Locale
}

// New 创建新游戏状态
func New() *State {
	u := user.New()
	nm := nutrition.NewManager()
	nm.InitPreviousLevels(u)

	return &State{
		World:     world.New(),
		User:      u,
		Nutrition: nm,
		Disease:   disease.NewManager(),
		Locale:    locale.NewWithDefaults(),
	}
}

// GetWorld 获取World（实现action.GameState接口）
func (s *State) GetWorld() interface{} {
	return s.World
}

// GetUser 获取User（实现action.GameState接口）
func (s *State) GetUser() interface{} {
	return s.User
}

// Reset 重置游戏状态
func (s *State) Reset() {
	s.World = world.New()
	s.User = user.New()
}
