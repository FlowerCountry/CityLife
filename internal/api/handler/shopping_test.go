package handler

import (
	"net/http"
	"testing"
)

func TestGetCommodities(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取商品列表", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到超市
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/shop/commodities", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Error("期望响应成功")
		}

		commodities, ok := parseDataList(resp)
		if !ok {
			t.Fatal("无法解析商品列表")
		}

		if len(commodities) == 0 {
			t.Error("商品列表不应为空")
		}

		// 验证商品结构
		for _, c := range commodities {
			commodity, ok := c.(map[string]interface{})
			if !ok {
				t.Error("无法解析商品结构")
				continue
			}

			requiredFields := []string{"name", "price", "type", "freshness"}
			for _, field := range requiredFields {
				if _, exists := commodity[field]; !exists {
					t.Errorf("商品缺少字段: %s", field)
				}
			}
		}
	})

	t.Run("商品列表包含面包", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到超市
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/shop/commodities", nil)
		resp, _ := parseResponse(w)
		commodities, _ := parseDataList(resp)

		hasBread := false
		for _, c := range commodities {
			commodity := c.(map[string]interface{})
			if commodity["name"] == "面包" {
				hasBread = true
				break
			}
		}

		if !hasBread {
			t.Error("商品列表应包含面包")
		}
	})
}

func TestBuyCommodity(t *testing.T) {
	env := setupTestEnv()

	t.Run("购买商品成功", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到超市
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		// 购买面包（价格15元）
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/shop/buy/面包", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("购买应该成功, message: %v", data["message"])
		}

		// 验证余额已减少
		walletResp := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/wallet", nil)
		walletData, _ := parseResponse(walletResp)
		wallet, _ := parseDataMap(walletData)

		if wallet["total"].(float64) >= 100 {
			t.Error("购买后余额应该减少")
		}
	})

	t.Run("购买不存在的商品", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到超市
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/shop/buy/不存在的商品", nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("余额不足无法购买", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到超市
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/1", nil)

		// 多次购买直到余额不足
		for i := 0; i < 10; i++ {
			env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/shop/buy/面包", nil)
		}

		// 尝试购买昂贵商品（牛排150元）
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/shop/buy/牛排", nil)
		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] == true {
			t.Error("余额不足时购买应该失败")
		}
	})
}
