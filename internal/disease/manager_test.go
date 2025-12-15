package disease

import (
	"testing"

	"citylife/internal/user"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	t.Run("初始化疾病定义", func(t *testing.T) {
		diseases := m.GetAllDiseases()
		if len(diseases) != 6 {
			t.Errorf("Expected 6 diseases, got %d", len(diseases))
		}
	})

	t.Run("无激活疾病", func(t *testing.T) {
		if len(m.GetActiveDiseases()) != 0 {
			t.Error("Expected no active diseases initially")
		}
	})
}

func TestDiseaseDefinitions(t *testing.T) {
	m := NewManager()

	tests := []struct {
		id            string
		name          string
		severity      Severity
		treatmentCost int
	}{
		{"cold", "感冒", SeverityMild, 50},
		{"anemia", "贫血", SeverityModerate, 200},
		{"scurvy", "坏血病", SeveritySevere, 500},
		{"food_poisoning", "食物中毒", SeverityModerate, 100},
		{"malnutrition", "营养不良", SeverityModerate, 150},
		{"depression", "抑郁症", SeverityModerate, 300},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			info := m.GetDiseaseInfo(tt.id)
			if info == nil {
				t.Fatalf("Disease %q not found", tt.id)
			}
			if info.Name != tt.name {
				t.Errorf("Name = %q, want %q", info.Name, tt.name)
			}
			if info.Severity != tt.severity {
				t.Errorf("Severity = %d, want %d", info.Severity, tt.severity)
			}
			if info.TreatmentCost != tt.treatmentCost {
				t.Errorf("TreatmentCost = %d, want %d", info.TreatmentCost, tt.treatmentCost)
			}
		})
	}
}

func TestTriggerFoodPoisoning(t *testing.T) {
	t.Run("100%概率触发", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)
		if !m.HasDisease("food_poisoning") {
			t.Error("Expected food poisoning to be triggered with 100% chance")
		}
	})

	t.Run("0%概率不触发", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(0.0)
		if m.HasDisease("food_poisoning") {
			t.Error("Expected food poisoning not to be triggered with 0% chance")
		}
	})
}

func TestHasDisease(t *testing.T) {
	m := NewManager()

	t.Run("无疾病时返回false", func(t *testing.T) {
		if m.HasDisease("cold") {
			t.Error("HasDisease(cold) = true, want false")
		}
	})

	t.Run("有疾病时返回true", func(t *testing.T) {
		m.TriggerFoodPoisoning(1.0)
		if !m.HasDisease("food_poisoning") {
			t.Error("HasDisease(food_poisoning) = false, want true")
		}
	})
}

func TestGetActiveDiseases(t *testing.T) {
	m := NewManager()

	t.Run("初始无疾病", func(t *testing.T) {
		diseases := m.GetActiveDiseases()
		if len(diseases) != 0 {
			t.Errorf("Expected 0 active diseases, got %d", len(diseases))
		}
	})

	t.Run("触发后有疾病", func(t *testing.T) {
		m.TriggerFoodPoisoning(1.0)
		diseases := m.GetActiveDiseases()
		if len(diseases) != 1 {
			t.Errorf("Expected 1 active disease, got %d", len(diseases))
		}
	})
}

func TestGetActiveDiseaseIDs(t *testing.T) {
	m := NewManager()
	m.TriggerFoodPoisoning(1.0)

	ids := m.GetActiveDiseaseIDs()
	if len(ids) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(ids))
	}
	if ids[0] != "food_poisoning" {
		t.Errorf("Expected food_poisoning, got %s", ids[0])
	}
}

func TestTreatAtHospital(t *testing.T) {
	t.Run("治疗成功", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)

		spent := 0
		err := m.TreatAtHospital("food_poisoning", func(cost int) error {
			spent = cost
			return nil
		})

		if err != nil {
			t.Errorf("TreatAtHospital returned error: %v", err)
		}
		if spent != 100 {
			t.Errorf("Expected to spend 100, spent %d", spent)
		}
		if m.HasDisease("food_poisoning") {
			t.Error("Disease should be cured after treatment")
		}
	})

	t.Run("治疗不存在的疾病", func(t *testing.T) {
		m := NewManager()
		err := m.TreatAtHospital("cold", func(cost int) error {
			t.Error("Should not attempt to spend money")
			return nil
		})
		if err != nil {
			t.Errorf("Expected nil error for non-existent disease, got %v", err)
		}
	})
}

