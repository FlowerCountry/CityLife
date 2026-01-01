// Package checkup 提供医院体检服务
package checkup

import (
	"fmt"
	"strings"

	"citylife/internal/user"
)

// CheckupType 检查类型
type CheckupType int

const (
	TypePackage CheckupType = iota // 套餐
	TypeSingle                     // 单项
)

// CheckupItem 检查项目
type CheckupItem struct {
	ID          string
	Name        string
	Description string
	Type        CheckupType
	Price       int
	Attributes  []string
}

// CheckupResult 检查结果
type CheckupResult struct {
	AttributeName string
	Value         int
	MinNormal     int
	MaxNormal     int
	Status        string
	Advice        string
}

// CheckupReport 检查报告
type CheckupReport struct {
	ItemName      string
	TotalCost     int
	Results       []CheckupResult
	OverallAdvice string
}

// Service 体检服务
type Service struct {
	items map[string]*CheckupItem
}

// 套餐项目
var packageItems = []*CheckupItem{
	{
		ID:          "p_core",
		Name:        "核心指标套餐",
		Description: "检测核心生存指标",
		Type:        TypePackage,
		Price:       50,
		Attributes:  []string{"饱腹感", "饥渴", "蛋白质", "碳水化合物"},
	},
	{
		ID:          "p_micro",
		Name:        "微量元素套餐",
		Description: "检测微量元素和矿物质",
		Type:        TypePackage,
		Price:       80,
		Attributes:  []string{"钙", "铁", "锌", "钾", "硒", "镁", "纤维素", "益生菌"},
	},
	{
		ID:          "p_vitamin",
		Name:        "维生素套餐",
		Description: "检测各类维生素水平",
		Type:        TypePackage,
		Price:       60,
		Attributes:  []string{"维生素A", "维生素B", "维生素C", "维生素D", "维生素E"},
	},
	{
		ID:          "p_mental",
		Name:        "心理检查套餐",
		Description: "评估心理健康状况",
		Type:        TypePackage,
		Price:       50,
		Attributes:  []string{"幸福感", "精神振奋", "糖分"},
	},
	{
		ID:          "p_full",
		Name:        "全身体检套餐",
		Description: "全面检查所有健康指标",
		Type:        TypePackage,
		Price:       150,
		Attributes: []string{
			"饱腹感", "饥渴", "蛋白质", "碳水化合物",
			"钙", "铁", "锌", "钾", "硒", "镁", "纤维素", "益生菌",
			"维生素A", "维生素B", "维生素C", "维生素D", "维生素E",
			"幸福感", "精神振奋", "糖分", "脂肪",
		},
	},
}

// 单项检查
var singleItems = []*CheckupItem{
	// 核心指标
	{ID: "s_satiety", Name: "饱腹感检查", Type: TypeSingle, Price: 15, Attributes: []string{"饱腹感"}},
	{ID: "s_thirst", Name: "饥渴度检查", Type: TypeSingle, Price: 15, Attributes: []string{"饥渴"}},
	{ID: "s_protein", Name: "蛋白质检查", Type: TypeSingle, Price: 15, Attributes: []string{"蛋白质"}},
	{ID: "s_carbs", Name: "碳水化合物检查", Type: TypeSingle, Price: 15, Attributes: []string{"碳水化合物"}},

	// 微量元素
	{ID: "s_calcium", Name: "钙元素检查", Type: TypeSingle, Price: 12, Attributes: []string{"钙"}},
	{ID: "s_sugar", Name: "糖分检查", Type: TypeSingle, Price: 12, Attributes: []string{"糖分"}},
	{ID: "s_fat", Name: "脂肪检查", Type: TypeSingle, Price: 12, Attributes: []string{"脂肪"}},
	{ID: "s_fiber", Name: "纤维素检查", Type: TypeSingle, Price: 10, Attributes: []string{"纤维素"}},
	{ID: "s_iron", Name: "铁元素检查", Type: TypeSingle, Price: 12, Attributes: []string{"铁"}},
	{ID: "s_potassium", Name: "钾元素检查", Type: TypeSingle, Price: 10, Attributes: []string{"钾"}},
	{ID: "s_selenium", Name: "硒元素检查", Type: TypeSingle, Price: 10, Attributes: []string{"硒"}},
	{ID: "s_zinc", Name: "锌元素检查", Type: TypeSingle, Price: 10, Attributes: []string{"锌"}},
	{ID: "s_probiotic", Name: "益生菌检查", Type: TypeSingle, Price: 10, Attributes: []string{"益生菌"}},
	{ID: "s_magnesium", Name: "镁元素检查", Type: TypeSingle, Price: 10, Attributes: []string{"镁"}},

	// 维生素
	{ID: "s_vita", Name: "维生素A检查", Type: TypeSingle, Price: 12, Attributes: []string{"维生素A"}},
	{ID: "s_vitb", Name: "维生素B检查", Type: TypeSingle, Price: 12, Attributes: []string{"维生素B"}},
	{ID: "s_vitc", Name: "维生素C检查", Type: TypeSingle, Price: 12, Attributes: []string{"维生素C"}},
	{ID: "s_vitd", Name: "维生素D检查", Type: TypeSingle, Price: 12, Attributes: []string{"维生素D"}},
	{ID: "s_vite", Name: "维生素E检查", Type: TypeSingle, Price: 12, Attributes: []string{"维生素E"}},

	// 心理评估
	{ID: "s_happiness", Name: "幸福感评估", Type: TypeSingle, Price: 20, Attributes: []string{"幸福感"}},
	{ID: "s_energy", Name: "精神状态评估", Type: TypeSingle, Price: 20, Attributes: []string{"精神振奋"}},
}

