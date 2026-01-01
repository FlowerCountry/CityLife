package v2

import (
	"errors"
	"io"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

type bankAmountRequest struct {
	Amount *int `json:"amount"`
}

// Deposit 存款（默认100元，支持amount）
func Deposit(c *gin.Context) {
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

	if state.World.Where != world.LocationBank {
		fail(c, response.ErrCodeInvalidLocation, "你需要在银行才能办理存款", actionData{State: buildState(state)})
		return
	}

	amount := 100
	var req bankAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.BadRequest(c, "无效的请求参数")
		return
	}
	if req.Amount != nil {
		amount = *req.Amount
	}
	if amount <= 0 {
		response.Error(c, 400, response.ErrCodeInvalidAmount, "amount 必须大于0")
		return
	}
	if amount%100 != 0 {
		response.Error(c, 400, response.ErrCodeInvalidAmount, "金额必须是100元的倍数")
		return
	}

	failCode := ""
	if state.World.GetWalletTotal() < amount {
		failCode = response.ErrCodeInsufficientFunds
	}

	executeAndRespond(c, state, &action.DepositAction{Amount: amount}, failCode)
}

// Withdraw 取款（默认100元，支持amount）
func Withdraw(c *gin.Context) {
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

	if state.World.Where != world.LocationBank {
		fail(c, response.ErrCodeInvalidLocation, "你需要在银行才能办理取款", actionData{State: buildState(state)})
		return
	}

	amount := 100
	var req bankAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.BadRequest(c, "无效的请求参数")
		return
	}
	if req.Amount != nil {
		amount = *req.Amount
	}
	if amount <= 0 {
		response.Error(c, 400, response.ErrCodeInvalidAmount, "amount 必须大于0")
		return
	}
	if amount%100 != 0 {
		response.Error(c, 400, response.ErrCodeInvalidAmount, "金额必须是100元的倍数")
		return
	}

	failCode := ""
	if state.World.BankDeposit < amount {
		failCode = response.ErrCodeInsufficientFunds
	}

	executeAndRespond(c, state, &action.WithdrawAction{Amount: amount}, failCode)
}

// DepositAll 存入全部现金
func DepositAll(c *gin.Context) {
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

	if state.World.Where != world.LocationBank {
		fail(c, response.ErrCodeInvalidLocation, "你需要在银行才能办理存款", actionData{State: buildState(state)})
		return
	}

	failCode := ""
	if state.World.GetWalletTotal() < 100 {
		failCode = response.ErrCodeInsufficientFunds
	}
	executeAndRespond(c, state, &action.DepositAllAction{}, failCode)
}

// WithdrawAll 取出全部存款
func WithdrawAll(c *gin.Context) {
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

	if state.World.Where != world.LocationBank {
		fail(c, response.ErrCodeInvalidLocation, "你需要在银行才能办理取款", actionData{State: buildState(state)})
		return
	}

	failCode := ""
	if state.World.BankDeposit < 100 {
		failCode = response.ErrCodeInsufficientFunds
	}
	executeAndRespond(c, state, &action.WithdrawAllAction{}, failCode)
}
