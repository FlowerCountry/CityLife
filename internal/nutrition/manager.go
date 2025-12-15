// Package nutrition 管理营养衰减和交互
package nutrition

import (
	"citylife/internal/user"
	"math/rand"
)

// Category 营养分类
type Category int

const (
	CategoryPrimary Category = iota // 核心属性
	CategoryMicro                   // 微量元素
	CategorySpecial                 // 特殊属性
)

// Config 营养配置
type Config struct {
	Name      string
	Category  Category
	BaseDecay float64 // 每次行动的基础衰减
	MinValue  int
	MaxValue  int
}

// Interaction 营养交互规则
type Interaction struct {
	Source     string  // 源营养
	Target     string  // 目标营养
	Ratio      float64 // 影响比例
	Threshold  int     // 阈值
	IsPositive bool    // true=源不足时影响目标，false=源充足时正向影响
}

// Manager 营养管理器
type Manager struct {
	configs        map[string]Config
	interactions   []Interaction
	previousLevels map[string]int
	pendingHints   []string
	hintMessages   map[string]map[int]string // name -> level -> message
}

// NewManager 创建营养管理器
func NewManager() *Manager {
	m := &Manager{
		configs:        make(map[string]Config),
		previousLevels: make(map[string]int),
		hintMessages:   make(map[string]map[int]string),
	}
	m.initConfigs()
	m.initInteractions()
	m.initHintMessages()
	return m
}

// initConfigs 初始化营养配置
func (m *Manager) initConfigs() {
	// 核心属性（Primary）- 衰减快，影响大
	m.configs["饱腹感"] = Config{"饱腹感", CategoryPrimary, 6.0, -30, 100}
	m.configs["饥渴"] = Config{"饥渴", CategoryPrimary, 5.0, -20, 100}
	m.configs["蛋白质"] = Config{"蛋白质", CategoryPrimary, 2.5, -15, 100}
	m.configs["碳水化合物"] = Config{"碳水化合物", CategoryPrimary, 3.0, -15, 100}

	// 微量元素（Micro）- 衰减慢，长期影响
	m.configs["钙"] = Config{"钙", CategoryMicro, 0.5, -10, 100}
	m.configs["糖分"] = Config{"糖分", CategorySpecial, 2.0, -10, 100}
	m.configs["脂肪"] = Config{"脂肪", CategorySpecial, 0.5, 0, 100}
	m.configs["纤维素"] = Config{"纤维素", CategoryMicro, 1.0, -10, 100}
	m.configs["铁"] = Config{"铁", CategoryMicro, 0.5, -10, 100}
	m.configs["维生素A"] = Config{"维生素A", CategoryMicro, 0.5, -10, 100}
	m.configs["维生素B"] = Config{"维生素B", CategoryMicro, 0.5, -10, 100}
	m.configs["维生素C"] = Config{"维生素C", CategoryMicro, 1.0, -10, 100}
	m.configs["维生素D"] = Config{"维生素D", CategoryMicro, 0.3, -10, 100}
	m.configs["维生素E"] = Config{"维生素E", CategoryMicro, 0.3, -10, 100}
	m.configs["钾"] = Config{"钾", CategoryMicro, 0.5, -10, 100}
	m.configs["硒"] = Config{"硒", CategoryMicro, 0.3, -10, 100}
	m.configs["锌"] = Config{"锌", CategoryMicro, 0.3, -10, 100}

	// 特殊属性
	m.configs["幸福感"] = Config{"幸福感", CategorySpecial, 1.5, -20, 100}
	m.configs["精神振奋"] = Config{"精神振奋", CategorySpecial, 3.0, -10, 100}
	m.configs["益生菌"] = Config{"益生菌", CategoryMicro, 0.5, -10, 100}
	m.configs["镁"] = Config{"镁", CategoryMicro, 0.3, -10, 100}

	// 兼容旧的"饥饿"属性
	m.configs["饥饿"] = Config{"饥饿", CategoryPrimary, 5.0, -20, 100}
}

// initInteractions 初始化营养交互规则
func (m *Manager) initInteractions() {
	m.interactions = []Interaction{
		// 铁不足影响蛋白质吸收
		{"铁", "蛋白质", 0.3, 30, true},
		// 维生素D不足影响钙吸收
		{"维生素D", "钙", 0.4, 30, true},
		// 维生素B不足影响碳水化合物代谢
		{"维生素B", "碳水化合物", 0.2, 30, true},
		// 糖分充足临时提升幸福感
		{"糖分", "幸福感", 0.5, 50, false},
		// 脂肪作为能量缓冲
		{"脂肪", "碳水化合物", 0.3, 20, true},
	}
}

