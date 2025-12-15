package handler

import (
	"net/http"
	"testing"
)

func TestDeposit(t *testing.T) {
	env := setupTestEnv()

	t.Run("在银行存款", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 先导航到银行
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/2", nil)

		// 初始余额100，存100需要先有足够的钱
		// 由于银行要求100元为单位，初始100元刚好可以存
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/deposit", map[string]int{
			"amount": 100,
		})

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("存款应该成功, message: %v", data["message"])
		}
	})
}

func TestWithdraw(t *testing.T) {
	env := setupTestEnv()

	t.Run("取款需要在银行位置", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/withdraw", map[string]int{
			"amount": 100,
		})

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] == true {
			t.Error("不在银行时取款应该失败")
		}
	})

	t.Run("取款余额不足", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/2", nil)

		// 银行初始余额为0，取款应该失败
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/withdraw", map[string]int{
			"amount": 100,
		})

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] == true {
			t.Error("余额不足时取款应该失败")
		}
	})
}

func TestDepositAll(t *testing.T) {
	env := setupTestEnv()

	t.Run("存入全部现金", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/2", nil)

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/deposit-all", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("存入全部现金应该成功, message: %v", data["message"])
		}
	})
}

func TestWithdrawAll(t *testing.T) {
	env := setupTestEnv()

	t.Run("取出全部存款", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到银行
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/2", nil)

		// 先存入全部
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/deposit-all", nil)

		// 再取出全部
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/bank/withdraw-all", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("取出全部存款应该成功, message: %v", data["message"])
		}
	})
}

func TestGetBankBalance(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取银行余额", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/bank", nil)

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
