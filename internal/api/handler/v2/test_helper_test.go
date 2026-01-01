package v2

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

type testEnv struct {
	Router     *gin.Engine
	SessionMgr *session.Manager
}

func setupTestEnv() *testEnv {
	return setupTestEnvWithMode(gin.TestMode)
}

func setupTestEnvWithMode(mode string) *testEnv {
	gin.SetMode(mode)

	sm := session.NewManager(30 * time.Minute)
	r := gin.New()

	registerTestRoutes(r, sm)

	return &testEnv{
		Router:     r,
		SessionMgr: sm,
	}
}

func registerTestRoutes(r *gin.Engine, sm *session.Manager) {
	r.Use(middleware.CORS())

	v2 := r.Group("/api/v2")
	v2.POST("/sessions", CreateSession(sm))

	sessionGroup := v2.Group("/sessions/:session_id")
	sessionGroup.Use(middleware.RequireSession(sm))
	{
		sessionGroup.GET("", GetSession(sm))
		sessionGroup.DELETE("", DeleteSession(sm))

		sessionGroup.GET("/state", GetState)
		sessionGroup.GET("/status", GetStatus)

		sessionGroup.GET("/actions", GetActions)
		sessionGroup.POST("/actions/:action_id", ExecuteAction)

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

		housing := sessionGroup.Group("/housing")
		{
			housing.GET("", GetHousing)
			housing.GET("/offers", GetHousingOffers)
			housing.POST("/rent/:level", RentHousing)
			housing.POST("/buy/:level", BuyHousing)
			housing.POST("/renew", RenewRent)
			housing.POST("/cancel", CancelRent)
			housing.POST("/go-home", GoHome)
			housing.POST("/leave-home", LeaveHome)
			housing.POST("/sleep", SleepAtHome)
		}

		jobs := sessionGroup.Group("/jobs")
		{
			jobs.GET("", GetJobs)
			jobs.POST("/:id", DoJob)
		}

		restaurant := sessionGroup.Group("/restaurant")
		{
			restaurant.GET("/menu", GetRestaurantMenu)
			restaurant.POST("/:id", DoRestaurant)
		}
		park := sessionGroup.Group("/park")
		{
			park.GET("/activities", GetParkActivities)
			park.POST("/:id", DoPark)
		}
		hotel := sessionGroup.Group("/hotel")
		{
			hotel.GET("/services", GetHotelServices)
			hotel.POST("/:id", DoHotel)
		}

		sessionGroup.GET("/wallet", GetWallet)
		sessionGroup.GET("/bank", GetBankBalance)
		sessionGroup.GET("/health", GetHealth)
		sessionGroup.GET("/diseases", GetDiseases)
		sessionGroup.GET("/map", GetMap)

		testTools := sessionGroup.Group("/test")
		{
			testTools.POST("/money", SetTestMoney)
		}
	}
}

func (e *testEnv) createTestSession() string {
	sess := e.SessionMgr.Create()
	return sess.ID
}

func (e *testEnv) setTestMoney(sessionID string, wallet, bank int) *httptest.ResponseRecorder {
	body := map[string]int{
		"wallet": wallet,
		"bank":   bank,
	}
	return e.doRequest("POST", "/api/v2/sessions/"+sessionID+"/test/money", body)
}

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

func parseResponse(w *httptest.ResponseRecorder) (*response.Response, error) {
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	return &resp, err
}

func parseDataMap(resp *response.Response) (map[string]interface{}, bool) {
	if resp.Data == nil {
		return nil, false
	}
	data, ok := resp.Data.(map[string]interface{})
	return data, ok
}

func parseDataList(resp *response.Response) ([]interface{}, bool) {
	if resp.Data == nil {
		return nil, false
	}
	data, ok := resp.Data.([]interface{})
	return data, ok
}
