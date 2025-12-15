package food

import (
	"testing"
)

func TestNewFood(t *testing.T) {
	health := []Health{
		{Name: "饱腹感", Amount: 30},
		{Name: "维生素C", Amount: 20},
	}
	f := NewFood("苹果", 5, FoodTypeFresh, health)

	t.Run("名称正确", func(t *testing.T) {
		if f.Name != "苹果" {
			t.Errorf("Name = %q, want 苹果", f.Name)
		}
	})

	t.Run("价格正确", func(t *testing.T) {
		if f.Price != 5 {
			t.Errorf("Price = %d, want 5", f.Price)
		}
	})

	t.Run("类型正确", func(t *testing.T) {
		if f.Type != FoodTypeFresh {
			t.Errorf("Type = %d, want FoodTypeFresh", f.Type)
		}
	})

	t.Run("初始新鲜度100", func(t *testing.T) {
		if f.Freshness != 100 {
			t.Errorf("Freshness = %d, want 100", f.Freshness)
		}
	})

	t.Run("营养效果正确", func(t *testing.T) {
		if len(f.BaseHealth) != 2 {
			t.Errorf("BaseHealth length = %d, want 2", len(f.BaseHealth))
		}
	})
}

func TestUpdateFreshness(t *testing.T) {
	tests := []struct {
		name        string
		foodType    FoodType
		gameMinutes int
		expected    int // 大约值
	}{
		{"新鲜食品1小时", FoodTypeFresh, 60, 96},       // 100 - 4.0 = 96
		{"新鲜食品6小时", FoodTypeFresh, 360, 76},      // 100 - 24 = 76
		{"饮料1小时", FoodTypeBeverage, 60, 98},        // 100 - 2.0 = 98
		{"加工食品1小时", FoodTypeProcessed, 60, 99},   // 100 - 1.0 = 99
		{"罐头1小时", FoodTypeCanned, 60, 99},          // 100 - 0.1 ≈ 99
		{"新鲜食品25小时完全腐烂", FoodTypeFresh, 1500, 0}, // 100 - 100 = 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFood("测试", 10, tt.foodType, nil)
			f.UpdateFreshness(tt.gameMinutes)
			// 允许±2的误差
			if f.Freshness < tt.expected-2 || f.Freshness > tt.expected+2 {
				t.Errorf("Freshness = %d, want ~%d", f.Freshness, tt.expected)
			}
		})
	}
}

func TestGetFreshnessLabel(t *testing.T) {
	tests := []struct {
		freshness int
		expected  string
	}{
		{100, "新鲜"},
		{76, "新鲜"},
		{75, "较新鲜"},
		{51, "较新鲜"},
		{50, "不太新鲜"},
		{26, "不太新鲜"},
		{25, "即将过期"},
		{1, "即将过期"},
		{0, "已过期"},
		{-10, "已过期"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			f := NewFood("测试", 10, FoodTypeFresh, nil)
			f.Freshness = tt.freshness
			if got := f.GetFreshnessLabel(); got != tt.expected {
				t.Errorf("GetFreshnessLabel() = %q, want %q (freshness=%d)",
					got, tt.expected, tt.freshness)
			}
		})
	}
}

func TestGetNutritionMultiplier(t *testing.T) {
	tests := []struct {
		freshness int
		expected  float64
	}{
		{100, 1.0},
		{76, 1.0},
		{75, 0.8},
		{51, 0.8},
		{50, 0.5},
		{26, 0.5},
		{25, 0.2},
		{1, 0.2},
		{0, -0.5},
		{-10, -0.5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f := NewFood("测试", 10, FoodTypeFresh, nil)
			f.Freshness = tt.freshness
			if got := f.GetNutritionMultiplier(); got != tt.expected {
				t.Errorf("GetNutritionMultiplier() = %f, want %f (freshness=%d)",
					got, tt.expected, tt.freshness)
			}
		})
	}
}

