package v2

import (
	"net/http"
	"testing"
)

func TestDeposit(t *testing.T) {
	env := setupTestEnv()

	t.Run("在银行存款", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 先导航到银行
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		// 初始余额100，存100需要先有足够的钱
		// 由于银行要求100元为单位，初始100元刚好可以存
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/deposit", map[string]int{
			"amount": 100,
		})

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("存款应该成功, error: %+v", resp.Error)
		}
	})

	t.Run("存款支持自定义金额", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 钱包加到500元
		env.setTestMoney(sessionID, 500, 0)

		// 导航到银行
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		// 存200元
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/deposit", map[string]int{
			"amount": 200,
		})

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("存款应该成功, error: %+v", resp.Error)
		}

		// 验证银行余额增加到200
		bankResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/bank", nil)
		bankData, _ := parseResponse(bankResp)
		bank, _ := parseDataMap(bankData)
		if bank["balance"].(float64) != 200 {
			t.Errorf("期望银行余额为200，得到 %v", bank["balance"])
		}
	})
}

func TestWithdraw(t *testing.T) {
	env := setupTestEnv()

	t.Run("取款需要在银行位置", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/withdraw", map[string]int{
			"amount": 100,
		})

		resp, _ := parseResponse(w)
		if resp.Success {
			t.Fatal("不在银行时取款应该失败")
		}
	})

	t.Run("取款余额不足", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		// 银行初始余额为0，取款应该失败
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/withdraw", map[string]int{
			"amount": 100,
		})

		resp, _ := parseResponse(w)
		if resp.Success {
			t.Fatal("余额不足时取款应该失败")
		}
	})
}

func TestDepositAll(t *testing.T) {
	env := setupTestEnv()

	t.Run("存入全部现金", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/deposit-all", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("存入全部现金应该成功, error: %+v", resp.Error)
		}
	})
}

func TestWithdrawAll(t *testing.T) {
	env := setupTestEnv()

	t.Run("取出全部存款", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		// 先存入全部
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/deposit-all", nil)

		// 再取出全部
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/withdraw-all", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("取出全部存款应该成功, error: %+v", resp.Error)
		}
	})
}

func TestGetBankBalance(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取银行余额", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/bank", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 初始银行余额应该为0
		if data["balance"].(float64) != 0 {
			t.Errorf("期望初始银行余额为0，得到 %v", data["balance"])
		}
	})
}
