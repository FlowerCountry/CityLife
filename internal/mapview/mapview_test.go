package mapview

import (
	"strings"
	"testing"

	"citylife/internal/world"
)

func TestGetGridPosition(t *testing.T) {
	tests := []struct {
		buildingID int
		wantX      int
		wantY      int
		wantOk     bool
	}{
		{world.LocationCityCenter, 1, 1, true},
		{world.LocationSupermarket, 0, 1, true},
		{world.LocationBank, 1, 0, true},
		{world.LocationHospital, 2, 1, true},
		{world.LocationSupermarketInner, -1, -1, false}, // 不在地图上
		{999, -1, -1, false},                            // 无效ID
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			x, y, ok := GetGridPosition(tt.buildingID)
			if x != tt.wantX || y != tt.wantY || ok != tt.wantOk {
				t.Errorf("GetGridPosition(%d) = (%d, %d, %v), want (%d, %d, %v)",
					tt.buildingID, x, y, ok, tt.wantX, tt.wantY, tt.wantOk)
			}
		})
	}
}

func TestIsAdjacent(t *testing.T) {
	tests := []struct {
		from     int
		to       int
		expected bool
	}{
		// 相邻
		{world.LocationCityCenter, world.LocationSupermarket, true},
		{world.LocationCityCenter, world.LocationBank, true},
		{world.LocationCityCenter, world.LocationHospital, true},

		// 不相邻
		{world.LocationSupermarket, world.LocationHospital, false}, // 对角
		{world.LocationBank, world.LocationSupermarket, false},     // 对角
		{world.LocationBank, world.LocationHospital, false},        // 对角

		// 同位置
		{world.LocationCityCenter, world.LocationCityCenter, false},

		// 无效位置
		{world.LocationSupermarketInner, world.LocationCityCenter, false},
		{999, world.LocationCityCenter, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := IsAdjacent(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("IsAdjacent(%d, %d) = %v, want %v",
					tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestRenderMap(t *testing.T) {
	t.Run("市中心", func(t *testing.T) {
		output := RenderMap(world.LocationCityCenter)

		// 检查标题
		if !strings.Contains(output, "城  市  地  图") {
			t.Error("Should contain map title")
		}

		// 检查当前位置标记
		if !strings.Contains(output, "★市中心") {
			t.Error("Should mark current location with ★")
		}

		// 检查边框
		if !strings.Contains(output, "╔") {
			t.Error("Should contain border")
		}

		// 检查方向指示
		if !strings.Contains(output, "北") {
			t.Error("Should contain North indicator")
		}
		if !strings.Contains(output, "南") {
			t.Error("Should contain South indicator")
		}
		if !strings.Contains(output, "东") {
			t.Error("Should contain East indicator")
		}
		if !strings.Contains(output, "西") {
			t.Error("Should contain West indicator")
		}

		// 检查可前往位置
		if !strings.Contains(output, "可前往") {
			t.Error("Should show adjacent locations")
		}
	})

	t.Run("银行", func(t *testing.T) {
		output := RenderMap(world.LocationBank)

		// 检查当前位置
		if !strings.Contains(output, "★银行") {
			t.Error("Should mark current location as Bank")
		}

		// 银行不应该被标记两次
		count := strings.Count(output, "★")
		if count != 2 { // 一次在地图中，一次在"当前位置"
			t.Errorf("★ should appear exactly 2 times, got %d", count)
		}
	})

	t.Run("超市", func(t *testing.T) {
		output := RenderMap(world.LocationSupermarket)

		if !strings.Contains(output, "★超市") {
			t.Error("Should mark current location as Supermarket")
		}
	})

	t.Run("医院", func(t *testing.T) {
		output := RenderMap(world.LocationHospital)

		if !strings.Contains(output, "★医院") {
			t.Error("Should mark current location as Hospital")
		}
	})
}

func TestRenderSimpleMap(t *testing.T) {
	output := RenderSimpleMap(world.LocationCityCenter)

	t.Run("包含所有位置", func(t *testing.T) {
		if !strings.Contains(output, "市中心") {
			t.Error("Should contain 市中心")
		}
		if !strings.Contains(output, "超市") {
			t.Error("Should contain 超市")
		}
		if !strings.Contains(output, "银行") {
			t.Error("Should contain 银行")
		}
		if !strings.Contains(output, "医院") {
			t.Error("Should contain 医院")
		}
	})

	t.Run("当前位置标记", func(t *testing.T) {
		if !strings.Contains(output, "*市中心*") {
			t.Error("Should mark current location with *")
		}
	})

	t.Run("无边框", func(t *testing.T) {
		if strings.Contains(output, "╔") {
			t.Error("Simple map should not have border")
		}
	})
}

func TestGetBuildingMarker(t *testing.T) {
	t.Run("当前位置", func(t *testing.T) {
		marker := getBuildingMarker(world.LocationCityCenter, world.LocationCityCenter)
		if !strings.Contains(marker, "★") {
			t.Error("Current location should have ★")
		}
		if !strings.Contains(marker, "市中心") {
			t.Error("Should contain building name")
		}
	})

	t.Run("非当前位置", func(t *testing.T) {
		marker := getBuildingMarker(world.LocationBank, world.LocationCityCenter)
		if strings.Contains(marker, "★") {
			t.Error("Non-current location should not have ★")
		}
		if !strings.Contains(marker, "银行") {
			t.Error("Should contain building name")
		}
	})
}

func TestGetBuildingName(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{world.LocationCityCenter, "市中心"},
		{world.LocationSupermarket, "超市"},
		{world.LocationBank, "银行"},
		{world.LocationHospital, "医院"},
		{999, "未知"},
		{-1, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBuildingName(tt.id)
			if result != tt.name {
				t.Errorf("getBuildingName(%d) = %s, want %s", tt.id, result, tt.name)
			}
		})
	}
}
