// Package disease 管理疾病系统
package disease

import (
	"math/rand"
	"strings"

	"citylife/internal/user"
)

// Severity 疾病严重程度
type Severity int

const (
	SeverityMild     Severity = iota // 轻度
	SeverityModerate                 // 中度
	SeveritySevere                   // 重度
)

// Trigger 触发/治愈条件
type Trigger struct {
	NutritionName string
	Threshold     int
}

// Disease 疾病定义
type Disease struct {
	ID                   string
	Name                 string
	Description          string
	Severity             Severity
	TriggerConditions    []Trigger
	TriggerDuration      int     // 满足条件多少次后检查触发
	TriggerChance        float64 // 触发概率
	DecayMultiplier      float64 // 营养衰减倍率
	ActionCostMultiplier float64 // 行动消耗倍率
	CureConditions       []Trigger
	CureDuration         int // 满足条件多少次后治愈
	TreatmentCost        int // 医院治疗费用
}

// ActiveDisease 激活的疾病状态
type ActiveDisease struct {
	ID              string
	TriggerProgress int
	CureProgress    int
	IsActive        bool
}

// Manager 疾病管理器
type Manager struct {
	definitions map[string]Disease
	active      map[string]*ActiveDisease
}

// NewManager 创建疾病管理器
func NewManager() *Manager {
	m := &Manager{
		definitions: make(map[string]Disease),
		active:      make(map[string]*ActiveDisease),
	}
	m.initDiseases()
	return m
}

// initDiseases 初始化疾病定义
func (m *Manager) initDiseases() {
	// 感冒 - 维生素C不足
	m.definitions["cold"] = Disease{
		ID:                   "cold",
		Name:                 "感冒",
		Description:          "维生素C不足导致的感冒",
		Severity:             SeverityMild,
		TriggerConditions:    []Trigger{{"维生素C", 20}},
		TriggerDuration:      10,
		TriggerChance:        0.15,
		DecayMultiplier:      1.2,
		ActionCostMultiplier: 1.2,
		CureConditions:       []Trigger{{"维生素C", 60}},
		CureDuration:         5,
		TreatmentCost:        50,
	}

	// 贫血 - 铁不足
	m.definitions["anemia"] = Disease{
		ID:                   "anemia",
		Name:                 "贫血",
		Description:          "铁元素不足导致的贫血",
		Severity:             SeverityModerate,
		TriggerConditions:    []Trigger{{"铁", 15}},
		TriggerDuration:      15,
		TriggerChance:        0.2,
		DecayMultiplier:      1.5,
		ActionCostMultiplier: 1.5,
		CureConditions:       []Trigger{{"铁", 70}, {"蛋白质", 60}},
		CureDuration:         10,
		TreatmentCost:        200,
	}

	// 坏血病 - 维生素C严重缺乏
	m.definitions["scurvy"] = Disease{
		ID:                   "scurvy",
		Name:                 "坏血病",
		Description:          "维生素C严重缺乏导致的坏血病",
		Severity:             SeveritySevere,
		TriggerConditions:    []Trigger{{"维生素C", -10}},
		TriggerDuration:      5,
		TriggerChance:        0.4,
		DecayMultiplier:      2.0,
		ActionCostMultiplier: 2.0,
		CureConditions:       []Trigger{{"维生素C", 80}},
		CureDuration:         15,
		TreatmentCost:        500,
	}

	// 食物中毒 - 特殊触发
	m.definitions["food_poisoning"] = Disease{
		ID:                   "food_poisoning",
		Name:                 "食物中毒",
		Description:          "食用不洁食物导致的食物中毒",
		Severity:             SeverityModerate,
		TriggerConditions:    []Trigger{}, // 特殊触发
		TriggerDuration:      1,
		TriggerChance:        1.0,
		DecayMultiplier:      2.0,
		ActionCostMultiplier: 1.8,
		CureConditions:       []Trigger{}, // 自然康复
		CureDuration:         8,
		TreatmentCost:        100,
	}

	// 营养不良 - 蛋白质不足
	m.definitions["malnutrition"] = Disease{
		ID:                   "malnutrition",
		Name:                 "营养不良",
		Description:          "蛋白质不足导致的营养不良",
		Severity:             SeverityModerate,
		TriggerConditions:    []Trigger{{"蛋白质", 20}},
		TriggerDuration:      12,
		TriggerChance:        0.2,
		DecayMultiplier:      1.4,
		ActionCostMultiplier: 1.4,
		CureConditions:       []Trigger{{"蛋白质", 70}, {"碳水化合物", 50}},
		CureDuration:         8,
		TreatmentCost:        150,
	}

	// 抑郁症 - 幸福感长期低
	m.definitions["depression"] = Disease{
		ID:                   "depression",
		Name:                 "抑郁症",
		Description:          "幸福感长期过低导致的抑郁症",
		Severity:             SeverityModerate,
		TriggerConditions:    []Trigger{{"幸福感", 10}},
		TriggerDuration:      20,
		TriggerChance:        0.15,
		DecayMultiplier:      1.3,
		ActionCostMultiplier: 1.3,
		CureConditions:       []Trigger{{"幸福感", 60}},
		CureDuration:         15,
		TreatmentCost:        300,
	}
}

