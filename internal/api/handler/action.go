package handler

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"

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
	actions := action.GetActionsForLocation(state.World.Where)
	sortedActions := action.SortActions(actions)

	result := make([]response.ActionInfo, len(sortedActions))
	for i, act := range sortedActions {
		result[i] = response.ActionInfo{
			ID:       act.ID(),
			Name:     act.Info(),
			Category: act.Category().String(),
		}
	}

	response.Success(c, result)
}

// ExecuteActionRequest 执行行动请求
type ExecuteActionRequest struct {
	ActionID string `json:"action_id" binding:"required"`
}

// ExecuteAction 执行指定行动
func ExecuteAction(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	var req ExecuteActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	state := sess.State

	// 检查游戏是否结束
	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	// 查找行动
	actions := action.GetActionsForLocation(state.World.Where)
	var targetAction action.Action
	for _, act := range actions {
		if act.ID() == req.ActionID {
			targetAction = act
			break
		}
	}

	if targetAction == nil {
		response.BadRequest(c, "无效的行动ID")
		return
	}

	// 执行行动
	result := executor.Execute(targetAction, state)

	// 构建响应
	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	if result.GameOver {
		resp.Message += "\n\n游戏结束: " + result.DeathCause
	}

	response.Success(c, resp)
}
