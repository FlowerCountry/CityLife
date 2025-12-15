package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

// testEnv 测试环境
type testEnv struct {
	Router     *gin.Engine
	SessionMgr *session.Manager
}

// setupTestEnv 创建测试环境
func setupTestEnv() *testEnv {
	gin.SetMode(gin.TestMode)

	sm := session.NewManager(30 * time.Minute)
	r := gin.New()

	// 注册路由（避免import cycle，直接在此注册）
	registerTestRoutes(r, sm)

	return &testEnv{
		Router:     r,
		SessionMgr: sm,
	}
}

// registerTestRoutes 为测试注册路由
func registerTestRoutes(r *gin.Engine, sm *session.Manager) {
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")

	v1.POST("/sessions", CreateSession(sm))

	sessionGroup := v1.Group("/sessions/:session_id")
	sessionGroup.Use(middleware.RequireSession(sm))
	{
		sessionGroup.GET("", GetSession(sm))
		sessionGroup.DELETE("", DeleteSession(sm))

		sessionGroup.GET("/state", GetState)
		sessionGroup.GET("/status", GetStatus)

		sessionGroup.GET("/actions", GetActions)
		sessionGroup.POST("/actions", ExecuteAction)

		sessionGroup.POST("/navigate/:location_id", Navigate)

		bank := sessionGroup.Group("/bank")
		{
			bank.POST("/deposit", Deposit)
			bank.POST("/withdraw", Withdraw)
			bank.POST("/deposit-all", DepositAll)
			bank.POST("/withdraw-all", WithdrawAll)
		}

		shop := sessionGroup.Group("/shop")
		{
			shop.GET("/commodities", GetCommodities)
			shop.POST("/buy/:commodity", BuyCommodity)
		}

		hospital := sessionGroup.Group("/hospital")
		{
			hospital.POST("/doctor", SeeDoctor)
			hospital.GET("/checkups", GetCheckups)
			hospital.POST("/checkup/:id", DoCheckup)
			hospital.GET("/medicines", GetMedicines)
			hospital.POST("/medicine/:id", BuyMedicine)
		}

		saves := sessionGroup.Group("/saves")
		{
			saves.GET("", GetSaveSlots)
			saves.POST("/:slot", SaveGame)
			saves.POST("/:slot/load", LoadGame)
		}

		sessionGroup.GET("/wallet", GetWallet)
		sessionGroup.GET("/bank", GetBankBalance)
		sessionGroup.GET("/health", GetHealth)
		sessionGroup.GET("/diseases", GetDiseases)
		sessionGroup.GET("/map", GetMap)
	}
}

// createTestSession 创建测试Session并返回ID
func (e *testEnv) createTestSession() string {
	sess := e.SessionMgr.Create()
	return sess.ID
}

// doRequest 发送HTTP请求
func (e *testEnv) doRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	e.Router.ServeHTTP(w, req)
	return w
}

// parseResponse 解析响应
func parseResponse(w *httptest.ResponseRecorder) (*response.Response, error) {
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	return &resp, err
}

// parseDataMap 将响应Data解析为map
func parseDataMap(resp *response.Response) (map[string]interface{}, bool) {
	if resp.Data == nil {
		return nil, false
	}
	data, ok := resp.Data.(map[string]interface{})
	return data, ok
}

// parseDataList 将响应Data解析为slice
func parseDataList(resp *response.Response) ([]interface{}, bool) {
	if resp.Data == nil {
		return nil, false
	}
	data, ok := resp.Data.([]interface{})
	return data, ok
}
