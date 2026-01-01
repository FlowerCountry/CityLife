package action

import "testing"

func TestWorkAction_AddsIncomeAndConsumesNutrition(t *testing.T) {
	state := newTestGameState()
	state.World.Wallet = [6]int{1, 0, 0, 0, 0, 0} // 100元

	act := &WorkAction{
		id:      "work_test",
		name:    "测试工作",
		seconds: 2 * 60 * 60,
		income:  300,
		consume: map[string]int{"饱腹感": 10, "饥渴": 12},
	}

	res := act.Execute(state)
	if !res.Success {
		t.Fatalf("工作应成功：%s", res.Message)
	}
	if state.World.GetWalletTotal() != 400 {
		t.Fatalf("工资发放错误：got=%d want=%d", state.World.GetWalletTotal(), 400)
	}
	if state.World.Hour != 13 || state.World.Minute != 3 {
		t.Fatalf("时间推进错误：got=%02d:%02d want=13:03", state.World.Hour, state.World.Minute)
	}
	if state.User.GetNutrition("饱腹感") != 90 || state.User.GetNutrition("饥渴") != 88 {
		t.Fatalf("工作消耗错误：饱腹感=%d 饥渴=%d", state.User.GetNutrition("饱腹感"), state.User.GetNutrition("饥渴"))
	}
}

func TestWorkAction_InvalidIncomeRejected(t *testing.T) {
	state := newTestGameState()
	beforeWallet := state.World.GetWalletTotal()
	beforeHour := state.World.Hour
	beforeMinute := state.World.Minute

	act := &WorkAction{id: "bad", name: "坏配置", seconds: 60, income: 350}
	res := act.Execute(state)

	if res.Success {
		t.Fatal("收入不是100倍数时应失败")
	}
	if state.World.GetWalletTotal() != beforeWallet {
		t.Fatal("失败时不应改动钱包")
	}
	if state.World.Hour != beforeHour || state.World.Minute != beforeMinute {
		t.Fatal("失败时不应推进时间")
	}
}

func TestPayAndApplyAction_InsufficientFundsNoTimeChange(t *testing.T) {
	state := newTestGameState()
	state.World.SpendMoney(state.World.GetWalletTotal()) // 清空

	beforeHour := state.World.Hour
	beforeMinute := state.World.Minute

	act := &PayAndApplyAction{
		id:      "test_pay",
		name:    "测试消费",
		price:   300,
		seconds: 60 * 60,
		add:     map[string]int{"幸福感": 10},
	}

	res := act.Execute(state)
	if res.Success {
		t.Fatal("现金不足应失败")
	}
	if state.World.Hour != beforeHour || state.World.Minute != beforeMinute {
		t.Fatal("失败时不应推进时间")
	}
}

func TestPayAndApplyAction_SpendsMoneyAndAppliesNutrition(t *testing.T) {
	state := newTestGameState()
	state.World.Wallet = [6]int{10, 0, 0, 0, 0, 0} // 1000元
	state.User.SetNutrition("饱腹感", 40)
	state.User.SetNutrition("幸福感", 0)

	act := &PayAndApplyAction{
		id:      "test_eat",
		name:    "测试吃饭",
		price:   300,
		seconds: 60 * 60,
		add:     map[string]int{"饱腹感": 55, "幸福感": 20},
	}

	res := act.Execute(state)
	if !res.Success {
		t.Fatalf("消费应成功：%s", res.Message)
	}
	if state.World.GetWalletTotal() != 700 {
		t.Fatalf("扣款错误：got=%d want=%d", state.World.GetWalletTotal(), 700)
	}
	if state.User.GetNutrition("饱腹感") != 95 || state.User.GetNutrition("幸福感") != 20 {
		t.Fatalf("营养/心情效果错误：饱腹感=%d 幸福感=%d", state.User.GetNutrition("饱腹感"), state.User.GetNutrition("幸福感"))
	}
	if res.TimeElapsed != 60*60 {
		t.Fatalf("TimeElapsed错误：got=%d", res.TimeElapsed)
	}
}
