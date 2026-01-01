// Package user 管理用户健康状态：营养、诊断历史
package user

// NutritionCategory 营养分类
type NutritionCategory int

const (
	CategoryPrimary NutritionCategory = iota // 核心属性（直接影响生存）
	CategoryMicro                            // 微量元素（长期健康）
	CategorySpecial                          // 特殊属性（心理/状态）
)

// NutritionConfig 营养配置
type NutritionConfig struct {
	Name     string
	Category NutritionCategory
	MinValue int // 负数下限
	MaxValue int // 上限（固定100）
}

// AllNutritionConfigs 所有营养类型配置
var AllNutritionConfigs = []NutritionConfig{
	// 核心属性（Primary）- 直接影响生存
	{Name: "饱腹感", Category: CategoryPrimary, MinValue: -30, MaxValue: 100},
	{Name: "饥渴", Category: CategoryPrimary, MinValue: -20, MaxValue: 100},
	{Name: "蛋白质", Category: CategoryPrimary, MinValue: -15, MaxValue: 100},
	{Name: "碳水化合物", Category: CategoryPrimary, MinValue: -15, MaxValue: 100},

	// 微量元素（Micro）- 长期健康
	{Name: "钙", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "糖分", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "脂肪", Category: CategoryMicro, MinValue: 0, MaxValue: 100}, // 脂肪不可为负
	{Name: "纤维素", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "铁", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "维生素A", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "维生素B", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "维生素C", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "维生素D", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "维生素E", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "钾", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "硒", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "锌", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "益生菌", Category: CategoryMicro, MinValue: -10, MaxValue: 100},
	{Name: "镁", Category: CategoryMicro, MinValue: -10, MaxValue: 100},

	// 特殊属性（Special）- 心理/状态
	{Name: "幸福感", Category: CategorySpecial, MinValue: -20, MaxValue: 100},
	{Name: "精神振奋", Category: CategorySpecial, MinValue: -10, MaxValue: 100},

	// 兼容旧食物定义
	{Name: "饥饿", Category: CategoryPrimary, MinValue: -20, MaxValue: 100}, // "饱腹感"的别名
}

// nutritionMinValues 营养最小值缓存
var nutritionMinValues = make(map[string]int)

func init() {
	for _, cfg := range AllNutritionConfigs {
		nutritionMinValues[cfg.Name] = cfg.MinValue
	}
}

// User 用户状态
type User struct {
	Nutrition map[string]int  // 营养值
	Diagnoses map[string]bool // 已诊断的疾病
}

// New 创建新用户
func New() *User {
	u := &User{
		Nutrition: make(map[string]int),
		Diagnoses: make(map[string]bool),
	}

	// 初始化所有营养值为100
	for _, cfg := range AllNutritionConfigs {
		u.Nutrition[cfg.Name] = 100
	}

	return u
}

// GetNutrition 获取营养值
func (u *User) GetNutrition(name string) int {
	if val, ok := u.Nutrition[name]; ok {
		return val
	}
	return 0
}

// SetNutrition 设置营养值（应用上下限）
func (u *User) SetNutrition(name string, value int) {
	minVal := GetMinValue(name)
	if value < minVal {
		value = minVal
	}
	if value > 100 {
		value = 100
	}
	u.Nutrition[name] = value
}

// ConsumeNutrition 消耗营养（减少）
func (u *User) ConsumeNutrition(name string, amount int) {
	current := u.GetNutrition(name)
	u.SetNutrition(name, current-amount)
}

// AddNutrition 补充营养（增加，上限100）
func (u *User) AddNutrition(name string, amount int) {
	current := u.GetNutrition(name)
	newVal := current + amount
	if newVal > 100 {
		newVal = 100
	}
	u.Nutrition[name] = newVal
}

// GetAllNutritionNames 获取所有营养名称
func (u *User) GetAllNutritionNames() []string {
	names := make([]string, 0, len(u.Nutrition))
	for name := range u.Nutrition {
		names = append(names, name)
	}
	return names
}

// IsAlive 检查是否存活
func (u *User) IsAlive() bool {
	// 检查核心属性是否低于死亡线
	if u.GetNutrition("饱腹感") <= -30 {
		return false
	}
	if u.GetNutrition("饥渴") <= -20 {
		return false
	}
	if u.GetNutrition("蛋白质") <= -15 {
		return false
	}
	if u.GetNutrition("碳水化合物") <= -15 {
		return false
	}
	return true
}

// GetDeathRisk 获取死亡风险（0.0-1.0）
func (u *User) GetDeathRisk() float64 {
	risk := 0.0

	// 检查核心属性的濒死区间
	satiety := u.GetNutrition("饱腹感")
	if satiety < -10 && satiety > -30 {
		risk += 0.05
	}

	thirst := u.GetNutrition("饥渴")
	if thirst < -10 && thirst > -20 {
		risk += 0.05
	}

	protein := u.GetNutrition("蛋白质")
	if protein < -10 && protein > -15 {
		risk += 0.03
	}

	carbs := u.GetNutrition("碳水化合物")
	if carbs < -10 && carbs > -15 {
		risk += 0.03
	}

	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// NutritionLevel 营养等级
const (
	LevelSafe      = 0 // 安全区 (50-100)
	LevelWarning   = 1 // 警告区 (20-50)
	LevelDanger    = 2 // 危险区 (0-20)
	LevelOverdraft = 3 // 透支区 (-10 到 0)
	LevelCritical  = 4 // 濒死区 (< -10)
)

// GetNutritionLevel 获取营养等级（0-4）
func (u *User) GetNutritionLevel(name string) int {
	value := u.GetNutrition(name)

	if value > 50 {
		return LevelSafe
	} else if value > 20 {
		return LevelWarning
	} else if value > 0 {
		return LevelDanger
	} else if value > -10 {
		return LevelOverdraft
	}
	return LevelCritical
}

// GetMinValue 获取营养的最小值
func GetMinValue(name string) int {
	if minVal, ok := nutritionMinValues[name]; ok {
		return minVal
	}
	return -10 // 默认下限
}

// AddDiagnosis 添加诊断记录
func (u *User) AddDiagnosis(diseaseID string) {
	u.Diagnoses[diseaseID] = true
}

// HasDiagnosis 检查是否有诊断记录
func (u *User) HasDiagnosis(diseaseID string) bool {
	return u.Diagnoses[diseaseID]
}

// ClearDiagnosis 清除诊断记录
func (u *User) ClearDiagnosis(diseaseID string) {
	delete(u.Diagnoses, diseaseID)
}

// GetDiagnoses 获取所有诊断记录
func (u *User) GetDiagnoses() []string {
	diagnoses := make([]string, 0, len(u.Diagnoses))
	for id := range u.Diagnoses {
		diagnoses = append(diagnoses, id)
	}
	return diagnoses
}

// GetAllNutrition 获取所有营养值（用于存档）
func (u *User) GetAllNutrition() map[string]int {
	result := make(map[string]int, len(u.Nutrition))
	for k, v := range u.Nutrition {
		result[k] = v
	}
	return result
}

// SetAllNutrition 设置所有营养值（用于加载存档）
func (u *User) SetAllNutrition(nutrition map[string]int) {
	for name, value := range nutrition {
		u.Nutrition[name] = value
	}
}
