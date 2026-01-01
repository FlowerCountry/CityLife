package action

import (
	"strings"
	"testing"

	"citylife/internal/disease"
	"citylife/internal/game"
	"citylife/internal/nutrition"
	"citylife/internal/user"
	"citylife/internal/world"
)

// 创建测试用的游戏状态
func newTestGameState() *game.State {
	return &game.State{
		World:     world.New(),
		User:      user.New(),
		Disease:   disease.NewManager(),
		Nutrition: nutrition.NewManager(),
	}
}

// ========== 看病行动测试 ==========

func TestSeeDoctorAction_Info(t *testing.T) {
	action := &SeeDoctorAction{}
	info := action.Info()

	if info != "看病" {
		t.Errorf("Info() = %s, want 看病", info)
	}
}

func TestSeeDoctorAction_ID(t *testing.T) {
	action := &SeeDoctorAction{}
	if action.ID() != "see_doctor" {
		t.Errorf("ID() = %s, want see_doctor", action.ID())
	}
}

func TestSeeDoctorAction_Category(t *testing.T) {
	action := &SeeDoctorAction{}
	if action.Category() != CategoryPrimary {
		t.Errorf("Category() = %d, want CategoryPrimary", action.Category())
	}
}

func TestSeeDoctorAction_Healthy(t *testing.T) {
	state := newTestGameState()
	action := &SeeDoctorAction{}

	result := action.Execute(state)

	if !result.Success {
		t.Error("Should succeed for healthy patient")
	}

	if !strings.Contains(result.Message, "健康") {
		t.Error("Message should mention healthy")
	}
}

func TestSeeDoctorAction_InsufficientFunds(t *testing.T) {
	state := newTestGameState()

	// 触发疾病
	state.User.SetNutrition("维生素C", 5)
	state.Disease.OnAction(state.User)
	for i := 0; i < 10; i++ {
		state.Disease.OnAction(state.User)
	}

	// 如果有疾病，清空钱包
	if len(state.Disease.GetActiveDiseaseIDs()) > 0 {
		// 花光所有钱
		total := state.World.GetWalletTotal()
		state.World.SpendMoney(total)

		action := &SeeDoctorAction{}
		result := action.Execute(state)

		if result.Success {
			t.Error("Should fail with insufficient funds")
		}

		if !strings.Contains(result.Message, "费用不足") {
			t.Error("Message should mention insufficient funds")
		}
	}
}

// ========== 体检菜单测试 ==========

func TestCheckupMenuAction_Info(t *testing.T) {
	action := &CheckupMenuAction{}
	info := action.Info()

	if info != "体检中心" {
		t.Errorf("Info() = %s, want 体检中心", info)
	}
}

func TestCheckupMenuAction_ID(t *testing.T) {
	action := &CheckupMenuAction{}
	if action.ID() != "checkup_menu" {
		t.Errorf("ID() = %s, want checkup_menu", action.ID())
	}
}

func TestCheckupMenuAction_Category(t *testing.T) {
	action := &CheckupMenuAction{}
	if action.Category() != CategoryPrimary {
		t.Errorf("Category() = %d, want CategoryPrimary", action.Category())
	}
}

func TestCheckupMenuAction_Execute(t *testing.T) {
	action := &CheckupMenuAction{}
	result := action.Execute(nil)

	if !result.Success {
		t.Error("Should always succeed")
	}

	if !strings.Contains(result.Message, "套餐检查") || !strings.Contains(result.Message, "单项检查") {
		t.Error("Message should mention checkup types")
	}
}

func TestPerformCheckupAction_Execute(t *testing.T) {
	state := newTestGameState()

	action := &PerformCheckupAction{ItemID: "p_core"}
	result := action.Execute(state)

	if !result.Success {
		t.Error("Checkup should succeed")
	}

	// 检查报告内容
	if !strings.Contains(result.Message, "饱腹感") && !strings.Contains(result.Message, "饥渴") {
		t.Error("Message should contain checkup results")
	}

	// 检查扣款
	if state.World.GetWalletTotal() >= 100 { // 初始100元 - 50元套餐费
		t.Error("Should have deducted checkup fee")
	}
}