// NewService 创建体检服务
func NewService() *Service {
	s := &Service{
		items: make(map[string]*CheckupItem),
	}

	// 注册所有项目
	for _, item := range packageItems {
		s.items[item.ID] = item
	}
	for _, item := range singleItems {
		s.items[item.ID] = item
	}

	return s
}

// GetAllItems 获取所有项目
func (s *Service) GetAllItems() []*CheckupItem {
	result := make([]*CheckupItem, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

// GetPackageItems 获取套餐项目
func (s *Service) GetPackageItems() []*CheckupItem {
	return packageItems
}

// GetSingleItems 获取单项检查
func (s *Service) GetSingleItems() []*CheckupItem {
	return singleItems
}

// GetItem 获取指定项目
func (s *Service) GetItem(id string) *CheckupItem {
	return s.items[id]
}

// PerformCheckup 执行检查并生成报告
func (s *Service) PerformCheckup(itemID string, u *user.User) *CheckupReport {
	item := s.items[itemID]
	if item == nil {
		return nil
	}

	report := &CheckupReport{
		ItemName:  item.Name,
		TotalCost: item.Price,
		Results:   make([]CheckupResult, 0, len(item.Attributes)),
	}

	criticalCount := 0
	dangerCount := 0
	warningCount := 0

	for _, attr := range item.Attributes {
		value := u.GetNutrition(attr)
		level := u.GetNutritionLevel(attr)

		result := CheckupResult{
			AttributeName: attr,
			Value:         value,
			MinNormal:     50,
			MaxNormal:     100,
			Status:        getStatusString(level),
			Advice:        getAdviceForAttribute(attr, level),
		}

		report.Results = append(report.Results, result)

		// 统计异常数量
		switch level {
		case user.LevelCritical, user.LevelOverdraft:
			criticalCount++
		case user.LevelDanger:
			dangerCount++
		case user.LevelWarning:
			warningCount++
		}
	}

	// 生成综合建议
	if criticalCount > 0 {
		report.OverallAdvice = "您有多项指标处于危急状态，请立即调整饮食或就医！"
	} else if dangerCount > 0 {
		report.OverallAdvice = "您有部分指标处于危险状态，请尽快补充营养。"
	} else if warningCount > 0 {
		report.OverallAdvice = "您有部分指标偏低，建议适当补充营养。"
	} else {
		report.OverallAdvice = "您的健康状况良好，继续保持！"
	}

	return report
}

// FormatReport 格式化报告
func (s *Service) FormatReport(report *CheckupReport) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("         城市生活医院体检报告\n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("检查项目：%s\n", report.ItemName))
	sb.WriteString(fmt.Sprintf("检查费用：¥%d\n\n", report.TotalCost))

	// 按类别分组显示结果
	categories := map[string][]CheckupResult{
		"核心指标": {},
		"微量元素": {},
		"维生素":  {},
		"心理评估": {},
	}

	for _, result := range report.Results {
		category := getCategoryName(result.AttributeName)
		categories[category] = append(categories[category], result)
	}

	categoryOrder := []string{"核心指标", "微量元素", "维生素", "心理评估"}
	for _, category := range categoryOrder {
		results := categories[category]
		if len(results) == 0 {
			continue
		}

		sb.WriteString("----------------------------------------\n")
		sb.WriteString(fmt.Sprintf("          %s\n", category))
		sb.WriteString("----------------------------------------\n")
		sb.WriteString("项目           数值    参考范围    状态\n")

		for _, result := range results {
			name := padRight(result.AttributeName, 12)
			sb.WriteString(fmt.Sprintf("%s %4d    %d-%d      %s\n",
				name, result.Value, result.MinNormal, result.MaxNormal, result.Status))
		}
		sb.WriteString("\n")
	}

	// 综合评估
	sb.WriteString("========================================\n")
	sb.WriteString("          综合评估\n")
	sb.WriteString("========================================\n")
	sb.WriteString(report.OverallAdvice)
	sb.WriteString("\n\n")

	// 注意事项
	hasAdvice := false
	for _, result := range report.Results {
		if result.Advice != "" {
			if !hasAdvice {
				sb.WriteString("注意事项：\n")
				hasAdvice = true
			}
			sb.WriteString(fmt.Sprintf("- %s：%s\n", result.AttributeName, result.Advice))
		}
	}

	if hasAdvice {
		sb.WriteString("\n")
	}
	sb.WriteString("医嘱：保持规律饮食，适当运动。\n")
	sb.WriteString("========================================\n")

	return sb.String()
}

// getStatusString 获取状态文字
func getStatusString(level int) string {
	switch level {
	case user.LevelSafe:
		return "正常"
	case user.LevelWarning:
		return "偏低"
	case user.LevelDanger:
		return "危险"
	case user.LevelOverdraft:
		return "透支"
	case user.LevelCritical:
		return "濒死"
	default:
		return "未知"
	}
}

// getAdviceForAttribute 获取属性建议
func getAdviceForAttribute(attr string, level int) string {
	if level == user.LevelSafe {
		return ""
	}

	adviceMap := map[string]string{
		"饱腹感":   "建议食用高能量食物，如面包、米饭",
		"饥渴":    "建议多补充水分，可饮用矿泉水、果汁",
		"蛋白质":   "建议食用高蛋白食物，如牛排、鸡蛋、牛奶",
		"碳水化合物": "建议食用主食类食物，如米饭、意大利面",
		"钙":     "建议食用奶制品，如牛奶、酸奶",
		"糖分":    "建议适量食用甜食补充糖分",
		"脂肪":    "建议食用含油脂食物，如坚果、肉类",
		"纤维素":   "建议多食用蔬菜水果",
		"铁":     "建议食用红肉、菠菜等富铁食物",
		"维生素A":  "建议食用沙拉、牛奶",
		"维生素B":  "建议食用谷物、肉类",
		"维生素C":  "建议食用水果，如苹果、橙汁、番茄",
		"维生素D":  "建议多晒太阳，食用鱼类",
		"维生素E":  "建议食用坚果、植物油",
		"钾":     "建议食用香蕉、土豆",
		"硒":     "建议食用海鲜、坚果",
		"锌":     "建议食用牡蛎、牛肉",
		"益生菌":   "建议食用酸奶、发酵食品",
		"镁":     "建议食用绿叶蔬菜、坚果",
		"幸福感":   "建议多进行社交活动，保持积极心态",
		"精神振奋":  "建议适当休息，保证充足睡眠",
	}

	advice := adviceMap[attr]
	if advice == "" {
		advice = "请咨询医生获取专业建议"
	}

	// 根据严重程度添加前缀
	switch level {
	case user.LevelCritical, user.LevelOverdraft:
		return "【紧急】" + advice
	case user.LevelDanger:
		return "【注意】" + advice
	default:
		return advice
	}
}

// getCategoryName 获取属性分类名称
func getCategoryName(attr string) string {
	coreAttrs := map[string]bool{
		"饱腹感": true, "饥渴": true, "蛋白质": true, "碳水化合物": true,
	}
	vitaminAttrs := map[string]bool{
		"维生素A": true, "维生素B": true, "维生素C": true, "维生素D": true, "维生素E": true,
	}
	mentalAttrs := map[string]bool{
		"幸福感": true, "精神振奋": true,
	}

	if coreAttrs[attr] {
		return "核心指标"
	}
	if vitaminAttrs[attr] {
		return "维生素"
	}
	if mentalAttrs[attr] {
		return "心理评估"
	}
	return "微量元素"
}

// padRight 右填充字符串
func padRight(s string, length int) string {
	runeCount := len([]rune(s))
	if runeCount >= length {
		return s
	}
	return s + strings.Repeat(" ", length-runeCount)
}
