// Package food 定义食物类型和属性
package food

// FoodType 食物类型
type FoodType int

const (
	FoodTypeFresh     FoodType = iota // 新鲜食品
	FoodTypeBeverage                  // 饮料
	FoodTypeProcessed                 // 加工食品
	FoodTypeCanned                    // 罐头食品
)

// DecayRates 衰减率（每小时）
var DecayRates = map[FoodType]float64{
	FoodTypeFresh:     4.0, // 25小时完全腐烂
	FoodTypeBeverage:  2.0, // 50小时
	FoodTypeProcessed: 1.0, // 100小时
	FoodTypeCanned:    0.1, // 1000小时
}

// FoodTypeNames 食物类型名称
var FoodTypeNames = map[FoodType]string{
	FoodTypeFresh:     "新鲜",
	FoodTypeBeverage:  "饮料",
	FoodTypeProcessed: "加工",
	FoodTypeCanned:    "罐头",
}

// Health 营养效果
type Health struct {
	Name   string // 营养名称
	Amount int    // 营养值
}

// Food 食物定义
type Food struct {
	Name         string
	BaseHealth   []Health // 基础营养效果
	Price        int
	Type         FoodType
	Freshness    int // 当前新鲜度（0-100，可为负）
	MaxFreshness int // 最大新鲜度
}

// NewFood 创建新食物
func NewFood(name string, price int, foodType FoodType, health []Health) *Food {
	return &Food{
		Name:         name,
		BaseHealth:   health,
		Price:        price,
		Type:         foodType,
		Freshness:    100,
		MaxFreshness: 100,
	}
}

// UpdateFreshness 更新新鲜度（基于游戏时间）
func (f *Food) UpdateFreshness(gameMinutes int) {
	hours := float64(gameMinutes) / 60.0
	decay := DecayRates[f.Type] * hours
	f.Freshness -= int(decay)
}

// GetFreshnessLabel 获取新鲜度标签
func (f *Food) GetFreshnessLabel() string {
	if f.Freshness > 75 {
		return "新鲜"
	} else if f.Freshness > 50 {
		return "较新鲜"
	} else if f.Freshness > 25 {
		return "不太新鲜"
	} else if f.Freshness > 0 {
		return "即将过期"
	}
	return "已过期"
}

// GetNutritionMultiplier 获取营养倍率
func (f *Food) GetNutritionMultiplier() float64 {
	if f.Freshness > 75 {
		return 1.0 // 100%
	} else if f.Freshness > 50 {
		return 0.8 // 80%
	} else if f.Freshness > 25 {
		return 0.5 // 50%
	} else if f.Freshness > 0 {
		return 0.2 // 20%
	}
	return -0.5 // 过期：负效果
}

// GetEffectiveHealth 获取实际营养效果
func (f *Food) GetEffectiveHealth() []Health {
	multiplier := f.GetNutritionMultiplier()

	if multiplier > 0 {
		result := make([]Health, 0, len(f.BaseHealth))
		for _, h := range f.BaseHealth {
			effectiveValue := int(float64(h.Amount) * multiplier)
			if effectiveValue > 0 {
				result = append(result, Health{Name: h.Name, Amount: effectiveValue})
			}
		}
		return result
	}

	// 过期食物：减少饱腹感和幸福感
	return []Health{
		{Name: "饱腹感", Amount: -10},
		{Name: "幸福感", Amount: -15},
		{Name: "饥渴", Amount: -5},
	}
}

// GetFoodPoisoningRisk 获取食物中毒风险
func (f *Food) GetFoodPoisoningRisk() float64 {
	if f.Freshness > 25 {
		return 0.0
	} else if f.Freshness > 0 {
		return 0.1
	} else if f.Freshness > -25 {
		return 0.3
	} else if f.Freshness > -50 {
		return 0.6
	}
	return 0.9
}

// GetInfo 获取显示信息
func (f *Food) GetInfo() string {
	info := f.Name + " " + formatPrice(f.Price)
	label := f.GetFreshnessLabel()
	if label != "新鲜" {
		info += " [" + label + "]"
	}
	return info
}

func formatPrice(price int) string {
	return intToString(price) + "元"
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	negative := n < 0
	if negative {
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
