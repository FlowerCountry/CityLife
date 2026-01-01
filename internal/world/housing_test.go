package world

import (
	"strings"
	"testing"
)

func TestHousingAdvance_RentExpiresViaWorldUpdateTime(t *testing.T) {
	w := New()
	w.Housing.SetRented(HomeBed, 60*60) // 1小时

	w.UpdateTime(60*60 - 1)
	if !w.Housing.HasHome() {
		t.Fatal("租期未到，不应失去住房")
	}
	if w.Housing.RentRemainingSeconds != 1 {
		t.Fatalf("剩余租期=%d，期望1", w.Housing.RentRemainingSeconds)
	}

	w.UpdateTime(2)
	if w.Housing.HasHome() {
		t.Fatal("租期到期后，应失去住房")
	}
	if w.Housing.Status != HousingStatusNone || w.Housing.Level != HomeNone || w.Housing.RentRemainingSeconds != 0 {
		t.Fatalf("到期后状态未清空：status=%d level=%d remaining=%d", w.Housing.Status, w.Housing.Level, w.Housing.RentRemainingSeconds)
	}
}

func TestHousingAdvance_OwnedNotAffected(t *testing.T) {
	w := New()
	w.Housing.SetOwned(HomeRoom)

	w.UpdateTime(30 * 24 * 60 * 60) // 30天
	if !w.Housing.IsOwned() {
		t.Fatal("自有住房不应被时间推进影响")
	}
	if w.Housing.Level != HomeRoom {
		t.Fatalf("自有住房档位被改变：got=%d want=%d", w.Housing.Level, HomeRoom)
	}
	if w.Housing.RentRemainingSeconds != 0 {
		t.Fatalf("自有住房不应有租期：got=%d", w.Housing.RentRemainingSeconds)
	}
}

func TestHousingSummary(t *testing.T) {
	t.Run("无住房", func(t *testing.T) {
		var h Housing
		if got := h.Summary(); got != "住房：无" {
			t.Fatalf("Summary()=%q，期望%q", got, "住房：无")
		}
	})

	t.Run("租房", func(t *testing.T) {
		h := Housing{Status: HousingStatusRented, Level: HomeApartment, RentRemainingSeconds: 3 * 24 * 60 * 60}
		got := h.Summary()
		if !strings.Contains(got, "小公寓") || !strings.Contains(got, "租") {
			t.Fatalf("Summary()=%q，期望包含租房与档位信息", got)
		}
	})

	t.Run("买房", func(t *testing.T) {
		h := Housing{Status: HousingStatusOwned, Level: HomeRoom}
		got := h.Summary()
		if !strings.Contains(got, "单间") || !strings.Contains(got, "自有") {
			t.Fatalf("Summary()=%q，期望包含自有与档位信息", got)
		}
	})
}
