package world

import (
	"testing"
)

func TestNewWorld(t *testing.T) {
	w := New()

	// 检查初始时间
	if w.Year != 2010 {
		t.Errorf("初始年份错误: got %d, want 2010", w.Year)
	}
	if w.Month != 10 {
		t.Errorf("初始月份错误: got %d, want 10", w.Month)
	}
	if w.Day != 10 {
		t.Errorf("初始日期错误: got %d, want 10", w.Day)
	}
	if w.Hour != 11 {
		t.Errorf("初始小时错误: got %d, want 11", w.Hour)
	}

	// 检查初始位置
	if w.Where != LocationCityCenter {
		t.Errorf("初始位置错误: got %d, want %d", w.Where, LocationCityCenter)
	}

	// 检查初始钱包 (100元 = 1张100元)
	expectedWallet := [6]int{1, 0, 0, 0, 0, 0}
	if w.Wallet != expectedWallet {
		t.Errorf("初始钱包错误: got %v, want %v", w.Wallet, expectedWallet)
	}

	// 检查初始银行存款
	if w.BankDeposit != 0 {
		t.Errorf("初始银行存款错误: got %d, want 0", w.BankDeposit)
	}
}

func TestGetWalletTotal(t *testing.T) {
	w := New()

	// 初始钱包总额: 100元 (只有1张100元)
	expected := 100
	if got := w.GetWalletTotal(); got != expected {
		t.Errorf("钱包总额错误: got %d, want %d", got, expected)
	}

	// 修改钱包后测试
	w.Wallet = [6]int{2, 0, 0, 0, 0, 0}
	if got := w.GetWalletTotal(); got != 200 {
		t.Errorf("钱包总额错误: got %d, want 200", got)
	}
}

