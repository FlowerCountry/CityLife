// Package center 管理市中心公告系统
package center

import (
	"math/rand"
	"strings"
	"time"
)

// Center 公告中心
type Center struct {
	announcements []string
}

// 默认公告内容
var defaultAnnouncements = []string{
	"于2009年7月,我市第一家银行正式完工",
	"于2009年12月,我市第一家银行正式营业,有需要者可以到银行办理业务",
	"于2010年5月,我市预计开始建设电信大楼",
	"于2010年8月,我市预计正式完工电信大楼",
}

// New 创建新的公告中心
func New() *Center {
	c := &Center{
		announcements: make([]string, len(defaultAnnouncements)),
	}

	// 复制公告
	copy(c.announcements, defaultAnnouncements)

	// 随机打乱顺序
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.announcements), func(i, j int) {
		c.announcements[i], c.announcements[j] = c.announcements[j], c.announcements[i]
	})

	return c
}

// NewWithSeed 创建新的公告中心（使用指定种子，用于测试）
func NewWithSeed(seed int64) *Center {
	c := &Center{
		announcements: make([]string, len(defaultAnnouncements)),
	}

	copy(c.announcements, defaultAnnouncements)

	rand.Seed(seed)
	rand.Shuffle(len(c.announcements), func(i, j int) {
		c.announcements[i], c.announcements[j] = c.announcements[j], c.announcements[i]
	})

	return c
}

// GetAnnouncements 获取所有公告
func (c *Center) GetAnnouncements() []string {
	return c.announcements
}

// FormatAnnouncements 格式化公告为显示字符串
func (c *Center) FormatAnnouncements() string {
	var sb strings.Builder

	sb.WriteString("========================================\n")
	sb.WriteString("           市 中 心 公 告\n")
	sb.WriteString("========================================\n\n")

	for _, announcement := range c.announcements {
		sb.WriteString("  • ")
		sb.WriteString(announcement)
		sb.WriteString("\n\n")
	}

	sb.WriteString("========================================\n")

	return sb.String()
}

// AddAnnouncement 添加新公告
func (c *Center) AddAnnouncement(announcement string) {
	c.announcements = append(c.announcements, announcement)
}

// Count 返回公告数量
func (c *Center) Count() int {
	return len(c.announcements)
}
