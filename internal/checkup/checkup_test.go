package checkup

import (
	"strings"
	"testing"

	"citylife/internal/user"
)

func TestNewService(t *testing.T) {
	s := NewService()

	t.Run("服务非空", func(t *testing.T) {
		if s == nil {
			t.Error("NewService() should not return nil")
		}
	})

	t.Run("项目数量", func(t *testing.T) {
		items := s.GetAllItems()
		// 5个套餐 + 21个单项 = 26
		if len(items) < 20 {
			t.Errorf("Should have at least 20 items, got %d", len(items))
		}
	})
}

func TestGetPackageItems(t *testing.T) {
	s := NewService()
	packages := s.GetPackageItems()

	t.Run("5个套餐", func(t *testing.T) {
		if len(packages) != 5 {
			t.Errorf("Should have 5 packages, got %d", len(packages))
		}
	})

	t.Run("套餐ID正确", func(t *testing.T) {
		expectedIDs := []string{"p_core", "p_micro", "p_vitamin", "p_mental", "p_full"}
		for i, pkg := range packages {
			if pkg.ID != expectedIDs[i] {
				t.Errorf("Package %d ID = %s, want %s", i, pkg.ID, expectedIDs[i])
			}
		}
	})

	t.Run("套餐价格正确", func(t *testing.T) {
		expectedPrices := []int{50, 80, 60, 50, 150}
		for i, pkg := range packages {
			if pkg.Price != expectedPrices[i] {
				t.Errorf("Package %s price = %d, want %d", pkg.ID, pkg.Price, expectedPrices[i])
			}
		}
	})
}

func TestGetSingleItems(t *testing.T) {
	s := NewService()
	singles := s.GetSingleItems()

	t.Run("单项数量", func(t *testing.T) {
		if len(singles) < 17 {
			t.Errorf("Should have at least 17 single items, got %d", len(singles))
		}
	})

	t.Run("单项类型正确", func(t *testing.T) {
		for _, item := range singles {
			if item.Type != TypeSingle {
				t.Errorf("Item %s should be TypeSingle", item.ID)
			}
		}
	})
}

func TestGetItem(t *testing.T) {
	s := NewService()

	t.Run("存在的项目", func(t *testing.T) {
		item := s.GetItem("p_core")
		if item == nil {
			t.Error("GetItem('p_core') should not return nil")
		}
		if item.Name != "核心指标套餐" {
			t.Errorf("Item name = %s, want 核心指标套餐", item.Name)
		}
	})

	t.Run("不存在的项目", func(t *testing.T) {
		item := s.GetItem("nonexistent")
		if item != nil {
			t.Error("GetItem for nonexistent ID should return nil")
		}
	})
}

func TestPerformCheckup(t *testing.T) {
	s := NewService()
	u := user.New()

	t.Run("核心套餐检查", func(t *testing.T) {
		report := s.PerformCheckup("p_core", u)

		if report == nil {
			t.Fatal("Report should not be nil")
		}

		if report.ItemName != "核心指标套餐" {
			t.Errorf("ItemName = %s, want 核心指标套餐", report.ItemName)
		}

		if report.TotalCost != 50 {
			t.Errorf("TotalCost = %d, want 50", report.TotalCost)
		}

		if len(report.Results) != 4 {
			t.Errorf("Should have 4 results, got %d", len(report.Results))
		}
	})

	t.Run("健康用户全部正常", func(t *testing.T) {
		report := s.PerformCheckup("p_core", u)

		for _, result := range report.Results {
			if result.Status != "正常" {
				t.Errorf("%s status = %s, want 正常", result.AttributeName, result.Status)
			}
			if result.Advice != "" {
				t.Errorf("%s should have no advice when normal", result.AttributeName)
			}
		}

		if !strings.Contains(report.OverallAdvice, "良好") {
			t.Errorf("OverallAdvice should mention 良好: %s", report.OverallAdvice)
		}
	})

	t.Run("低营养显示警告", func(t *testing.T) {
		u := user.New()
		u.SetNutrition("饱腹感", 40) // 警告区 (20-50)

		report := s.PerformCheckup("p_core", u)

		found := false
		for _, result := range report.Results {
			if result.AttributeName == "饱腹感" {
				found = true
				if result.Status != "偏低" {
					t.Errorf("饱腹感 status = %s, want 偏低", result.Status)
				}
				if result.Advice == "" {
					t.Error("Should have advice for low nutrition")
				}
			}
		}
		if !found {
			t.Error("Should have 饱腹感 result")
		}
	})

	t.Run("危险状态", func(t *testing.T) {
		u := user.New()
		u.SetNutrition("蛋白质", 10) // 危险区 (0-20)

		report := s.PerformCheckup("p_core", u)

		for _, result := range report.Results {
			if result.AttributeName == "蛋白质" {
				if result.Status != "危险" {
					t.Errorf("蛋白质 status = %s, want 危险", result.Status)
				}
				if !strings.Contains(result.Advice, "【注意】") {
					t.Error("Danger advice should have 【注意】 prefix")
				}
			}
		}
	})

	t.Run("濒死状态", func(t *testing.T) {
		u := user.New()
		u.SetNutrition("饥渴", -15) // 濒死区

		report := s.PerformCheckup("p_core", u)

		for _, result := range report.Results {
			if result.AttributeName == "饥渴" {
				if result.Status != "濒死" {
					t.Errorf("饥渴 status = %s, want 濒死", result.Status)
				}
				if !strings.Contains(result.Advice, "【紧急】") {
					t.Error("Critical advice should have 【紧急】 prefix")
				}
			}
		}

		if !strings.Contains(report.OverallAdvice, "危急") {
			t.Errorf("OverallAdvice should mention 危急: %s", report.OverallAdvice)
		}
	})

	t.Run("无效项目ID", func(t *testing.T) {
		report := s.PerformCheckup("invalid", u)
		if report != nil {
			t.Error("Invalid item ID should return nil report")
		}
	})
}

