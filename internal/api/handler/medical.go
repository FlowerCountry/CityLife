package handler

import (
	"fmt"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"

	"github.com/gin-gonic/gin"
)

// SeeDoctor 看病
func SeeDoctor(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	act := &action.SeeDoctorAction{}
	result := executor.Execute(act, state)

	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	response.Success(c, resp)
}

// GetCheckups 获取体检项目列表
func GetCheckups(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	packages := action.GetCheckupPackages()
	singles := action.GetCheckupSingles()

	packageList := make([]response.CheckupInfo, len(packages))
	for i, pkg := range packages {
		packageList[i] = response.CheckupInfo{
			ID:    pkg.ID,
			Name:  pkg.Name,
			Price: pkg.Price,
			Items: pkg.Attributes,
		}
	}

	singleList := make([]response.CheckupInfo, len(singles))
	for i, single := range singles {
		singleList[i] = response.CheckupInfo{
			ID:    single.ID,
			Name:  single.Name,
			Price: single.Price,
			Items: single.Attributes,
		}
	}

	response.Success(c, gin.H{
		"packages": packageList,
		"singles":  singleList,
	})
}

// DoCheckup 执行体检
func DoCheckup(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	checkupID := c.Param("id")

	act := &action.PerformCheckupAction{ItemID: checkupID}
	result := executor.Execute(act, state)

	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	response.Success(c, resp)
}

// GetMedicines 获取可购买药品列表
func GetMedicines(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	// OTC药品
	otcMeds := action.GetOTCMedicines()
	otcList := make([]response.MedicineInfo, len(otcMeds))
	for i, med := range otcMeds {
		var effects []string
		for _, eff := range med.Effects {
			effects = append(effects, fmt.Sprintf("%s+%d", eff.Attribute, eff.Amount))
		}
		otcList[i] = response.MedicineInfo{
			ID:      med.ID,
			Name:    med.Name,
			Price:   med.Price,
			Type:    "OTC",
			Effects: effects,
		}
	}

	// 处方药（根据诊断记录）
	prescriptionMeds := action.GetAvailablePrescriptions(state)
	prescriptionList := make([]response.MedicineInfo, len(prescriptionMeds))
	for i, med := range prescriptionMeds {
		var effects []string
		for _, eff := range med.Effects {
			effects = append(effects, fmt.Sprintf("%s+%d", eff.Attribute, eff.Amount))
		}
		prescriptionList[i] = response.MedicineInfo{
			ID:              med.ID,
			Name:            med.Name,
			Price:           med.Price,
			Type:            "处方药",
			RequiredDisease: med.RequiredDisease,
			Effects:         effects,
		}
	}

	response.Success(c, gin.H{
		"otc":          otcList,
		"prescription": prescriptionList,
	})
}

// BuyMedicine 购买药品
func BuyMedicine(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	medicineID := c.Param("id")

	act := &action.BuyMedicineAction{MedicineID: medicineID}
	result := executor.Execute(act, state)

	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	response.Success(c, resp)
}
