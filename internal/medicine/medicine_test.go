package medicine

import (
	"strings"
	"testing"

	"citylife/internal/user"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	t.Run("管理器非空", func(t *testing.T) {
		if m == nil {
			t.Error("NewManager() should not return nil")
		}
	})

	t.Run("药品数量", func(t *testing.T) {
		all := m.GetAllMedicines()
		// 10 OTC + 6 处方 = 16
		if len(all) != 16 {
			t.Errorf("Should have 16 medicines, got %d", len(all))
		}
	})
}

func TestGetOTCMedicines(t *testing.T) {
	m := NewManager()
	otc := m.GetOTCMedicines()

	t.Run("10种OTC药品", func(t *testing.T) {
		if len(otc) != 10 {
			t.Errorf("Should have 10 OTC medicines, got %d", len(otc))
		}
	})

	t.Run("所有都是OTC类型", func(t *testing.T) {
		for _, med := range otc {
			if med.Type != TypeOTC {
				t.Errorf("%s should be TypeOTC", med.Name)
			}
		}
	})

	t.Run("包含复合维生素片", func(t *testing.T) {
		found := false
		for _, med := range otc {
			if med.Name == "复合维生素片" {
				found = true
				if med.Price != 30 {
					t.Errorf("复合维生素片 price = %d, want 30", med.Price)
				}
				break
			}
		}
		if !found {
			t.Error("Should contain 复合维生素片")
		}
	})
}

func TestGetPrescriptionMedicines(t *testing.T) {
	m := NewManager()
	prescriptions := m.GetPrescriptionMedicines()

	t.Run("6种处方药", func(t *testing.T) {
		if len(prescriptions) != 6 {
			t.Errorf("Should have 6 prescription medicines, got %d", len(prescriptions))
		}
	})

	t.Run("所有都是处方药类型", func(t *testing.T) {
		for _, med := range prescriptions {
			if med.Type != TypePrescription {
				t.Errorf("%s should be TypePrescription", med.Name)
			}
		}
	})

	t.Run("都有RequiredDisease", func(t *testing.T) {
		for _, med := range prescriptions {
			if med.RequiredDisease == "" {
				t.Errorf("%s should have RequiredDisease", med.Name)
			}
		}
	})
}

func TestGetAvailablePrescriptions(t *testing.T) {
	m := NewManager()
	u := user.New()

	t.Run("无诊断时无可用处方药", func(t *testing.T) {
		available := m.GetAvailablePrescriptions(u)
		if len(available) != 0 {
			t.Errorf("Should have 0 available, got %d", len(available))
		}
	})

	t.Run("有诊断时可购买对应处方药", func(t *testing.T) {
		u.AddDiagnosis("cold")
		available := m.GetAvailablePrescriptions(u)

		if len(available) != 1 {
			t.Errorf("Should have 1 available, got %d", len(available))
		}

		if available[0].RequiredDisease != "cold" {
			t.Errorf("Available medicine should require cold diagnosis")
		}
	})

	t.Run("多个诊断时可购买多种处方药", func(t *testing.T) {
		u.AddDiagnosis("anemia")
		available := m.GetAvailablePrescriptions(u)

		if len(available) != 2 {
			t.Errorf("Should have 2 available, got %d", len(available))
		}
	})
}

func TestGetMedicine(t *testing.T) {
	m := NewManager()

	t.Run("存在的药品", func(t *testing.T) {
		med := m.GetMedicine("multivitamin")
		if med == nil {
			t.Fatal("Should find multivitamin")
		}
		if med.Name != "复合维生素片" {
			t.Errorf("Name = %s, want 复合维生素片", med.Name)
		}
	})

	t.Run("不存在的药品", func(t *testing.T) {
		med := m.GetMedicine("nonexistent")
		if med != nil {
			t.Error("Should return nil for nonexistent medicine")
		}
	})
}

func TestCanPurchase(t *testing.T) {
	m := NewManager()
	u := user.New()

	t.Run("OTC药品始终可购买", func(t *testing.T) {
		med := m.GetMedicine("multivitamin")
		if !m.CanPurchase(med, u) {
			t.Error("OTC medicine should always be purchasable")
		}
	})

	t.Run("处方药无诊断不可购买", func(t *testing.T) {
		med := m.GetMedicine("cold_medicine")
		if m.CanPurchase(med, u) {
			t.Error("Prescription without diagnosis should not be purchasable")
		}
	})

	t.Run("处方药有诊断可购买", func(t *testing.T) {
		u.AddDiagnosis("cold")
		med := m.GetMedicine("cold_medicine")
		if !m.CanPurchase(med, u) {
			t.Error("Prescription with diagnosis should be purchasable")
		}
	})

	t.Run("处方药诊断不匹配不可购买", func(t *testing.T) {
		med := m.GetMedicine("antidepressant") // 需要 depression 诊断
		if m.CanPurchase(med, u) {
			t.Error("Prescription with wrong diagnosis should not be purchasable")
		}
	})
}

