package handler

import (
	"net/http"
	"testing"
)

func TestCreateSession(t *testing.T) {
	env := setupTestEnv()

	t.Run("创建Session成功", func(t *testing.T) {
		w := env.doRequest("POST", "/api/v1/sessions", nil)

		if w.Code != http.StatusCreated {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusCreated, w.Code)
		}

		resp, err := parseResponse(w)
		if err != nil {
			t.Fatalf("解析响应失败: %v", err)
		}

		if !resp.Success {
			t.Error("期望响应成功")
		}

		data, ok := parseDataMap(resp)
		if !ok {
			t.Fatal("无法解析响应数据")
		}

		if _, exists := data["session_id"]; !exists {
			t.Error("响应中缺少 session_id")
		}

		if _, exists := data["created_at"]; !exists {
			t.Error("响应中缺少 created_at")
		}

		if _, exists := data["expires_at"]; !exists {
			t.Error("响应中缺少 expires_at")
		}
	})

	t.Run("Session计数增加", func(t *testing.T) {
		initialCount := env.SessionMgr.Count()

		env.doRequest("POST", "/api/v1/sessions", nil)

		newCount := env.SessionMgr.Count()
		if newCount != initialCount+1 {
			t.Errorf("期望Session数量 %d, 得到 %d", initialCount+1, newCount)
		}
	})
}

func TestGetSession(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取存在的Session", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID, nil)

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

		data, ok := parseDataMap(resp)
		if !ok {
			t.Fatal("无法解析响应数据")
		}

		if data["session_id"] != sessionID {
			t.Errorf("期望 session_id %s, 得到 %v", sessionID, data["session_id"])
		}
	})

	t.Run("获取不存在的Session", func(t *testing.T) {
		w := env.doRequest("GET", "/api/v1/sessions/non-existent-id", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}

		resp, err := parseResponse(w)
		if err != nil {
			t.Fatalf("解析响应失败: %v", err)
		}

		if resp.Success {
			t.Error("期望响应失败")
		}

		if resp.Error == nil {
			t.Error("期望有错误信息")
		}
	})
}

func TestDeleteSession(t *testing.T) {
	env := setupTestEnv()

	t.Run("删除存在的Session", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 确认Session存在
		_, exists := env.SessionMgr.Get(sessionID)
		if !exists {
			t.Fatal("Session应该存在")
		}

		w := env.doRequest("DELETE", "/api/v1/sessions/"+sessionID, nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		// 确认Session已删除
		_, exists = env.SessionMgr.Get(sessionID)
		if exists {
			t.Error("Session应该已被删除")
		}
	})

	t.Run("删除不存在的Session", func(t *testing.T) {
		w := env.doRequest("DELETE", "/api/v1/sessions/non-existent-id", nil)

		// 删除不存在的Session返回404（因为中间件会先验证Session）
		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}
	})
}
