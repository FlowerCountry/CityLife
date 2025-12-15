// Package action 包含行动注册表
package action

import (
	"citylife/internal/world"
)

// GetActionsForLocation 获取指定位置的行动列表
func GetActionsForLocation(location int) []Action {
	actions := []Action{}

	// 添加导航行动
	adjacentBuildings := world.GetAdjacentBuildings(location)
	for _, targetID := range adjacentBuildings {
		if targetID >= 0 && targetID < len(world.BuildingNames) {
			actions = append(actions, &GoWhereAction{
				TargetID:   targetID,
				TargetName: world.BuildingNames[targetID],
			})
		}
	}

	// 根据位置添加特定行动
	switch location {
	case world.LocationCityCenter:
		actions = append(actions, &InformationAction{})
		actions = append(actions, &ViewMapAction{})
		actions = append(actions, &ViewHealthAction{})
		actions = append(actions, &ViewSaveSlotsAction{})
		actions = append(actions, &SaveMenuAction{})
		actions = append(actions, &LoadMenuAction{})
		actions = append(actions, &CheckCashAction{})

	case world.LocationSupermarket:
		actions = append(actions, &EnterSupermarketAction{})
		actions = append(actions, &CheckCashAction{})

	case world.LocationBank:
		actions = append(actions, &DepositAction{})
		actions = append(actions, &DepositAllAction{})
		actions = append(actions, &WithdrawAction{})
		actions = append(actions, &WithdrawAllAction{})
		actions = append(actions, &CheckBankAction{})
		actions = append(actions, &CheckCashAction{})

	case world.LocationSupermarketInner:
		// 添加所有商品
		actions = append(actions, GetCommodityActions()...)
		actions = append(actions, &LeaveSupermarketAction{})
		actions = append(actions, &CheckCashAction{})

	case world.LocationHospital:
		// 医疗功能
		actions = append(actions, &SeeDoctorAction{})
		actions = append(actions, &CheckupMenuAction{})
		actions = append(actions, &MedicineMenuAction{})
		actions = append(actions, &ViewDiseaseAction{})
		actions = append(actions, &ViewHealthAction{})
		actions = append(actions, &CheckCashAction{})
	}

	return actions
}
