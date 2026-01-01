package v2

import (
	"fmt"
	"net/http"

	"citylife/internal/action"
	"citylife/internal/api/response"
	"citylife/internal/disease"
	"citylife/internal/executor"
	"citylife/internal/game"

	"github.com/gin-gonic/gin"
)

type actionData struct {
	Message    string                 `json:"message"`
	Hints      []string               `json:"hints,omitempty"`
	GameOver   bool                   `json:"game_over,omitempty"`
	DeathCause string                 `json:"death_cause,omitempty"`
	State      response.StateResponse `json:"state"`
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    data,
	})
}

func fail(c *gin.Context, code, message string, data interface{}) {
	c.JSON(statusForCode(code), response.Response{
		Success: false,
		Data:    data,
		Error: &response.ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func buildState(state *game.State) response.StateResponse {
	locationName := "未知"
	if current := state.World.CurrentLocation(); current != nil {
		locationName = current.Name
	}

	housingStatus := "none"
	rentRemaining := 0
	housingName := ""
	housingLevel := 0
	if state.World.Housing.HasHome() {
		housingName = state.World.Housing.HomeName()
		housingLevel = int(state.World.Housing.Level)
		if state.World.Housing.IsOwned() {
			housingStatus = "owned"
		} else if state.World.Housing.IsRented() {
			housingStatus = "rented"
			rentRemaining = state.World.Housing.RentRemainingSeconds
		}
	}

	return response.StateResponse{
		Location: response.LocationInfo{
			ID:   state.World.Where,
			Name: locationName,
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
		Housing: response.HousingInfo{
			Status:               housingStatus,
			Level:                housingLevel,
			Name:                 housingName,
			RentRemainingSeconds: rentRemaining,
		},
		Health:   state.User.GetAllNutrition(),
		Diseases: getDiseases(state),
		IsAlive:  state.User.IsAlive(),
	}
}

func getDiseases(state *game.State) []response.DiseaseInfo {
	diseaseIDs := state.Disease.GetActiveDiseaseIDs()
	result := make([]response.DiseaseInfo, 0, len(diseaseIDs))

	for _, id := range diseaseIDs {
		info := state.Disease.GetDiseaseInfo(id)
		if info == nil {
			continue
		}
		result = append(result, response.DiseaseInfo{
			ID:            id,
			Name:          info.Name,
			Description:   info.Description,
			Severity:      disease.SeverityToString(info.Severity),
			TreatmentCost: info.TreatmentCost,
		})
	}

	return result
}

func executeAndRespond(c *gin.Context, state *game.State, act action.Action, failCode string) {
	result := executor.Execute(act, state)

	data := actionData{
		Message:    result.ActionResult.Message,
		Hints:      result.Hints,
		GameOver:   result.GameOver,
		DeathCause: result.DeathCause,
		State:      buildState(state),
	}

	if result.ActionResult.Success {
		ok(c, data)
		return
	}

	if failCode == "" {
		failCode = response.ErrCodeInvalidAction
	}
	fail(c, failCode, result.ActionResult.Message, data)
}

func failGameOver(c *gin.Context, state *game.State) {
	fail(c, response.ErrCodeGameOver, "游戏已结束，角色已死亡", actionData{
		Message:  "游戏已结束，角色已死亡",
		GameOver: true,
		State:    buildState(state),
	})
}

func statusForCode(code string) int {
	switch code {
	case response.ErrCodeSessionNotFound, response.ErrCodeSaveNotFound:
		return http.StatusNotFound
	case response.ErrCodeInvalidRequest,
		response.ErrCodeInvalidAction,
		response.ErrCodeInvalidAmount,
		response.ErrCodeInvalidLocation,
		response.ErrCodeInvalidSlot:
		return http.StatusBadRequest
	case response.ErrCodeInsufficientFunds,
		response.ErrCodeGameOver,
		response.ErrCodePrescriptionRequired:
		return http.StatusConflict
	case response.ErrCodeInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