func TestPerformCheckupAction_InsufficientFunds(t *testing.T) {
	state := newTestGameState()

	// 花光所有钱
	total := state.World.GetWalletTotal()
	state.World.SpendMoney(total)

	action := &PerformCheckupAction{ItemID: "p_full"} // 全身体检 150元
	result := action.Execute(state)

	if result.Success {
		t.Error("Should fail with insufficient funds")
	}
}

func TestPerformCheckupAction_InvalidItem(t *testing.T) {
	state := newTestGameState()

	action := &PerformCheckupAction{ItemID: "invalid_item"}
	result := action.Execute(state)

	if result.Success {
		t.Error("Should fail with invalid item")
	}
}

// ========== 药品菜单测试 ==========

func TestMedicineMenuAction_Info(t *testing.T) {
	action := &MedicineMenuAction{}
	info := action.Info()

	if info != "购买药品" {
		t.Errorf("Info() = %s, want 购买药品", info)
	}
}

func TestMedicineMenuAction_ID(t *testing.T) {
	action := &MedicineMenuAction{}
	if action.ID() != "medicine_menu" {
		t.Errorf("ID() = %s, want medicine_menu", action.ID())
	}
}

func TestMedicineMenuAction_Execute(t *testing.T) {
	action := &MedicineMenuAction{}
	result := action.Execute(nil)

	if !result.Success {
		t.Error("Should always succeed")
	}

	if !strings.Contains(result.Message, "普通药物") || !strings.Contains(result.Message, "处方药物") {
		t.Error("Message should mention medicine types")
	}
}

// ========== 购买药品测试 ==========

func TestBuyMedicineAction_OTC(t *testing.T) {
	state := newTestGameState()
	// 降低初始值以测试效果增加
	state.User.SetNutrition("维生素C", 30)
	initialVitC := state.User.GetNutrition("维生素C")

	action := &BuyMedicineAction{MedicineID: "vitamin_c"}
	result := action.Execute(state)

	if !result.Success {
		t.Errorf("Should succeed buying OTC medicine, got: %s", result.Message)
	}

	// 检查效果应用
	newVitC := state.User.GetNutrition("维生素C")
	if newVitC <= initialVitC {
		t.Errorf("Vitamin C should increase, was %d, now %d", initialVitC, newVitC)
	}
}

func TestBuyMedicineAction_PrescriptionWithDiagnosis(t *testing.T) {
	state := newTestGameState()
	state.User.AddDiagnosis("cold")

	action := &BuyMedicineAction{MedicineID: "cold_medicine"}
	result := action.Execute(state)

	if !result.Success {
		t.Errorf("Should succeed buying prescription with diagnosis, got: %s", result.Message)
	}

	// 检查诊断记录被清除
	if state.User.HasDiagnosis("cold") {
		t.Error("Diagnosis should be cleared after purchase")
	}
}

func TestBuyMedicineAction_PrescriptionWithoutDiagnosis(t *testing.T) {
	state := newTestGameState()

	action := &BuyMedicineAction{MedicineID: "cold_medicine"}
	result := action.Execute(state)

	if result.Success {
		t.Error("Should fail buying prescription without diagnosis")
	}
}

func TestBuyMedicineAction_InsufficientFunds(t *testing.T) {
	state := newTestGameState()

	// 花光所有钱
	total := state.World.GetWalletTotal()
	state.World.SpendMoney(total)

	action := &BuyMedicineAction{MedicineID: "vitamin_c"}
	result := action.Execute(state)

	if result.Success {
		t.Error("Should fail with insufficient funds")
	}
}

func TestBuyMedicineAction_InvalidMedicine(t *testing.T) {
	state := newTestGameState()

	action := &BuyMedicineAction{MedicineID: "invalid_medicine"}
	result := action.Execute(state)

	if result.Success {
		t.Error("Should fail with invalid medicine")
	}
}

// ========== 疾病查看测试 ==========

func TestViewDiseaseAction_Info(t *testing.T) {
	action := &ViewDiseaseAction{}
	info := action.Info()

	if info != "查看疾病状态" {
		t.Errorf("Info() = %s, want 查看疾病状态", info)
	}
}

func TestViewDiseaseAction_ID(t *testing.T) {
	action := &ViewDiseaseAction{}
	if action.ID() != "view_disease" {
		t.Errorf("ID() = %s, want view_disease", action.ID())
	}
}

