package center

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()

	t.Run("初始化公告数量", func(t *testing.T) {
		if c.Count() != 4 {
			t.Errorf("Count() = %d, want 4", c.Count())
		}
	})

	t.Run("公告非空", func(t *testing.T) {
		for i, a := range c.GetAnnouncements() {
			if a == "" {
				t.Errorf("Announcement %d is empty", i)
			}
		}
	})
}

func TestNewWithSeed(t *testing.T) {
	t.Run("相同种子相同顺序", func(t *testing.T) {
		c1 := NewWithSeed(42)
		c2 := NewWithSeed(42)

		a1 := c1.GetAnnouncements()
		a2 := c2.GetAnnouncements()

		for i := range a1 {
			if a1[i] != a2[i] {
				t.Errorf("Announcement %d differs: %s vs %s", i, a1[i], a2[i])
			}
		}
	})

	t.Run("不同种子可能不同顺序", func(t *testing.T) {
		c1 := NewWithSeed(42)
		c2 := NewWithSeed(123)

		a1 := c1.GetAnnouncements()
		a2 := c2.GetAnnouncements()

		// 至少应该有一个不同（概率上几乎必然）
		allSame := true
		for i := range a1 {
			if a1[i] != a2[i] {
				allSame = false
				break
			}
		}
		if allSame {
			t.Log("Different seeds produced same order (unlikely but possible)")
		}
	})
}

func TestGetAnnouncements(t *testing.T) {
	c := NewWithSeed(42)
	announcements := c.GetAnnouncements()

	t.Run("返回4条公告", func(t *testing.T) {
		if len(announcements) != 4 {
			t.Errorf("len(GetAnnouncements()) = %d, want 4", len(announcements))
		}
	})

	t.Run("包含银行公告", func(t *testing.T) {
		found := false
		for _, a := range announcements {
			if strings.Contains(a, "银行") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should contain announcement about 银行")
		}
	})

	t.Run("包含电信大楼公告", func(t *testing.T) {
		found := false
		for _, a := range announcements {
			if strings.Contains(a, "电信大楼") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should contain announcement about 电信大楼")
		}
	})
}

func TestFormatAnnouncements(t *testing.T) {
	c := NewWithSeed(42)
	formatted := c.FormatAnnouncements()

	t.Run("包含标题", func(t *testing.T) {
		if !strings.Contains(formatted, "市 中 心 公 告") {
			t.Error("Should contain title '市 中 心 公 告'")
		}
	})

	t.Run("包含分隔线", func(t *testing.T) {
		if !strings.Contains(formatted, "========") {
			t.Error("Should contain separator line")
		}
	})

	t.Run("包含所有公告", func(t *testing.T) {
		for _, a := range c.GetAnnouncements() {
			if !strings.Contains(formatted, a) {
				t.Errorf("Formatted output should contain: %s", a)
			}
		}
	})

	t.Run("包含项目符号", func(t *testing.T) {
		if !strings.Contains(formatted, "•") {
			t.Error("Should contain bullet point '•'")
		}
	})
}

func TestAddAnnouncement(t *testing.T) {
	c := NewWithSeed(42)
	initialCount := c.Count()

	newAnnouncement := "于2011年1月,我市图书馆正式开放"
	c.AddAnnouncement(newAnnouncement)

	t.Run("公告数量增加", func(t *testing.T) {
		if c.Count() != initialCount+1 {
			t.Errorf("Count() = %d, want %d", c.Count(), initialCount+1)
		}
	})

	t.Run("新公告存在", func(t *testing.T) {
		found := false
		for _, a := range c.GetAnnouncements() {
			if a == newAnnouncement {
				found = true
				break
			}
		}
		if !found {
			t.Error("New announcement should be in list")
		}
	})
}

func TestCount(t *testing.T) {
	c := NewWithSeed(42)

	if c.Count() != 4 {
		t.Errorf("Initial Count() = %d, want 4", c.Count())
	}

	c.AddAnnouncement("测试公告")
	if c.Count() != 5 {
		t.Errorf("After add, Count() = %d, want 5", c.Count())
	}
}

func TestRandomization(t *testing.T) {
	// 测试随机化确实发生
	t.Run("多次创建顺序不同", func(t *testing.T) {
		// 创建多个实例，统计第一条公告
		firstItems := make(map[string]int)
		for i := 0; i < 100; i++ {
			c := New()
			firstItems[c.GetAnnouncements()[0]]++
		}

		// 应该有多种不同的第一条公告
		if len(firstItems) < 2 {
			t.Error("Randomization should produce different orderings")
		}
	})
}
