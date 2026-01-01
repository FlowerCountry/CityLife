// Package testutil 提供自动化测试支持
package testutil

import (
	"flag"
	"os"
)

// Config 测试配置
type Config struct {
	InputPath  string
	OutputPath string
	Enabled    bool
}

// ParseArgs 解析命令行参数
// 支持 --test input output 格式
func ParseArgs() Config {
	cfg := Config{}

	// 检查是否有 --test 参数
	for i, arg := range os.Args {
		if arg == "--test" {
			cfg.Enabled = true
			if i+1 < len(os.Args) {
				cfg.InputPath = os.Args[i+1]
			}
			if i+2 < len(os.Args) {
				cfg.OutputPath = os.Args[i+2]
			}
			break
		}
	}

	// 如果没有使用--test，使用flag包解析
	if !cfg.Enabled {
		inputPtr := flag.String("input", "", "测试输入文件路径")
		outputPtr := flag.String("output", "", "测试输出文件路径")
		testPtr := flag.Bool("test-mode", false, "启用测试模式")
		flag.Parse()

		if *testPtr || (*inputPtr != "" && *outputPtr != "") {
			cfg.Enabled = true
			cfg.InputPath = *inputPtr
			cfg.OutputPath = *outputPtr
		}
	}

	return cfg
}

// IsTestMode 检查是否在测试模式下运行
func IsTestMode() bool {
	for _, arg := range os.Args {
		if arg == "--test" || arg == "-test-mode" {
			return true
		}
	}
	return false
}
