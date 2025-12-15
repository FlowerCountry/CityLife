package handler

import (
	"strconv"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"

	"github.com/gin-gonic/gin"
)

// GetSaveSlots 获取存档槽位列表
func GetSaveSlots(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	slots, err := action.GetSaveSlots()
	if err != nil {
		response.InternalError(c, "获取存档列表失败: "+err.Error())
		return
	}

	result := make([]response.SaveSlotInfo, len(slots))
	for i, slot := range slots {
		result[i] = response.SaveSlotInfo{
			Slot:        slot.Slot,
			Exists:      slot.Exists,
			Description: slot.Description,
		}
	}

	response.Success(c, result)
}

// SaveGame 保存游戏到指定槽位
func SaveGame(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	slotStr := c.Param("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil || slot < 1 || slot > 3 {
		response.BadRequest(c, "无效的槽位号，有效范围: 1-3")
		return
	}

	act := &action.SaveGameAction{Slot: slot}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}

// LoadGame 从指定槽位加载游戏
func LoadGame(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	slotStr := c.Param("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil || slot < 1 || slot > 3 {
		response.BadRequest(c, "无效的槽位号，有效范围: 1-3")
		return
	}

	act := &action.LoadGameAction{Slot: slot}
	result := executor.ExecuteWithoutDecay(act, sess.State)

	response.Success(c, response.ActionResponse{
		Message: result.ActionResult.Message,
		Success: result.ActionResult.Success,
	})
}
