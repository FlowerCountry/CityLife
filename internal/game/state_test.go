package game

import (
	"testing"

	"citylife/internal/user"
	"citylife/internal/world"
)

func TestNew(t *testing.T) {
	s := New()

	t.Run("World初始化", func(t *testing.T) {
		if s.World == nil {
			t.Error("World should not be nil")
		}
	})

	t.Run("User初始化", func(t *testing.T) {
		if s.User == nil {
			t.Error("User should not be nil")
		}
	})

	t.Run("Nutrition初始化", func(t *testing.T) {
		if s.Nutrition == nil {
			t.Error("Nutrition should not be nil")
		}
	})

	t.Run("Disease初始化", func(t *testing.T) {
		if s.Disease == nil {
			t.Error("Disease should not be nil")
		}
	})
}

func TestGetWorld(t *testing.T) {
	s := New()

	result := s.GetWorld()

	t.Run("返回World", func(t *testing.T) {
		w, ok := result.(*world.World)
		if !ok {
			t.Error("GetWorld() should return *world.World")
		}
		if w != s.World {
			t.Error("GetWorld() should return the same World instance")
		}
	})
}

func TestGetUser(t *testing.T) {
	s := New()

	result := s.GetUser()

	t.Run("返回User", func(t *testing.T) {
		u, ok := result.(*user.User)
		if !ok {
			t.Error("GetUser() should return *user.User")
		}
		if u != s.User {
			t.Error("GetUser() should return the same User instance")
		}
	})
}

func TestReset(t *testing.T) {
	s := New()

	// 修改World状态
	s.World.Year = 2020
	s.World.BankDeposit = 1000
	s.World.Where = 2

	// 修改User状态
	s.User.SetNutrition("饱腹感", 50)

	// 记录旧实例
	oldWorld := s.World
	oldUser := s.User

	// 重置
	s.Reset()

	t.Run("创建新World实例", func(t *testing.T) {
		if s.World == oldWorld {
			t.Error("Reset() should create a new World instance")
		}
	})

	t.Run("创建新User实例", func(t *testing.T) {
		if s.User == oldUser {
			t.Error("Reset() should create a new User instance")
		}
	})

	t.Run("World状态重置", func(t *testing.T) {
		if s.World.Year != 2010 {
			t.Errorf("Year = %d, want 2010", s.World.Year)
		}
		if s.World.BankDeposit != 0 {
			t.Errorf("BankDeposit = %d, want 0", s.World.BankDeposit)
		}
		if s.World.Where != 0 {
			t.Errorf("Where = %d, want 0", s.World.Where)
		}
	})

	t.Run("User状态重置", func(t *testing.T) {
		if s.User.GetNutrition("饱腹感") != 100 {
			t.Errorf("饱腹感 = %d, want 100", s.User.GetNutrition("饱腹感"))
		}
	})
}

func TestStateIntegration(t *testing.T) {
	s := New()

	t.Run("World建筑数量", func(t *testing.T) {
		if len(s.World.Buildings) == 0 {
			t.Error("World should have buildings")
		}
	})

	t.Run("User营养初始值", func(t *testing.T) {
		// 验证User已正确初始化
		if s.User.GetNutrition("饱腹感") != 100 {
			t.Errorf("Initial 饱腹感 = %d, want 100", s.User.GetNutrition("饱腹感"))
		}
	})

	t.Run("初始钱包金额", func(t *testing.T) {
		total := s.World.GetWalletTotal()
		if total != 100 {
			t.Errorf("Initial wallet = %d, want 100", total)
		}
	})

	t.Run("User存活状态", func(t *testing.T) {
		if !s.User.IsAlive() {
			t.Error("New user should be alive")
		}
	})
}

func TestMultipleStates(t *testing.T) {
	s1 := New()
	s2 := New()

	t.Run("独立World实例", func(t *testing.T) {
		if s1.World == s2.World {
			t.Error("Different states should have different World instances")
		}
	})

	t.Run("独立User实例", func(t *testing.T) {
		if s1.User == s2.User {
			t.Error("Different states should have different User instances")
		}
	})

	t.Run("状态隔离", func(t *testing.T) {
		s1.World.Year = 2025
		if s2.World.Year == 2025 {
			t.Error("Changes to one state should not affect another")
		}
	})
}
