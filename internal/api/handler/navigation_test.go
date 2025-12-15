package handler

import (
	"net/http"
	"testing"
)

func TestNavigate(t *testing.T) {
	env := setupTestEnv()

	t.Run("导航到超市成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}

		data, _ := parseDataMap(resp)
		if _, exists := data["message"]; !exists {
			t.Error("响应中缺少message字段")
		}

		// 验证位置已更改
		stateResp := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/state", nil)
		stateData, _ := parseResponse(stateResp)
		state, _ := parseDataMap(stateData)
		location := state["location"].(map[string]interface{})

		if location["id"].(float64) != 1 {
			t.Errorf("期望位置ID为1，得到 %v", location["id"])
		}
	})

	t.Run("导航到银行成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/2", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}
	})

	t.Run("导航到医院成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}
	})

	t.Run("导航到无效位置", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/999", nil)

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 无效位置应该返回失败
		if data["success"] == true {
			t.Error("导航到无效位置应该失败")
		}
	})

	t.Run("无效Session返回404", func(t *testing.T) {
		w := env.doRequest("POST", "/api/v1/sessions/invalid-id/navigate/1", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}
	})
}
