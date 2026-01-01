package user

import (
	"testing"
)

func TestNew(t *testing.T) {
	u := New()

	t.Run("营养值初始化为100", func(t *testing.T) {
		// 检查核心营养
		coreNutrients := []string{"饱腹感", "饥渴", "蛋白质", "碳水化合物"}
		for _, name := range coreNutrients {
			if got := u.GetNutrition(name); got != 100 {
				t.Errorf("GetNutrition(%q) = %d, want 100", name, got)
			}
		}
	})

	t.Run("微量元素初始化为100", func(t *testing.T) {
		microNutrients := []string{"钙", "铁", "锌", "维生素C", "维生素B"}
		for _, name := range microNutrients {
			if got := u.GetNutrition(name); got != 100 {
				t.Errorf("GetNutrition(%q) = %d, want 100", name, got)
			}
		}
	})

	t.Run("诊断记录为空", func(t *testing.T) {
		if len(u.Diagnoses) != 0 {
			t.Errorf("Diagnoses should be empty, got %d items", len(u.Diagnoses))
		}
	})
}

func TestGetSetNutrition(t *testing.T) {
	u := New()

	t.Run("获取已存在的营养值", func(t *testing.T) {
		if got := u.GetNutrition("饱腹感"); got != 100 {
			t.Errorf("GetNutrition(饱腹感) = %d, want 100", got)
		}
	})

	t.Run("获取不存在的营养值返回0", func(t *testing.T) {
		if got := u.GetNutrition("不存在的营养"); got != 0 {
			t.Errorf("GetNutrition(不存在的营养) = %d, want 0", got)
		}
	})

	t.Run("设置营养值", func(t *testing.T) {
		u.SetNutrition("饱腹感", 50)
		if got := u.GetNutrition("饱腹感"); got != 50 {
			t.Errorf("GetNutrition(饱腹感) = %d, want 50", got)
		}
	})
}

func TestNutritionBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		nutrition string
		setValue  int
		expected  int
	}{
		{"上限100", "饱腹感", 150, 100},
		{"正常值", "饱腹感", 50, 50},
		{"下限-30(饱腹感)", "饱腹感", -50, -30},
		{"下限-20(饥渴)", "饥渴", -50, -20},
		{"下限-15(蛋白质)", "蛋白质", -50, -15},
		{"下限-15(碳水化合物)", "碳水化合物", -50, -15},
		{"下限-10(钙)", "钙", -50, -10},
		{"下限0(脂肪)", "脂肪", -50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New()
			u.SetNutrition(tt.nutrition, tt.setValue)
			if got := u.GetNutrition(tt.nutrition); got != tt.expected {
				t.Errorf("SetNutrition(%q, %d) resulted in %d, want %d",
					tt.nutrition, tt.setValue, got, tt.expected)
			}
		})
	}
}

