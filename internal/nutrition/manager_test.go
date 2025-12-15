package nutrition

import (
	"testing"

	"citylife/internal/user"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	t.Run("配置初始化", func(t *testing.T) {
		if len(m.configs) == 0 {
			t.Error("Expected configs to be initialized")
		}
	})

	t.Run("交互规则初始化", func(t *testing.T) {
		if len(m.interactions) == 0 {
			t.Error("Expected interactions to be initialized")
		}
	})

	t.Run("提示消息初始化", func(t *testing.T) {
		if len(m.hintMessages) == 0 {
			t.Error("Expected hint messages to be initialized")
		}
	})
}

func TestOnAction(t *testing.T) {
	m := NewManager()
	u := user.New()
	m.InitPreviousLevels(u)

	initialSatiety := u.GetNutrition("饱腹感")

	// 执行一次行动
	m.OnAction(u, 1.0)

	newSatiety := u.GetNutrition("饱腹感")

	t.Run("营养衰减", func(t *testing.T) {
		if newSatiety >= initialSatiety {
			t.Errorf("Satiety should decrease: %d -> %d", initialSatiety, newSatiety)
		}
	})
}

func TestOnActionWithMultiplier(t *testing.T) {
	m := NewManager()
	u1 := user.New()
	u2 := user.New()
	m.InitPreviousLevels(u1)
	m.InitPreviousLevels(u2)

	// 正常衰减
	m.OnAction(u1, 1.0)
	// 加速衰减
	m.OnAction(u2, 2.0)

	s1 := u1.GetNutrition("饱腹感")
	s2 := u2.GetNutrition("饱腹感")

	t.Run("倍率影响衰减", func(t *testing.T) {
		if s2 >= s1 {
			t.Errorf("Higher multiplier should cause more decay: %d (1x) vs %d (2x)", s1, s2)
		}
	})
}

func TestGetDecayMultiplier(t *testing.T) {
	m := NewManager()

	tests := []struct {
		name     string
		value    int
		expected float64
	}{
		{"正常区", 50, 1.0},
		{"危险区", 10, 1.2},
		{"透支区", -5, 1.5},
		{"濒死区", -15, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.New()
			u.SetNutrition("饱腹感", tt.value)
			if got := m.GetDecayMultiplier(u, "饱腹感"); got != tt.expected {
				t.Errorf("GetDecayMultiplier() = %f, want %f (value=%d)", got, tt.expected, tt.value)
			}
		})
	}
}

func TestCheckLevelChanges(t *testing.T) {
	m := NewManager()
	u := user.New()
	m.InitPreviousLevels(u)

	// 清空提示
	m.ClearPendingHints()

	// 将饱腹感降到警告区
	u.SetNutrition("饱腹感", 30)
	m.CheckLevelChanges(u)

	hints := m.GetPendingHints()
	t.Run("等级变化生成提示", func(t *testing.T) {
		if len(hints) == 0 {
			t.Error("Expected hints when nutrition level changes")
		}
	})
}

func TestGetPendingHints(t *testing.T) {
	m := NewManager()

	t.Run("初始无提示", func(t *testing.T) {
		if len(m.GetPendingHints()) != 0 {
			t.Error("Expected no pending hints initially")
		}
	})
}

func TestClearPendingHints(t *testing.T) {
	m := NewManager()
	u := user.New()
	m.InitPreviousLevels(u)

	// 触发一些提示
	u.SetNutrition("饱腹感", 10)
	m.CheckLevelChanges(u)

	m.ClearPendingHints()

	t.Run("清空后无提示", func(t *testing.T) {
		if len(m.GetPendingHints()) != 0 {
			t.Error("Expected no pending hints after clear")
		}
	})
}

func TestCheckRandomDeath(t *testing.T) {
	m := NewManager()

	t.Run("健康状态无死亡", func(t *testing.T) {
		u := user.New()
		death := false
		for i := 0; i < 100; i++ {
			if m.CheckRandomDeath(u) {
				death = true
				break
			}
		}
		if death {
			t.Error("Healthy user should not die randomly")
		}
	})

	t.Run("濒死状态有死亡风险", func(t *testing.T) {
		u := user.New()
		u.SetNutrition("饱腹感", -15)
		// 只验证不会panic
		_ = m.CheckRandomDeath(u)
	})
}

func TestGetWarnings(t *testing.T) {
	m := NewManager()
	u := user.New()

	t.Run("健康状态无警告", func(t *testing.T) {
		warnings := m.GetWarnings(u)
		if len(warnings) != 0 {
			t.Errorf("Expected no warnings, got %d", len(warnings))
		}
	})

	t.Run("警告区有警告", func(t *testing.T) {
		u.SetNutrition("饱腹感", 30)
		warnings := m.GetWarnings(u)
		if len(warnings) == 0 {
			t.Error("Expected warnings for low nutrition")
		}
	})
}

func TestGetCriticals(t *testing.T) {
	m := NewManager()
	u := user.New()

	t.Run("健康状态无危险", func(t *testing.T) {
		criticals := m.GetCriticals(u)
		if len(criticals) != 0 {
			t.Errorf("Expected no criticals, got %d", len(criticals))
		}
	})

	t.Run("危险区有危险", func(t *testing.T) {
		u.SetNutrition("饱腹感", 10)
		criticals := m.GetCriticals(u)
		if len(criticals) == 0 {
			t.Error("Expected criticals for dangerous nutrition")
		}
	})
}

func TestProcessInteractions(t *testing.T) {
	m := NewManager()
	u := user.New()

	// 测试交互规则不会panic
	t.Run("交互处理不panic", func(t *testing.T) {
		// 设置各种边界条件
		u.SetNutrition("铁", 0)
		u.SetNutrition("维生素D", 0)
		u.SetNutrition("糖分", 100)

		// 只验证不会panic
		m.ProcessInteractions(u)
	})

	// 测试极端情况
	t.Run("极端低值处理", func(t *testing.T) {
		u.SetNutrition("铁", -10)
		u.SetNutrition("蛋白质", 100)

		initialProtein := u.GetNutrition("蛋白质")
		m.ProcessInteractions(u)
		newProtein := u.GetNutrition("蛋白质")

		// 铁极低时应该影响蛋白质
		// penalty = (30-(-10))/100 * 0.3 = 0.12
		// extraDecay = int(100 * 0.12 * 0.1) = int(1.2) = 1
		if newProtein > initialProtein {
			t.Errorf("Iron deficiency should not increase protein: %d -> %d", initialProtein, newProtein)
		}
	})
}

func TestInitPreviousLevels(t *testing.T) {
	m := NewManager()
	u := user.New()

	m.InitPreviousLevels(u)

	t.Run("记录初始等级", func(t *testing.T) {
		if len(m.previousLevels) == 0 {
			t.Error("Expected previous levels to be recorded")
		}
	})
}

func TestNutritionConfigs(t *testing.T) {
	m := NewManager()

	tests := []struct {
		name     string
		category Category
	}{
		{"饱腹感", CategoryPrimary},
		{"饥渴", CategoryPrimary},
		{"钙", CategoryMicro},
		{"幸福感", CategorySpecial},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, ok := m.configs[tt.name]
			if !ok {
				t.Fatalf("Config for %q not found", tt.name)
			}
			if cfg.Category != tt.category {
				t.Errorf("Category = %d, want %d", cfg.Category, tt.category)
			}
		})
	}
}
