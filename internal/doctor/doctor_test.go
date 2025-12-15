package doctor

import (
	"strings"
	"testing"

	"citylife/internal/disease"
	"citylife/internal/user"
)

func TestDiagnose(t *testing.T) {
	t.Run("健康用户", func(t *testing.T) {
		u := user.New()
		dm := disease.NewManager()

		report := Diagnose(u, dm)

		if report == nil {
			t.Fatal("Report should not be nil")
		}

		if len(report.Results) != 0 {
			t.Errorf("Healthy user should have 0 results, got %d", len(report.Results))
		}

		if report.TotalCost != DiagnosisCost {
			t.Errorf("TotalCost = %d, want %d", report.TotalCost, DiagnosisCost)
		}

		if !strings.Contains(report.OverallStatus, "健康") {
			t.Errorf("OverallStatus should mention 健康: %s", report.OverallStatus)
		}
	})

	t.Run("患病用户", func(t *testing.T) {
		u := user.New()
		dm := disease.NewManager()

		// 手动触发食物中毒
		dm.TriggerFoodPoisoning(1.0) // 100%触发

		report := Diagnose(u, dm)

		if len(report.Results) != 1 {
			t.Errorf("Should have 1 result, got %d", len(report.Results))
		}

		if report.Results[0].DiseaseID != "food_poisoning" {
			t.Errorf("DiseaseID = %s, want food_poisoning", report.Results[0].DiseaseID)
		}

		if report.Results[0].DiseaseName != "食物中毒" {
			t.Errorf("DiseaseName = %s, want 食物中毒", report.Results[0].DiseaseName)
		}

		// 检查诊断记录已添加
		if !u.HasDiagnosis("food_poisoning") {
			t.Error("User should have food_poisoning diagnosis")
		}
	})

	t.Run("多疾病诊断", func(t *testing.T) {
		u := user.New()
		dm := disease.NewManager()

		// 触发多个疾病
		dm.TriggerFoodPoisoning(1.0)
		dm.SetActiveDisease("cold", 0, 0, true)

		report := Diagnose(u, dm)

		if len(report.Results) != 2 {
			t.Errorf("Should have 2 results, got %d", len(report.Results))
		}

		if !strings.Contains(report.OverallStatus, "多种疾病") {
			t.Errorf("OverallStatus should mention 多种疾病: %s", report.OverallStatus)
		}

		// 检查所有诊断记录
		if !u.HasDiagnosis("food_poisoning") {
			t.Error("Should have food_poisoning diagnosis")
		}
		if !u.HasDiagnosis("cold") {
			t.Error("Should have cold diagnosis")
		}
	})
}

func TestFormatReport(t *testing.T) {
	u := user.New()
	dm := disease.NewManager()
	dm.TriggerFoodPoisoning(1.0)

	report := Diagnose(u, dm)
	formatted := FormatReport(report)

	t.Run("包含标题", func(t *testing.T) {
		if !strings.Contains(formatted, "诊断报告") {
			t.Error("Should contain report title")
		}
	})

	t.Run("包含费用", func(t *testing.T) {
		if !strings.Contains(formatted, "¥50") {
			t.Error("Should contain cost")
		}
	})

	t.Run("包含疾病信息", func(t *testing.T) {
		if !strings.Contains(formatted, "食物中毒") {
			t.Error("Should contain disease name")
		}
	})

	t.Run("包含建议", func(t *testing.T) {
		if !strings.Contains(formatted, "建议") {
			t.Error("Should contain advice")
		}
	})

	t.Run("包含处方药提示", func(t *testing.T) {
		if !strings.Contains(formatted, "处方药") {
			t.Error("Should mention prescription drugs")
		}
	})
}

func TestFormatReportHealthy(t *testing.T) {
	u := user.New()
	dm := disease.NewManager()

	report := Diagnose(u, dm)
	formatted := FormatReport(report)

	if !strings.Contains(formatted, "身体健康") {
		t.Error("Should indicate healthy status")
	}

	if strings.Contains(formatted, "处方药") {
		t.Error("Should not mention prescription for healthy user")
	}
}

func TestGetDiseaseName(t *testing.T) {
	tests := []struct {
		id   string
		name string
	}{
		{"cold", "感冒"},
		{"anemia", "贫血"},
		{"scurvy", "坏血病"},
		{"food_poisoning", "食物中毒"},
		{"malnutrition", "营养不良"},
		{"depression", "抑郁症"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := GetDiseaseName(tt.id)
		if result != tt.name {
			t.Errorf("GetDiseaseName(%s) = %s, want %s", tt.id, result, tt.name)
		}
	}
}

func TestGetCost(t *testing.T) {
	if GetCost() != 50 {
		t.Errorf("GetCost() = %d, want 50", GetCost())
	}
}

func TestHasDiagnosedDisease(t *testing.T) {
	u := user.New()

	t.Run("未诊断", func(t *testing.T) {
		if HasDiagnosedDisease(u, "cold") {
			t.Error("Should not have cold diagnosis initially")
		}
	})

	t.Run("已诊断", func(t *testing.T) {
		u.AddDiagnosis("cold")
		if !HasDiagnosedDisease(u, "cold") {
			t.Error("Should have cold diagnosis after adding")
		}
	})
}

func TestGetUserDiagnoses(t *testing.T) {
	u := user.New()

	t.Run("无诊断", func(t *testing.T) {
		diagnoses := GetUserDiagnoses(u)
		if len(diagnoses) != 0 {
			t.Errorf("Should have 0 diagnoses, got %d", len(diagnoses))
		}
	})

	t.Run("有诊断", func(t *testing.T) {
		u.AddDiagnosis("cold")
		u.AddDiagnosis("anemia")

		diagnoses := GetUserDiagnoses(u)
		if len(diagnoses) != 2 {
			t.Errorf("Should have 2 diagnoses, got %d", len(diagnoses))
		}
	})
}

func TestClearDiagnosis(t *testing.T) {
	u := user.New()
	u.AddDiagnosis("cold")

	if !u.HasDiagnosis("cold") {
		t.Fatal("Should have cold diagnosis before clearing")
	}

	ClearDiagnosis(u, "cold")

	if u.HasDiagnosis("cold") {
		t.Error("Should not have cold diagnosis after clearing")
	}
}

func TestDiagnosisResultFields(t *testing.T) {
	u := user.New()
	dm := disease.NewManager()
	dm.SetActiveDisease("cold", 0, 0, true)

	report := Diagnose(u, dm)

	if len(report.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(report.Results))
	}

	result := report.Results[0]

	t.Run("DiseaseID", func(t *testing.T) {
		if result.DiseaseID != "cold" {
			t.Errorf("DiseaseID = %s, want cold", result.DiseaseID)
		}
	})

	t.Run("DiseaseName", func(t *testing.T) {
		if result.DiseaseName != "感冒" {
			t.Errorf("DiseaseName = %s, want 感冒", result.DiseaseName)
		}
	})

	t.Run("Severity", func(t *testing.T) {
		if result.Severity == "" {
			t.Error("Severity should not be empty")
		}
	})

	t.Run("Description", func(t *testing.T) {
		if result.Description == "" {
			t.Error("Description should not be empty")
		}
	})

	t.Run("Advice", func(t *testing.T) {
		if result.Advice == "" {
			t.Error("Advice should not be empty")
		}
		if !strings.Contains(result.Advice, "维生素C") {
			t.Errorf("Cold advice should mention 维生素C: %s", result.Advice)
		}
	})
}
