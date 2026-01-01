// Package action 包含所有游戏行动实现
package action

import (
	"fmt"
	"math/rand"
	"strings"

	"citylife/internal/game"
	"citylife/internal/payment"
	"citylife/internal/world"
)

// GoWhereAction 前往行动
type GoWhereAction struct {
	TargetID   int
	TargetName string
}

func (a *GoWhereAction) ID() string {
	return fmt.Sprintf("go_where_%d", a.TargetID)
}

func (a *GoWhereAction) Info() string {
	return "前往" + a.TargetName
}

func (a *GoWhereAction) Execute(state *game.State) *Result {
	// 计算距离和时间
	distance := world.GetDistance(state.World.Where, a.TargetID)
	timeMinutes := distance

	// 更新位置和时间
	state.World.ChangeLocation(a.TargetID)
	state.World.UpdateTime(timeMinutes * 60)

	// 消耗营养（根据距离）
	state.User.ConsumeNutrition("饱腹感", distance*2)
	state.User.ConsumeNutrition("饥渴", int(float64(distance)*1.5))
	state.User.ConsumeNutrition("蛋白质", distance/2)
	state.User.ConsumeNutrition("碳水化合物", distance/2)

	return &Result{
		Message:     fmt.Sprintf("你到达了%s（用时%d分钟）", a.TargetName, timeMinutes),
		Success:     true,
		TimeElapsed: timeMinutes * 60,
	}
}

func (a *GoWhereAction) Category() EventCategory {
	return CategoryNavigation
}

// CheckCashAction 查看现金
type CheckCashAction struct{}

func (a *CheckCashAction) ID() string {
	return "check_cash"
}

func (a *CheckCashAction) Info() string {
	return "查看钱包"
}

func (a *CheckCashAction) Execute(state *game.State) *Result {
	w := state.World

	// 使用WalletInspector生成钱包信息
	inspector := payment.NewWalletInspector(state.Locale)
	total := w.GetWalletTotal()

	var msg strings.Builder

	// 检查是否为空钱包
	if total == 0 {
		msg.WriteString(inspector.PickZeroWalletMessage())
	} else {
		msg.WriteString(inspector.PickSummaryLead())
		msg.WriteString("\n")

		// 生成摘要
		summary := inspector.BuildSummary(w.Wallet)
		for _, line := range summary.Lines {
			msg.WriteString("  " + line + "\n")
		}

		if summary.HasMore {
			msg.WriteString("\n")
			msg.WriteString(inspector.PickHasMoreMessage())
		}

		msg.WriteString("\n")
		msg.WriteString(fmt.Sprintf("总计: %d元\n", total))
		msg.WriteString("\n")
		msg.WriteString(inspector.PickCloseMessage())
	}

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *CheckCashAction) Category() EventCategory {
	return CategoryInsight
}

// CheckBankAction 查看银行余额
type CheckBankAction struct{}

func (a *CheckBankAction) ID() string {
	return "check_bank"
}

func (a *CheckBankAction) Info() string {
	return "查看银行余额"
}

func (a *CheckBankAction) Execute(state *game.State) *Result {
	return &Result{
		Message: fmt.Sprintf("══════ 银行账户 ══════\n  存款余额: %d元\n══════════════════════", state.World.BankDeposit),
		Success: true,
	}
}

func (a *CheckBankAction) Category() EventCategory {
	return CategoryInsight
}

// EnterSupermarketAction 进入超市
type EnterSupermarketAction struct{}

func (a *EnterSupermarketAction) ID() string {
	return "enter_supermarket"
}

func (a *EnterSupermarketAction) Info() string {
	return "进入超市购物"
}

func (a *EnterSupermarketAction) Execute(state *game.State) *Result {
	state.World.ChangeLocation(world.LocationSupermarketInner)

	return &Result{
		Message: "欢迎光临超市！请选择您要购买的商品。",
		Success: true,
	}
}

func (a *EnterSupermarketAction) Category() EventCategory {
	return CategoryPrimary
}

// LeaveSupermarketAction 离开超市
type LeaveSupermarketAction struct{}

func (a *LeaveSupermarketAction) ID() string {
	return "leave_supermarket"
}

func (a *LeaveSupermarketAction) Info() string {
	return "离开超市"
}

func (a *LeaveSupermarketAction) Execute(state *game.State) *Result {
	state.World.ChangeLocation(world.LocationSupermarket)

	return &Result{
		Message: "你离开了超市",
		Success: true,
	}
}

func (a *LeaveSupermarketAction) Category() EventCategory {
	return CategoryNavigation
}

// InformationAction 查看公告
type InformationAction struct{}

