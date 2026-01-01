package v2

import (
	"net/http"
	"testing"

	"citylife/internal/api/response"
)

func TestCreateSession_V2(t *testing.T) {
	env := setupTestEnv()

	w := env.doRequest("POST", "/api/v2/sessions", nil)
	if w.Code != http.StatusCreated {
		t.Fatalf("期望状态码 %d, 得到 %d", http.StatusCreated, w.Code)
	}

	resp, err := parseResponse(w)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if !resp.Success {
		t.Fatalf("期望success=true, got: %+v", resp)
	}

	data, ok := parseDataMap(resp)
	if !ok {
		t.Fatal("无法解析data")
	}
	if data["session_id"] == "" {
		t.Fatal("session_id不应为空")
	}
}

func TestExecuteAction_DepositAmount(t *testing.T) {
	env := setupTestEnv()
	sessionID := env.createTestSession()

	env.setTestMoney(sessionID, 500, 0)

	// 先到银行
	env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

	w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/actions/deposit", map[string]int{
		"amount": 200,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
	}

	resp, _ := parseResponse(w)
	if !resp.Success {
		t.Fatalf("存款应该成功, error: %+v", resp.Error)
	}

	// 验证银行余额为200
	stateResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/state", nil)
	state, _ := parseResponse(stateResp)
	stateData, _ := parseDataMap(state)
	money := stateData["money"].(map[string]interface{})
	if money["bank_deposit"].(float64) != 200 {
		t.Errorf("期望bank_deposit为200，得到 %v", money["bank_deposit"])
	}
	if money["wallet_total"].(float64) != 300 {
		t.Errorf("期望wallet_total为300，得到 %v", money["wallet_total"])
	}
}

func TestBankDeposit_WrongLocation(t *testing.T) {
	env := setupTestEnv()
	sessionID := env.createTestSession()

	w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/bank/deposit", map[string]int{
		"amount": 100,
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
	}

	resp, _ := parseResponse(w)
	if resp.Success {
		t.Fatal("不在银行时存款应该失败")
	}
	if resp.Error == nil || resp.Error.Code != response.ErrCodeInvalidLocation {
		t.Fatalf("期望error.code=%s, got: %+v", response.ErrCodeInvalidLocation, resp.Error)
	}
}
