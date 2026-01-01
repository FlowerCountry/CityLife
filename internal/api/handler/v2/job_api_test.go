package v2

import "testing"

func TestDoJob(t *testing.T) {
	env := setupTestEnv()

	t.Run("不在人才市场无法工作", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/jobs/work_delivery", nil)
		resp, _ := parseResponse(w)
		if resp.Success {
			t.Fatal("不在人才市场时工作应该失败")
		}
	})

	t.Run("工作成功后钱包增加", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 去人才市场
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/8", nil)

		// 送外卖 +300
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/jobs/work_delivery", nil)
		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("工作应该成功, error: %+v", resp.Error)
		}

		// 验证钱包为400
		walletResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/wallet", nil)
		walletData, _ := parseResponse(walletResp)
		wallet, _ := parseDataMap(walletData)
		if wallet["total"].(float64) != 400 {
			t.Errorf("期望钱包为400，得到 %v", wallet["total"])
		}
	})
}
