package v2

import (
	"net/http"
	"testing"
)

func TestGetHousingOffers(t *testing.T) {
	env := setupTestEnv()

	t.Run("初始住房可选项", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/housing/offers", nil)
		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["can_renew"] != false || data["can_cancel"] != false {
			t.Errorf("初始不应支持续租/退租, data: %+v", data)
		}

		rent, ok := data["rent"].([]interface{})
		if !ok || len(rent) != 3 {
			t.Errorf("期望3个租房选项，得到 %v", data["rent"])
		}

		buy, ok := data["buy"].([]interface{})
		if !ok || len(buy) != 2 {
			t.Errorf("期望2个买房选项，得到 %v", data["buy"])
		}
	})
}

func TestRentHousing(t *testing.T) {
	env := setupTestEnv()

	t.Run("不在房产中介无法租房", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/housing/rent/bed", nil)
		resp, _ := parseResponse(w)
		if resp.Success {
			t.Fatal("不在房产中介时租房应该失败")
		}
	})

	t.Run("租房成功后住房状态更新", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 钱包加到5000元
		env.setTestMoney(sessionID, 5000, 0)

		// 去房产中介
		env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/navigate/7", nil)

		// 租合租床位
		w := env.doRequest("POST", "/api/v2/sessions/"+sessionID+"/housing/rent/bed", nil)
		resp, _ := parseResponse(w)
		if !resp.Success {
			t.Fatalf("租房应该成功, error: %+v", resp.Error)
		}

		// 验证住房状态
		housingResp := env.doRequest("GET", "/api/v2/sessions/"+sessionID+"/housing", nil)
		housingData, _ := parseResponse(housingResp)
		housing, _ := parseDataMap(housingData)
		if housing["status"] != "rented" {
			t.Errorf("期望status=rented，得到 %v", housing["status"])
		}
		if housing["name"] == "" {
			t.Error("租房后name不应为空")
		}
	})
}
