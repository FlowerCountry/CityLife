// Package action 包含医疗相关行动
package action

import (
	"fmt"
	"strings"

	"citylife/internal/checkup"
	"citylife/internal/disease"
	"citylife/internal/game"
	"citylife/internal/medicine"
)

// ========== 看病行动 ==========

// SeeDoctorAction 看病行动
type SeeDoctorAction struct{}

func (a *SeeDoctorAction) ID() string {
	return "see_doctor"
}

func (a *SeeDoctorAction) Info() string {
	return "看病"
}

func (a *SeeDoctorAction) Execute(state *game.State) *Result {
	// 获取当前疾病
	diseaseIDs := state.Disease.GetActiveDiseaseIDs()
	if len(diseaseIDs) == 0 {
		return &Result{
			Message: "医生检查后说：你很健康，不需要治疗。",
			Success: true,
		}
	}

	// 计算总费用
	totalCost := 0
	var treatments []string
	for _, id := range diseaseIDs {
		info := state.Disease.GetDiseaseInfo(id)
		if info != nil {
			totalCost += info.TreatmentCost
			treatments = append(treatments, fmt.Sprintf("%s(%d元)", info.Name, info.TreatmentCost))
		}
	}

	// 检查钱包
	if state.World.GetWalletTotal() < totalCost {
		return &Result{
			Message: fmt.Sprintf("治疗费用不足！需要 %d 元，当前只有 %d 元\n需要治疗: %s",
				totalCost, state.World.GetWalletTotal(), strings.Join(treatments, ", ")),
			Success: false,
		}
	}

	// 治疗所有疾病并添加诊断记录
	var cured []string
	for _, id := range diseaseIDs {
		// 添加诊断记录（用于购买处方药）
		state.User.AddDiagnosis(id)

		err := state.Disease.TreatAtHospital(id, func(cost int) error {
			return state.World.SpendMoney(cost)
		})
		if err == nil {
			info := state.Disease.GetDiseaseInfo(id)
			if info != nil {
				cured = append(cured, info.Name)
			}
		}
	}

	// 更新时间（看病消耗30分钟）
	state.World.UpdateTime(30 * 60)

	return &Result{
		Message: fmt.Sprintf("治疗完成！花费 %d 元\n已治愈: %s\n\n提示：您现在可以购买相关处方药预防复发",
			totalCost, strings.Join(cured, ", ")),
		Success:     true,
		TimeElapsed: 30 * 60,
	}
}

func (a *SeeDoctorAction) Category() EventCategory {
	return CategoryPrimary
}

// ========== 体检系统 ==========

// 体检服务单例
var checkupService = checkup.NewService()

// CheckupMenuAction 体检中心入口
type CheckupMenuAction struct{}

func (a *CheckupMenuAction) ID() string {
	return "checkup_menu"
}

func (a *CheckupMenuAction) Info() string {
	return "体检中心"
}

func (a *CheckupMenuAction) Execute(state *game.State) *Result {
	return &Result{
		Message: "请选择检查类型：\n\n套餐检查：多项指标组合，优惠价格\n单项检查：针对特定指标的精确检查",
		Success: true,
	}
}

func (a *CheckupMenuAction) Category() EventCategory {
	return CategoryPrimary
}

// GetCheckupPackages 获取体检套餐列表
func GetCheckupPackages() []*checkup.CheckupItem {
	return checkupService.GetPackageItems()
}

// GetCheckupSingles 获取单项检查列表
func GetCheckupSingles() []*checkup.CheckupItem {
	return checkupService.GetSingleItems()
}

// PerformCheckupAction 执行体检
type PerformCheckupAction struct {
	ItemID string
}

func (a *PerformCheckupAction) ID() string {
	return "checkup_" + a.ItemID
}

func (a *PerformCheckupAction) Info() string {
	item := checkupService.GetItem(a.ItemID)
	if item != nil {
		return fmt.Sprintf("%s（¥%d）", item.Name, item.Price)
	}
	return "体检"
}

