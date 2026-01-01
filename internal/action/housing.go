package action

import (
	"fmt"
	"strings"

	"citylife/internal/game"
	"citylife/internal/world"
)

const (
	rentDurationSeconds   = 7 * 24 * 60 * 60
	rentServiceTimeSecond = 30 * 60
)

// ViewHousingAction 查看住房信息
type ViewHousingAction struct{}

func (a *ViewHousingAction) ID() string { return "view_housing" }

func (a *ViewHousingAction) Info() string { return "查看住房" }

func (a *ViewHousingAction) Execute(state *game.State) *Result {
	var msg strings.Builder
	msg.WriteString("══════ 住房状态 ══════\n\n")
	msg.WriteString("  " + state.World.Housing.Summary() + "\n")
	msg.WriteString("\n══════════════════════")

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *ViewHousingAction) Category() EventCategory { return CategoryInsight }

// GoHomeAction 回家（住宅区 -> 家）
type GoHomeAction struct{}

func (a *GoHomeAction) ID() string { return "go_home" }

func (a *GoHomeAction) Info() string { return "回家" }

func (a *GoHomeAction) Execute(state *game.State) *Result {
	if !state.World.Housing.HasHome() {
		return &Result{
			Message: "你现在没有可用的住所（可去房产中介租房/买房）",
			Success: false,
		}
	}

	state.World.ChangeLocation(world.LocationHome)
	return &Result{
		Message: "你回到了家。",
		Success: true,
	}
}

func (a *GoHomeAction) Category() EventCategory { return CategoryNavigation }

// LeaveHomeAction 出门（家 -> 住宅区）
type LeaveHomeAction struct{}

func (a *LeaveHomeAction) ID() string { return "leave_home" }

func (a *LeaveHomeAction) Info() string { return "出门" }

func (a *LeaveHomeAction) Execute(state *game.State) *Result {
	state.World.ChangeLocation(world.LocationResidentialArea)
	return &Result{
		Message: "你出了门，来到住宅区。",
		Success: true,
	}
}

func (a *LeaveHomeAction) Category() EventCategory { return CategoryNavigation }

// SleepAtHomeAction 在家睡觉（8小时）
type SleepAtHomeAction struct{}

func (a *SleepAtHomeAction) ID() string { return "sleep_home" }

func (a *SleepAtHomeAction) Info() string { return "睡觉（8小时）" }

func (a *SleepAtHomeAction) Execute(state *game.State) *Result {
	if !state.World.Housing.HasHome() {
		return &Result{
			Message: "你现在没有可用的住所，无法在家睡觉。",
			Success: false,
		}
	}

	beforeRent := state.World.Housing.RentRemainingSeconds
	level := state.World.Housing.Level

	// 睡眠时间
	sleepSeconds := 8 * 60 * 60
	state.World.UpdateTime(sleepSeconds)

	spiritBoost, happinessBoost := homeSleepBenefits(level)
	applyHomeSleepEffects(state, spiritBoost, happinessBoost)

	rentExpired := beforeRent > 0 && !state.World.Housing.HasHome()
	msg := buildHomeSleepMessage(spiritBoost, happinessBoost, rentExpired)

	return &Result{
		Message:     msg,
		Success:     true,
		TimeElapsed: sleepSeconds,
	}
}

func (a *SleepAtHomeAction) Category() EventCategory { return CategoryPrimary }

// RentBedAction 租合租床位
type RentBedAction struct{}

func (a *RentBedAction) ID() string { return "rent_bed" }

func (a *RentBedAction) Info() string { return "租合租床位（7天）- 1200元" }

func (a *RentBedAction) Execute(state *game.State) *Result {
	return rentHome(state, world.HomeBed, 1200)
}

func (a *RentBedAction) Category() EventCategory { return CategoryPrimary }

// RentRoomAction 租单间
type RentRoomAction struct{}

func (a *RentRoomAction) ID() string { return "rent_room" }

func (a *RentRoomAction) Info() string { return "租单间（7天）- 2400元" }

func (a *RentRoomAction) Execute(state *game.State) *Result {
	return rentHome(state, world.HomeRoom, 2400)
}

func (a *RentRoomAction) Category() EventCategory { return CategoryPrimary }

// RentApartmentAction 租小公寓
type RentApartmentAction struct{}

func (a *RentApartmentAction) ID() string { return "rent_apartment" }

func (a *RentApartmentAction) Info() string { return "租小公寓（7天）- 4000元" }

func (a *RentApartmentAction) Execute(state *game.State) *Result {
	return rentHome(state, world.HomeApartment, 4000)
}

func (a *RentApartmentAction) Category() EventCategory { return CategoryPrimary }

// RenewRentAction 续租当前住房
type RenewRentAction struct{}

func (a *RenewRentAction) ID() string { return "renew_rent" }

func (a *RenewRentAction) Info() string { return "续租当前住房（7天）" }

