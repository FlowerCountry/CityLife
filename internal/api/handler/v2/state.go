package v2

import (
	"fmt"

	"citylife/internal/api/middleware"
	"citylife/internal/api/response"

	"github.com/gin-gonic/gin"
)

// GetState 获取完整游戏状态
func GetState(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, buildState(sess.State))
}

// GetStatus 获取状态摘要
func GetStatus(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	locationName := "未知"
	if current := state.World.CurrentLocation(); current != nil {
		locationName = current.Name
	}
	ok(c, gin.H{
		"location":    locationName,
		"location_id": state.World.Where,
		"time":        fmt.Sprintf("%d年%d月%d日 %02d:%02d", state.World.Year, state.World.Month, state.World.Day, state.World.Hour, state.World.Minute),
		"wallet":      state.World.GetWalletTotal(),
		"bank":        state.World.BankDeposit,
		"housing":     state.World.Housing.Summary(),
		"is_alive":    state.User.IsAlive(),
		"diseases":    len(state.Disease.GetActiveDiseaseIDs()),
	})
}

// GetWallet 获取钱包详情
func GetWallet(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	w := sess.State.World
	denominations := map[int]int{
		100: w.Wallet[0],
		50:  w.Wallet[1],
		20:  w.Wallet[2],
		10:  w.Wallet[3],
		5:   w.Wallet[4],
		1:   w.Wallet[5],
	}

	ok(c, response.WalletInfo{
		Total:         w.GetWalletTotal(),
		Denominations: denominations,
	})
}

// GetBankBalance 获取银行余额
func GetBankBalance(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, gin.H{
		"balance": sess.State.World.BankDeposit,
	})
}

// GetHealth 获取健康状态
func GetHealth(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	coreNutrients := []string{"饱腹感", "饥渴", "蛋白质", "碳水化合物"}
	core := make(map[string]response.NutritionInfo)
	for _, name := range coreNutrients {
		value := state.User.GetNutrition(name)
		level := state.User.GetNutritionLevel(name)
		core[name] = response.NutritionInfo{
			Value:  value,
			Level:  level,
			Status: getLevelStatus(level),
		}
	}

	ok(c, gin.H{
		"core":     core,
		"all":      state.User.GetAllNutrition(),
		"is_alive": state.User.IsAlive(),
	})
}

// GetDiseases 获取疾病状态
func GetDiseases(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	ok(c, gin.H{
		"diseases":  getDiseases(sess.State),
		"diagnoses": sess.State.User.GetDiagnoses(),
	})
}

// GetMap 获取地图
func GetMap(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	locations := make([]gin.H, 0, len(state.World.Buildings))
	for _, b := range state.World.Buildings {
		if b == nil || b.IsInterior {
			continue
		}
		locations = append(locations, gin.H{
			"id":   b.ID,
			"name": b.Name,
			"x":    b.X,
			"y":    b.Y,
		})
	}

	ok(c, gin.H{
		"current_location": state.World.Where,
		"locations":        locations,
	})
}

func getLevelStatus(level int) string {
	switch level {
	case 0:
		return "安全"
	case 1:
		return "警告"
	case 2:
		return "危险"
	case 3:
		return "透支"
	case 4:
		return "濒死"
	default:
		return "未知"
	}
}