// 公告列表
var announcements = []string{
	"于2009年7月,我市第一家银行正式完工",
	"于2009年12月,我市第一家银行正式营业,有需要者可以到银行办理业务",
	"于2010年5月,我市预计开始建设电信大楼",
	"于2010年8月,我市预计正式完工电信大楼",
	"温馨提示：保持良好的饮食习惯，注意补充各种营养",
	"健康小贴士：多喝水，保持身体水分充足",
	"市民须知：银行存取款只能以100元为单位",
}

func (a *InformationAction) ID() string {
	return "information"
}

func (a *InformationAction) Info() string {
	return "查看公告"
}

func (a *InformationAction) Execute(state *game.State) *Result {
	// 随机打乱公告顺序
	shuffled := make([]string, len(announcements))
	copy(shuffled, announcements)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	var msg strings.Builder
	msg.WriteString("══════ 市中心公告 ══════\n\n")
	for i, ann := range shuffled {
		msg.WriteString(fmt.Sprintf("  %d. %s\n", i+1, ann))
	}
	msg.WriteString("\n════════════════════════")

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *InformationAction) Category() EventCategory {
	return CategoryInsight
}

// ViewMapAction 查看地图
type ViewMapAction struct{}

func (a *ViewMapAction) ID() string {
	return "view_map"
}

func (a *ViewMapAction) Info() string {
	return "查看地图"
}

func (a *ViewMapAction) Execute(state *game.State) *Result {
	current := state.World.CurrentLocation()
	if current == nil {
		return &Result{Message: "当前位置未知", Success: false}
	}

	adjacentIDs := world.GetAdjacentBuildings(state.World.Where)
	adjacentNames := make([]string, 0, len(adjacentIDs))
	for _, id := range adjacentIDs {
		if id >= 0 && id < len(world.BuildingNames) {
			adjacentNames = append(adjacentNames, world.BuildingNames[id])
		}
	}

	publicNames := make([]string, 0, len(state.World.Buildings))
	for _, b := range state.World.Buildings {
		if b == nil || b.IsInterior {
			continue
		}
		publicNames = append(publicNames, b.Name)
	}

	var msg strings.Builder
	msg.WriteString("══════ 城市地图 ══════\n\n")
	msg.WriteString(fmt.Sprintf("  当前位置: 【%s】\n", current.Name))
	if len(adjacentNames) > 0 {
		msg.WriteString(fmt.Sprintf("  可前往: %s\n", strings.Join(adjacentNames, "、")))
	} else {
		msg.WriteString("  可前往: 无\n")
	}
	msg.WriteString("\n  公共地点: " + strings.Join(publicNames, "、") + "\n")
	msg.WriteString("\n  提示: 部分室内地点需要通过特定行动进入。\n")
	msg.WriteString("\n══════════════════════")

	return &Result{Message: msg.String(), Success: true}
}

func (a *ViewMapAction) Category() EventCategory {
	return CategoryInsight
}

// ViewHealthAction 查看健康状态
type ViewHealthAction struct{}

func (a *ViewHealthAction) ID() string {
	return "view_health"
}

func (a *ViewHealthAction) Info() string {
	return "查看健康状态"
}

func (a *ViewHealthAction) Execute(state *game.State) *Result {
	u := state.User

	var msg strings.Builder
	msg.WriteString("══════ 健康状态 ══════\n\n")

	// 核心属性
	msg.WriteString("【核心指标】\n")
	coreNutrients := []string{"饱腹感", "饥渴", "蛋白质", "碳水化合物"}
	for _, name := range coreNutrients {
		value := u.GetNutrition(name)
		level := u.GetNutritionLevel(name)
		status := getNutritionStatus(level)
		msg.WriteString(fmt.Sprintf("  %s: %d %s\n", name, value, status))
	}

	// 微量元素
	msg.WriteString("\n【微量元素】\n")
	microNutrients := []string{"钙", "铁", "锌", "维生素C", "维生素B"}
	for _, name := range microNutrients {
		value := u.GetNutrition(name)
		msg.WriteString(fmt.Sprintf("  %s: %d\n", name, value))
	}

	// 特殊属性
	msg.WriteString("\n【心理状态】\n")
	msg.WriteString(fmt.Sprintf("  幸福感: %d\n", u.GetNutrition("幸福感")))
	msg.WriteString(fmt.Sprintf("  精神振奋: %d\n", u.GetNutrition("精神振奋")))

	msg.WriteString("\n══════════════════════")

	return &Result{
		Message: msg.String(),
		Success: true,
	}
}

func (a *ViewHealthAction) Category() EventCategory {
	return CategoryInsight
}

func getNutritionStatus(level int) string {
	switch level {
	case 0:
		return "[安全]"
	case 1:
		return "[警告]"
	case 2:
		return "[危险]"
	case 3:
		return "[透支]"
	case 4:
		return "[濒死]"
	default:
		return ""
	}
}
