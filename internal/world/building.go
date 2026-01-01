package world

// LocationID 位置ID常量
const (
	LocationCityCenter       = 0 // 市中心
	LocationSupermarket      = 1 // 超市（外部）
	LocationBank             = 2 // 银行
	LocationSupermarketInner = 3 // 超市内部
	LocationHospital         = 4 // 医院
	// 新地点只允许追加，严禁重排（保证旧存档Where不炸）
	LocationResidentialArea  = 5  // 住宅区（外部）
	LocationHome             = 6  // 家（室内）
	LocationRealEstateAgency = 7  // 房产中介
	LocationJobMarket        = 8  // 人才市场
	LocationRestaurant       = 9  // 餐馆
	LocationPark             = 10 // 公园
	LocationHotel            = 11 // 旅馆
)

// Building 建筑
type Building struct {
	ID         int
	X          int
	Y          int
	Name       string
	IsInterior bool // 室内地点不在地图上展示，且通常不允许直接导航
}

// BuildingNames 建筑名称
var BuildingNames = []string{
	"市中心",
	"超市",
	"银行",
	"超市内部",
	"医院",
	"住宅区",
	"家",
	"房产中介",
	"人才市场",
	"餐馆",
	"公园",
	"旅馆",
}

// initBuildings 初始化建筑列表
func initBuildings() []*Building {
	return []*Building{
		{ID: LocationCityCenter, X: 0, Y: 0, Name: "市中心"},
		{ID: LocationSupermarket, X: 1, Y: 0, Name: "超市"},
		{ID: LocationBank, X: 0, Y: 1, Name: "银行"},
		{ID: LocationSupermarketInner, X: 1, Y: 0, Name: "超市内部", IsInterior: true},
		{ID: LocationHospital, X: 1, Y: 1, Name: "医院"},
		{ID: LocationResidentialArea, X: -1, Y: 0, Name: "住宅区"},
		{ID: LocationHome, X: -1, Y: 0, Name: "家", IsInterior: true},
		{ID: LocationRealEstateAgency, X: -1, Y: -1, Name: "房产中介"},
		{ID: LocationJobMarket, X: 0, Y: 2, Name: "人才市场"},
		{ID: LocationRestaurant, X: 0, Y: -1, Name: "餐馆"},
		{ID: LocationPark, X: -1, Y: 1, Name: "公园"},
		{ID: LocationHotel, X: 2, Y: 0, Name: "旅馆"},
	}
}

// Adjacency 邻接关系表：从位置A可以到达的位置列表
var Adjacency = map[int][]int{
	LocationCityCenter:       {LocationSupermarket, LocationBank, LocationHospital, LocationResidentialArea, LocationRealEstateAgency, LocationJobMarket, LocationRestaurant, LocationPark, LocationHotel},
	LocationSupermarket:      {LocationCityCenter, LocationBank, LocationHospital},
	LocationBank:             {LocationCityCenter, LocationSupermarket, LocationHospital},
	LocationSupermarketInner: {}, // 从超市内部只能通过"离开"返回超市外部
	LocationHospital:         {LocationCityCenter, LocationSupermarket, LocationBank},
	LocationResidentialArea:  {LocationCityCenter, LocationPark, LocationRealEstateAgency},
	LocationHome:             {}, // 室内地点：只能通过"回家/出门"切换
	LocationRealEstateAgency: {LocationCityCenter, LocationResidentialArea},
	LocationJobMarket:        {LocationCityCenter},
	LocationRestaurant:       {LocationCityCenter},
	LocationPark:             {LocationCityCenter, LocationResidentialArea},
	LocationHotel:            {LocationCityCenter},
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