func (a *RenewRentAction) Execute(state *game.State) *Result {
	if !state.World.Housing.IsRented() {
		return &Result{
			Message: "你目前没有租房。",
			Success: false,
		}
	}

	level := state.World.Housing.Level
	remaining := state.World.Housing.RentRemainingSeconds
	price := rentPriceForLevel(level)
	if err := state.World.SpendMoney(price); err != nil {
		return &Result{
			Message: fmt.Sprintf("续租失败：现金不足，需要 %d 元，当前只有 %d 元", price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	state.World.UpdateTime(rentServiceTimeSecond)

	remainingAfter := remaining - rentServiceTimeSecond
	if remainingAfter < 0 {
		remainingAfter = 0
	}
	state.World.Housing.SetRented(level, remainingAfter+rentDurationSeconds)

	return &Result{
		Message:     fmt.Sprintf("续租成功：%s（+7天）", state.World.Housing.HomeName()),
		Success:     true,
		TimeElapsed: rentServiceTimeSecond,
	}
}

func (a *RenewRentAction) Category() EventCategory { return CategoryPrimary }

// CancelRentAction 退租
type CancelRentAction struct{}

func (a *CancelRentAction) ID() string { return "cancel_rent" }

func (a *CancelRentAction) Info() string { return "退租（无退款）" }

func (a *CancelRentAction) Execute(state *game.State) *Result {
	if !state.World.Housing.IsRented() {
		return &Result{
			Message: "你目前没有租房。",
			Success: false,
		}
	}

	state.World.Housing.Clear()
	state.World.UpdateTime(10 * 60)

	return &Result{
		Message:     "你退租了，现在没有可用的住所。",
		Success:     true,
		TimeElapsed: 10 * 60,
	}
}

func (a *CancelRentAction) Category() EventCategory { return CategoryPrimary }

// BuyRoomAction 买单间
type BuyRoomAction struct{}

func (a *BuyRoomAction) ID() string { return "buy_room" }

func (a *BuyRoomAction) Info() string { return "买单间 - 20000元" }

func (a *BuyRoomAction) Execute(state *game.State) *Result {
	return buyHome(state, world.HomeRoom, 20000)
}

func (a *BuyRoomAction) Category() EventCategory { return CategoryPrimary }

// BuyApartmentAction 买小公寓
type BuyApartmentAction struct{}

func (a *BuyApartmentAction) ID() string { return "buy_apartment" }

func (a *BuyApartmentAction) Info() string { return "买小公寓 - 50000元" }

func (a *BuyApartmentAction) Execute(state *game.State) *Result {
	return buyHome(state, world.HomeApartment, 50000)
}

func (a *BuyApartmentAction) Category() EventCategory { return CategoryPrimary }

func homeSleepBenefits(level world.HomeLevel) (spiritBoost, happinessBoost int) {
	// 租房收益明显：档位越高，恢复越多
	switch level {
	case world.HomeBed:
		return 35, 18
	case world.HomeRoom:
		return 55, 28
	case world.HomeApartment:
		return 75, 40
	default:
		return 20, 10
	}
}

func applyHomeSleepEffects(state *game.State, spiritBoost, happinessBoost int) {
	state.User.AddNutrition("精神振奋", spiritBoost)
	state.User.AddNutrition("幸福感", happinessBoost)

	// 睡觉也会饿/渴
	state.User.ConsumeNutrition("饱腹感", 18)
	state.User.ConsumeNutrition("饥渴", 22)
}

func buildHomeSleepMessage(spiritBoost, happinessBoost int, rentExpired bool) string {
	var msg strings.Builder
	msg.WriteString("你睡了一个好觉。\n")
	msg.WriteString(fmt.Sprintf("精神振奋+%d，幸福感+%d", spiritBoost, happinessBoost))
	if rentExpired {
		msg.WriteString("\n\n⚠️ 租期已到，你现在没有可用的住所了。")
	}
	return msg.String()
}

func rentHome(state *game.State, level world.HomeLevel, price int) *Result {
	if state.World.Housing.IsOwned() {
		return &Result{
			Message: "你已经有自有住房，不需要租房。",
			Success: false,
		}
	}

	if err := state.World.SpendMoney(price); err != nil {
		return &Result{
			Message: fmt.Sprintf("租房失败：现金不足，需要 %d 元，当前只有 %d 元", price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	state.World.UpdateTime(rentServiceTimeSecond)
	state.World.Housing.SetRented(level, rentDurationSeconds)

	return &Result{
		Message:     fmt.Sprintf("租房成功：%s（7天），租金 %d 元", state.World.Housing.HomeName(), price),
		Success:     true,
		TimeElapsed: rentServiceTimeSecond,
	}
}

func buyHome(state *game.State, level world.HomeLevel, price int) *Result {
	if err := state.World.SpendMoney(price); err != nil {
		return &Result{
			Message: fmt.Sprintf("购房失败：现金不足，需要 %d 元，当前只有 %d 元", price, state.World.GetWalletTotal()),
			Success: false,
		}
	}

	state.World.UpdateTime(60 * 60)
	state.World.Housing.SetOwned(level)

	return &Result{
		Message:     fmt.Sprintf("购房成功：%s（自有住房）", state.World.Housing.HomeName()),
		Success:     true,
		TimeElapsed: 60 * 60,
	}
}

func rentPriceForLevel(level world.HomeLevel) int {
	switch level {
	case world.HomeBed:
		return 1200
	case world.HomeRoom:
		return 2400
	case world.HomeApartment:
		return 4000
	default:
		return 999999
	}
}

func buildHousingMarketActions(state *game.State) []Action {
	actions := make([]Action, 0, 8)
	h := state.World.Housing

	if h.IsRented() {
		actions = append(actions, &RenewRentAction{}, &CancelRentAction{})
	}

	// 自有住房不提供租房入口（避免无意义的状态切换）
	if !h.IsOwned() {
		actions = append(actions, &RentBedAction{}, &RentRoomAction{}, &RentApartmentAction{})
	}

	// 买房/升级：允许从无房/租房买入；已自有单间可升级小公寓
	if !h.IsOwned() {
		actions = append(actions, &BuyRoomAction{}, &BuyApartmentAction{})
	} else if h.Level == world.HomeRoom {
		actions = append(actions, &BuyApartmentAction{})
	}

	return actions
}