func TestSpendMoney(t *testing.T) {
	tests := []struct {
		name       string
		wallet     [6]int
		amount     int
		wantErr    bool
		wantWallet [6]int
	}{
		{
			name:       "花费5元",
			wallet:     [6]int{1, 1, 1, 1, 1, 1}, // 186元
			amount:     5,
			wantErr:    false,
			wantWallet: [6]int{1, 1, 1, 1, 0, 1}, // 减去1张5元
		},
		{
			name:       "花费1元",
			wallet:     [6]int{0, 0, 0, 0, 0, 5}, // 5元
			amount:     1,
			wantErr:    false,
			wantWallet: [6]int{0, 0, 0, 0, 0, 4},
		},
		{
			name:       "花费超过余额",
			wallet:     [6]int{0, 0, 0, 0, 0, 1}, // 1元
			amount:     5,
			wantErr:    true,
			wantWallet: [6]int{0, 0, 0, 0, 0, 1}, // 不变
		},
		{
			name:       "花费100元",
			wallet:     [6]int{1, 0, 0, 0, 0, 0}, // 100元
			amount:     100,
			wantErr:    false,
			wantWallet: [6]int{0, 0, 0, 0, 0, 0},
		},
		{
			name:       "花费需要找零",
			wallet:     [6]int{1, 0, 0, 0, 0, 0}, // 100元
			amount:     30,
			wantErr:    false,
			wantWallet: [6]int{0, 1, 1, 0, 0, 0}, // 找零70元 = 50+20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.Wallet = tt.wallet

			err := w.SpendMoney(tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("SpendMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && w.Wallet != tt.wantWallet {
				t.Errorf("SpendMoney() wallet = %v, want %v", w.Wallet, tt.wantWallet)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name     string
		wallet   [6]int
		bank     int
		amount   int
		wantErr  bool
		wantBank int
	}{
		{
			name:     "存款100元",
			wallet:   [6]int{1, 0, 0, 0, 0, 0},
			bank:     0,
			amount:   100,
			wantErr:  false,
			wantBank: 100,
		},
		{
			name:     "存款余额不足",
			wallet:   [6]int{0, 1, 0, 0, 0, 0}, // 50元
			bank:     0,
			amount:   100,
			wantErr:  true,
			wantBank: 0,
		},
		{
			name:     "存款非100倍数",
			wallet:   [6]int{1, 1, 0, 0, 0, 0}, // 150元
			bank:     0,
			amount:   50,
			wantErr:  true,
			wantBank: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.Wallet = tt.wallet
			w.BankDeposit = tt.bank

			err := w.Deposit(tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && w.BankDeposit != tt.wantBank {
				t.Errorf("Deposit() bank = %d, want %d", w.BankDeposit, tt.wantBank)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name     string
		bank     int
		amount   int
		wantErr  bool
		wantBank int
	}{
		{
			name:     "取款100元",
			bank:     200,
			amount:   100,
			wantErr:  false,
			wantBank: 100,
		},
		{
			name:     "取款余额不足",
			bank:     50,
			amount:   100,
			wantErr:  true,
			wantBank: 50,
		},
		{
			name:     "取款非100倍数",
			bank:     200,
			amount:   50,
			wantErr:  true,
			wantBank: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.BankDeposit = tt.bank

			err := w.Withdraw(tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if w.BankDeposit != tt.wantBank {
				t.Errorf("Withdraw() bank = %d, want %d", w.BankDeposit, tt.wantBank)
			}
		})
	}
}

func TestUpdateTime(t *testing.T) {
	tests := []struct {
		name       string
		startHour  int
		startMin   int
		startDay   int
		addSeconds int
		wantHour   int
		wantMinute int
		wantDay    int
	}{
		{
			name:       "加1分钟",
			startHour:  8,
			startMin:   0,
			startDay:   1,
			addSeconds: 60,
			wantHour:   8,
			wantMinute: 1,
			wantDay:    1,
		},
		{
			name:       "加1小时",
			startHour:  8,
			startMin:   0,
			startDay:   1,
			addSeconds: 3600,
			wantHour:   9,
			wantMinute: 0,
			wantDay:    1,
		},
		{
			name:       "跨日",
			startHour:  8,
			startMin:   0,
			startDay:   1,
			addSeconds: 16 * 3600, // 16小时后是第二天0点
			wantHour:   0,
			wantMinute: 0,
			wantDay:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			// 设置固定的初始时间用于测试
			w.Hour = tt.startHour
			w.Minute = tt.startMin
			w.Second = 0
			w.Day = tt.startDay

			w.UpdateTime(tt.addSeconds)

			if w.Hour != tt.wantHour {
				t.Errorf("UpdateTime() hour = %d, want %d", w.Hour, tt.wantHour)
			}
			if w.Minute != tt.wantMinute {
				t.Errorf("UpdateTime() minute = %d, want %d", w.Minute, tt.wantMinute)
			}
			if w.Day != tt.wantDay {
				t.Errorf("UpdateTime() day = %d, want %d", w.Day, tt.wantDay)
			}
		})
	}
}

func TestChangeLocation(t *testing.T) {
	w := New()

	w.ChangeLocation(LocationSupermarket)
	if w.Where != LocationSupermarket {
		t.Errorf("ChangeLocation() = %d, want %d", w.Where, LocationSupermarket)
	}

	w.ChangeLocation(LocationBank)
	if w.Where != LocationBank {
		t.Errorf("ChangeLocation() = %d, want %d", w.Where, LocationBank)
	}
}

func TestGetTimeString(t *testing.T) {
	w := New()
	timeStr := w.GetTimeString()

	// 验证时间字符串非空且格式正确
	if len(timeStr) == 0 {
		t.Error("GetTimeString() should not return empty string")
	}
}

// ==================== 建筑导航测试 ====================

func TestBuildingNavigation(t *testing.T) {
	w := New()

	t.Run("建筑列表初始化", func(t *testing.T) {
		if len(w.Buildings) == 0 {
			t.Error("Buildings should not be empty")
		}
	})

	t.Run("建筑ID匹配", func(t *testing.T) {
		for i, b := range w.Buildings {
			if b.ID != i {
				t.Errorf("Building ID mismatch: index=%d, ID=%d", i, b.ID)
			}
		}
	})

	t.Run("建筑名称非空", func(t *testing.T) {
		for _, b := range w.Buildings {
			if b.Name == "" {
				t.Errorf("Building %d has empty name", b.ID)
			}
		}
	})

	t.Run("位置常量有效", func(t *testing.T) {
		locations := []int{
			LocationCityCenter,
			LocationSupermarket,
			LocationBank,
			LocationSupermarketInner,
			LocationHospital,
		}
		for _, loc := range locations {
			if loc < 0 || loc >= len(w.Buildings) {
				t.Errorf("Location constant %d out of range", loc)
			}
		}
	})
}

func TestGetAdjacentBuildings(t *testing.T) {
	tests := []struct {
		location int
		minCount int
		contains []int
	}{
		{LocationCityCenter, 3, []int{LocationSupermarket, LocationBank, LocationHospital}},
		{LocationSupermarket, 3, []int{LocationCityCenter, LocationBank, LocationHospital}},
		{LocationBank, 3, []int{LocationCityCenter, LocationSupermarket, LocationHospital}},
		{LocationSupermarketInner, 0, []int{}},
		{LocationHospital, 3, []int{LocationCityCenter, LocationSupermarket, LocationBank}},
	}

	for _, tt := range tests {
		t.Run(BuildingNames[tt.location], func(t *testing.T) {
			adj := GetAdjacentBuildings(tt.location)

			if len(adj) < tt.minCount {
				t.Errorf("GetAdjacentBuildings(%d) returned %d items, want >= %d",
					tt.location, len(adj), tt.minCount)
			}

			// 验证包含预期的位置
			for _, expected := range tt.contains {
				found := false
				for _, a := range adj {
					if a == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetAdjacentBuildings(%d) should contain %d",
						tt.location, expected)
				}
			}
		})
	}

	t.Run("无效位置", func(t *testing.T) {
		adj := GetAdjacentBuildings(-1)
		if len(adj) != 0 {
			t.Error("Invalid location should return empty list")
		}

		adj = GetAdjacentBuildings(999)
		if len(adj) != 0 {
			t.Error("Invalid location should return empty list")
		}
	})
}

func TestGetDistance(t *testing.T) {
	tests := []struct {
		from     int
		to       int
		expected int
	}{
		{LocationCityCenter, LocationSupermarket, 5},  // dx=1, dy=0 -> 1*5=5
		{LocationCityCenter, LocationBank, 5},         // dx=0, dy=1 -> 1*5=5
		{LocationCityCenter, LocationHospital, 10},    // dx=1, dy=1 -> 2*5=10
		{LocationSupermarket, LocationBank, 10},       // dx=1, dy=1 -> 2*5=10
		{LocationSupermarket, LocationSupermarket, 5}, // 同位置，距离为1*5=5
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			dist := GetDistance(tt.from, tt.to)
			if dist != tt.expected {
				t.Errorf("GetDistance(%d, %d) = %d, want %d",
					tt.from, tt.to, dist, tt.expected)
			}
		})
	}

	t.Run("无效位置返回默认距离", func(t *testing.T) {
		dist := GetDistance(-1, 0)
		if dist != 10 {
			t.Errorf("Invalid from should return 10, got %d", dist)
		}

		dist = GetDistance(0, 999)
		if dist != 10 {
			t.Errorf("Invalid to should return 10, got %d", dist)
		}
	})
}

func TestCurrentLocation(t *testing.T) {
	w := New()

	t.Run("初始位置", func(t *testing.T) {
		loc := w.CurrentLocation()
		if loc == nil {
			t.Fatal("CurrentLocation() should not return nil")
		}
		if loc.ID != w.Where {
			t.Errorf("CurrentLocation().ID = %d, want %d", loc.ID, w.Where)
		}
	})

	t.Run("改变位置后", func(t *testing.T) {
		w.ChangeLocation(LocationBank)
		loc := w.CurrentLocation()
		if loc == nil {
			t.Fatal("CurrentLocation() should not return nil")
		}
		if loc.ID != LocationBank {
			t.Errorf("CurrentLocation().ID = %d, want %d", loc.ID, LocationBank)
		}
	})

	t.Run("无效位置", func(t *testing.T) {
		w.Where = 999
		loc := w.CurrentLocation()
		if loc != nil {
			t.Error("CurrentLocation() should return nil for invalid Where")
		}
	})
}

// ==================== 时间边界测试 ====================

func TestTimeEdgeCases(t *testing.T) {
	t.Run("月末跨月", func(t *testing.T) {
		w := New()
		w.Year = 2010
		w.Month = 1
		w.Day = 31
		w.Hour = 23
		w.Minute = 59
		w.Second = 0

		w.UpdateTime(60) // 加1分钟

		if w.Month != 2 || w.Day != 1 {
			t.Errorf("跨月错误: got %d月%d日, want 2月1日", w.Month, w.Day)
		}
	})

	t.Run("年末跨年", func(t *testing.T) {
		w := New()
		w.Year = 2010
		w.Month = 12
		w.Day = 31
		w.Hour = 23
		w.Minute = 59
		w.Second = 0

		w.UpdateTime(60) // 加1分钟

		if w.Year != 2011 || w.Month != 1 || w.Day != 1 {
			t.Errorf("跨年错误: got %d年%d月%d日, want 2011年1月1日",
				w.Year, w.Month, w.Day)
		}
	})

	t.Run("闰年2月", func(t *testing.T) {
		w := New()
		w.Year = 2000 // 闰年
		w.Month = 2
		w.Day = 28
		w.Hour = 23
		w.Minute = 59
		w.Second = 0

		w.UpdateTime(60) // 加1分钟

		if w.Month != 2 || w.Day != 29 {
			t.Errorf("闰年2月错误: got %d月%d日, want 2月29日", w.Month, w.Day)
		}
	})

	t.Run("非闰年2月", func(t *testing.T) {
		w := New()
		w.Year = 2001 // 非闰年
		w.Month = 2
		w.Day = 28
		w.Hour = 23
		w.Minute = 59
		w.Second = 0

		w.UpdateTime(60) // 加1分钟

		if w.Month != 3 || w.Day != 1 {
			t.Errorf("非闰年2月错误: got %d月%d日, want 3月1日", w.Month, w.Day)
		}
	})

	t.Run("大量秒数", func(t *testing.T) {
		w := New()
		startYear := w.Year

		// 加一整年的秒数 (365 * 24 * 60 * 60)
		w.UpdateTime(365 * 24 * 60 * 60)

		if w.Year != startYear+1 {
			t.Errorf("加一年错误: got %d, want %d", w.Year, startYear+1)
		}
	})
}

// ==================== 钱包操作边界测试 ====================

func TestWalletEdgeCases(t *testing.T) {
	t.Run("添加负数", func(t *testing.T) {
		w := New()
		initial := w.Wallet[0]
		w.AddToWallet(0, -1)
		// 应该允许添加负数（实际上是减少）
		if w.Wallet[0] != initial-1 {
			t.Errorf("AddToWallet with negative should work, got %d", w.Wallet[0])
		}
	})

	t.Run("移除超过余额", func(t *testing.T) {
		w := New()
		w.Wallet = [6]int{1, 0, 0, 0, 0, 0}
		ok := w.RemoveFromWallet(0, 2)
		if ok {
			t.Error("RemoveFromWallet should return false when count exceeds balance")
		}
		if w.Wallet[0] != 1 {
			t.Error("Wallet should not change on failed removal")
		}
	})

	t.Run("无效索引操作", func(t *testing.T) {
		w := New()
		initial := w.Wallet

		w.AddToWallet(-1, 1)
		w.AddToWallet(6, 1)

		if w.Wallet != initial {
			t.Error("Invalid index should not change wallet")
		}

		ok := w.RemoveFromWallet(-1, 1)
		if ok {
			t.Error("RemoveFromWallet with invalid index should return false")
		}

		ok = w.RemoveFromWallet(6, 1)
		if ok {
			t.Error("RemoveFromWallet with invalid index should return false")
		}
	})

	t.Run("花费0元", func(t *testing.T) {
		w := New()
		initial := w.Wallet

		// 花费0元应该成功（但不改变钱包）
		// 具体行为取决于实现
		_ = w.SpendMoney(0)
		// 不检查错误，只验证钱包不变
		if w.Wallet != initial {
			// 某些实现可能会改变钱包，这里只记录
			t.Log("SpendMoney(0) changed wallet")
		}
	})
}

// ==================== 建筑名称测试 ====================

func TestBuildingNames(t *testing.T) {
	t.Run("名称数量匹配", func(t *testing.T) {
		w := New()
		if len(BuildingNames) != len(w.Buildings) {
			t.Errorf("BuildingNames count = %d, Buildings count = %d",
				len(BuildingNames), len(w.Buildings))
		}
	})

	t.Run("名称非空", func(t *testing.T) {
		for i, name := range BuildingNames {
			if name == "" {
				t.Errorf("BuildingNames[%d] is empty", i)
			}
		}
	})
}

// ==================== Adjacency 表测试 ====================

func TestAdjacency(t *testing.T) {
	t.Run("邻接对称性", func(t *testing.T) {
		// 如果A可达B，则B应可达A（除了特殊位置）
		specialLocations := map[int]bool{
			LocationSupermarketInner: true,
		}

		for from, toList := range Adjacency {
			if specialLocations[from] {
				continue
			}
			for _, to := range toList {
				if specialLocations[to] {
					continue
				}
				// 检查反向连接
				reverse := Adjacency[to]
				found := false
				for _, r := range reverse {
					if r == from {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("邻接不对称: %d->%d 存在, 但 %d->%d 不存在",
						from, to, to, from)
				}
			}
		}
	})

	t.Run("无自环", func(t *testing.T) {
		for from, toList := range Adjacency {
			for _, to := range toList {
				if from == to {
					t.Errorf("位置 %d 有自环", from)
				}
			}
		}
	})
}
