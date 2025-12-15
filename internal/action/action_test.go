package action

import (
	"testing"

	"citylife/internal/world"
)

func TestSortActions(t *testing.T) {
	// 创建不同类别的行动
	actions := []Action{
		&CheckCashAction{},            // Insight
		&GoWhereAction{TargetID: 1},   // Navigation
		&EnterSupermarketAction{},     // Primary
		&ViewMapAction{},              // Insight
		&InformationAction{},          // Insight
	}

	sorted := SortActions(actions)

	t.Run("Primary在前", func(t *testing.T) {
		if sorted[0].Category() != CategoryPrimary {
			t.Errorf("First action should be Primary, got %d", sorted[0].Category())
		}
	})

	t.Run("Navigation在后", func(t *testing.T) {
		last := sorted[len(sorted)-1]
		if last.Category() != CategoryNavigation {
			t.Errorf("Last action should be Navigation, got %d", last.Category())
		}
	})
}

func TestEventCategory(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		expected EventCategory
	}{
		{"EnterSupermarket", &EnterSupermarketAction{}, CategoryPrimary},
		{"CheckCash", &CheckCashAction{}, CategoryInsight},
		{"CheckBank", &CheckBankAction{}, CategoryInsight},
		{"Information", &InformationAction{}, CategoryInsight},
		{"ViewMap", &ViewMapAction{}, CategoryInsight},
		{"ViewHealth", &ViewHealthAction{}, CategoryInsight},
		{"GoWhere", &GoWhereAction{}, CategoryNavigation},
		{"LeaveSupermarket", &LeaveSupermarketAction{}, CategoryNavigation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.Category(); got != tt.expected {
				t.Errorf("Category() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestActionInfo(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		contains string
	}{
		{"GoWhere", &GoWhereAction{TargetName: "超市"}, "超市"},
		{"CheckCash", &CheckCashAction{}, "钱包"},
		{"CheckBank", &CheckBankAction{}, "银行"},
		{"EnterSupermarket", &EnterSupermarketAction{}, "超市"},
		{"LeaveSupermarket", &LeaveSupermarketAction{}, "离开"},
		{"Information", &InformationAction{}, "公告"},
		{"ViewMap", &ViewMapAction{}, "地图"},
		{"ViewHealth", &ViewHealthAction{}, "健康"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := tt.action.Info()
			if info == "" {
				t.Error("Info() should not return empty string")
			}
		})
	}
}

func TestGetActionsForLocation(t *testing.T) {
	tests := []struct {
		location int
		minCount int
	}{
		{world.LocationCityCenter, 5},
		{world.LocationSupermarket, 2},
		{world.LocationBank, 5},
		{world.LocationSupermarketInner, 10},
		{world.LocationHospital, 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actions := GetActionsForLocation(tt.location)
			if len(actions) < tt.minCount {
				t.Errorf("Location %d: got %d actions, want >= %d",
					tt.location, len(actions), tt.minCount)
			}
		})
	}
}

func TestGoWhereActionInfo(t *testing.T) {
	action := &GoWhereAction{
		TargetID:   1,
		TargetName: "超市",
	}

	info := action.Info()
	if info != "前往超市" {
		t.Errorf("Info() = %q, want 前往超市", info)
	}
}

func TestActionInterface(t *testing.T) {
	// 验证所有行动都实现了Action接口
	actions := []Action{
		&GoWhereAction{},
		&CheckCashAction{},
		&CheckBankAction{},
		&EnterSupermarketAction{},
		&LeaveSupermarketAction{},
		&InformationAction{},
		&ViewMapAction{},
		&ViewHealthAction{},
	}

	for i, a := range actions {
		t.Run("", func(t *testing.T) {
			_ = a.Info()
			_ = a.Category()
			// Execute需要game.State，这里只验证接口
			if a == nil {
				t.Errorf("Action %d is nil", i)
			}
		})
	}
}

func TestCategoryConstants(t *testing.T) {
	// 验证类别常量的值
	if CategoryPrimary != 0 {
		t.Errorf("CategoryPrimary = %d, want 0", CategoryPrimary)
	}
	if CategoryInsight != 1 {
		t.Errorf("CategoryInsight = %d, want 1", CategoryInsight)
	}
	if CategoryNavigation != 2 {
		t.Errorf("CategoryNavigation = %d, want 2", CategoryNavigation)
	}
}