// OnAction 每次行动时调用
func (m *Manager) OnAction(u *user.User) {
	m.checkTriggerConditions(u)
	m.checkCureConditions(u)
}

// checkTriggerConditions 检查触发条件
func (m *Manager) checkTriggerConditions(u *user.User) {
	for id, disease := range m.definitions {
		// 跳过特殊触发的疾病
		if len(disease.TriggerConditions) == 0 {
			continue
		}

		// 检查是否已经患病
		ad, exists := m.active[id]
		if exists && ad.IsActive {
			continue
		}

		// 检查触发条件
		if m.checkSingleTrigger(u, disease) {
			if !exists {
				ad = &ActiveDisease{ID: id}
				m.active[id] = ad
			}
			ad.TriggerProgress++

			// 检查是否达到触发阈值
			if ad.TriggerProgress >= disease.TriggerDuration {
				if rand.Float64() < disease.TriggerChance {
					ad.IsActive = true
				}
			}
		} else {
			// 不满足条件，重置触发进度
			if exists && !ad.IsActive {
				ad.TriggerProgress = 0
			}
		}
	}
}

// checkCureConditions 检查治愈条件
func (m *Manager) checkCureConditions(u *user.User) {
	for _, ad := range m.active {
		if !ad.IsActive {
			continue
		}

		disease, ok := m.definitions[ad.ID]
		if !ok {
			continue
		}

		// 食物中毒自然康复
		if len(disease.CureConditions) == 0 {
			ad.CureProgress++
			if ad.CureProgress >= disease.CureDuration {
				ad.IsActive = false
				ad.CureProgress = 0
				ad.TriggerProgress = 0
			}
			continue
		}

		// 检查康复条件
		if m.checkSingleCure(u, disease) {
			ad.CureProgress++
			if ad.CureProgress >= disease.CureDuration {
				ad.IsActive = false
				ad.CureProgress = 0
				ad.TriggerProgress = 0
			}
		} else {
			ad.CureProgress = 0
		}
	}
}

// checkSingleTrigger 检查单个疾病的触发条件
func (m *Manager) checkSingleTrigger(u *user.User, disease Disease) bool {
	for _, trigger := range disease.TriggerConditions {
		value := u.GetNutrition(trigger.NutritionName)
		if value >= trigger.Threshold {
			return false
		}
	}
	return true
}

// checkSingleCure 检查单个疾病的治愈条件
func (m *Manager) checkSingleCure(u *user.User, disease Disease) bool {
	for _, cure := range disease.CureConditions {
		value := u.GetNutrition(cure.NutritionName)
		if value < cure.Threshold {
			return false
		}
	}
	return true
}

