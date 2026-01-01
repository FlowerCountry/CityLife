package v2

import "testing"

func TestRestaurant(t *testing.T) {
	env := setupTestEnv()

	t.Run("菜单不为空", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/restaurant/menu", nil)
		resp, _ := parseResponse(w)
		items, ok := parseDataList(resp)
		if !ok || len(items) == 0 {
			t.Fatal("菜单不应为空")
		}
	})

	t.Run("不在餐馆无法消费", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/restaurant/eat_fastfood", nil)
		resp, _ := parseResponse(w)
		if resp.Success {
			t.Fatal("不在餐馆时消费应该失败")
		}
	})

	t.Run("餐馆消费成功后钱包减少", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 去餐馆
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/9", nil)

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/restaurant/eat_fastfood", nil)
		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("消费应该成功, error: %+v", resp.Error)
		}

		walletResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/wallet", nil)
		walletData, _ := parseResponse(walletResp)
		wallet, _ := parseDataMap(walletData)
		if wallet["total"].(float64) != 50 {
			t.Errorf("期望钱包为50，得到 %v", wallet["total"])
		}
	})
}
