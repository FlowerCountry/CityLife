package v2

import (
	"net/http"
	"testing"
)

func TestGetSaveSlots(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取存档槽位列表", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/saves", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}

		slots, ok := parseDataList(resp)
		if !ok {
			t.Fatal("无法解析存档槽位列表")
		}

		// 应该有3个槽位
		if len(slots) != 3 {
			t.Errorf("期望3个存档槽位，得到 %d", len(slots))
		}

		// 验证槽位结构
		for _, s := range slots {
			slot, ok := s.(map[string]interface{})
			if !ok {
				t.Error("无法解析槽位结构")
				continue
			}

			if _, exists := slot["slot"]; !exists {
				t.Error("槽位缺少slot字段")
			}
			if _, exists := slot["exists"]; !exists {
				t.Error("槽位缺少exists字段")
			}
		}
	})
}

func TestSaveGame(t *testing.T) {
	env := setupTestEnv()

	t.Run("保存游戏到槽位1", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/1", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("保存应该成功, error: %+v", resp.Error)
		}
	})

	t.Run("无效槽位号", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 槽位只有1-3
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/0", nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
		}

		w = env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/4", nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestLoadGame(t *testing.T) {
	env := setupTestEnv()

	t.Run("保存后加载游戏", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 先导航到超市改变状态
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/1", nil)

		// 保存游戏
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/1", nil)

		// 导航到其他位置
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/2", nil)

		// 验证位置已变
		stateResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/state", nil)
		stateData, _ := parseResponse(stateResp)
		state, _ := parseDataMap(stateData)
		location := state["location"].(map[string]interface{})
		if location["id"].(float64) != 2 {
			t.Error("位置应该已变为银行")
		}

		// 加载游戏
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/1/load", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("加载应该成功, error: %+v", resp.Error)
		}

		// 验证位置已恢复
		stateResp = env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/state", nil)
		stateData, _ = parseResponse(stateResp)
		state, _ = parseDataMap(stateData)
		location = state["location"].(map[string]interface{})
		if location["id"].(float64) != 1 {
			t.Error("位置应该已恢复为超市")
		}
	})

	t.Run("加载不存在的存档", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 使用槽位3（假设为空）
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/saves/3/load", nil)

		resp, _ := parseResponse(w)
		// 加载空存档应该失败
		if resp.Success {
			t.Fatal("加载不存在的存档应该失败")
		}
	})
}
