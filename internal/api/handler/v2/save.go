package v2

import (
	"strconv"

	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/executor"
	"citylife/internal/save"

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

	ok(c, result)
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
		response.Error(c, 400, response.ErrCodeInvalidSlot, "无效的槽位号，有效范围: 1-3")
		return
	}

	state := sess.State
	result := executor.ExecuteWithoutDecay(&action.SaveGameAction{Slot: slot}, state)
	if !result.ActionResult.Success {
		response.InternalError(c, result.ActionResult.Message)
		return
	}

	ok(c, gin.H{
		"message": result.ActionResult.Message,
		"state":   buildState(state),
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
		response.Error(c, 400, response.ErrCodeInvalidSlot, "无效的槽位号，有效范围: 1-3")
		return
	}

	mgr, err := save.NewManager()
	if err != nil {
		response.InternalError(c, "创建存档管理器失败: "+err.Error())
		return
	}
	if !mgr.Exists(slot) {
		fail(c, response.ErrCodeSaveNotFound, "存档不存在", gin.H{
			"state": buildState(sess.State),
		})
		return
	}

	state := sess.State
	result := executor.ExecuteWithoutDecay(&action.LoadGameAction{Slot: slot}, state)
	if !result.ActionResult.Success {
		response.InternalError(c, result.ActionResult.Message)
		return
	}

	ok(c, gin.H{
		"message": result.ActionResult.Message,
		"state":   buildState(state),
	})
}