// TriggerFoodPoisoning 触发食物中毒
func (m *Manager) TriggerFoodPoisoning(severity float64) {
	if rand.Float64() < severity {
		m.active["food_poisoning"] = &ActiveDisease{
			ID:       "food_poisoning",
			IsActive: true,
		}
	}
}

// GetActiveDiseases 获取当前患病列表
func (m *Manager) GetActiveDiseases() []string {
	var result []string
	for _, ad := range m.active {
		if ad.IsActive {
			if disease, ok := m.definitions[ad.ID]; ok {
				result = append(result, disease.Name)
			}
		}
	}
	return result
}

// GetActiveDiseaseIDs 获取当前患病ID列表
func (m *Manager) GetActiveDiseaseIDs() []string {
	var result []string
	for id, ad := range m.active {
		if ad.IsActive {
			result = append(result, id)
		}
	}
	return result
}

// HasDisease 检查是否患有指定疾病
func (m *Manager) HasDisease(diseaseID string) bool {
	ad, exists := m.active[diseaseID]
	return exists && ad.IsActive
}

// GetActionCostMultiplier 获取行动消耗倍率
func (m *Manager) GetActionCostMultiplier() float64 {
	multiplier := 1.0
	for _, ad := range m.active {
		if ad.IsActive {
			if disease, ok := m.definitions[ad.ID]; ok {
				multiplier *= disease.ActionCostMultiplier
			}
		}
	}
	return multiplier
}

// GetNutritionDecayMultiplier 获取营养衰减倍率
func (m *Manager) GetNutritionDecayMultiplier() float64 {
	multiplier := 1.0
	for _, ad := range m.active {
		if ad.IsActive {
			if disease, ok := m.definitions[ad.ID]; ok {
				multiplier *= disease.DecayMultiplier
			}
		}
	}
	return multiplier
}

// TreatAtHospital 在医院治疗
func (m *Manager) TreatAtHospital(diseaseID string, spend func(int) error) error {
	ad, exists := m.active[diseaseID]
	if !exists || !ad.IsActive {
		return nil
	}

	disease, ok := m.definitions[diseaseID]
	if !ok {
		return nil
	}

	// 扣款
	err := spend(disease.TreatmentCost)
	if err != nil {
		return err
	}

	// 立即治愈
	ad.IsActive = false
	ad.CureProgress = 0
	ad.TriggerProgress = 0
	return nil
}

// GetDiseaseInfo 获取疾病信息
func (m *Manager) GetDiseaseInfo(diseaseID string) *Disease {
	if disease, ok := m.definitions[diseaseID]; ok {
		return &disease
	}
	return nil
}

// GetAllDiseases 获取所有疾病定义
func (m *Manager) GetAllDiseases() map[string]Disease {
	return m.definitions
}

// GetStatusSummary 获取状态摘要
func (m *Manager) GetStatusSummary() string {
	diseases := m.GetActiveDiseases()
	if len(diseases) == 0 {
		return "健康"
	}
	return strings.Join(diseases, ", ")
}

// SeverityToString 严重程度转字符串
func SeverityToString(severity Severity) string {
	switch severity {
	case SeverityMild:
		return "轻度"
	case SeverityModerate:
		return "中度"
	case SeveritySevere:
		return "重度"
	default:
		return "未知"
	}
}

// GetActiveDisease 获取指定疾病的活动状态
func (m *Manager) GetActiveDisease(diseaseID string) *ActiveDisease {
	if ad, exists := m.active[diseaseID]; exists {
		return ad
	}
	return nil
}

// ClearAllDiseases 清除所有疾病状态
func (m *Manager) ClearAllDiseases() {
	m.active = make(map[string]*ActiveDisease)
}

// SetActiveDisease 设置疾病状态（用于加载存档）
func (m *Manager) SetActiveDisease(diseaseID string, triggerProgress, cureProgress int, isActive bool) {
	m.active[diseaseID] = &ActiveDisease{
		ID:              diseaseID,
		TriggerProgress: triggerProgress,
		CureProgress:    cureProgress,
		IsActive:        isActive,
	}
}
