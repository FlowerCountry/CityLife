package handler

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"

	"github.com/gin-gonic/gin"
)

// Deposit 存款100元
func Deposit(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	act := &action.DepositAction{}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}

// DepositAll 存入全部现金
func DepositAll(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	act := &action.DepositAllAction{}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}

// Withdraw 取款100元
func Withdraw(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	act := &action.WithdrawAction{}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}

// WithdrawAll 取出全部存款
func WithdrawAll(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	act := &action.WithdrawAllAction{}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}
