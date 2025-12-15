package handler

import (
	"net/http"
	"testing"
)

func TestGetState(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取游戏状态成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/state", nil)

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

		// 验证关键字段存在
		requiredFields := []string{"location", "time", "money", "health", "is_alive"}
		for _, field := range requiredFields {
			if _, exists := data[field]; !exists {
				t.Errorf("响应中缺少字段: %s", field)
			}
		}
	})

	t.Run("初始位置为市中心", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/state", nil)
		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		location, ok := data["location"].(map[string]interface{})
		if !ok {
			t.Fatal("无法解析位置信息")
		}

		if location["id"].(float64) != 0 {
			t.Errorf("期望初始位置ID为0，得到 %v", location["id"])
		}

		if location["name"] != "市中心" {
			t.Errorf("期望初始位置为市中心，得到 %v", location["name"])
		}
	})

	t.Run("初始钱包余额为100", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/state", nil)
		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		money, ok := data["money"].(map[string]interface{})
		if !ok {
			t.Fatal("无法解析金钱信息")
		}

		if money["wallet_total"].(float64) != 100 {
			t.Errorf("期望初始钱包余额为100，得到 %v", money["wallet_total"])
		}
	})

	t.Run("无效Session返回404", func(t *testing.T) {
		w := env.doRequest("GET", "/api/v1/sessions/invalid-id/state", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestGetStatus(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取状态摘要成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/status", nil)

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

		// 验证关键字段
		if data["location"] != "市中心" {
			t.Errorf("期望位置为市中心，得到 %v", data["location"])
		}

		if data["wallet"].(float64) != 100 {
			t.Errorf("期望钱包余额为100，得到 %v", data["wallet"])
		}

		if data["is_alive"] != true {
			t.Error("期望角色存活")
		}
	})
}

func TestGetWallet(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取钱包详情", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/wallet", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["total"].(float64) != 100 {
			t.Errorf("期望钱包总额为100，得到 %v", data["total"])
		}

		// 验证面额存在
		if _, exists := data["denominations"]; !exists {
			t.Error("响应中缺少面额信息")
		}
	})
}

func TestGetHealth(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取健康状态", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/health", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["is_alive"] != true {
			t.Error("期望角色存活")
		}

		// 验证核心营养存在
		if _, exists := data["core"]; !exists {
			t.Error("响应中缺少核心营养信息")
		}
	})
}

func TestGetMap(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取地图信息", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/map", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 验证当前位置
		if data["current_location"].(float64) != 0 {
			t.Errorf("期望当前位置为0，得到 %v", data["current_location"])
		}

		// 验证位置列表存在
		locations, ok := data["locations"].([]interface{})
		if !ok || len(locations) == 0 {
			t.Error("响应中缺少或位置列表为空")
		}
	})
}
