package v2

import (
	"net/http"
	"testing"

	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

func TestSetTestMoney(t *testing.T) {
	env := setupTestEnv()

	t.Run("设置钱包与银行余额", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.setTestMoney(sessionID, 1200, 300)
		if w.Code != http.StatusOK {
			t.Fatalf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatal("期望响应成功")
		}

		walletResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/wallet", nil)
		walletData, _ := parseResponse(walletResp)
		wallet, _ := parseDataMap(walletData)
		if wallet["total"].(float64) != 1200 {
			t.Errorf("期望钱包余额为1200，得到 %v", wallet["total"])
		}

		bankResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/bank", nil)
		bankData, _ := parseResponse(bankResp)
		bank, _ := parseDataMap(bankData)
		if bank["balance"].(float64) != 300 {
			t.Errorf("期望银行余额为300，得到 %v", bank["balance"])
		}
	})
}

func TestSetTestMoney_DisabledOutsideTestMode(t *testing.T) {
	originalMode := gin.Mode()
	t.Cleanup(func() {
		gin.SetMode(originalMode)
	})

	env := setupTestEnvWithMode(gin.ReleaseMode)

	t.Run("非测试模式禁用", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.setTestMoney(sessionID, 1000, 0)
		if w.Code != http.StatusForbidden {
			t.Fatalf("期望状态码 %d, 得到 %d", http.StatusForbidden, w.Code)
		}

		resp, _ := parseResponse(w)
		if resp.Success {
			t.Error("非测试模式下接口不应成功")
		}
	})
}

// TestMoneyRequest 测试专用余额设置请求
type TestMoneyRequest struct {
	Wallet *int `json:"wallet"`
	Bank   *int `json:"bank"`
}

// SetTestMoney 测试专用：设置钱包与银行余额
func SetTestMoney(c *gin.Context) {
	if gin.Mode() != gin.TestMode {
		response.Error(c, http.StatusForbidden, response.ErrCodeInvalidRequest, "仅测试模式可用")
		return
	}

	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	var req TestMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	if req.Wallet == nil && req.Bank == nil {
		response.BadRequest(c, "请至少提供 wallet 或 bank")
		return
	}

	if req.Wallet != nil && *req.Wallet < 0 {
		response.BadRequest(c, "wallet 不能为负数")
		return
	}

	if req.Bank != nil && *req.Bank < 0 {
		response.BadRequest(c, "bank 不能为负数")
		return
	}

	if req.Wallet != nil {
		setWalletTotal(sess.State.World, *req.Wallet)
	}

	if req.Bank != nil {
		sess.State.World.BankDeposit = *req.Bank
	}

	response.Success(c, gin.H{
		"wallet_total": sess.State.World.GetWalletTotal(),
		"bank_deposit": sess.State.World.BankDeposit,
	})
}

func setWalletTotal(w *world.World, total int) {
	w.Wallet = [6]int{}
	remaining := total
	for i, denom := range world.Denominations {
		if remaining <= 0 {
			break
		}
		count := remaining / denom
		if count > 0 {
			w.Wallet[i] = count
			remaining -= count * denom
		}
	}
}