func TestFormatReport(t *testing.T) {
	s := NewService()
	u := user.New()
	report := s.PerformCheckup("p_core", u)
	formatted := s.FormatReport(report)

	t.Run("包含标题", func(t *testing.T) {
		if !strings.Contains(formatted, "城市生活医院体检报告") {
			t.Error("Should contain report title")
		}
	})

	t.Run("包含项目名称", func(t *testing.T) {
		if !strings.Contains(formatted, "核心指标套餐") {
			t.Error("Should contain item name")
		}
	})

	t.Run("包含费用", func(t *testing.T) {
		if !strings.Contains(formatted, "¥50") {
			t.Error("Should contain cost")
		}
	})

	t.Run("包含分类标题", func(t *testing.T) {
		if !strings.Contains(formatted, "核心指标") {
			t.Error("Should contain category header")
		}
	})

	t.Run("包含综合评估", func(t *testing.T) {
		if !strings.Contains(formatted, "综合评估") {
			t.Error("Should contain overall assessment section")
		}
	})

	t.Run("包含医嘱", func(t *testing.T) {
		if !strings.Contains(formatted, "医嘱") {
			t.Error("Should contain doctor's advice")
		}
	})
}

func TestGetStatusString(t *testing.T) {
	tests := []struct {
		level  int
		status string
	}{
		{user.LevelSafe, "正常"},
		{user.LevelWarning, "偏低"},
		{user.LevelDanger, "危险"},
		{user.LevelOverdraft, "透支"},
		{user.LevelCritical, "濒死"},
		{99, "未知"},
	}

	for _, tt := range tests {
		result := getStatusString(tt.level)
		if result != tt.status {
			t.Errorf("getStatusString(%d) = %s, want %s", tt.level, result, tt.status)
		}
	}
}

func TestGetCategoryName(t *testing.T) {
	tests := []struct {
		attr     string
		category string
	}{
		{"饱腹感", "核心指标"},
		{"饥渴", "核心指标"},
		{"维生素A", "维生素"},
		{"维生素C", "维生素"},
		{"幸福感", "心理评估"},
		{"钙", "微量元素"},
		{"铁", "微量元素"},
	}

	for _, tt := range tests {
		result := getCategoryName(tt.attr)
		if result != tt.category {
			t.Errorf("getCategoryName(%s) = %s, want %s", tt.attr, result, tt.category)
		}
	}
}

func TestFullBodyCheckup(t *testing.T) {
	s := NewService()
	u := user.New()

	report := s.PerformCheckup("p_full", u)

	t.Run("包含所有属性", func(t *testing.T) {
		if len(report.Results) < 20 {
			t.Errorf("Full checkup should have at least 20 results, got %d", len(report.Results))
		}
	})

	t.Run("费用正确", func(t *testing.T) {
		if report.TotalCost != 150 {
			t.Errorf("Full checkup cost = %d, want 150", report.TotalCost)
		}
	})
}
