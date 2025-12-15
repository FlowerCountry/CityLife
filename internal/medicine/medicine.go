// Package medicine 提供药品系统
package medicine

import (
	"citylife/internal/user"
)

// MedicineType 药品类型
type MedicineType int

const (
	TypeOTC          MedicineType = iota // 非处方药
	TypePrescription                     // 处方药
)

// HealthEffect 健康效果
type HealthEffect struct {
	Attribute string
	Amount    int
}

// Medicine 药品
type Medicine struct {
	ID              string
	Name            string
	Price           int
	Effects         []HealthEffect
	Type            MedicineType
	RequiredDisease string // 处方药需要的诊断（空为OTC）
}

// Manager 药品管理器
type Manager struct {
	medicines map[string]*Medicine
}

// OTC药品列表
var otcMedicines = []*Medicine{
	{
		ID:    "multivitamin",
		Name:  "复合维生素片",
		Price: 30,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"维生素A", 45},
			{"维生素B", 45},
			{"维生素C", 45},
			{"维生素D", 45},
			{"维生素E", 45},
		},
	},
	{
		ID:    "vitamin_c",
		Name:  "维生素C片",
		Price: 20,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"维生素C", 60},
		},
	},
	{
		ID:    "calcium",
		Name:  "钙片",
		Price: 25,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"钙", 55},
			{"维生素D", 20},
		},
	},
	{
		ID:    "iron_supplement",
		Name:  "铁剂",
		Price: 35,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"铁", 55},
			{"维生素C", 15},
		},
	},
	{
		ID:    "protein_powder",
		Name:  "蛋白粉",
		Price: 50,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"蛋白质", 50},
			{"碳水化合物", 20},
		},
	},
	{
		ID:    "fiber",
		Name:  "膳食纤维片",
		Price: 22,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"纤维素", 50},
			{"益生菌", 25},
		},
	},
	{
		ID:    "probiotic",
		Name:  "益生菌胶囊",
		Price: 28,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"益生菌", 55},
			{"纤维素", 15},
		},
	},
	{
		ID:    "sports_drink",
		Name:  "运动饮料",
		Price: 15,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"饥渴", 50},
			{"钾", 25},
			{"糖分", 15},
		},
	},
	{
		ID:    "energy_capsule",
		Name:  "能量胶囊",
		Price: 40,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"精神振奋", 55},
			{"碳水化合物", 25},
		},
	},
	{
		ID:    "happy_candy",
		Name:  "快乐糖",
		Price: 18,
		Type:  TypeOTC,
		Effects: []HealthEffect{
			{"幸福感", 40},
			{"糖分", 30},
		},
	},
}

// 处方药列表
var prescriptionMedicines = []*Medicine{
	{
		ID:              "cold_medicine",
		Name:            "感冒特效药",
		Price:           40,
		Type:            TypePrescription,
		RequiredDisease: "cold",
		Effects: []HealthEffect{
			{"维生素C", 60},
			{"精神振奋", 30},
		},
	},
	{
		ID:              "blood_tonic",
		Name:            "补血口服液",
		Price:           80,
		Type:            TypePrescription,
		RequiredDisease: "anemia",
		Effects: []HealthEffect{
			{"铁", 60},
			{"蛋白质", 40},
			{"维生素B", 20},
		},
	},
	{
		ID:              "vitamin_c_injection",
		Name:            "维C注射液",
		Price:           120,
		Type:            TypePrescription,
		RequiredDisease: "scurvy",
		Effects: []HealthEffect{
			{"维生素C", 80},
		},
	},
	{
		ID:              "antidiarrheal",
		Name:            "止泻灵",
		Price:           45,
		Type:            TypePrescription,
		RequiredDisease: "food_poisoning",
		Effects: []HealthEffect{
			{"益生菌", 55},
			{"饱腹感", 30},
			{"饥渴", 25},
		},
	},
	{
		ID:              "nutrition_booster",
		Name:            "营养强化液",
		Price:           100,
		Type:            TypePrescription,
		RequiredDisease: "malnutrition",
		Effects: []HealthEffect{
			{"蛋白质", 55},
			{"碳水化合物", 45},
			{"脂肪", 30},
		},
	},
	{
		ID:              "antidepressant",
		Name:            "抗抑郁药",
		Price:           150,
		Type:            TypePrescription,
		RequiredDisease: "depression",
		Effects: []HealthEffect{
			{"幸福感", 60},
			{"精神振奋", 40},
		},
	},
}

// NewManager 创建药品管理器
func NewManager() *Manager {
	m := &Manager{
		medicines: make(map[string]*Medicine),
	}

	// 注册所有药品
	for _, med := range otcMedicines {
		m.medicines[med.ID] = med
	}
	for _, med := range prescriptionMedicines {
		m.medicines[med.ID] = med
	}

	return m
}

// GetAllMedicines 获取所有药品
func (m *Manager) GetAllMedicines() []*Medicine {
	result := make([]*Medicine, 0, len(m.medicines))
	for _, med := range m.medicines {
		result = append(result, med)
	}
	return result
}

// GetOTCMedicines 获取OTC药品
func (m *Manager) GetOTCMedicines() []*Medicine {
	return otcMedicines
}

// GetPrescriptionMedicines 获取处方药
func (m *Manager) GetPrescriptionMedicines() []*Medicine {
	return prescriptionMedicines
}

// GetAvailablePrescriptions 获取用户可购买的处方药
func (m *Manager) GetAvailablePrescriptions(u *user.User) []*Medicine {
	result := make([]*Medicine, 0)
	for _, med := range prescriptionMedicines {
		if u.HasDiagnosis(med.RequiredDisease) {
			result = append(result, med)
		}
	}
	return result
}

// GetMedicine 获取指定药品
func (m *Manager) GetMedicine(id string) *Medicine {
	return m.medicines[id]
}

// CanPurchase 检查是否可以购买
func (m *Manager) CanPurchase(med *Medicine, u *user.User) bool {
	if med.Type == TypeOTC {
		return true
	}
	return u.HasDiagnosis(med.RequiredDisease)
}

// ApplyEffects 应用药品效果到用户
func ApplyEffects(med *Medicine, u *user.User) {
	for _, effect := range med.Effects {
		u.AddNutrition(effect.Attribute, effect.Amount)
	}
}

// GetInfo 获取药品显示信息
func (med *Medicine) GetInfo() string {
	if med.Type == TypePrescription {
		return med.Name + " ¥" + itoa(med.Price) + " [处方]"
	}
	return med.Name + " ¥" + itoa(med.Price)
}

// GetTypeString 获取类型字符串
func (med *Medicine) GetTypeString() string {
	if med.Type == TypePrescription {
		return "处方药"
	}
	return "非处方药"
}

// GetEffectsDescription 获取效果描述
func (med *Medicine) GetEffectsDescription() string {
	result := ""
	for i, effect := range med.Effects {
		if i > 0 {
			result += ", "
		}
		result += effect.Attribute + "+" + itoa(effect.Amount)
	}
	return result
}

// 简单的整数转字符串
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
