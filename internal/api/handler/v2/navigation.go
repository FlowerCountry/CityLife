package v2

import (
	"strconv"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
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
	if !state.User.IsAlive() {
		failGameOver(c, state)
		return
	}

	locationIDStr := c.Param("location_id")
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		response.BadRequest(c, "无效的位置ID")
		return
	}

	if locationID < 0 || locationID >= len(world.BuildingNames) {
		response.BadRequest(c, "位置不存在")
		return
	}

	adjacentBuildings := world.GetAdjacentBuildings(state.World.Where)
	canReach := false
	for _, adj := range adjacentBuildings {
		if adj == locationID {
			canReach = true
			break
		}
	}
	if !canReach {
		fail(c, response.ErrCodeInvalidLocation, "无法从当前位置到达目标位置", actionData{
			State: buildState(state),
		})
		return
	}

	act := &action.GoWhereAction{
		TargetID:   locationID,
		TargetName: world.BuildingNames[locationID],
	}

	executeAndRespond(c, state, act, "")
}
