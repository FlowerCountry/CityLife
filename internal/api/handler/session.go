// Package handler 包含API处理器
package handler

import (
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

// CreateSession 创建新游戏会话
func CreateSession(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sm.Create()
		response.Created(c, gin.H{
			"session_id": sess.ID,
			"created_at": sess.CreatedAt,
			"expires_at": sess.ExpiresAt,
		})
	}
}

// GetSession 获取会话信息
func GetSession(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := middleware.GetSession(c)
		if sess == nil {
			response.SessionNotFound(c)
			return
		}

		response.Success(c, gin.H{
			"session_id": sess.ID,
			"created_at": sess.CreatedAt,
			"updated_at": sess.UpdatedAt,
			"expires_at": sess.ExpiresAt,
		})
	}
}

// DeleteSession 删除会话
func DeleteSession(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")
		sm.Delete(sessionID)
		response.Success(c, gin.H{
			"message": "会话已删除",
		})
	}
}
