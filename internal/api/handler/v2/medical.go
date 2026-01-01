package v2

import (
	"fmt"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

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
		failGameOver(c, state)
		return
	}

	if state.World.Where != world.LocationHospital {
		fail(c, response.ErrCodeInvalidLocation, "你需要在医院才能看病", actionData{State: buildState(state)})
		return
	}

	failCode := ""
	active := state.Disease.GetActiveDiseaseIDs()
	if len(active) > 0 {
		totalCost := 0
		for _, id := range active {
			info := state.Disease.GetDiseaseInfo(id)
			if info == nil {
				continue
			}
			totalCost += info.TreatmentCost
		}
		if state.World.GetWalletTotal() < totalCost {
			failCode = response.ErrCodeInsufficientFunds
		}
	}

	executeAndRespond(c, state, &action.SeeDoctorAction{}, failCode)
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

	ok(c, gin.H{
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
		failGameOver(c, state)
		return
	}

	if state.World.Where != world.LocationHospital {
		fail(c, response.ErrCodeInvalidLocation, "你需要在医院才能体检", actionData{State: buildState(state)})
		return
	}

	checkupID := c.Param("id")
	item := action.GetCheckupService().GetItem(checkupID)
	if item == nil {
		response.BadRequest(c, "未知的检查项目")
		return
	}

	failCode := ""
	if state.World.GetWalletTotal() < item.Price {
		failCode = response.ErrCodeInsufficientFunds
	}

	executeAndRespond(c, state, &action.PerformCheckupAction{ItemID: checkupID}, failCode)
}

// GetMedicines 获取可购买药品列表
func GetMedicines(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	otcMeds := action.GetOTCMedicines()
	otcList := make([]response.MedicineInfo, len(otcMeds))
	for i, med := range otcMeds {
		var effects []string
		for _, eff := range med.Effects {
			effects = append(effects, fmt.Sprintf("%s%+d", eff.Attribute, eff.Amount))
		}
		otcList[i] = response.MedicineInfo{
			ID:      med.ID,
			Name:    med.Name,
			Price:   med.Price,
			Type:    "OTC",
			Effects: effects,
		}
	}

	prescriptionMeds := action.GetAvailablePrescriptions(state)
	prescriptionList := make([]response.MedicineInfo, len(prescriptionMeds))
	for i, med := range prescriptionMeds {
		var effects []string
		for _, eff := range med.Effects {
			effects = append(effects, fmt.Sprintf("%s%+d", eff.Attribute, eff.Amount))
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

	ok(c, gin.H{
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
		failGameOver(c, state)
		return
	}

	if state.World.Where != world.LocationHospital {
		fail(c, response.ErrCodeInvalidLocation, "你需要在医院才能买药", actionData{State: buildState(state)})
		return
	}

	medicineID := c.Param("id")
	medMgr := action.GetMedicineManager()
	med := medMgr.GetMedicine(medicineID)
	if med == nil {
		response.NotFound(c, response.ErrCodeInvalidRequest, "未知的药品")
		return
	}

	if !medMgr.CanPurchase(med, state.User) {
		fail(c, response.ErrCodePrescriptionRequired, "处方药需要先经医生诊断对应疾病", actionData{State: buildState(state)})
		return
	}

	failCode := ""
	if state.World.GetWalletTotal() < med.Price {
		failCode = response.ErrCodeInsufficientFunds
	}

	executeAndRespond(c, state, &action.BuyMedicineAction{MedicineID: medicineID}, failCode)
}
