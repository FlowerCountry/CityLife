// Package food 包含所有商品定义
package food

// AllCommodities 所有超市商品定义
var AllCommodities = []*Food{
	NewFood("面包", 15, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 18},
		{Name: "饥饿", Amount: 18},
		{Name: "蛋白质", Amount: 9},
		{Name: "维生素B", Amount: 9},
	}),
	NewFood("牛奶", 25, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 10},
		{Name: "饥饿", Amount: 16},
		{Name: "钙", Amount: 24},
		{Name: "蛋白质", Amount: 14},
		{Name: "维生素A", Amount: 12},
	}),
	NewFood("蛋糕", 40, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 24},
		{Name: "饥饿", Amount: 10},
		{Name: "糖分", Amount: 30},
		{Name: "脂肪", Amount: 20},
		{Name: "维生素E", Amount: 14},
	}),
	NewFood("苹果", 10, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 16},
		{Name: "维生素C", Amount: 20},
		{Name: "纤维素", Amount: 14},
		{Name: "钾", Amount: 10},
	}),
	NewFood("牛排", 150, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 20},
		{Name: "蛋白质", Amount: 25},
		{Name: "铁", Amount: 15},
		{Name: "脂肪", Amount: 8},
		{Name: "锌", Amount: 32},
	}),
	NewFood("橙汁", 18, FoodTypeBeverage, []Health{
		{Name: "饱腹感", Amount: 7},
		{Name: "维生素C", Amount: 22},
		{Name: "饥渴", Amount: 14},
		{Name: "糖分", Amount: 7},
		{Name: "钙", Amount: 10},
	}),
	NewFood("沙拉", 35, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 9},
		{Name: "维生素A", Amount: 30},
		{Name: "纤维素", Amount: 20},
		{Name: "钾", Amount: 14},
	}),
	NewFood("鸡蛋", 12, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 16},
		{Name: "蛋白质", Amount: 24},
		{Name: "硒", Amount: 12},
	}),
	NewFood("意大利面", 30, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 24},
		{Name: "碳水化合物", Amount: 30},
		{Name: "维生素B", Amount: 12},
	}),
	NewFood("咖啡", 20, FoodTypeBeverage, []Health{
		{Name: "饱腹感", Amount: 4},
		{Name: "精神振奋", Amount: 20},
		{Name: "饥渴", Amount: 10},
		{Name: "钾", Amount: 12},
	}),
	NewFood("巧克力", 25, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 6},
		{Name: "幸福感", Amount: 24},
		{Name: "糖分", Amount: 30},
		{Name: "脂肪", Amount: 16},
		{Name: "镁", Amount: 12},
	}),
	NewFood("番茄", 7, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 8},
		{Name: "维生素C", Amount: 16},
		{Name: "纤维素", Amount: 12},
		{Name: "钾", Amount: 8},
	}),
	NewFood("鸡胸肉", 70, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 24},
		{Name: "蛋白质", Amount: 36},
		{Name: "脂肪", Amount: 6},
		{Name: "铁", Amount: 26},
	}),
	NewFood("矿泉水", 3, FoodTypeBeverage, []Health{
		{Name: "饱腹感", Amount: 2},
		{Name: "饥渴", Amount: 45},
		{Name: "钙", Amount: 5},
	}),
	NewFood("米饭", 15, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 24},
		{Name: "碳水化合物", Amount: 36},
		{Name: "维生素B", Amount: 12},
	}),
	NewFood("燕麦片", 22, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 28},
		{Name: "纤维素", Amount: 20},
		{Name: "钾", Amount: 14},
	}),
	NewFood("酸奶", 18, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 12},
		{Name: "钙", Amount: 24},
		{Name: "蛋白质", Amount: 16},
		{Name: "益生菌", Amount: 16},
	}),
	NewFood("三明治", 40, FoodTypeProcessed, []Health{
		{Name: "饱腹感", Amount: 32},
		{Name: "蛋白质", Amount: 20},
		{Name: "铁", Amount: 18},
	}),
	NewFood("薯片", 12, FoodTypeCanned, []Health{
		{Name: "饱腹感", Amount: 6},
		{Name: "脂肪", Amount: 24},
		{Name: "碳水化合物", Amount: 16},
		{Name: "维生素C", Amount: 12},
	}),
	NewFood("冰淇淋", 35, FoodTypeFresh, []Health{
		{Name: "饱腹感", Amount: 6},
		{Name: "幸福感", Amount: 30},
		{Name: "糖分", Amount: 36},
		{Name: "脂肪", Amount: 20},
		{Name: "维生素D", Amount: 14},
	}),
}

// GetCommodityByName 根据名称获取商品
func GetCommodityByName(name string) *Food {
	for _, c := range AllCommodities {
		if c.Name == name {
			return c
		}
	}
	return nil
}
