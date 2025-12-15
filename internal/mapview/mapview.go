// Package mapview 提供城市地图渲染功能
package mapview

import (
	"fmt"
	"strings"

	"citylife/internal/world"
)

// GridPosition 网格坐标
type GridPosition struct {
	X int
	Y int
}

// GridPositions 位置ID到网格坐标的映射
// 布局:
//
//	      (1,0) 银行
//	        |
//	(0,1) 超市 -- (1,1) 市中心 -- (2,1) 医院
var GridPositions = map[int]GridPosition{
	world.LocationCityCenter:  {1, 1}, // 中心
	world.LocationSupermarket: {0, 1}, // 西边
	world.LocationBank:        {1, 0}, // 北边
	world.LocationHospital:    {2, 1}, // 东边
}

// GetGridPosition 获取建筑的网格坐标
func GetGridPosition(buildingID int) (int, int, bool) {
	if pos, ok := GridPositions[buildingID]; ok {
		return pos.X, pos.Y, true
	}
	return -1, -1, false
}

// IsAdjacent 判断两个位置是否相邻（曼哈顿距离为1）
func IsAdjacent(from, to int) bool {
	pos1, ok1 := GridPositions[from]
	pos2, ok2 := GridPositions[to]

	if !ok1 || !ok2 {
		return false
	}

	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y

	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	return dx+dy == 1
}

// RenderMap 渲染ASCII地图
func RenderMap(currentLocation int) string {
	var sb strings.Builder

	// 边框顶部
	sb.WriteString("╔══════════════════════════════════════════════╗\n")
	sb.WriteString("║             城  市  地  图                   ║\n")
	sb.WriteString("╠══════════════════════════════════════════════╣\n")

	// 北方向指示
	sb.WriteString("║                     北                       ║\n")
	sb.WriteString("║                     ↑                        ║\n")
	sb.WriteString("║                                              ║\n")

	// 银行 (北边)
	bankMarker := getBuildingMarker(world.LocationBank, currentLocation)
	sb.WriteString(fmt.Sprintf("║                 %s                       ║\n", bankMarker))
	sb.WriteString("║                     ║                        ║\n")

	// 主横线: 超市 - 市中心 - 医院
	superMarker := getBuildingMarker(world.LocationSupermarket, currentLocation)
	centerMarker := getBuildingMarker(world.LocationCityCenter, currentLocation)
	hospitalMarker := getBuildingMarker(world.LocationHospital, currentLocation)

	sb.WriteString(fmt.Sprintf("║  西 ← %s═══%s═══%s → 东 ║\n",
		superMarker, centerMarker, hospitalMarker))

	sb.WriteString("║                                              ║\n")

	// 南方向指示
	sb.WriteString("║                     ↓                        ║\n")
	sb.WriteString("║                     南                       ║\n")

	// 边框底部
	sb.WriteString("╚══════════════════════════════════════════════╝\n")

	// 当前位置信息
	sb.WriteString("\n")
	currentName := getBuildingName(currentLocation)
	sb.WriteString(fmt.Sprintf("  当前位置：★%s\n", currentName))

	// 可前往位置
	adjacent := world.GetAdjacentBuildings(currentLocation)
	if len(adjacent) > 0 {
		names := make([]string, 0, len(adjacent))
		for _, id := range adjacent {
			names = append(names, getBuildingName(id))
		}
		sb.WriteString(fmt.Sprintf("  可前往：%s\n", strings.Join(names, "、")))
	}

	return sb.String()
}

// getBuildingMarker 获取建筑标记（当前位置加★）
func getBuildingMarker(buildingID, currentLocation int) string {
	name := getBuildingName(buildingID)
	if buildingID == currentLocation {
		return fmt.Sprintf("[★%s]", name)
	}
	return fmt.Sprintf("[%s]", name)
}

// getBuildingName 获取建筑名称
func getBuildingName(buildingID int) string {
	if buildingID >= 0 && buildingID < len(world.BuildingNames) {
		return world.BuildingNames[buildingID]
	}
	return "未知"
}

// RenderSimpleMap 渲染简化地图（无边框）
func RenderSimpleMap(currentLocation int) string {
	var sb strings.Builder

	// 银行
	bankMarker := getSimpleMarker(world.LocationBank, currentLocation)
	sb.WriteString(fmt.Sprintf("        %s\n", bankMarker))
	sb.WriteString("          |\n")

	// 超市 - 市中心 - 医院
	superMarker := getSimpleMarker(world.LocationSupermarket, currentLocation)
	centerMarker := getSimpleMarker(world.LocationCityCenter, currentLocation)
	hospitalMarker := getSimpleMarker(world.LocationHospital, currentLocation)
	sb.WriteString(fmt.Sprintf("%s---%s---%s\n", superMarker, centerMarker, hospitalMarker))

	return sb.String()
}

// getSimpleMarker 获取简单标记
func getSimpleMarker(buildingID, currentLocation int) string {
	name := getBuildingName(buildingID)
	if buildingID == currentLocation {
		return fmt.Sprintf("*%s*", name)
	}
	return name
}
