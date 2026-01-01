package save

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// testManager 创建使用临时目录的测试管理器
func testManager(t *testing.T) (*Manager, func()) {
	tmpDir, err := os.MkdirTemp("", "citylife-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	m := &Manager{saveDir: tmpDir}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return m, cleanup
}

func TestNewManager(t *testing.T) {
	m, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() error: %v", err)
	}
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	if m.saveDir == "" {
		t.Error("saveDir should not be empty")
	}
}

func TestGetSavePath(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	path := m.GetSavePath(1)
	expected := filepath.Join(m.saveDir, "save_1.json")
	if path != expected {
		t.Errorf("GetSavePath(1) = %q, want %q", path, expected)
	}
}

func TestSaveAndLoad(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	// 创建测试数据
	data := &SaveData{
		World: WorldData{
			Year:        2009,
			Month:       7,
			Day:         1,
			Hour:        12,
			Minute:      30,
			Second:      0,
			Where:       1,
			Wallet:      [6]int{1, 2, 3, 4, 5, 6},
			BankDeposit: 1000,
		},
		User: UserData{
			Nutrition: map[string]int{
				"饱腹感": 80,
				"饥渴":  70,
			},
		},
		Diseases: []DiseaseData{
			{ID: "cold", IsActive: true},
		},
	}

	// 保存
	err := m.Save(1, data)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// 加载
	loaded, err := m.Load(1)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// 验证数据
	t.Run("World数据", func(t *testing.T) {
		if loaded.World.Year != data.World.Year {
			t.Errorf("Year = %d, want %d", loaded.World.Year, data.World.Year)
		}
		if loaded.World.BankDeposit != data.World.BankDeposit {
			t.Errorf("BankDeposit = %d, want %d", loaded.World.BankDeposit, data.World.BankDeposit)
		}
		if loaded.World.Wallet != data.World.Wallet {
			t.Errorf("Wallet = %v, want %v", loaded.World.Wallet, data.World.Wallet)
		}
	})

	t.Run("User数据", func(t *testing.T) {
		if loaded.User.Nutrition["饱腹感"] != 80 {
			t.Errorf("Nutrition[饱腹感] = %d, want 80", loaded.User.Nutrition["饱腹感"])
		}
	})

	t.Run("Disease数据", func(t *testing.T) {
		if len(loaded.Diseases) != 1 {
			t.Fatalf("Diseases length = %d, want 1", len(loaded.Diseases))
		}
		if loaded.Diseases[0].ID != "cold" {
			t.Errorf("Disease ID = %q, want cold", loaded.Diseases[0].ID)
		}
	})

	t.Run("元数据", func(t *testing.T) {
		if loaded.Version != CurrentVersion {
			t.Errorf("Version = %d, want %d", loaded.Version, CurrentVersion)
		}
		if loaded.Slot != 1 {
			t.Errorf("Slot = %d, want 1", loaded.Slot)
		}
		if loaded.Timestamp.IsZero() {
			t.Error("Timestamp should not be zero")
		}
	})
}

func TestSlotValidation(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	tests := []struct {
		slot    int
		wantErr bool
	}{
		{0, true},
		{1, false},
		{2, false},
		{3, false},
		{4, true},
		{-1, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := m.Save(tt.slot, &SaveData{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Save(slot=%d) error = %v, wantErr = %v", tt.slot, err, tt.wantErr)
			}
		})
	}
}

func TestLoadNonExistent(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	_, err := m.Load(1)
	if err == nil {
		t.Error("Load() should return error for non-existent save")
	}
}

func TestDelete(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	// 创建存档
	m.Save(1, &SaveData{})

	// 删除
	err := m.Delete(1)
	if err != nil {
		t.Errorf("Delete() error: %v", err)
	}

	// 验证已删除
	if m.Exists(1) {
		t.Error("Save should not exist after delete")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	// 删除不存在的存档应该成功
	err := m.Delete(1)
	if err != nil {
		t.Errorf("Delete() should not error for non-existent save: %v", err)
	}
}

func TestExists(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	t.Run("不存在", func(t *testing.T) {
		if m.Exists(1) {
			t.Error("Exists() should return false for non-existent save")
		}
	})

	t.Run("存在", func(t *testing.T) {
		m.Save(1, &SaveData{})
		if !m.Exists(1) {
			t.Error("Exists() should return true for existing save")
		}
	})

	t.Run("无效槽位", func(t *testing.T) {
		if m.Exists(0) {
			t.Error("Exists() should return false for invalid slot")
		}
		if m.Exists(4) {
			t.Error("Exists() should return false for invalid slot")
		}
	})
}

func TestGetSlotInfo(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	t.Run("空槽位", func(t *testing.T) {
		info := m.GetSlotInfo(1)
		if info != "空" {
			t.Errorf("GetSlotInfo() = %q, want 空", info)
		}
	})

	t.Run("有存档", func(t *testing.T) {
		data := &SaveData{
			World: WorldData{
				Year:   2009,
				Month:  7,
				Day:    15,
				Hour:   10,
				Minute: 30,
			},
		}
		m.Save(1, data)

		info := m.GetSlotInfo(1)
		if info == "空" || info == "读取失败" {
			t.Errorf("GetSlotInfo() should return save info, got %q", info)
		}
	})

	t.Run("无效槽位", func(t *testing.T) {
		info := m.GetSlotInfo(0)
		if info != "无效槽位" {
			t.Errorf("GetSlotInfo(0) = %q, want 无效槽位", info)
		}
	})
}

func TestGetAvailableSlots(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	// 创建两个存档
	m.Save(1, &SaveData{})
	m.Save(3, &SaveData{})

	slots := m.GetAvailableSlots()

	if len(slots) != MaxSlots {
		t.Fatalf("GetAvailableSlots() returned %d slots, want %d", len(slots), MaxSlots)
	}

	// 验证槽位状态
	if !slots[0].Exists {
		t.Error("Slot 1 should exist")
	}
	if slots[1].Exists {
		t.Error("Slot 2 should not exist")
	}
	if !slots[2].Exists {
		t.Error("Slot 3 should exist")
	}
}

func TestSaveTimestamp(t *testing.T) {
	m, cleanup := testManager(t)
	defer cleanup()

	before := time.Now()
	m.Save(1, &SaveData{})
	after := time.Now()

	loaded, _ := m.Load(1)

	if loaded.Timestamp.Before(before) || loaded.Timestamp.After(after) {
		t.Errorf("Timestamp = %v, should be between %v and %v",
			loaded.Timestamp, before, after)
	}
}