// initHintMessages 初始化提示消息
func (m *Manager) initHintMessages() {
	m.hintMessages["饱腹感"] = map[int]string{
		1: "你感到有些饿了",
		2: "你的肚子咕咕叫",
		3: "你饿得头昏眼花",
		4: "你快要饿晕了！",
	}
	m.hintMessages["饥渴"] = map[int]string{
		1: "你感到有些口渴",
		2: "你的嘴唇干裂",
		3: "你渴得喉咙冒烟",
		4: "你快要渴死了！",
	}
	m.hintMessages["蛋白质"] = map[int]string{
		1: "你感到肌肉有些无力",
		2: "你的身体开始虚弱",
		3: "你严重缺乏蛋白质",
		4: "你的身体已经极度虚弱！",
	}
	m.hintMessages["碳水化合物"] = map[int]string{
		1: "你感到有些疲劳",
		2: "你的能量不足",
		3: "你严重缺乏能量",
		4: "你已经筋疲力尽！",
	}
	m.hintMessages["幸福感"] = map[int]string{
		1: "你的心情有些低落",
		2: "你感到郁郁寡欢",
		3: "你非常沮丧",
		4: "你感到极度绝望！",
	}
	m.hintMessages["维生素C"] = map[int]string{
		1: "你可能需要补充维生素C",
		2: "你严重缺乏维生素C",
	}
}

// OnAction 每次行动时调用，处理营养衰减
func (m *Manager) OnAction(u *user.User, actionCost float64) {
	for name, config := range m.configs {
		// 计算衰减倍率
		multiplier := m.GetDecayMultiplier(u, name)

		// 计算实际衰减量
		decayAmount := config.BaseDecay * actionCost * multiplier

		// 应用衰减
		current := u.GetNutrition(name)
		newValue := current - int(decayAmount)

		// 限制在最小值以上
		if newValue < config.MinValue {
			newValue = config.MinValue
		}
		u.SetNutrition(name, newValue)
	}
}

// ProcessInteractions 处理营养交互
func (m *Manager) ProcessInteractions(u *user.User) {
	for _, interaction := range m.interactions {
		sourceValue := u.GetNutrition(interaction.Source)

		if interaction.IsPositive {
			// 源营养不足时，影响目标营养的衰减
			if sourceValue < interaction.Threshold {
				penalty := float64(interaction.Threshold-sourceValue) / 100.0 * interaction.Ratio
				targetValue := u.GetNutrition(interaction.Target)
				extraDecay := int(float64(targetValue) * penalty * 0.1)
				if extraDecay > 0 {
					u.ConsumeNutrition(interaction.Target, extraDecay)
				}
			}
		} else {
			// 源营养充足时，正向影响目标
			if sourceValue > interaction.Threshold {
				boost := int(float64(sourceValue-interaction.Threshold) * interaction.Ratio * 0.05)
				if boost > 0 {
					u.AddNutrition(interaction.Target, boost)
				}
			}
		}
	}
}

// GetDecayMultiplier 获取衰减倍率
func (m *Manager) GetDecayMultiplier(u *user.User, name string) float64 {
	value := u.GetNutrition(name)

	if value > 20 {
		return 1.0 // 正常衰减
	} else if value > 0 {
		return 1.2 // 危险区：加速20%
	} else if value > -10 {
		return 1.5 // 透支区：加速50%
	}
	return 2.0 // 濒死区：加速100%
}

// InitPreviousLevels 初始化之前的等级记录
func (m *Manager) InitPreviousLevels(u *user.User) {
	for name := range m.configs {
		m.previousLevels[name] = u.GetNutritionLevel(name)
	}
}

// CheckLevelChanges 检查等级变化并生成提示
func (m *Manager) CheckLevelChanges(u *user.User) {
	for name := range m.configs {
		currentLevel := u.GetNutritionLevel(name)
		previousLevel := m.previousLevels[name]

		// 等级数值增大 = 状态恶化
		if currentLevel > previousLevel {
			if messages, ok := m.hintMessages[name]; ok {
				if msg, ok := messages[currentLevel]; ok {
					m.pendingHints = append(m.pendingHints, msg)
				}
			}
		}

		m.previousLevels[name] = currentLevel
	}

	// 限制最多3条提示
	if len(m.pendingHints) > 3 {
		m.pendingHints = m.pendingHints[:3]
	}
}

// GetPendingHints 获取待显示的提示
func (m *Manager) GetPendingHints() []string {
	return m.pendingHints
}

// ClearPendingHints 清除待显示的提示
func (m *Manager) ClearPendingHints() {
	m.pendingHints = nil
}

// CheckRandomDeath 检查随机死亡
func (m *Manager) CheckRandomDeath(u *user.User) bool {
	risk := u.GetDeathRisk()
	if risk <= 0 {
		return false
	}
	return rand.Float64() < risk
}

// GetWarnings 获取警告列表
func (m *Manager) GetWarnings(u *user.User) []string {
	var warnings []string
	for name := range m.configs {
		level := u.GetNutritionLevel(name)
		if level == 1 {
			warnings = append(warnings, name+"偏低")
		}
	}
	return warnings
}

// GetCriticals 获取危险列表
func (m *Manager) GetCriticals(u *user.User) []string {
	var criticals []string
	levelNames := map[int]string{
		2: "危险",
		3: "透支",
		4: "濒死",
	}
	for name := range m.configs {
		level := u.GetNutritionLevel(name)
		if level >= 2 {
			criticals = append(criticals, name+levelNames[level])
		}
	}
	return criticals
}