func TestApplyEffects(t *testing.T) {
	u := user.New()
	u.SetNutrition("维生素C", 50)

	m := NewManager()
	med := m.GetMedicine("vitamin_c")

	ApplyEffects(med, u)

	// 维生素C +60，但上限100
	if u.GetNutrition("维生素C") != 100 {
		t.Errorf("维生素C = %d, want 100 (capped)", u.GetNutrition("维生素C"))
	}
}

func TestApplyEffectsMultiple(t *testing.T) {
	u := user.New()
	u.SetNutrition("铁", 30)
	u.SetNutrition("维生素C", 40)

	m := NewManager()
	med := m.GetMedicine("iron_supplement") // 铁+55, 维C+15

	ApplyEffects(med, u)

	if u.GetNutrition("铁") != 85 {
		t.Errorf("铁 = %d, want 85", u.GetNutrition("铁"))
	}
	if u.GetNutrition("维生素C") != 55 {
		t.Errorf("维生素C = %d, want 55", u.GetNutrition("维生素C"))
	}
}

func TestMedicineGetInfo(t *testing.T) {
	m := NewManager()

	t.Run("OTC药品显示", func(t *testing.T) {
		med := m.GetMedicine("vitamin_c")
		info := med.GetInfo()

		if !strings.Contains(info, "维生素C片") {
			t.Error("Should contain medicine name")
		}
		if !strings.Contains(info, "¥20") {
			t.Error("Should contain price")
		}
		if strings.Contains(info, "[处方]") {
			t.Error("OTC should not have [处方] tag")
		}
	})

	t.Run("处方药显示", func(t *testing.T) {
		med := m.GetMedicine("cold_medicine")
		info := med.GetInfo()

		if !strings.Contains(info, "感冒特效药") {
			t.Error("Should contain medicine name")
		}
		if !strings.Contains(info, "[处方]") {
			t.Error("Prescription should have [处方] tag")
		}
	})
}

func TestMedicineGetTypeString(t *testing.T) {
	m := NewManager()

	otc := m.GetMedicine("multivitamin")
	if otc.GetTypeString() != "非处方药" {
		t.Errorf("OTC type = %s, want 非处方药", otc.GetTypeString())
	}

	prescription := m.GetMedicine("cold_medicine")
	if prescription.GetTypeString() != "处方药" {
		t.Errorf("Prescription type = %s, want 处方药", prescription.GetTypeString())
	}
}

func TestMedicineGetEffectsDescription(t *testing.T) {
	m := NewManager()
	med := m.GetMedicine("vitamin_c")

	desc := med.GetEffectsDescription()

	if !strings.Contains(desc, "维生素C") {
		t.Error("Should contain effect name")
	}
	if !strings.Contains(desc, "+60") {
		t.Error("Should contain effect amount")
	}
}

func TestPrescriptionDiseaseMapping(t *testing.T) {
	m := NewManager()

	mapping := map[string]string{
		"cold_medicine":        "cold",
		"blood_tonic":          "anemia",
		"vitamin_c_injection":  "scurvy",
		"antidiarrheal":        "food_poisoning",
		"nutrition_booster":    "malnutrition",
		"antidepressant":       "depression",
	}

	for medID, diseaseID := range mapping {
		med := m.GetMedicine(medID)
		if med == nil {
			t.Errorf("Medicine %s not found", medID)
			continue
		}
		if med.RequiredDisease != diseaseID {
			t.Errorf("%s RequiredDisease = %s, want %s",
				medID, med.RequiredDisease, diseaseID)
		}
	}
}

func TestMedicinePrices(t *testing.T) {
	m := NewManager()

	// 所有药品价格应该大于0
	for _, med := range m.GetAllMedicines() {
		if med.Price <= 0 {
			t.Errorf("%s has invalid price: %d", med.Name, med.Price)
		}
	}

	// 处方药通常更贵
	cheapestPrescription := 1000
	for _, med := range m.GetPrescriptionMedicines() {
		if med.Price < cheapestPrescription {
			cheapestPrescription = med.Price
		}
	}

	// 感冒特效药是最便宜的处方药 (¥40)
	if cheapestPrescription != 40 {
		t.Errorf("Cheapest prescription = %d, want 40", cheapestPrescription)
	}
}

func TestMedicineEffects(t *testing.T) {
	m := NewManager()

	// 所有药品都应该有至少一个效果
	for _, med := range m.GetAllMedicines() {
		if len(med.Effects) == 0 {
			t.Errorf("%s has no effects", med.Name)
		}
	}

	// 效果值应该大于0
	for _, med := range m.GetAllMedicines() {
		for _, effect := range med.Effects {
			if effect.Amount <= 0 {
				t.Errorf("%s has invalid effect amount for %s: %d",
					med.Name, effect.Attribute, effect.Amount)
			}
		}
	}
}
