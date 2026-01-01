package session

import (
	"sync"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	t.Run("创建管理器", func(t *testing.T) {
		m := NewManager(30 * time.Minute)
		if m == nil {
			t.Fatal("管理器不应为nil")
		}

		if m.Count() != 0 {
			t.Errorf("初始Session数量应为0，得到 %d", m.Count())
		}
	})
}

func TestManagerCreate(t *testing.T) {
	t.Run("创建Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)

		sess := m.Create()
		if sess == nil {
			t.Fatal("Session不应为nil")
		}

		if sess.ID == "" {
			t.Error("Session ID不应为空")
		}

		if sess.State == nil {
			t.Error("Session State不应为nil")
		}

		if m.Count() != 1 {
			t.Errorf("Session数量应为1，得到 %d", m.Count())
		}
	})

	t.Run("创建多个Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)

		sess1 := m.Create()
		sess2 := m.Create()
		sess3 := m.Create()

		if sess1.ID == sess2.ID || sess2.ID == sess3.ID {
			t.Error("Session ID应该唯一")
		}

		if m.Count() != 3 {
			t.Errorf("Session数量应为3，得到 %d", m.Count())
		}
	})
}

func TestManagerGet(t *testing.T) {
	t.Run("获取存在的Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)
		created := m.Create()

		sess, exists := m.Get(created.ID)
		if !exists {
			t.Fatal("Session应该存在")
		}

		if sess.ID != created.ID {
			t.Errorf("Session ID不匹配")
		}
	})

	t.Run("获取不存在的Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)

		sess, exists := m.Get("non-existent-id")
		if exists {
			t.Error("Session不应该存在")
		}
		if sess != nil {
			t.Error("Session应该为nil")
		}
	})

	t.Run("获取过期的Session", func(t *testing.T) {
		// 使用1毫秒的TTL
		m := NewManager(1 * time.Millisecond)
		created := m.Create()

		// 等待过期
		time.Sleep(10 * time.Millisecond)

		_, exists := m.Get(created.ID)
		if exists {
			t.Error("过期Session不应该被获取到")
		}
	})
}

func TestManagerTouch(t *testing.T) {
	t.Run("刷新Session过期时间", func(t *testing.T) {
		m := NewManager(100 * time.Millisecond)
		sess := m.Create()
		originalExpiry := sess.ExpiresAt

		// 等待一点时间
		time.Sleep(20 * time.Millisecond)

		// 刷新
		m.Touch(sess.ID)

		// 重新获取
		refreshed, exists := m.Get(sess.ID)
		if !exists {
			t.Fatal("Session应该存在")
		}

		if !refreshed.ExpiresAt.After(originalExpiry) {
			t.Error("过期时间应该被延长")
		}
	})
}

func TestManagerDelete(t *testing.T) {
	t.Run("删除Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)
		sess := m.Create()

		m.Delete(sess.ID)

		_, exists := m.Get(sess.ID)
		if exists {
			t.Error("删除后Session不应该存在")
		}

		if m.Count() != 0 {
			t.Errorf("删除后Session数量应为0，得到 %d", m.Count())
		}
	})

	t.Run("删除不存在的Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)

		// 不应该panic
		m.Delete("non-existent-id")
	})
}

func TestManagerConcurrency(t *testing.T) {
	t.Run("并发创建Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)
		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.Create()
			}()
		}

		wg.Wait()

		if m.Count() != numGoroutines {
			t.Errorf("Session数量应为 %d，得到 %d", numGoroutines, m.Count())
		}
	})

	t.Run("并发读写Session", func(t *testing.T) {
		m := NewManager(30 * time.Minute)
		sess := m.Create()
		var wg sync.WaitGroup

		// 并发读取
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.Get(sess.ID)
			}()
		}

		// 并发刷新
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.Touch(sess.ID)
			}()
		}

		wg.Wait()

		// 验证Session仍然存在
		_, exists := m.Get(sess.ID)
		if !exists {
			t.Error("Session应该仍然存在")
		}
	})
}

func TestManagerCleanup(t *testing.T) {
	t.Run("手动清理测试", func(t *testing.T) {
		m := NewManager(1 * time.Millisecond)

		// 创建几个Session
		m.Create()
		m.Create()
		m.Create()

		if m.Count() != 3 {
			t.Fatalf("创建后Session数量应为3，得到 %d", m.Count())
		}

		// 等待过期
		time.Sleep(10 * time.Millisecond)

		// 手动调用cleanup
		m.cleanup()

		if m.Count() != 0 {
			t.Errorf("清理后Session数量应为0，得到 %d", m.Count())
		}
	})
}