func TestConsumeNutrition(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		consume  int
		expected int
	}{
		{"正常消耗", 100, 30, 70},
		{"消耗到0", 100, 100, 0},
		{"消耗到负值", 50, 80, -30}, // 饱腹感下限-30
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New()
			u.SetNutrition("饱腹感", tt.initial)
			u.ConsumeNutrition("饱腹感", tt.consume)
			if got := u.GetNutrition("饱腹感"); got != tt.expected {
				t.Errorf("ConsumeNutrition: got %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestAddNutrition(t *testing.T) {
	tests := []struct {
		name     string
		initial  int
		add      int
		expected int
	}{
		{"正常增加", 50, 30, 80},
		{"增加到上限", 80, 50, 100},
		{"从负值增加", -10, 50, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New()
			u.SetNutrition("饱腹感", tt.initial)
			u.AddNutrition("饱腹感", tt.add)
			if got := u.GetNutrition("饱腹感"); got != tt.expected {
				t.Errorf("AddNutrition: got %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestIsAlive(t *testing.T) {
	tests := []struct {
		name      string
		nutrition string
		value     int
		alive     bool
	}{
		{"饱腹感安全", "饱腹感", 50, true},
		{"饱腹感临界", "饱腹感", -29, true},
		{"饱腹感死亡", "饱腹感", -30, false},
		{"饥渴安全", "饥渴", 50, true},
		{"饥渴临界", "饥渴", -19, true},
		{"饥渴死亡", "饥渴", -20, false},
		{"蛋白质安全", "蛋白质", 50, true},
		{"蛋白质临界", "蛋白质", -14, true},
		{"蛋白质死亡", "蛋白质", -15, false},
		{"碳水化合物安全", "碳水化合物", 50, true},
		{"碳水化合物临界", "碳水化合物", -14, true},
		{"碳水化合物死亡", "碳水化合物", -15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New()
			u.SetNutrition(tt.nutrition, tt.value)
			if got := u.IsAlive(); got != tt.alive {
				t.Errorf("IsAlive() = %v, want %v (nutrition=%s, value=%d)",
					got, tt.alive, tt.nutrition, tt.value)
			}
		})
	}
}

func TestGetDeathRisk(t *testing.T) {
	t.Run("健康状态无风险", func(t *testing.T) {
		u := New()
		if risk := u.GetDeathRisk(); risk != 0.0 {
			t.Errorf("GetDeathRisk() = %f, want 0.0", risk)
		}
	})

	t.Run("饱腹感濒死区有风险", func(t *testing.T) {
		u := New()
		u.SetNutrition("饱腹感", -15)
		if risk := u.GetDeathRisk(); risk <= 0.0 {
			t.Errorf("GetDeathRisk() = %f, want > 0.0", risk)
		}
	})

	t.Run("多个营养濒死累加风险", func(t *testing.T) {
		u := New()
		u.SetNutrition("饱腹感", -15)
		risk1 := u.GetDeathRisk()

		u.SetNutrition("饥渴", -15)
		risk2 := u.GetDeathRisk()

		if risk2 <= risk1 {
			t.Errorf("Multiple critical nutrients should increase risk: %f <= %f", risk2, risk1)
		}
	})
}

func TestGetNutritionLevel(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"安全区高", 100, LevelSafe},
		{"安全区低", 51, LevelSafe},
		{"警告区高", 50, LevelWarning},
		{"警告区低", 21, LevelWarning},
		{"危险区高", 20, LevelDanger},
		{"危险区低", 1, LevelDanger},
		{"透支区高", 0, LevelOverdraft},
		{"透支区低", -9, LevelOverdraft},
		{"濒死区", -10, LevelCritical},
		{"濒死区深", -20, LevelCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New()
			u.SetNutrition("饱腹感", tt.value)
			if got := u.GetNutritionLevel("饱腹感"); got != tt.expected {
				t.Errorf("GetNutritionLevel() = %d, want %d (value=%d)",
					got, tt.expected, tt.value)
			}
		})
	}
}

func TestDiagnosis(t *testing.T) {
	u := New()

	t.Run("添加诊断", func(t *testing.T) {
		u.AddDiagnosis("cold")
		if !u.HasDiagnosis("cold") {
			t.Error("HasDiagnosis(cold) = false, want true")
		}
	})

	t.Run("检查不存在的诊断", func(t *testing.T) {
		if u.HasDiagnosis("anemia") {
			t.Error("HasDiagnosis(anemia) = true, want false")
		}
	})

	t.Run("清除诊断", func(t *testing.T) {
		u.ClearDiagnosis("cold")
		if u.HasDiagnosis("cold") {
			t.Error("HasDiagnosis(cold) = true after clear, want false")
		}
	})

	t.Run("获取所有诊断", func(t *testing.T) {
		u.AddDiagnosis("cold")
		u.AddDiagnosis("anemia")
		diagnoses := u.GetDiagnoses()
		if len(diagnoses) != 2 {
			t.Errorf("GetDiagnoses() returned %d items, want 2", len(diagnoses))
		}
	})
}

func TestGetAllNutrition(t *testing.T) {
	u := New()
	u.SetNutrition("饱腹感", 50)
	u.SetNutrition("饥渴", 60)

	t.Run("包含修改后的值", func(t *testing.T) {
		all := u.GetAllNutrition()
		if all["饱腹感"] != 50 {
			t.Errorf("all[饱腹感] = %d, want 50", all["饱腹感"])
		}
		if all["饥渴"] != 60 {
			t.Errorf("all[饥渴] = %d, want 60", all["饥渴"])
		}
	})

	t.Run("返回副本", func(t *testing.T) {
		all := u.GetAllNutrition()
		all["饱腹感"] = 999
		if u.GetNutrition("饱腹感") == 999 {
			t.Error("GetAllNutrition should return a copy, not a reference")
		}
	})
}

func TestSetAllNutrition(t *testing.T) {
	u := New()

	newValues := map[string]int{
		"饱腹感": 50,
		"饥渴":  60,
		"蛋白质": 70,
	}

	u.SetAllNutrition(newValues)

	for name, expected := range newValues {
		if got := u.GetNutrition(name); got != expected {
			t.Errorf("GetNutrition(%q) = %d, want %d", name, got, expected)
		}
	}
}

func TestGetMinValue(t *testing.T) {
	tests := []struct {
		nutrition string
		expected  int
	}{
		{"饱腹感", -30},
		{"饥渴", -20},
		{"蛋白质", -15},
		{"碳水化合物", -15},
		{"钙", -10},
		{"脂肪", 0},
		{"不存在", -10}, // 默认值
	}

	for _, tt := range tests {
		t.Run(tt.nutrition, func(t *testing.T) {
			if got := GetMinValue(tt.nutrition); got != tt.expected {
				t.Errorf("GetMinValue(%q) = %d, want %d", tt.nutrition, got, tt.expected)
			}
		})
	}
}