func TestGetEffectiveHealth(t *testing.T) {
	health := []Health{
		{Name: "饱腹感", Amount: 100},
		{Name: "维生素C", Amount: 50},
	}

	t.Run("新鲜食物100%效果", func(t *testing.T) {
		f := NewFood("测试", 10, FoodTypeFresh, health)
		f.Freshness = 100
		effective := f.GetEffectiveHealth()

		if len(effective) != 2 {
			t.Fatalf("Expected 2 effects, got %d", len(effective))
		}
		for _, e := range effective {
			if e.Name == "饱腹感" && e.Amount != 100 {
				t.Errorf("饱腹感 Amount = %d, want 100", e.Amount)
			}
			if e.Name == "维生素C" && e.Amount != 50 {
				t.Errorf("维生素C Amount = %d, want 50", e.Amount)
			}
		}
	})

	t.Run("较新鲜食物80%效果", func(t *testing.T) {
		f := NewFood("测试", 10, FoodTypeFresh, health)
		f.Freshness = 60
		effective := f.GetEffectiveHealth()

		for _, e := range effective {
			if e.Name == "饱腹感" && e.Amount != 80 {
				t.Errorf("饱腹感 Amount = %d, want 80", e.Amount)
			}
			if e.Name == "维生素C" && e.Amount != 40 {
				t.Errorf("维生素C Amount = %d, want 40", e.Amount)
			}
		}
	})

	t.Run("过期食物负效果", func(t *testing.T) {
		f := NewFood("测试", 10, FoodTypeFresh, health)
		f.Freshness = -10
		effective := f.GetEffectiveHealth()

		// 过期食物应该返回负面效果
		hasNegative := false
		for _, e := range effective {
			if e.Amount < 0 {
				hasNegative = true
			}
		}
		if !hasNegative {
			t.Error("Expected negative effects for spoiled food")
		}
	})
}

func TestGetFoodPoisoningRisk(t *testing.T) {
	tests := []struct {
		freshness int
		minRisk   float64
		maxRisk   float64
	}{
		{100, 0.0, 0.0},
		{50, 0.0, 0.0},
		{26, 0.0, 0.0},
		{25, 0.1, 0.1},
		{1, 0.1, 0.1},
		{0, 0.3, 0.3},
		{-24, 0.3, 0.3},
		{-25, 0.6, 0.6},
		{-49, 0.6, 0.6},
		{-50, 0.9, 0.9},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f := NewFood("测试", 10, FoodTypeFresh, nil)
			f.Freshness = tt.freshness
			risk := f.GetFoodPoisoningRisk()
			if risk < tt.minRisk || risk > tt.maxRisk {
				t.Errorf("GetFoodPoisoningRisk() = %f, want %f-%f (freshness=%d)",
					risk, tt.minRisk, tt.maxRisk, tt.freshness)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	t.Run("新鲜食物不显示标签", func(t *testing.T) {
		f := NewFood("苹果", 5, FoodTypeFresh, nil)
		info := f.GetInfo()
		if info != "苹果 5元" {
			t.Errorf("GetInfo() = %q, want 苹果 5元", info)
		}
	})

	t.Run("不新鲜食物显示标签", func(t *testing.T) {
		f := NewFood("苹果", 5, FoodTypeFresh, nil)
		f.Freshness = 50
		info := f.GetInfo()
		if info != "苹果 5元 [不太新鲜]" {
			t.Errorf("GetInfo() = %q, want 苹果 5元 [不太新鲜]", info)
		}
	})
}

func TestDecayRates(t *testing.T) {
	tests := []struct {
		foodType FoodType
		expected float64
	}{
		{FoodTypeFresh, 4.0},
		{FoodTypeBeverage, 2.0},
		{FoodTypeProcessed, 1.0},
		{FoodTypeCanned, 0.1},
	}

	for _, tt := range tests {
		t.Run(FoodTypeNames[tt.foodType], func(t *testing.T) {
			if rate := DecayRates[tt.foodType]; rate != tt.expected {
				t.Errorf("DecayRates[%d] = %f, want %f", tt.foodType, rate, tt.expected)
			}
		})
	}
}

func TestFoodTypeNames(t *testing.T) {
	tests := []struct {
		foodType FoodType
		expected string
	}{
		{FoodTypeFresh, "新鲜"},
		{FoodTypeBeverage, "饮料"},
		{FoodTypeProcessed, "加工"},
		{FoodTypeCanned, "罐头"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if name := FoodTypeNames[tt.foodType]; name != tt.expected {
				t.Errorf("FoodTypeNames[%d] = %q, want %q", tt.foodType, name, tt.expected)
			}
		})
	}
}

func TestAllCommodities(t *testing.T) {
	t.Run("商品数量", func(t *testing.T) {
		if len(AllCommodities) < 10 {
			t.Errorf("Expected at least 10 commodities, got %d", len(AllCommodities))
		}
	})

	t.Run("每个商品有效", func(t *testing.T) {
		for i, c := range AllCommodities {
			if c.Name == "" {
				t.Errorf("Commodity %d has empty name", i)
			}
			if c.Price <= 0 {
				t.Errorf("Commodity %q has invalid price %d", c.Name, c.Price)
			}
			if len(c.BaseHealth) == 0 {
				t.Errorf("Commodity %q has no health effects", c.Name)
			}
		}
	})
}
