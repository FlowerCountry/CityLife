// Package path 统一路径管理器 - 管理应用程序所有数据路径
//
// 目录结构:
//
//	bin/
//	├── main              # 可执行文件
//	└── data/             # 统一数据目录
//	    ├── locale/       # 语言资源
//	    └── saves/        # 存档文件
package path

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager 路径管理器（单例）
type Manager struct {
	basePath    string
	initialized bool
	mu          sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{}
	})
	return instance
}

// Initialize 初始化路径管理器
// executableDir 为可执行文件所在目录
func (m *Manager) Initialize(executableDir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.basePath = executableDir
	if m.basePath == "" {
		m.basePath = "."
	}
	m.initialized = true
}

// InitializeFromExecutable 从可执行文件路径自动初始化
// 会自动获取当前可执行文件所在目录
func (m *Manager) InitializeFromExecutable() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %w", err)
	}

	// 解析符号链接获取真实路径
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("无法解析符号链接: %w", err)
	}

	m.Initialize(filepath.Dir(exe))
	return nil
}

// IsInitialized 检查是否已初始化
func (m *Manager) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}

// GetBasePath 获取可执行文件所在目录
func (m *Manager) GetBasePath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.basePath
}

// GetDataDirectory 获取数据根目录
// 返回如 "/path/to/bin/data"
func (m *Manager) GetDataDirectory() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return filepath.Join(m.basePath, "data")
}

// GetSaveDirectory 获取存档目录
// 返回如 "/path/to/bin/data/saves"
func (m *Manager) GetSaveDirectory() string {
	return filepath.Join(m.GetDataDirectory(), "saves")
}

// GetLocaleDirectory 获取语言资源目录
// 返回如 "/path/to/bin/data/locale"
func (m *Manager) GetLocaleDirectory() string {
	return filepath.Join(m.GetDataDirectory(), "locale")
}

// GetSavePath 获取指定存档槽位的完整路径
// slot 为存档槽位 (1-3)
// 返回如 "/path/to/bin/data/saves/save_1.json"
func (m *Manager) GetSavePath(slot int) string {
	return filepath.Join(m.GetSaveDirectory(), fmt.Sprintf("save_%d.json", slot))
}

// GetLocalePath 获取指定语言文件的完整路径
// locale 为语言代码 (如 "zh_CN")
// 返回如 "/path/to/bin/data/locale/zh_CN.ini"
func (m *Manager) GetLocalePath(locale string) string {
	return filepath.Join(m.GetLocaleDirectory(), locale+".ini")
}

// EnsureDirectories 确保所有数据目录存在
// 返回创建过程中的错误
func (m *Manager) EnsureDirectories() error {
	if !m.IsInitialized() {
		return fmt.Errorf("路径管理器未初始化")
	}

	// 创建数据根目录
	if err := os.MkdirAll(m.GetDataDirectory(), 0755); err != nil {
		return fmt.Errorf("无法创建数据目录 %s: %w", m.GetDataDirectory(), err)
	}

	// 创建存档目录
	if err := os.MkdirAll(m.GetSaveDirectory(), 0755); err != nil {
		return fmt.Errorf("无法创建存档目录 %s: %w", m.GetSaveDirectory(), err)
	}

	// 创建语言资源目录
	if err := os.MkdirAll(m.GetLocaleDirectory(), 0755); err != nil {
		return fmt.Errorf("无法创建语言目录 %s: %w", m.GetLocaleDirectory(), err)
	}

	return nil
}

// DirectoryExists 检查目录是否存在
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// GetHomeDirectory 获取用户主目录
func GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// GetOldSaveDirectory 获取旧版存档目录 (~/.citylife-go/saves)
// 用于兼容性检查和数据迁移
func GetOldSaveDirectory() string {
	return filepath.Join(GetHomeDirectory(), ".citylife-go", "saves")
}
