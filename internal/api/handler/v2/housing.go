package v2

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

// GetHousing 获取住房状态（结构化）
func GetHousing(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	h := state.World.Housing

	status := "none"
	rentRemaining := 0
	name := ""
	level := 0
	if h.HasHome() {
		name = h.HomeName()
		level = int(h.Level)
		if h.IsOwned() {
			status = "owned"
		} else if h.IsRented() {
			status = "rented"
			rentRemaining = h.RentRemainingSeconds
		}
	}

	ok(c, response.HousingInfo{
		Status:               status,
		Level:                level,
		Name:                 name,
		RentRemainingSeconds: rentRemaining,
	})
}

type HousingOffer struct {
	Level           string `json:"level"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
}

type HousingOffersResponse struct {
	CanRenew  bool           `json:"can_renew"`
	CanCancel bool           `json:"can_cancel"`
	Rent      []HousingOffer `json:"rent"`
	Buy       []HousingOffer `json:"buy"`
}

// GetHousingOffers 获取租房/买房可选项（结构化）
func GetHousingOffers(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	h := sess.State.World.Housing
	resp := HousingOffersResponse{
		CanRenew:  h.IsRented(),
		CanCancel: h.IsRented(),
	}

	if !h.IsOwned() {
		const rentDurationSeconds = 7 * 24 * 60 * 60
		resp.Rent = []HousingOffer{
			{Level: "bed", Name: "合租床位", Price: 1200, DurationSeconds: rentDurationSeconds},
			{Level: "room", Name: "单间", Price: 2400, DurationSeconds: rentDurationSeconds},
			{Level: "apartment", Name: "小公寓", Price: 4000, DurationSeconds: rentDurationSeconds},
		}
	}

	if !h.IsOwned() {
		resp.Buy = []HousingOffer{
			{Level: "room", Name: "单间", Price: 20000},
			{Level: "apartment", Name: "小公寓", Price: 50000},
		}
	} else if h.Level == world.HomeRoom {
		resp.Buy = []HousingOffer{
			{Level: "apartment", Name: "小公寓", Price: 50000},
		}
	}

	ok(c, resp)
}

// RentHousing 租房（需要在房产中介）
func RentHousing(c *gin.Context) {
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

	if state.World.Where != world.LocationRealEstateAgency {
		fail(c, response.ErrCodeInvalidLocation, "你需要在房产中介才能办理租房", actionData{State: buildState(state)})
		return
	}

	level := c.Param("level")
	var act action.Action
	switch level {
	case "bed":
		act = &action.RentBedAction{}
	case "room":
		act = &action.RentRoomAction{}
	case "apartment":
		act = &action.RentApartmentAction{}
	default:
		response.BadRequest(c, "无效的level参数，可选: bed|room|apartment")
		return
	}

	executeAndRespond(c, state, act, "")
}

// BuyHousing 买房/升级（需要在房产中介）
func BuyHousing(c *gin.Context) {
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

	if state.World.Where != world.LocationRealEstateAgency {
		fail(c, response.ErrCodeInvalidLocation, "你需要在房产中介才能办理购房", actionData{State: buildState(state)})
		return
	}

	level := c.Param("level")
	var act action.Action
	switch level {
	case "room":
		act = &action.BuyRoomAction{}
	case "apartment":
		act = &action.BuyApartmentAction{}
	default:
		response.BadRequest(c, "无效的level参数，可选: room|apartment")
		return
	}

	executeAndRespond(c, state, act, "")
}

// RenewRent 续租（需要在房产中介）
func RenewRent(c *gin.Context) {
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

	if state.World.Where != world.LocationRealEstateAgency {
		fail(c, response.ErrCodeInvalidLocation, "你需要在房产中介才能办理续租", actionData{State: buildState(state)})
		return
	}

	executeAndRespond(c, state, &action.RenewRentAction{}, "")
}

// CancelRent 退租（需要在房产中介）
func CancelRent(c *gin.Context) {
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

	if state.World.Where != world.LocationRealEstateAgency {
		fail(c, response.ErrCodeInvalidLocation, "你需要在房产中介才能办理退租", actionData{State: buildState(state)})
		return
	}

	executeAndRespond(c, state, &action.CancelRentAction{}, "")
}

// GoHome 回家（需要在住宅区）
func GoHome(c *gin.Context) {
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

	if state.World.Where != world.LocationResidentialArea {
		fail(c, response.ErrCodeInvalidLocation, "你需要先到住宅区才能回家", actionData{State: buildState(state)})
		return
	}

	executeAndRespond(c, state, &action.GoHomeAction{}, "")
}

// LeaveHome 出门（需要在家）
func LeaveHome(c *gin.Context) {
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

	if state.World.Where != world.LocationHome {
		fail(c, response.ErrCodeInvalidLocation, "你需要在家里才能出门", actionData{State: buildState(state)})
		return
	}

	executeAndRespond(c, state, &action.LeaveHomeAction{}, "")
}

// SleepAtHome 睡觉（需要在家）
func SleepAtHome(c *gin.Context) {
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

	if state.World.Where != world.LocationHome {
		fail(c, response.ErrCodeInvalidLocation, "你需要在家里才能睡觉", actionData{State: buildState(state)})
		return
	}

	executeAndRespond(c, state, &action.SleepAtHomeAction{}, "")
}
