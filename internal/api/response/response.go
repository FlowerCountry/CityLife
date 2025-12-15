// Package response 定义API响应结构
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 错误代码
const (
	ErrCodeSessionNotFound      = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired       = "SESSION_EXPIRED"
	ErrCodeInvalidAction        = "INVALID_ACTION"
	ErrCodeInsufficientFunds    = "INSUFFICIENT_FUNDS"
	ErrCodeInvalidAmount        = "INVALID_AMOUNT"
	ErrCodeInvalidLocation      = "INVALID_LOCATION"
	ErrCodeGameOver             = "GAME_OVER"
	ErrCodeSaveNotFound         = "SAVE_NOT_FOUND"
	ErrCodeInvalidSlot          = "INVALID_SLOT"
	ErrCodePrescriptionRequired = "PRESCRIPTION_REQUIRED"
	ErrCodeInvalidRequest       = "INVALID_REQUEST"
	ErrCodeInternalError        = "INTERNAL_ERROR"
)

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created 返回创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpStatus int, code, message string) {
	c.JSON(httpStatus, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// BadRequest 返回400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrCodeInvalidRequest, message)
}

// NotFound 返回404错误
func NotFound(c *gin.Context, code, message string) {
	Error(c, http.StatusNotFound, code, message)
}

// Conflict 返回409错误
func Conflict(c *gin.Context, code, message string) {
	Error(c, http.StatusConflict, code, message)
}

// InternalError 返回500错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeInternalError, message)
}

// SessionNotFound 返回Session不存在错误
func SessionNotFound(c *gin.Context) {
	NotFound(c, ErrCodeSessionNotFound, "游戏会话不存在或已过期")
}

// GameOver 返回游戏结束错误
func GameOver(c *gin.Context, message string) {
	Error(c, http.StatusConflict, ErrCodeGameOver, message)
}

// LocationInfo 位置信息
type LocationInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TimeInfo 时间信息
type TimeInfo struct {
	Year    int    `json:"year"`
	Month   int    `json:"month"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
	Second  int    `json:"second"`
	Display string `json:"display"`
}

// MoneyInfo 金钱信息
type MoneyInfo struct {
	WalletTotal int    `json:"wallet_total"`
	BankDeposit int    `json:"bank_deposit"`
	Wallet      [6]int `json:"wallet"`
}

// NutritionInfo 营养详情
type NutritionInfo struct {
	Value  int    `json:"value"`
	Level  int    `json:"level"`
	Status string `json:"status"`
}

// DiseaseInfo 疾病信息
type DiseaseInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	TreatmentCost int    `json:"treatment_cost"`
}

// StateResponse 游戏状态响应
type StateResponse struct {
	Location LocationInfo          `json:"location"`
	Time     TimeInfo              `json:"time"`
	Money    MoneyInfo             `json:"money"`
	Health   map[string]int        `json:"health"`
	Diseases []DiseaseInfo         `json:"diseases"`
	IsAlive  bool                  `json:"is_alive"`
}

// ActionInfo 可用行动信息
type ActionInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

// ActionResponse 行动响应
type ActionResponse struct {
	Message  string   `json:"message"`
	Success  bool     `json:"success"`
	Hints    []string `json:"hints,omitempty"`
	GameOver bool     `json:"game_over,omitempty"`
}

// CommodityInfo 商品信息
type CommodityInfo struct {
	Name           string   `json:"name"`
	Price          int      `json:"price"`
	Type           string   `json:"type"`
	Freshness      int      `json:"freshness"`
	FreshnessLabel string   `json:"freshness_label"`
	Effects        []string `json:"effects"`
}

// MedicineInfo 药品信息
type MedicineInfo struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Price           int      `json:"price"`
	Type            string   `json:"type"`
	RequiredDisease string   `json:"required_disease,omitempty"`
	Effects         []string `json:"effects"`
}

// CheckupInfo 体检项目信息
type CheckupInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Items []string `json:"items"`
}

// SaveSlotInfo 存档槽位信息
type SaveSlotInfo struct {
	Slot        int    `json:"slot"`
	Exists      bool   `json:"exists"`
	Description string `json:"description"`
}

// WalletInfo 钱包详情
type WalletInfo struct {
	Total       int            `json:"total"`
	Denominations map[int]int  `json:"denominations"`
}
