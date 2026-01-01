package v2

import (
	"net/http"
	"testing"
)

func TestGetActions(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取可用行动列表", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/actions", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, err := parseResponse(w)
		if err != nil {
			t.Fatalf("解析响应失败: %v", err)
		}

		if !resp.Success {
			t.Error("期望响应成功")
		}

		actions, ok := parseDataList(resp)
		if !ok {
			t.Fatal("无法解析行动列表")
		}

		if len(actions) == 0 {
			t.Error("行动列表不应为空")
		}

		// 验证行动结构
		for _, a := range actions {
			action, ok := a.(map[string]interface{})
			if !ok {
				t.Error("无法解析行动结构")
				continue
			}

			if _, exists := action["id"]; !exists {
				t.Error("行动缺少id字段")
			}
			if _, exists := action["name"]; !exists {
				t.Error("行动缺少name字段")
			}
			if _, exists := action["category"]; !exists {
				t.Error("行动缺少category字段")
			}
		}
	})

	t.Run("初始位置包含导航行动", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/actions", nil)
		resp, _ := parseResponse(w)
		actions, _ := parseDataList(resp)

		hasNavigation := false
		for _, a := range actions {
			action := a.(map[string]interface{})
			if action["category"] == "navigation" {
				hasNavigation = true
				break
			}
		}

		if !hasNavigation {
			t.Error("初始位置应包含导航行动")
		}
	})

	t.Run("无效Session返回404", func(t *testing.T) {
		w := env.doRequest("GET", "/api/v2/sessions/invalid-id/actions", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestExecuteAction(t *testing.T) {
	env := setupTestEnv()

	t.Run("执行查看钱包行动", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/actions/check_cash", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}

		data, ok := parseDataMap(resp)
		if !ok {
			t.Fatal("无法解析响应数据")
		}

		if _, exists := data["message"]; !exists {
			t.Error("响应中缺少message字段")
		}
	})

	t.Run("执行无效行动", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/actions/invalid_action", nil)

		// 无效行动返回404
		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}
	})
}
