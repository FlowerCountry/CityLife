package v2

import (
	"errors"
	"io"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"

	"github.com/gin-gonic/gin"
)

// GetActions 获取当前位置可用的行动列表
func GetActions(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	actions := action.GetActionsForLocation(state, state.World.Where)
	sortedActions := action.SortActions(actions)

	result := make([]response.ActionInfo, len(sortedActions))
	for i, act := range sortedActions {
		result[i] = response.ActionInfo{
			ID:       act.ID(),
			Name:     act.Info(),
			Category: act.Category().String(),
		}
	}

	ok(c, result)
}

type executeActionRequest struct {
	Amount *int `json:"amount"`
}

// ExecuteAction 执行指定行动（v2：使用path参数，不再需要action_id字段）
func ExecuteAction(c *gin.Context) {
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

	actionID := c.Param("action_id")
	if actionID == "" {
		response.BadRequest(c, "缺少action_id参数")
		return
	}

	actions := action.GetActionsForLocation(state, state.World.Where)
	var target action.Action
	for _, act := range actions {
		if act.ID() == actionID {
			target = act
			break
		}
	}

	if target == nil {
		response.NotFound(c, response.ErrCodeInvalidAction, "当前地点不可用的行动")
		return
	}

	var req executeActionRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	failCode := ""
	if req.Amount != nil {
		amount := *req.Amount
		if amount <= 0 {
			response.Error(c, 400, response.ErrCodeInvalidAmount, "amount 必须大于0")
			return
		}
		if amount%100 != 0 {
			response.Error(c, 400, response.ErrCodeInvalidAmount, "金额必须是100元的倍数")
			return
		}

		switch a := target.(type) {
		case *action.DepositAction:
			walletTotal := state.World.GetWalletTotal()
			if walletTotal < amount {
				failCode = response.ErrCodeInsufficientFunds
			}
			a.Amount = amount
		case *action.WithdrawAction:
			if state.World.BankDeposit < amount {
				failCode = response.ErrCodeInsufficientFunds
			}
			a.Amount = amount
		default:
			response.BadRequest(c, "该行动不支持amount参数")
			return
		}
	}

	executeAndRespond(c, state, target, failCode)
}
