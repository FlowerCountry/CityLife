package v2

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

type leisureInfo struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Price           int            `json:"price"`
	DurationSeconds int            `json:"duration_seconds"`
	Add             map[string]int `json:"add,omitempty"`
	Consume         map[string]int `json:"consume,omitempty"`
}

func buildLeisureInfos(actions []action.Action) []leisureInfo {
	items := make([]leisureInfo, 0, len(actions))
	for _, act := range actions {
		item, ok := act.(*action.PayAndApplyAction)
		if !ok {
			continue
		}
		items = append(items, leisureInfo{
			ID:              item.ID(),
			Name:            item.Name(),
			Price:           item.Price(),
			DurationSeconds: item.DurationSeconds(),
			Add:             item.Add(),
			Consume:         item.Consume(),
		})
	}
	return items
}

// GetRestaurantMenu 获取餐馆菜单
func GetRestaurantMenu(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, buildLeisureInfos(action.GetRestaurantActions()))
}

// DoRestaurant 执行餐馆消费（需要在餐馆）
func DoRestaurant(c *gin.Context) {
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

	if state.World.Where != world.LocationRestaurant {
		fail(c, response.ErrCodeInvalidLocation, "你需要在餐馆才能消费", actionData{State: buildState(state)})
		return
	}

	id := c.Param("id")
	var target action.Action
	for _, act := range action.GetRestaurantActions() {
		if act.ID() == id {
			target = act
			break
		}
	}
	if target == nil {
		response.BadRequest(c, "无效的id参数")
		return
	}

	executeAndRespond(c, state, target, "")
}

// GetParkActivities 获取公园活动列表
func GetParkActivities(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, buildLeisureInfos(action.GetParkActions()))
}

// DoPark 执行公园活动（需要在公园）
func DoPark(c *gin.Context) {
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

	if state.World.Where != world.LocationPark {
		fail(c, response.ErrCodeInvalidLocation, "你需要在公园才能进行该活动", actionData{State: buildState(state)})
		return
	}

	id := c.Param("id")
	var target action.Action
	for _, act := range action.GetParkActions() {
		if act.ID() == id {
			target = act
			break
		}
	}
	if target == nil {
		response.BadRequest(c, "无效的id参数")
		return
	}

	executeAndRespond(c, state, target, "")
}

// GetHotelServices 获取旅馆服务列表
func GetHotelServices(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, buildLeisureInfos(action.GetHotelActions()))
}

// DoHotel 执行旅馆服务（需要在旅馆）
func DoHotel(c *gin.Context) {
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

	if state.World.Where != world.LocationHotel {
		fail(c, response.ErrCodeInvalidLocation, "你需要在旅馆才能使用该服务", actionData{State: buildState(state)})
		return
	}

	id := c.Param("id")
	var target action.Action
	for _, act := range action.GetHotelActions() {
		if act.ID() == id {
			target = act
			break
		}
	}
	if target == nil {
		response.BadRequest(c, "无效的id参数")
		return
	}

	executeAndRespond(c, state, target, "")
}
