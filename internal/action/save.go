// Package action 包含存档相关行动
package action

import (
	"fmt"
	"strings"

	"citylife/internal/disease"
	"citylife/internal/game"
	"citylife/internal/save"
	"citylife/internal/world"
)

// SaveGameAction 保存游戏行动
type SaveGameAction struct {
	Slot int
}

func (a *SaveGameAction) ID() string {
	return fmt.Sprintf("save_%d", a.Slot)
}

func (a *SaveGameAction) Info() string {
	return fmt.Sprintf("保存到槽位 %d", a.Slot)
}

func (a *SaveGameAction) Execute(state *game.State) *Result {
	// 创建存档管理器
	mgr, err := save.NewManager()
	if err != nil {
		return &Result{
			Message: "创建存档管理器失败: " + err.Error(),
			Success: false,
		}
	}

	// 构建存档数据
	saveData := &save.SaveData{
		World: save.WorldData{
			Year:                 state.World.Year,
			Month:                state.World.Month,
			Day:                  state.World.Day,
			Hour:                 state.World.Hour,
			Minute:               state.World.Minute,
			Second:               state.World.Second,
			Where:                state.World.Where,
			Wallet:               state.World.Wallet,
			BankDeposit:          state.World.BankDeposit,
			HousingStatus:        int(state.World.Housing.Status),
			HousingLevel:         int(state.World.Housing.Level),
			RentRemainingSeconds: state.World.Housing.RentRemainingSeconds,
		},
		User: save.UserData{
			Nutrition: state.User.GetAllNutrition(),
		},
		Diseases: getDiseaseData(state.Disease),
	}

	// 保存
	if err := mgr.Save(a.Slot, saveData); err != nil {
		return &Result{
			Message: "保存失败: " + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("游戏已保存到槽位 %d", a.Slot),
		Success: true,
	}
}

func (a *SaveGameAction) Category() EventCategory {
	return CategoryInsight
}

// LoadGameAction 加载游戏行动
type LoadGameAction struct {
	Slot int
}

func (a *LoadGameAction) ID() string {
	return fmt.Sprintf("load_%d", a.Slot)
}

func (a *LoadGameAction) Info() string {
	return fmt.Sprintf("加载槽位 %d", a.Slot)
}

func (a *LoadGameAction) Execute(state *game.State) *Result {
	// 创建存档管理器
	mgr, err := save.NewManager()
	if err != nil {
		return &Result{
			Message: "创建存档管理器失败: " + err.Error(),
			Success: false,
		}
	}

	// 加载
	saveData, err := mgr.Load(a.Slot)
	if err != nil {
		return &Result{
			Message: "加载失败: " + err.Error(),
			Success: false,
		}
	}

	// 恢复世界状态
	state.World.Year = saveData.World.Year
	state.World.Month = saveData.World.Month
	state.World.Day = saveData.World.Day
	state.World.Hour = saveData.World.Hour
	state.World.Minute = saveData.World.Minute
	state.World.Second = saveData.World.Second
	state.World.Where = saveData.World.Where
	state.World.Wallet = saveData.World.Wallet
	state.World.BankDeposit = saveData.World.BankDeposit
	state.World.Housing.Status = world.HousingStatus(saveData.World.HousingStatus)
	state.World.Housing.Level = world.HomeLevel(saveData.World.HousingLevel)
	state.World.Housing.RentRemainingSeconds = saveData.World.RentRemainingSeconds
	if !state.World.Housing.HasHome() {
		state.World.Housing.Clear()
	}

	// 恢复用户状态
	state.User.SetAllNutrition(saveData.User.Nutrition)

	// 恢复疾病状态
	loadDiseaseData(state.Disease, saveData.Diseases)

	return &Result{
		Message: fmt.Sprintf("游戏已从槽位 %d 加载\n%d年%d月%d日 %02d:%02d",
			a.Slot,
			saveData.World.Year, saveData.World.Month, saveData.World.Day,
			saveData.World.Hour, saveData.World.Minute),
		Success: true,
	}
}

func (a *LoadGameAction) Category() EventCategory {
	return CategoryInsight
}

// ViewSaveSlotsAction 查看存档槽位
type ViewSaveSlotsAction struct{}

func (a *ViewSaveSlotsAction) ID() string {
	return "view_saves"
}

func (a *ViewSaveSlotsAction) Info() string {
	return "查看存档"
}

func (a *ViewSaveSlotsAction) Execute(state *game.State) *Result {
	// 创建存档管理器
	mgr, err := save.NewManager()
	if err != nil {
		return &Result{
			Message: "创建存档管理器失败: " + err.Error(),
			Success: false,
		}
	}

	// 获取所有槽位信息
	slots := mgr.GetAvailableSlots()

	var msg strings.Builder
	msg.WriteString("══════ 存档槽位 ══════\n\n")
	for _, slot := range slots {
		status := "空"
		if slot.Exists {
			status = slot.Description
		}
		msg.WriteString(fmt.Sprintf("  槽位 %d: %s\n", slot.Slot, status))
	}
	msg.WriteString("\n══════════════════════")

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *ViewSaveSlotsAction) Category() EventCategory {
	return CategoryInsight
}

// SaveMenuAction 保存菜单行动（自动保存）
type SaveMenuAction struct{}

func (a *SaveMenuAction) ID() string {
	return "save_auto"
}

func (a *SaveMenuAction) Info() string {
	return "保存游戏"
}

func (a *SaveMenuAction) Execute(state *game.State) *Result {
	// 创建存档管理器
	mgr, err := save.NewManager()
	if err != nil {
		return &Result{
			Message: "创建存档管理器失败: " + err.Error(),
			Success: false,
		}
	}

	// 自动保存到第一个可用槽位或最早的槽位
	var targetSlot int = 1
	slots := mgr.GetAvailableSlots()
	for _, slot := range slots {
		if !slot.Exists {
			targetSlot = slot.Slot
			break
		}
	}

	// 构建存档数据
	saveData := &save.SaveData{
		World: save.WorldData{
			Year:                 state.World.Year,
			Month:                state.World.Month,
			Day:                  state.World.Day,
			Hour:                 state.World.Hour,
			Minute:               state.World.Minute,
			Second:               state.World.Second,
			Where:                state.World.Where,
			Wallet:               state.World.Wallet,
			BankDeposit:          state.World.BankDeposit,
			HousingStatus:        int(state.World.Housing.Status),
			HousingLevel:         int(state.World.Housing.Level),
			RentRemainingSeconds: state.World.Housing.RentRemainingSeconds,
		},
		User: save.UserData{
			Nutrition: state.User.GetAllNutrition(),
		},
		Diseases: getDiseaseData(state.Disease),
	}

	// 保存
	if err := mgr.Save(targetSlot, saveData); err != nil {
		return &Result{
			Message: "保存失败: " + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("游戏已保存到槽位 %d", targetSlot),
		Success: true,
	}
}

func (a *SaveMenuAction) Category() EventCategory {
	return CategoryInsight
}

// LoadMenuAction 加载菜单行动（自动加载）
type LoadMenuAction struct{}

func (a *LoadMenuAction) ID() string {
	return "load_auto"
}

func (a *LoadMenuAction) Info() string {
	return "加载游戏"
}

func (a *LoadMenuAction) Execute(state *game.State) *Result {
	// 创建存档管理器
	mgr, err := save.NewManager()
	if err != nil {
		return &Result{
			Message: "创建存档管理器失败: " + err.Error(),
			Success: false,
		}
	}

	// 查找第一个存在的存档
	slots := mgr.GetAvailableSlots()
	var targetSlot int = 0
	for _, slot := range slots {
		if slot.Exists {
			targetSlot = slot.Slot
			break
		}
	}

	if targetSlot == 0 {
		return &Result{
			Message: "没有可用的存档",
			Success: false,
		}
	}

	// 加载
	saveData, err := mgr.Load(targetSlot)
	if err != nil {
		return &Result{
			Message: "加载失败: " + err.Error(),
			Success: false,
		}
	}

	// 恢复世界状态
	state.World.Year = saveData.World.Year
	state.World.Month = saveData.World.Month
	state.World.Day = saveData.World.Day
	state.World.Hour = saveData.World.Hour
	state.World.Minute = saveData.World.Minute
	state.World.Second = saveData.World.Second
	state.World.Where = saveData.World.Where
	state.World.Wallet = saveData.World.Wallet
	state.World.BankDeposit = saveData.World.BankDeposit
	state.World.Housing.Status = world.HousingStatus(saveData.World.HousingStatus)
	state.World.Housing.Level = world.HomeLevel(saveData.World.HousingLevel)
	state.World.Housing.RentRemainingSeconds = saveData.World.RentRemainingSeconds
	if !state.World.Housing.HasHome() {
		state.World.Housing.Clear()
	}

	// 恢复用户状态
	state.User.SetAllNutrition(saveData.User.Nutrition)

	// 恢复疾病状态
	loadDiseaseData(state.Disease, saveData.Diseases)

	return &Result{
		Message: fmt.Sprintf("游戏已从槽位 %d 加载\n%d年%d月%d日 %02d:%02d",
			targetSlot,
			saveData.World.Year, saveData.World.Month, saveData.World.Day,
			saveData.World.Hour, saveData.World.Minute),
		Success: true,
	}
}

func (a *LoadMenuAction) Category() EventCategory {
	return CategoryInsight
}

// getDiseaseData 获取疾病存档数据
func getDiseaseData(dm *disease.Manager) []save.DiseaseData {
	diseaseIDs := dm.GetActiveDiseaseIDs()
	data := make([]save.DiseaseData, 0, len(diseaseIDs))

	for _, id := range diseaseIDs {
		ad := dm.GetActiveDisease(id)
		if ad != nil {
			data = append(data, save.DiseaseData{
				ID:              ad.ID,
				TriggerProgress: ad.TriggerProgress,
				CureProgress:    ad.CureProgress,
				IsActive:        ad.IsActive,
			})
		}
	}

	return data
}

// loadDiseaseData 加载疾病存档数据
func loadDiseaseData(dm *disease.Manager, data []save.DiseaseData) {
	// 清除当前疾病状态
	dm.ClearAllDiseases()

	// 恢复疾病状态
	for _, d := range data {
		dm.SetActiveDisease(d.ID, d.TriggerProgress, d.CureProgress, d.IsActive)
	}
}

// GetSaveSlots 获取存档槽位信息
func GetSaveSlots() ([]save.SlotInfo, error) {
	mgr, err := save.NewManager()
	if err != nil {
		return nil, err
	}
	return mgr.GetAvailableSlots(), nil
}
