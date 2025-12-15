// Package save 管理游戏存档
package save

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"citylife/internal/path"
)

// SaveData 存档数据结构
type SaveData struct {
	// 元数据
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Slot      int       `json:"slot"`

	// 世界状态
	World WorldData `json:"world"`

	// 用户状态
	User UserData `json:"user"`

	// 疾病状态
	Diseases []DiseaseData `json:"diseases"`
}

// WorldData 世界存档数据
type WorldData struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	Hour        int    `json:"hour"`
	Minute      int    `json:"minute"`
	Second      int    `json:"second"`
	Where       int    `json:"where"`
	Wallet      [6]int `json:"wallet"`
	BankDeposit int    `json:"bank_deposit"`
}

// UserData 用户存档数据
type UserData struct {
	Nutrition map[string]int `json:"nutrition"`
}

// DiseaseData 疾病存档数据
type DiseaseData struct {
	ID              string `json:"id"`
	TriggerProgress int    `json:"trigger_progress"`
	CureProgress    int    `json:"cure_progress"`
	IsActive        bool   `json:"is_active"`
}

const (
	CurrentVersion = 1
	MaxSlots       = 3
)

// Manager 存档管理器
type Manager struct {
	saveDir string
}

// NewManager 创建存档管理器
// 使用统一的路径管理器获取存档目录
func NewManager() (*Manager, error) {
	pm := path.GetInstance()

	// 如果PathManager已初始化，使用新路径
	if pm.IsInitialized() {
		saveDir := pm.GetSaveDirectory()

		// 创建存档目录
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return nil, fmt.Errorf("无法创建存档目录: %w", err)
		}

		mgr := &Manager{saveDir: saveDir}

		// 迁移旧数据（如果存在）
		mgr.migrateOldData()

		return mgr, nil
	}

	// 回退到旧路径（兼容模式）
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户目录: %w", err)
	}

	saveDir := filepath.Join(homeDir, ".citylife-go", "saves")

	// 创建存档目录
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建存档目录: %w", err)
	}

	return &Manager{saveDir: saveDir}, nil
}

// migrateOldData 迁移旧版存档数据到新位置
func (m *Manager) migrateOldData() {
	oldSaveDir := path.GetOldSaveDirectory()

	// 检查旧目录是否存在
	if !path.DirectoryExists(oldSaveDir) {
		return
	}

	// 迁移每个存档槽位
	for slot := 1; slot <= MaxSlots; slot++ {
		oldPath := filepath.Join(oldSaveDir, fmt.Sprintf("save_%d.json", slot))
		newPath := m.GetSavePath(slot)

		// 如果旧文件存在且新文件不存在，则迁移
		if path.FileExists(oldPath) && !path.FileExists(newPath) {
			if err := copyFile(oldPath, newPath); err == nil {
				fmt.Printf("已迁移存档槽位 %d 到新位置\n", slot)
			}
		}
	}
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// GetSavePath 获取存档路径
func (m *Manager) GetSavePath(slot int) string {
	return filepath.Join(m.saveDir, fmt.Sprintf("save_%d.json", slot))
}

// Save 保存游戏
func (m *Manager) Save(slot int, data *SaveData) error {
	if slot < 1 || slot > MaxSlots {
		return fmt.Errorf("无效的存档槽位: %d (有效范围: 1-%d)", slot, MaxSlots)
	}

	data.Version = CurrentVersion
	data.Timestamp = time.Now()
	data.Slot = slot

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	// 写入文件
	savePath := m.GetSavePath(slot)
	if err := os.WriteFile(savePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入存档失败: %w", err)
	}

	return nil
}

// Load 加载游戏
func (m *Manager) Load(slot int) (*SaveData, error) {
	if slot < 1 || slot > MaxSlots {
		return nil, fmt.Errorf("无效的存档槽位: %d (有效范围: 1-%d)", slot, MaxSlots)
	}

	savePath := m.GetSavePath(slot)

	// 读取文件
	jsonData, err := os.ReadFile(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("存档不存在")
		}
		return nil, fmt.Errorf("读取存档失败: %w", err)
	}

	// 反序列化
	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("解析存档失败: %w", err)
	}

	// 版本检查
	if data.Version > CurrentVersion {
		return nil, fmt.Errorf("存档版本过高: %d (当前支持: %d)", data.Version, CurrentVersion)
	}

	return &data, nil
}

// Delete 删除存档
func (m *Manager) Delete(slot int) error {
	if slot < 1 || slot > MaxSlots {
		return fmt.Errorf("无效的存档槽位: %d", slot)
	}

	savePath := m.GetSavePath(slot)
	if err := os.Remove(savePath); err != nil {
		if os.IsNotExist(err) {
			return nil // 不存在也算成功
		}
		return fmt.Errorf("删除存档失败: %w", err)
	}

	return nil
}

// Exists 检查存档是否存在
func (m *Manager) Exists(slot int) bool {
	if slot < 1 || slot > MaxSlots {
		return false
	}

	savePath := m.GetSavePath(slot)
	_, err := os.Stat(savePath)
	return err == nil
}

// GetSlotInfo 获取存档槽位信息
func (m *Manager) GetSlotInfo(slot int) string {
	if slot < 1 || slot > MaxSlots {
		return "无效槽位"
	}

	if !m.Exists(slot) {
		return "空"
	}

	data, err := m.Load(slot)
	if err != nil {
		return "读取失败"
	}

	return fmt.Sprintf("%d年%d月%d日 %02d:%02d (保存于: %s)",
		data.World.Year, data.World.Month, data.World.Day,
		data.World.Hour, data.World.Minute,
		data.Timestamp.Format("2006-01-02 15:04"))
}

// GetAvailableSlots 获取所有存档槽位状态
func (m *Manager) GetAvailableSlots() []SlotInfo {
	slots := make([]SlotInfo, MaxSlots)
	for i := 1; i <= MaxSlots; i++ {
		slots[i-1] = SlotInfo{
			Slot:        i,
			Exists:      m.Exists(i),
			Description: m.GetSlotInfo(i),
		}
	}
	return slots
}

// SlotInfo 存档槽位信息
type SlotInfo struct {
	Slot        int
	Exists      bool
	Description string
}
