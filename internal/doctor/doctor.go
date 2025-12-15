// Package doctor 提供医院看病诊断服务
package doctor

import (
	"fmt"
	"strings"

	"citylife/internal/disease"
	"citylife/internal/user"
)

// DiagnosisResult 诊断结果
type DiagnosisResult struct {
	DiseaseID   string
	DiseaseName string
	Severity    string
	Description string
	Advice      string
}

// DiagnosisReport 诊断报告
type DiagnosisReport struct {
	Results       []DiagnosisResult
	TotalCost     int
	OverallStatus string
}

// 诊断费用
const DiagnosisCost = 50

// 疾病ID到中文名的映射
var diseaseNames = map[string]string{
	"cold":           "感冒",
	"anemia":         "贫血",
	"scurvy":         "坏血病",
	"food_poisoning": "食物中毒",
	"malnutrition":   "营养不良",
	"depression":     "抑郁症",
}

// 疾病建议
var diseaseAdvice = map[string]string{
	"cold":           "多休息，补充维生素C，可购买感冒特效药",
	"anemia":         "补充铁元素和蛋白质，可购买补血口服液",
	"scurvy":         "大量补充维生素C，可购买维C注射液",
	"food_poisoning": "多喝水，注意饮食卫生，可购买止泻灵",
	"malnutrition":   "均衡饮食，补充蛋白质，可购买营养强化液",
	"depression":     "多参与社交活动，保持心情愉悦，可购买抗抑郁药",
}

// Diagnose 执行诊断
func Diagnose(u *user.User, dm *disease.Manager) *DiagnosisReport {
	report := &DiagnosisReport{
		Results:   make([]DiagnosisResult, 0),
		TotalCost: DiagnosisCost,
	}

	// 获取当前疾病
	diseaseIDs := dm.GetActiveDiseaseIDs()

	if len(diseaseIDs) == 0 {
		report.OverallStatus = "您目前身体健康，没有发现任何疾病。"
		return report
	}

	// 为每个疾病生成诊断结果
	for _, id := range diseaseIDs {
		info := dm.GetDiseaseInfo(id)
		if info == nil {
			continue
		}

		result := DiagnosisResult{
			DiseaseID:   id,
			DiseaseName: info.Name,
			Severity:    disease.SeverityToString(info.Severity),
			Description: info.Description,
			Advice:      getAdvice(id),
		}

		report.Results = append(report.Results, result)

		// 添加到用户诊断记录
		u.AddDiagnosis(id)
	}

	// 生成总体状态
	if len(report.Results) == 1 {
		report.OverallStatus = fmt.Sprintf("诊断发现您患有%s，请及时治疗。", report.Results[0].DiseaseName)
	} else {
		names := make([]string, len(report.Results))
		for i, r := range report.Results {
			names[i] = r.DiseaseName
		}
		report.OverallStatus = fmt.Sprintf("诊断发现您患有多种疾病：%s，请及时治疗。", strings.Join(names, "、"))
	}

	return report
}

// getAdvice 获取疾病建议
func getAdvice(diseaseID string) string {
	if advice, ok := diseaseAdvice[diseaseID]; ok {
		return advice
	}
	return "请遵医嘱进行治疗"
}

// GetDiseaseName 获取疾病中文名
func GetDiseaseName(diseaseID string) string {
	if name, ok := diseaseNames[diseaseID]; ok {
		return name
	}
	return diseaseID
}

// FormatReport 格式化诊断报告
func FormatReport(report *DiagnosisReport) string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("         城市生活医院诊断报告\n")
	sb.WriteString("========================================\n")
	sb.WriteString(fmt.Sprintf("诊断费用：¥%d\n\n", report.TotalCost))

	if len(report.Results) == 0 {
		sb.WriteString("诊断结果：\n")
		sb.WriteString("  ✓ 身体健康，未发现疾病\n\n")
	} else {
		sb.WriteString("诊断结果：\n")
		sb.WriteString("----------------------------------------\n")

		for _, result := range report.Results {
			sb.WriteString(fmt.Sprintf("疾病名称：%s\n", result.DiseaseName))
			sb.WriteString(fmt.Sprintf("严重程度：%s\n", result.Severity))
			sb.WriteString(fmt.Sprintf("病情描述：%s\n", result.Description))
			sb.WriteString(fmt.Sprintf("医生建议：%s\n", result.Advice))
			sb.WriteString("----------------------------------------\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString("========================================\n")
	sb.WriteString("          综合评估\n")
	sb.WriteString("========================================\n")
	sb.WriteString(report.OverallStatus)
	sb.WriteString("\n\n")

	if len(report.Results) > 0 {
		sb.WriteString("提示：您可以到医院药房购买处方药进行治疗。\n")
		sb.WriteString("      诊断记录已保存，凭此记录可购买对应处方药。\n")
	}

	sb.WriteString("========================================\n")

	return sb.String()
}

// GetCost 获取诊断费用
func GetCost() int {
	return DiagnosisCost
}

// HasDiagnosedDisease 检查用户是否已诊断某疾病
func HasDiagnosedDisease(u *user.User, diseaseID string) bool {
	return u.HasDiagnosis(diseaseID)
}

// GetUserDiagnoses 获取用户所有诊断记录
func GetUserDiagnoses(u *user.User) []string {
	return u.GetDiagnoses()
}

// ClearDiagnosis 清除诊断记录（疾病治愈后调用）
func ClearDiagnosis(u *user.User, diseaseID string) {
	u.ClearDiagnosis(diseaseID)
}
