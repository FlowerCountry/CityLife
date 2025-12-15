package handler

import (
	"net/http"
	"testing"
)

func TestSeeDoctor(t *testing.T) {
	env := setupTestEnv()

	t.Run("看病成功（无疾病）", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/hospital/doctor", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("看病应该成功, message: %v", data["message"])
		}

		// 无疾病时应该提示健康
		message, ok := data["message"].(string)
		if !ok {
			t.Error("响应中缺少message字段")
		}
		if message == "" {
			t.Error("message不应为空")
		}
	})
}

func TestGetCheckups(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取体检项目列表", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/hospital/checkups", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 验证套餐和单项检查存在
		if _, exists := data["packages"]; !exists {
			t.Error("响应中缺少packages字段")
		}

		if _, exists := data["singles"]; !exists {
			t.Error("响应中缺少singles字段")
		}

		// 验证套餐不为空
		packages, ok := data["packages"].([]interface{})
		if !ok || len(packages) == 0 {
			t.Error("套餐列表不应为空")
		}
	})
}

func TestDoCheckup(t *testing.T) {
	env := setupTestEnv()

	t.Run("执行体检", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		// 执行核心指标套餐体检
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/hospital/checkup/p_core", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("体检应该成功, message: %v", data["message"])
		}
	})

	t.Run("执行无效体检项目", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/hospital/checkup/invalid_checkup", nil)

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] == true {
			t.Error("无效体检项目应该失败")
		}
	})
}

func TestGetMedicines(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取药品列表", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/hospital/medicines", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 验证OTC和处方药列表存在
		if _, exists := data["otc"]; !exists {
			t.Error("响应中缺少otc字段")
		}

		if _, exists := data["prescription"]; !exists {
			t.Error("响应中缺少prescription字段")
		}

		// OTC药品不应为空
		otc, ok := data["otc"].([]interface{})
		if !ok || len(otc) == 0 {
			t.Error("OTC药品列表不应为空")
		}
	})
}

func TestBuyMedicine(t *testing.T) {
	env := setupTestEnv()

	t.Run("购买OTC药品", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		// 购买维生素C片
		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/hospital/medicine/vitamin_c", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] != true {
			t.Errorf("购买OTC药品应该成功, message: %v", data["message"])
		}
	})

	t.Run("购买无效药品", func(t *testing.T) {
		sessionID := env.createTestSession()

		// 导航到医院
		env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/navigate/4", nil)

		w := env.doRequest("POST", "/api/v1/sessions/"+sessionID+"/hospital/medicine/invalid_medicine", nil)

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		if data["success"] == true {
			t.Error("购买无效药品应该失败")
		}
	})
}

func TestGetDiseases(t *testing.T) {
	env := setupTestEnv()

	t.Run("获取疾病状态", func(t *testing.T) {
		sessionID := env.createTestSession()

		w := env.doRequest("GET", "/api/v1/sessions/"+sessionID+"/diseases", nil)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}

		resp, _ := parseResponse(w)
		data, _ := parseDataMap(resp)

		// 初始应该没有疾病
		diseases, ok := data["diseases"].([]interface{})
		if !ok {
			t.Fatal("无法解析疾病列表")
		}

		if len(diseases) != 0 {
			t.Error("初始应该没有疾病")
		}
	})
}