func TestViewDiseaseAction_Category(t *testing.T) {
	action := &ViewDiseaseAction{}
	if action.Category() != CategoryInsight {
		t.Errorf("Category() = %d, want CategoryInsight", action.Category())
	}
}

func TestViewDiseaseAction_Healthy(t *testing.T) {
	state := newTestGameState()

	action := &ViewDiseaseAction{}
	result := action.Execute(state)

	if !result.Success {
		t.Error("Should always succeed")
	}

	if !strings.Contains(result.Message, "健康") {
		t.Error("Message should mention healthy")
	}
}

func TestViewDiseaseAction_WithDiagnosis(t *testing.T) {
	state := newTestGameState()
	state.User.AddDiagnosis("cold")

	action := &ViewDiseaseAction{}
	result := action.Execute(state)

	if !strings.Contains(result.Message, "诊断记录") {
		t.Error("Message should show diagnosis records")
	}
}

// ========== 辅助函数测试 ==========

func TestGetCheckupPackages(t *testing.T) {
	packages := GetCheckupPackages()

	if len(packages) == 0 {
		t.Error("Should have checkup packages")
	}

	// 检查包含核心指标套餐
	found := false
	for _, pkg := range packages {
		if pkg.ID == "p_core" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should contain core checkup package")
	}
}

func TestGetCheckupSingles(t *testing.T) {
	singles := GetCheckupSingles()

	if len(singles) == 0 {
		t.Error("Should have single checkup items")
	}
}

func TestGetOTCMedicines(t *testing.T) {
	meds := GetOTCMedicines()

	if len(meds) == 0 {
		t.Error("Should have OTC medicines")
	}

	// 检查包含维生素C
	found := false
	for _, med := range meds {
		if med.ID == "vitamin_c" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should contain vitamin_c")
	}
}

func TestGetAvailablePrescriptions(t *testing.T) {
	state := newTestGameState()

	// 无诊断时应该为空
	meds := GetAvailablePrescriptions(state)
	if len(meds) != 0 {
		t.Error("Should have no prescriptions without diagnosis")
	}

	// 有诊断时应该返回对应药品
	state.User.AddDiagnosis("cold")
	meds = GetAvailablePrescriptions(state)
	if len(meds) == 0 {
		t.Error("Should have prescriptions with diagnosis")
	}
}

// ========== 医院行动集成测试 ==========

func TestHospitalActionsRegistered(t *testing.T) {
	state := newTestGameState()
	actions := GetActionsForLocation(state, world.LocationHospital)

	// 检查必要的行动存在
	hasSeeDoctorAction := false
	hasCheckupMenuAction := false
	hasMedicineMenuAction := false
	hasViewDiseaseAction := false

	for _, a := range actions {
		switch a.(type) {
		case *SeeDoctorAction:
			hasSeeDoctorAction = true
		case *CheckupMenuAction:
			hasCheckupMenuAction = true
		case *MedicineMenuAction:
			hasMedicineMenuAction = true
		case *ViewDiseaseAction:
			hasViewDiseaseAction = true
		}
	}

	if !hasSeeDoctorAction {
		t.Error("Hospital should have SeeDoctorAction")
	}
	if !hasCheckupMenuAction {
		t.Error("Hospital should have CheckupMenuAction")
	}
	if !hasMedicineMenuAction {
		t.Error("Hospital should have MedicineMenuAction")
	}
	if !hasViewDiseaseAction {
		t.Error("Hospital should have ViewDiseaseAction")
	}
}

// ========== 时间消耗测试 ==========

func TestMedicalActionsTimeConsumption(t *testing.T) {
	tests := []struct {
		name           string
		action         Action
		expectedMinute int
	}{
		{"PerformCheckup", &PerformCheckupAction{ItemID: "p_core"}, 15},
		{"BuyMedicine", &BuyMedicineAction{MedicineID: "vitamin_c"}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := newTestGameState()
			initialMinute := state.World.Minute

			result := tt.action.Execute(state)

			if result.Success {
				expectedMinute := (initialMinute + tt.expectedMinute) % 60
				if state.World.Minute != expectedMinute {
					t.Errorf("Minute = %d, want %d", state.World.Minute, expectedMinute)
				}
			}
		})
	}
}