func (a *PerformCheckupAction) Execute(state *game.State) *Result {
	item := checkupService.GetItem(a.ItemID)

	if item == nil {
		return &Result{
			Message: "未知的检查项目",
			Success: false,
		}
	}

	// 检查钱包
	if state.World.GetWalletTotal() < item.Price {
		return &Result{
			Message: fmt.Sprintf("费用不足！%s 需要 ¥%d，当前只有 ¥%d",
				item.Name, item.Price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	// 扣款
	err := state.World.SpendMoney(item.Price)
	if err != nil {
		return &Result{
			Message: "扣款失败：" + err.Error(),
			Success: false,
		}
	}

	// 执行检查并生成报告
	report := checkupService.PerformCheckup(a.ItemID, state.User)
	if report == nil {
		return &Result{
			Message: "检查失败",
			Success: false,
		}
	}

	// 格式化报告
	reportStr := checkupService.FormatReport(report)

	// 更新时间（体检消耗15分钟）
	state.World.UpdateTime(15 * 60)

	return &Result{
		Message:     reportStr,
		Success:     true,
		TimeElapsed: 15 * 60,
	}
}

func (a *PerformCheckupAction) Category() EventCategory {
	return CategoryPrimary
}

// ========== 药品系统 ==========

// 药品管理器单例
var medicineManager = medicine.NewManager()

// MedicineMenuAction 药品购买入口
type MedicineMenuAction struct{}

func (a *MedicineMenuAction) ID() string {
	return "medicine_menu"
}

func (a *MedicineMenuAction) Info() string {
	return "购买药品"
}

func (a *MedicineMenuAction) Execute(state *game.State) *Result {
	return &Result{
		Message: "请选择药品类型：\n\n普通药物：无需处方，可随时购买\n处方药物：需要先经医生诊断",
		Success: true,
	}
}

func (a *MedicineMenuAction) Category() EventCategory {
	return CategoryPrimary
}

// GetOTCMedicines 获取OTC药品列表
func GetOTCMedicines() []*medicine.Medicine {
	return medicineManager.GetOTCMedicines()
}

// GetAvailablePrescriptions 获取可用的处方药列表
func GetAvailablePrescriptions(state *game.State) []*medicine.Medicine {
	return medicineManager.GetAvailablePrescriptions(state.User)
}

// BuyMedicineAction 购买药品行动
type BuyMedicineAction struct {
	MedicineID string
}

func (a *BuyMedicineAction) ID() string {
	return "buy_medicine_" + a.MedicineID
}

func (a *BuyMedicineAction) Info() string {
	med := medicineManager.GetMedicine(a.MedicineID)
	if med != nil {
		return med.GetInfo()
	}
	return "购买药品"
}

func (a *BuyMedicineAction) Execute(state *game.State) *Result {
	med := medicineManager.GetMedicine(a.MedicineID)

	if med == nil {
		return &Result{
			Message: "未知的药品",
			Success: false,
		}
	}

	// 检查处方药是否可购买
	if !medicineManager.CanPurchase(med, state.User) {
		return &Result{
			Message: fmt.Sprintf("购买失败！%s 是处方药，需要先经医生诊断对应疾病才能购买。", med.Name),
			Success: false,
		}
	}

	// 检查钱包
	if state.World.GetWalletTotal() < med.Price {
		return &Result{
			Message: fmt.Sprintf("购买失败！现金不足，需要 ¥%d，当前只有 ¥%d",
				med.Price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	// 扣款
	err := state.World.SpendMoney(med.Price)
	if err != nil {
		return &Result{
			Message: "购买失败：" + err.Error(),
			Success: false,
		}
	}

	// 应用药品效果
	medicine.ApplyEffects(med, state.User)

	// 如果是处方药，购买后清除诊断记录
	if med.Type == medicine.TypePrescription && med.RequiredDisease != "" {
		state.User.ClearDiagnosis(med.RequiredDisease)
	}

	// 更新时间（购药消耗5分钟）
	state.World.UpdateTime(5 * 60)

	return &Result{
		Message: fmt.Sprintf("购买成功！%s（¥%d）\n效果: %s",
			med.Name, med.Price, med.GetEffectsDescription()),
		Success:     true,
		TimeElapsed: 5 * 60,
	}
}

func (a *BuyMedicineAction) Category() EventCategory {
	return CategoryPrimary
}

// ========== 疾病查看 ==========

// ViewDiseaseAction 查看当前疾病
type ViewDiseaseAction struct{}

func (a *ViewDiseaseAction) ID() string {
	return "view_disease"
}

func (a *ViewDiseaseAction) Info() string {
	return "查看疾病状态"
}

func (a *ViewDiseaseAction) Execute(state *game.State) *Result {
	var msg strings.Builder
	msg.WriteString("══════ 疾病状态 ══════\n\n")

	diseaseIDs := state.Disease.GetActiveDiseaseIDs()
	if len(diseaseIDs) == 0 {
		msg.WriteString("当前没有任何疾病，身体健康！\n")
	} else {
		for _, id := range diseaseIDs {
			info := state.Disease.GetDiseaseInfo(id)
			if info != nil {
				severity := disease.SeverityToString(info.Severity)
				msg.WriteString(fmt.Sprintf("【%s】(%s)\n", info.Name, severity))
				msg.WriteString(fmt.Sprintf("  描述: %s\n", info.Description))
				msg.WriteString(fmt.Sprintf("  治疗费用: %d 元\n", info.TreatmentCost))
				msg.WriteString("\n")
			}
		}
	}

	// 显示诊断记录
	diagnoses := state.User.GetDiagnoses()
	if len(diagnoses) > 0 {
		msg.WriteString("\n【诊断记录】\n")
		msg.WriteString("以下诊断记录可用于购买处方药：\n")
		for _, d := range diagnoses {
			info := state.Disease.GetDiseaseInfo(d)
			if info != nil {
				msg.WriteString(fmt.Sprintf("  - %s\n", info.Name))
			} else {
				msg.WriteString(fmt.Sprintf("  - %s\n", d))
			}
		}
	}

	msg.WriteString("\n══════════════════════")

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *ViewDiseaseAction) Category() EventCategory {
	return CategoryInsight
}

// ========== 辅助函数 ==========

// GetMedicineManager 获取药品管理器
func GetMedicineManager() *medicine.Manager {
	return medicineManager
}

// GetCheckupService 获取体检服务
func GetCheckupService() *checkup.Service {
	return checkupService
}
