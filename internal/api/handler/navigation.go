package handler

import (
	"strconv"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

// Navigate 导航到指定位置
func Navigate(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	// 检查游戏是否结束
	if !state.User.IsAlive() {
		response.GameOver(c, "游戏已结束，角色已死亡")
		return
	}

	// 获取目标位置
	locationIDStr := c.Param("location_id")
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		response.BadRequest(c, "无效的位置ID")
		return
	}

	// 检查位置是否有效
	if locationID < 0 || locationID >= len(world.BuildingNames) {
		response.BadRequest(c, "位置不存在")
		return
	}

	// 检查是否可以到达
	adjacentBuildings := world.GetAdjacentBuildings(state.World.Where)
	canReach := false
	for _, adj := range adjacentBuildings {
		if adj == locationID {
			canReach = true
			break
		}
	}

	if !canReach {
		response.Conflict(c, response.ErrCodeInvalidLocation, "无法从当前位置到达目标位置")
		return
	}

	// 创建并执行行动
	act := &action.GoWhereAction{
		TargetID:   locationID,
		TargetName: world.BuildingNames[locationID],
	}

	result := executor.Execute(act, state)

	resp := response.ActionResponse{
		Message:  result.ActionResult.Message,
		Success:  result.ActionResult.Success,
		Hints:    result.Hints,
		GameOver: result.GameOver,
	}

	response.Success(c, resp)
}
