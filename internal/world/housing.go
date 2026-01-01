package world

import "fmt"

// HousingStatus 住房状态
type HousingStatus int

const (
	HousingStatusNone HousingStatus = iota
	HousingStatusRented
	HousingStatusOwned
)

// HomeLevel 住房档位（3档）
type HomeLevel int

const (
	HomeNone HomeLevel = iota
	HomeBed
	HomeRoom
	HomeApartment
)

// Housing 住房数据（单槽位：要么租要么买）
type Housing struct {
	Status               HousingStatus
	Level                HomeLevel
	RentRemainingSeconds int // 仅租房使用
}

func (h Housing) HasHome() bool {
	switch h.Status {
	case HousingStatusOwned:
		return h.Level != HomeNone
	case HousingStatusRented:
		return h.Level != HomeNone && h.RentRemainingSeconds > 0
	default:
		return false
	}
}

func (h Housing) IsRented() bool { return h.Status == HousingStatusRented && h.HasHome() }

func (h Housing) IsOwned() bool { return h.Status == HousingStatusOwned && h.HasHome() }

func (h Housing) StatusName() string {
	switch h.Status {
	case HousingStatusRented:
		return "租住"
	case HousingStatusOwned:
		return "自有"
	default:
		return "无"
	}
}

func (h Housing) HomeName() string {
	switch h.Level {
	case HomeBed:
		return "合租床位"
	case HomeRoom:
		return "单间"
	case HomeApartment:
		return "小公寓"
	default:
		return "无"
	}
}

func (h Housing) Summary() string {
	if !h.HasHome() {
		return "住房：无"
	}

	if h.IsOwned() {
		return fmt.Sprintf("住房：%s（自有）", h.HomeName())
	}

	return fmt.Sprintf("住房：%s（租住，剩余%s）", h.HomeName(), formatDurationChinese(h.RentRemainingSeconds))
}

func (h *Housing) Clear() {
	h.Status = HousingStatusNone
	h.Level = HomeNone
	h.RentRemainingSeconds = 0
}

func (h *Housing) SetOwned(level HomeLevel) {
	h.Status = HousingStatusOwned
	h.Level = level
	h.RentRemainingSeconds = 0
}

func (h *Housing) SetRented(level HomeLevel, rentSeconds int) {
	h.Status = HousingStatusRented
	h.Level = level
	h.RentRemainingSeconds = rentSeconds
	if rentSeconds <= 0 {
		h.Clear()
	}
}

// Advance 随时间推进（仅影响租房）
func (h *Housing) Advance(seconds int) (expired bool) {
	if seconds <= 0 {
		return false
	}
	if h.Status != HousingStatusRented {
		return false
	}
	if h.RentRemainingSeconds <= 0 {
		h.Clear()
		return true
	}

	h.RentRemainingSeconds -= seconds
	if h.RentRemainingSeconds <= 0 {
		h.Clear()
		return true
	}
	return false
}

func formatDurationChinese(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	const (
		day  = 24 * 60 * 60
		hour = 60 * 60
		min  = 60
	)

	d := seconds / day
	seconds %= day
	h := seconds / hour
	seconds %= hour
	m := seconds / min

	if d > 0 {
		if h > 0 {
			return fmt.Sprintf("%d天%d小时", d, h)
		}
		return fmt.Sprintf("%d天", d)
	}
	if h > 0 {
		if m > 0 {
			return fmt.Sprintf("%d小时%d分钟", h, m)
		}
		return fmt.Sprintf("%d小时", h)
	}
	return fmt.Sprintf("%d分钟", m)
}
