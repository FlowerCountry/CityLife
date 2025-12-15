package middleware

import (
	"citylife/internal/api/response"
	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

// SessionKey 上下文中Session的Key
const SessionKey = "session"

// RequireSession Session验证中间件
func RequireSession(sm *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")
		if sessionID == "" {
			response.BadRequest(c, "缺少session_id参数")
			c.Abort()
			return
		}

		sess, exists := sm.Get(sessionID)
		if !exists {
			response.SessionNotFound(c)
			c.Abort()
			return
		}

		// 刷新Session过期时间
		sm.Touch(sessionID)

		// 将Session存入上下文
		c.Set(SessionKey, sess)
		c.Next()
	}
}

// GetSession 从上下文获取Session
func GetSession(c *gin.Context) *session.Session {
	if sess, exists := c.Get(SessionKey); exists {
		return sess.(*session.Session)
	}
	return nil
}
