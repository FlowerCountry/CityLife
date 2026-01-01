package action

import (
	"strings"
	"testing"

	"citylife/internal/world"
)

func TestGoHomeAction_NoHome(t *testing.T) {
	state := newTestGameState()
	state.World.ChangeLocation(world.LocationResidentialArea)

	act := &GoHomeAction{}
	res := act.Execute(state)

	if res.Success {
		t.Fatal("无住房时不应允许回家")
	}
	if state.World.Where != world.LocationResidentialArea {
		t.Fatalf("位置不应变化：got=%d", state.World.Where)
	}
}

func TestGoHomeAction_WithHome(t *testing.T) {
	state := newTestGameState()
	state.World.Housing.SetOwned(world.HomeRoom)
	state.World.ChangeLocation(world.LocationResidentialArea)

	act := &GoHomeAction{}
	res := act.Execute(state)

	if !res.Success {
		t.Fatal("有住房时应允许回家")
	}
	if state.World.Where != world.LocationHome {
		t.Fatalf("应到达家：got=%d want=%d", state.World.Where, world.LocationHome)
	}
}

func TestRentBedAction_Success(t *testing.T) {
	state := newTestGameState()
	state.World.Wallet = [6]int{50, 0, 0, 0, 0, 0} // 5000元

	act := &RentBedAction{}
	res := act.Execute(state)

	if !res.Success {
		t.Fatalf("租房应成功：%s", res.Message)
	}
	if !state.World.Housing.IsRented() || state.World.Housing.Level != world.HomeBed {
		t.Fatalf("租房状态错误：status=%d level=%d", state.World.Housing.Status, state.World.Housing.Level)
	}
	if state.World.Housing.RentRemainingSeconds != 7*24*60*60 {
		t.Fatalf("租期错误：got=%d", state.World.Housing.RentRemainingSeconds)
	}
	if state.World.GetWalletTotal() != 5000-1200 {
		t.Fatalf("租金扣款错误：got=%d want=%d", state.World.GetWalletTotal(), 5000-1200)
	}
}

func TestResidentialAreaActions_GoHomeOnlyWhenHasHome(t *testing.T) {
	state := newTestGameState()

	actions := GetActionsForLocation(state, world.LocationResidentialArea)
	if containsActionID(actions, "go_home") {
		t.Fatal("无住房时住宅区不应出现回家行动")
	}

	state.World.Housing.SetOwned(world.HomeRoom)
	actions = GetActionsForLocation(state, world.LocationResidentialArea)
	if !containsActionID(actions, "go_home") {
		t.Fatal("有住房时住宅区应出现回家行动")
	}
}

func TestSleepAtHomeAction_ExpiresRentAndAppliesEffects(t *testing.T) {
	state := newTestGameState()
	state.World.Housing.SetRented(world.HomeBed, 60*60) // 1小时租期

	state.User.SetNutrition("精神振奋", 0)
	state.User.SetNutrition("幸福感", 0)

	act := &SleepAtHomeAction{}
	res := act.Execute(state)

	if !res.Success {
		t.Fatalf("睡觉应成功：%s", res.Message)
	}
	if res.TimeElapsed != 8*60*60 {
		t.Fatalf("TimeElapsed=%d，期望8小时", res.TimeElapsed)
	}
	if state.World.Housing.HasHome() {
		t.Fatal("睡觉跨过租期后，应失去住房")
	}
	if !strings.Contains(res.Message, "租期已到") {
		t.Fatalf("应提示租期到期：%s", res.Message)
	}
	if state.User.GetNutrition("精神振奋") != 35 || state.User.GetNutrition("幸福感") != 18 {
		t.Fatalf("睡觉恢复值错误：精神振奋=%d 幸福感=%d", state.User.GetNutrition("精神振奋"), state.User.GetNutrition("幸福感"))
	}
	if state.User.GetNutrition("饱腹感") != 82 || state.User.GetNutrition("饥渴") != 78 {
		t.Fatalf("睡觉消耗错误：饱腹感=%d 饥渴=%d", state.User.GetNutrition("饱腹感"), state.User.GetNutrition("饥渴"))
	}
}

func TestRenewRentAction_ExtendsRent(t *testing.T) {
	state := newTestGameState()
	state.World.Wallet = [6]int{50, 0, 0, 0, 0, 0} // 5000元
	state.World.Housing.SetRented(world.HomeBed, 24*60*60)

	act := &RenewRentAction{}
	res := act.Execute(state)

	if !res.Success {
		t.Fatalf("续租应成功：%s", res.Message)
	}
	wantRemaining := (24*60*60 - 30*60) + 7*24*60*60
	if state.World.Housing.RentRemainingSeconds != wantRemaining {
		t.Fatalf("续租后租期错误：got=%d want=%d", state.World.Housing.RentRemainingSeconds, wantRemaining)
	}
	if state.World.GetWalletTotal() != 5000-1200 {
		t.Fatalf("续租扣款错误：got=%d want=%d", state.World.GetWalletTotal(), 5000-1200)
	}
	if res.TimeElapsed != 30*60 {
		t.Fatalf("续租耗时错误：got=%d", res.TimeElapsed)
	}
}

func TestRealEstateAgencyActions_ByHousingState(t *testing.T) {
	t.Run("无住房", func(t *testing.T) {
		state := newTestGameState()
		actions := GetActionsForLocation(state, world.LocationRealEstateAgency)
		if containsActionID(actions, "renew_rent") || containsActionID(actions, "cancel_rent") {
			t.Fatal("无住房不应出现续租/退租")
		}
		if !containsAllActionIDs(actions, []string{"rent_bed", "rent_room", "rent_apartment", "buy_room", "buy_apartment"}) {
			t.Fatal("无住房时应提供租/买选项")
		}
	})

	t.Run("租房中", func(t *testing.T) {
		state := newTestGameState()
		state.World.Housing.SetRented(world.HomeRoom, 2*24*60*60)
		actions := GetActionsForLocation(state, world.LocationRealEstateAgency)
		if !containsAllActionIDs(actions, []string{"renew_rent", "cancel_rent"}) {
			t.Fatal("租房中应提供续租/退租")
		}
	})

	t.Run("自有单间可升级", func(t *testing.T) {
		state := newTestGameState()
		state.World.Housing.SetOwned(world.HomeRoom)
		actions := GetActionsForLocation(state, world.LocationRealEstateAgency)
		if !containsActionID(actions, "buy_apartment") {
			t.Fatal("自有单间应允许升级小公寓")
		}
		if containsActionID(actions, "rent_bed") || containsActionID(actions, "buy_room") {
			t.Fatal("自有住房不应提供租房/买同档位")
		}
	})
}

func containsActionID(actions []Action, id string) bool {
	for _, a := range actions {
		if a.ID() == id {
			return true
		}
	}
	return false
}

func containsAllActionIDs(actions []Action, ids []string) bool {
	for _, id := range ids {
		if !containsActionID(actions, id) {
			return false
		}
	}
	return true
}
