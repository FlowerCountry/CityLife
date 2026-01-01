// Package action 包含行动注册表
package action

import (
	"citylife/internal/game"
	"citylife/internal/world"
)

// GetActionsForLocation 获取指定位置的行动列表
func GetActionsForLocation(state *game.State, location int) []Action {
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
		actions = append(actions, &ViewHousingAction{})
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

	case world.LocationResidentialArea:
		actions = append(actions, &ViewHousingAction{})
		if state.World.Housing.HasHome() {
			actions = append(actions, &GoHomeAction{})
		}
		actions = append(actions, &CheckCashAction{})

	case world.LocationHome:
		actions = append(actions, &ViewHousingAction{})
		actions = append(actions, &SleepAtHomeAction{})
		actions = append(actions, &LeaveHomeAction{})
		actions = append(actions, &CheckCashAction{})

	case world.LocationRealEstateAgency:
		actions = append(actions, &ViewHousingAction{})
		actions = append(actions, buildHousingMarketActions(state)...)
		actions = append(actions, &CheckCashAction{})

	case world.LocationJobMarket:
		actions = append(actions, GetJobActions()...)
		actions = append(actions, &CheckCashAction{})

	case world.LocationRestaurant:
		actions = append(actions, GetRestaurantActions()...)
		actions = append(actions, &CheckCashAction{})

	case world.LocationPark:
		actions = append(actions, GetParkActions()...)
		actions = append(actions, &ViewHealthAction{})

	case world.LocationHotel:
		actions = append(actions, GetHotelActions()...)
		actions = append(actions, &CheckCashAction{})
	}

	return actions
}
