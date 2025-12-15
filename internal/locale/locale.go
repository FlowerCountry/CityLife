// Package locale 提供多语言支持
package locale

import (
	"math/rand"
	"os"
	"path/filepath"
)

// Config 语言配置
type Config struct {
	Locale   string // 语言代码，如 "zh_CN"
	BasePath string // 资源文件基础路径
}

// Locale 语言管理器
type Locale struct {
	sections      map[string]map[string]string // section -> key -> value
	currentLocale string
	basePath      string
	initialized   bool
}

// New 创建语言管理器
func New(cfg Config) (*Locale, error) {
	l := &Locale{
		sections:      make(map[string]map[string]string),
		currentLocale: cfg.Locale,
		basePath:      cfg.BasePath,
		initialized:   false,
	}

	if l.currentLocale == "" {
		l.currentLocale = "zh_CN"
	}

	if l.basePath == "" {
		// 默认使用 ~/.citylife-go/locale/
		home, _ := os.UserHomeDir()
		l.basePath = filepath.Join(home, ".citylife-go", "locale")
	}

	// 尝试加载语言文件
	err := l.load()
	if err != nil {
		// 加载失败时使用默认值
		l.loadDefaults()
	}

	l.initialized = true
	return l, nil
}

// NewWithDefaults 创建使用默认值的语言管理器
func NewWithDefaults() *Locale {
	l := &Locale{
		sections:      make(map[string]map[string]string),
		currentLocale: "zh_CN",
		basePath:      "",
		initialized:   true,
	}
	l.loadDefaults()
	return l
}

// load 加载语言文件
func (l *Locale) load() error {
	filepath := l.getLocalePath(l.currentLocale)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return l.parseINI(string(content))
}

// getLocalePath 获取语言文件路径
func (l *Locale) getLocalePath(locale string) string {
	return filepath.Join(l.basePath, locale+".ini")
}

// Get 获取指定section和key的字符串
// 找不到时返回 "[section.key]" 作为 fallback
func (l *Locale) Get(section, key string) string {
	if sec, ok := l.sections[section]; ok {
		if val, ok := sec[key]; ok {
			return val
		}
	}
	return "[" + section + "." + key + "]"
}

// GetSection 获取整个section的所有key-value对
func (l *Locale) GetSection(section string) map[string]string {
	if sec, ok := l.sections[section]; ok {
		// 返回副本
		result := make(map[string]string)
		for k, v := range sec {
			result[k] = v
		}
		return result
	}
	return make(map[string]string)
}

// GetSectionValues 获取section中所有值的列表（用于随机选择）
func (l *Locale) GetSectionValues(section string) []string {
	if sec, ok := l.sections[section]; ok {
		values := make([]string, 0, len(sec))
		for _, v := range sec {
			values = append(values, v)
		}
		return values
	}
	return nil
}

// GetRandom 从section中随机获取一个值
func (l *Locale) GetRandom(section string) string {
	values := l.GetSectionValues(section)
	if len(values) == 0 {
		return "[" + section + "]"
	}
	return values[rand.Intn(len(values))]
}

// HasSection 检查section是否存在
func (l *Locale) HasSection(section string) bool {
	_, ok := l.sections[section]
	return ok
}

// HasKey 检查key是否存在
func (l *Locale) HasKey(section, key string) bool {
	if sec, ok := l.sections[section]; ok {
		_, ok := sec[key]
		return ok
	}
	return false
}

// CurrentLocale 获取当前语言代码
func (l *Locale) CurrentLocale() string {
	return l.currentLocale
}

// IsInitialized 检查是否已初始化
func (l *Locale) IsInitialized() bool {
	return l.initialized
}

// SetLocale 切换语言（重新加载）
func (l *Locale) SetLocale(locale string) error {
	l.currentLocale = locale
	l.sections = make(map[string]map[string]string)
	err := l.load()
	if err != nil {
		l.loadDefaults()
	}
	return nil
}

// set 内部方法：设置一个键值对
func (l *Locale) set(section, key, value string) {
	if _, ok := l.sections[section]; !ok {
		l.sections[section] = make(map[string]string)
	}
	l.sections[section][key] = value
}