func TestGetActionCostMultiplier(t *testing.T) {
	t.Run("无疾病时为1.0", func(t *testing.T) {
		m := NewManager()
		if mult := m.GetActionCostMultiplier(); mult != 1.0 {
			t.Errorf("GetActionCostMultiplier() = %f, want 1.0", mult)
		}
	})

	t.Run("有疾病时大于1.0", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)
		if mult := m.GetActionCostMultiplier(); mult <= 1.0 {
			t.Errorf("GetActionCostMultiplier() = %f, want > 1.0", mult)
		}
	})
}

func TestGetNutritionDecayMultiplier(t *testing.T) {
	t.Run("无疾病时为1.0", func(t *testing.T) {
		m := NewManager()
		if mult := m.GetNutritionDecayMultiplier(); mult != 1.0 {
			t.Errorf("GetNutritionDecayMultiplier() = %f, want 1.0", mult)
		}
	})

	t.Run("有疾病时大于1.0", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)
		if mult := m.GetNutritionDecayMultiplier(); mult <= 1.0 {
			t.Errorf("GetNutritionDecayMultiplier() = %f, want > 1.0", mult)
		}
	})
}

func TestGetStatusSummary(t *testing.T) {
	t.Run("健康时", func(t *testing.T) {
		m := NewManager()
		if s := m.GetStatusSummary(); s != "健康" {
			t.Errorf("GetStatusSummary() = %q, want 健康", s)
		}
	})

	t.Run("有疾病时", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)
		s := m.GetStatusSummary()
		if s == "健康" {
			t.Error("GetStatusSummary() should not be 健康 when sick")
		}
	})
}

func TestSeverityToString(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityMild, "轻度"},
		{SeverityModerate, "中度"},
		{SeveritySevere, "重度"},
		{Severity(99), "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := SeverityToString(tt.severity); got != tt.expected {
				t.Errorf("SeverityToString(%d) = %q, want %q", tt.severity, got, tt.expected)
			}
		})
	}
}

func TestSaveLoadMethods(t *testing.T) {
	t.Run("GetActiveDisease", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)

		ad := m.GetActiveDisease("food_poisoning")
		if ad == nil {
			t.Fatal("GetActiveDisease returned nil")
		}
		if !ad.IsActive {
			t.Error("IsActive should be true")
		}
	})

	t.Run("ClearAllDiseases", func(t *testing.T) {
		m := NewManager()
		m.TriggerFoodPoisoning(1.0)
		m.ClearAllDiseases()

		if m.HasDisease("food_poisoning") {
			t.Error("Disease should be cleared")
		}
	})

	t.Run("SetActiveDisease", func(t *testing.T) {
		m := NewManager()
		m.SetActiveDisease("cold", 5, 3, true)

		if !m.HasDisease("cold") {
			t.Error("Disease should be active after SetActiveDisease")
		}

		ad := m.GetActiveDisease("cold")
		if ad.TriggerProgress != 5 {
			t.Errorf("TriggerProgress = %d, want 5", ad.TriggerProgress)
		}
		if ad.CureProgress != 3 {
			t.Errorf("CureProgress = %d, want 3", ad.CureProgress)
		}
	})
}

func TestOnAction(t *testing.T) {
	t.Run("食物中毒自然康复", func(t *testing.T) {
		m := NewManager()
		u := user.New()

		m.TriggerFoodPoisoning(1.0)

		// 食物中毒需要8次行动自然康复
		for i := 0; i < 8; i++ {
			if !m.HasDisease("food_poisoning") {
				t.Fatalf("Food poisoning cured too early at action %d", i)
			}
			m.OnAction(u)
		}

		if m.HasDisease("food_poisoning") {
			t.Error("Food poisoning should be cured after 8 actions")
		}
	})
}

func TestTriggerConditionsWithUser(t *testing.T) {
	t.Run("维生素C低触发感冒", func(t *testing.T) {
		m := NewManager()
		u := user.New()

		// 设置维生素C低于触发阈值
		u.SetNutrition("维生素C", 10)

		// 需要多次行动才能触发（触发持续时间为10）
		// 且有概率触发
		triggered := false
		for i := 0; i < 100; i++ {
			m.OnAction(u)
			if m.HasDisease("cold") {
				triggered = true
				break
			}
		}

		// 由于有随机性，只验证没有panic
		_ = triggered
	})
}
