package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"citylife/internal/session"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCORS(t *testing.T) {
	t.Run("设置CORS头", func(t *testing.T) {
		r := gin.New()
		r.Use(CORS())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		// 验证CORS头
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("缺少Access-Control-Allow-Origin头")
		}

		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("缺少Access-Control-Allow-Methods头")
		}

		if w.Header().Get("Access-Control-Allow-Headers") == "" {
			t.Error("缺少Access-Control-Allow-Headers头")
		}
	})

	t.Run("处理OPTIONS预检请求", func(t *testing.T) {
		r := gin.New()
		r.Use(CORS())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != 204 {
			t.Errorf("期望状态码 204, 得到 %d", w.Code)
		}
	})
}

func TestRequireSession(t *testing.T) {
	t.Run("有效Session通过验证", func(t *testing.T) {
		sm := session.NewManager(30 * time.Minute)
		sess := sm.Create()

		r := gin.New()
		r.GET("/sessions/:session_id/test", RequireSession(sm), func(c *gin.Context) {
			// 验证Session已存入上下文
			s := GetSession(c)
			if s == nil {
				c.String(500, "Session not found in context")
				return
			}
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sessions/"+sess.ID+"/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("期望状态码 200, 得到 %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("无效Session返回404", func(t *testing.T) {
		sm := session.NewManager(30 * time.Minute)

		r := gin.New()
		r.GET("/sessions/:session_id/test", RequireSession(sm), func(c *gin.Context) {
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sessions/invalid-id/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Errorf("期望状态码 404, 得到 %d", w.Code)
		}
	})

	t.Run("过期Session返回404", func(t *testing.T) {
		sm := session.NewManager(1 * time.Millisecond)
		sess := sm.Create()

		// 等待过期
		time.Sleep(10 * time.Millisecond)

		r := gin.New()
		r.GET("/sessions/:session_id/test", RequireSession(sm), func(c *gin.Context) {
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sessions/"+sess.ID+"/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Errorf("期望状态码 404, 得到 %d", w.Code)
		}
	})

	t.Run("Session访问刷新过期时间", func(t *testing.T) {
		sm := session.NewManager(100 * time.Millisecond)
		sess := sm.Create()
		originalExpiry := sess.ExpiresAt

		r := gin.New()
		r.GET("/sessions/:session_id/test", RequireSession(sm), func(c *gin.Context) {
			c.String(200, "OK")
		})

		// 等待一点时间
		time.Sleep(20 * time.Millisecond)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sessions/"+sess.ID+"/test", nil)
		r.ServeHTTP(w, req)

		// 重新获取Session检查过期时间
		refreshed, _ := sm.Get(sess.ID)
		if !refreshed.ExpiresAt.After(originalExpiry) {
			t.Error("访问后Session过期时间应该延长")
		}
	})
}

func TestGetSession(t *testing.T) {
	t.Run("从上下文获取Session", func(t *testing.T) {
		sm := session.NewManager(30 * time.Minute)
		sess := sm.Create()

		r := gin.New()
		r.GET("/sessions/:session_id/test", RequireSession(sm), func(c *gin.Context) {
			s := GetSession(c)
			if s == nil {
				t.Error("应该能从上下文获取Session")
				c.String(500, "Failed")
				return
			}
			if s.ID != sess.ID {
				t.Error("获取的Session ID不匹配")
			}
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sessions/"+sess.ID+"/test", nil)
		r.ServeHTTP(w, req)
	})

	t.Run("无Session时返回nil", func(t *testing.T) {
		r := gin.New()
		r.GET("/test", func(c *gin.Context) {
			s := GetSession(c)
			if s != nil {
				t.Error("无Session时应该返回nil")
			}
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
	})
}
