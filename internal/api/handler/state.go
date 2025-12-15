package handler

import (
	"fmt"

	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/disease"
	"citylife/internal/game"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

// GetState 获取完整游戏状态
func GetState(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	// 构建响应
	resp := response.StateResponse{
		Location: response.LocationInfo{
			ID:   state.World.Where,
			Name: state.World.CurrentLocation().Name,
		},
		Time: response.TimeInfo{
			Year:    state.World.Year,
			Month:   state.World.Month,
			Day:     state.World.Day,
			Hour:    state.World.Hour,
			Minute:  state.World.Minute,
			Second:  state.World.Second,
			Display: fmt.Sprintf("%d年%d月%d日 %02d:%02d:%02d", state.World.Year, state.World.Month, state.World.Day, state.World.Hour, state.World.Minute, state.World.Second),
		},
		Money: response.MoneyInfo{
			WalletTotal: state.World.GetWalletTotal(),
			BankDeposit: state.World.BankDeposit,
			Wallet:      state.World.Wallet,
		},
		Health:   state.User.GetAllNutrition(),
		Diseases: getDiseases(state),
		IsAlive:  state.User.IsAlive(),
	}

	response.Success(c, resp)
}

// GetStatus 获取状态摘要
func GetStatus(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State

	response.Success(c, gin.H{
		"location":    state.World.CurrentLocation().Name,
		"location_id": state.World.Where,
		"time":        fmt.Sprintf("%d年%d月%d日 %02d:%02d", state.World.Year, state.World.Month, state.World.Day, state.World.Hour, state.World.Minute),
		"wallet":      state.World.GetWalletTotal(),
		"bank":        state.World.BankDeposit,
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

	state := sess.State
	w := state.World

	denominations := map[int]int{
		100: w.Wallet[0],
		50:  w.Wallet[1],
		20:  w.Wallet[2],
		10:  w.Wallet[3],
		5:   w.Wallet[4],
		1:   w.Wallet[5],
	}

	response.Success(c, response.WalletInfo{
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

	response.Success(c, gin.H{
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

	// 核心营养
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

	// 其他营养
	allNutrition := state.User.GetAllNutrition()

	response.Success(c, gin.H{
		"core":     core,
		"all":      allNutrition,
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

	response.Success(c, gin.H{
		"diseases":   getDiseases(sess.State),
		"diagnoses":  sess.State.User.GetDiagnoses(),
	})
}

// GetMap 获取地图
func GetMap(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	locations := []gin.H{
		{"id": world.LocationCityCenter, "name": "市中心"},
		{"id": world.LocationSupermarket, "name": "超市"},
		{"id": world.LocationBank, "name": "银行"},
		{"id": world.LocationHospital, "name": "医院"},
	}

	response.Success(c, gin.H{
		"current_location": sess.State.World.Where,
		"locations":        locations,
	})
}

// 辅助函数

func getDiseases(state *game.State) []response.DiseaseInfo {
	diseaseIDs := state.Disease.GetActiveDiseaseIDs()
	result := make([]response.DiseaseInfo, 0, len(diseaseIDs))

	for _, id := range diseaseIDs {
		info := state.Disease.GetDiseaseInfo(id)
		if info != nil {
			result = append(result, response.DiseaseInfo{
				ID:            id,
				Name:          info.Name,
				Description:   info.Description,
				Severity:      disease.SeverityToString(info.Severity),
				TreatmentCost: info.TreatmentCost,
			})
		}
	}

	return result
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
