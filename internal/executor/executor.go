// Package executor 封装行动执行逻辑
package executor

import (
	"citylife/internal/action"
	"citylife/internal/game"
)

// ExecutionResult 执行结果
type ExecutionResult struct {
	ActionResult *action.Result
	Hints        []string
	GameOver     bool
	DeathCause   string
}

// Execute 执行行动并处理后续效果
func Execute(act action.Action, state *game.State) *ExecutionResult {
	// 1. 执行行动
	result := act.Execute(state)

	// 2. 获取疾病影响的倍率
	decayMultiplier := state.Disease.GetNutritionDecayMultiplier()

	// 3. 执行营养衰减
	state.Nutrition.OnAction(state.User, decayMultiplier)
	state.Nutrition.ProcessInteractions(state.User)
	state.Nutrition.CheckLevelChanges(state.User)

	// 4. 疾病系统检查
	state.Disease.OnAction(state.User)

	// 5. 收集营养提示
	hints := state.Nutrition.GetPendingHints()
	state.Nutrition.ClearPendingHints()

	// 6. 额外营养消耗（疾病惩罚）
	actionCostMultiplier := state.Disease.GetActionCostMultiplier()
	if actionCostMultiplier > 1.0 {
		extraCost := int((actionCostMultiplier - 1.0) * 5)
		if extraCost > 0 {
			state.User.ConsumeNutrition("饱腹感", extraCost)
			state.User.ConsumeNutrition("饥渴", extraCost)
		}
	}

	// 7. 检查死亡
	execResult := &ExecutionResult{
		ActionResult: result,
		Hints:        hints,
		GameOver:     false,
	}

	if !state.User.IsAlive() {
		execResult.GameOver = true
		execResult.DeathCause = "营养不足"
	} else if state.Nutrition.CheckRandomDeath(state.User) {
		execResult.GameOver = true
		execResult.DeathCause = "随机死亡"
	}

	return execResult
}

// ExecuteWithoutDecay 执行行动但不进行营养衰减（用于查看类行动）
func ExecuteWithoutDecay(act action.Action, state *game.State) *ExecutionResult {
	result := act.Execute(state)

	return &ExecutionResult{
		ActionResult: result,
		Hints:        nil,
		GameOver:     false,
	}
}
