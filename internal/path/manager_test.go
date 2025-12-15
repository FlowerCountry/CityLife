package path

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPathManager_Initialize(t *testing.T) {
	pm := GetInstance()

	// 测试未初始化状态
	// 注意：由于单例模式，如果之前已初始化，这个测试可能不准确
	// 在实际测试中应该使用新实例

	// 初始化
	testDir := "/tmp/citylife-test"
	pm.Initialize(testDir)

	if !pm.IsInitialized() {
		t.Error("PathManager 应该已初始化")
	}

	if pm.GetBasePath() != testDir {
		t.Errorf("基础路径不匹配: got %s, want %s", pm.GetBasePath(), testDir)
	}
}

func TestPathManager_GetDirectories(t *testing.T) {
	pm := GetInstance()
	testDir := "/tmp/citylife-test"
	pm.Initialize(testDir)

	// 测试数据目录
	dataDir := pm.GetDataDirectory()
	if !strings.HasSuffix(dataDir, "data") {
		t.Errorf("数据目录应该以 'data' 结尾: %s", dataDir)
	}

	// 测试存档目录
	saveDir := pm.GetSaveDirectory()
	if !strings.HasSuffix(saveDir, filepath.Join("data", "saves")) {
		t.Errorf("存档目录应该以 'data/saves' 结尾: %s", saveDir)
	}

	// 测试语言目录
	localeDir := pm.GetLocaleDirectory()
	if !strings.HasSuffix(localeDir, filepath.Join("data", "locale")) {
		t.Errorf("语言目录应该以 'data/locale' 结尾: %s", localeDir)
	}
}

func TestPathManager_GetSavePath(t *testing.T) {
	pm := GetInstance()
	testDir := "/tmp/citylife-test"
	pm.Initialize(testDir)

	// 测试存档路径
	savePath := pm.GetSavePath(1)
	expectedSuffix := filepath.Join("data", "saves", "save_1.json")
	if !strings.HasSuffix(savePath, expectedSuffix) {
		t.Errorf("存档路径不正确: got %s, want suffix %s", savePath, expectedSuffix)
	}

	// 测试多个槽位
	for slot := 1; slot <= 3; slot++ {
		savePath := pm.GetSavePath(slot)
		if !strings.Contains(savePath, "save_") {
			t.Errorf("槽位 %d 的存档路径应该包含 'save_': %s", slot, savePath)
		}
	}
}

func TestPathManager_GetLocalePath(t *testing.T) {
	pm := GetInstance()
	testDir := "/tmp/citylife-test"
	pm.Initialize(testDir)

	localePath := pm.GetLocalePath("zh_CN")
	expectedSuffix := filepath.Join("data", "locale", "zh_CN.ini")
	if !strings.HasSuffix(localePath, expectedSuffix) {
		t.Errorf("语言文件路径不正确: got %s, want suffix %s", localePath, expectedSuffix)
	}
}

func TestPathManager_EnsureDirectories(t *testing.T) {
	pm := GetInstance()
	testDir := "/tmp/citylife-pathmanager-test"
	pm.Initialize(testDir)

	// 清理测试目录
	os.RemoveAll(testDir)

	// 确保目录创建
	err := pm.EnsureDirectories()
	if err != nil {
		t.Errorf("EnsureDirectories 失败: %v", err)
	}

	// 验证目录存在
	if !DirectoryExists(pm.GetDataDirectory()) {
		t.Error("数据目录应该已创建")
	}
	if !DirectoryExists(pm.GetSaveDirectory()) {
		t.Error("存档目录应该已创建")
	}
	if !DirectoryExists(pm.GetLocaleDirectory()) {
		t.Error("语言目录应该已创建")
	}

	// 清理
	os.RemoveAll(testDir)
}

func TestHelperFunctions(t *testing.T) {
	// 测试 DirectoryExists
	if !DirectoryExists("/tmp") {
		t.Error("/tmp 应该存在")
	}
	if DirectoryExists("/nonexistent-dir-12345") {
		t.Error("不存在的目录不应该返回 true")
	}

	// 测试 GetHomeDirectory
	home := GetHomeDirectory()
	if home == "" || home == "." {
		t.Log("警告: 无法获取用户主目录")
	}

	// 测试 GetOldSaveDirectory
	oldSaveDir := GetOldSaveDirectory()
	if !strings.Contains(oldSaveDir, ".citylife-go") {
		t.Errorf("旧存档目录应该包含 '.citylife-go': %s", oldSaveDir)
	}
}
