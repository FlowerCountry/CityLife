package world

// LocationID 位置ID常量
const (
	LocationCityCenter       = 0 // 市中心
	LocationSupermarket      = 1 // 超市（外部）
	LocationBank             = 2 // 银行
	LocationSupermarketInner = 3 // 超市内部
	LocationHospital         = 4 // 医院
)

// Building 建筑
type Building struct {
	ID   int
	X    int
	Y    int
	Name string
}

// BuildingNames 建筑名称
var BuildingNames = []string{
	"市中心",
	"超市",
	"银行",
	"超市内部",
	"医院",
}

// initBuildings 初始化建筑列表
func initBuildings() []*Building {
	return []*Building{
		{ID: LocationCityCenter, X: 0, Y: 0, Name: "市中心"},
		{ID: LocationSupermarket, X: 1, Y: 0, Name: "超市"},
		{ID: LocationBank, X: 0, Y: 1, Name: "银行"},
		{ID: LocationSupermarketInner, X: 1, Y: 0, Name: "超市内部"},
		{ID: LocationHospital, X: 1, Y: 1, Name: "医院"},
	}
}

// Adjacency 邻接关系表：从位置A可以到达的位置列表
var Adjacency = map[int][]int{
	LocationCityCenter:       {LocationSupermarket, LocationBank, LocationHospital},
	LocationSupermarket:      {LocationCityCenter, LocationBank, LocationHospital},
	LocationBank:             {LocationCityCenter, LocationSupermarket, LocationHospital},
	LocationSupermarketInner: {}, // 从超市内部只能通过"离开"返回超市外部
	LocationHospital:         {LocationCityCenter, LocationSupermarket, LocationBank},
}

// GetAdjacentBuildings 获取从当前位置可到达的建筑ID列表
func GetAdjacentBuildings(currentLocation int) []int {
	if adj, ok := Adjacency[currentLocation]; ok {
		return adj
	}
	return []int{}
}

// GetDistance 获取两个位置之间的距离（用于计算时间消耗）
func GetDistance(from, to int) int {
	// 简化距离计算：基于位置ID差异
	buildings := initBuildings()
	if from < 0 || from >= len(buildings) || to < 0 || to >= len(buildings) {
		return 10 // 默认距离
	}

	b1 := buildings[from]
	b2 := buildings[to]

	// 曼哈顿距离
	dx := b1.X - b2.X
	dy := b1.Y - b2.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	distance := dx + dy
	if distance == 0 {
		distance = 1
	}
	return distance * 5 // 每单位距离5分钟
}
