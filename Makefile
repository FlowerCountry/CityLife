.PHONY: build run clean test fmt

# 默认目标
all: build

# 构建
build:
	go build -o bin/citylife ./cmd/citylife

# 运行
run: build
	./bin/citylife

# 清理
clean:
	rm -rf bin/
	go clean

# 测试
test:
	go test -v ./...

# 格式化
fmt:
	go fmt ./...

# 检查
vet:
	go vet ./...

# 下载依赖
deps:
	go mod tidy

# 跨平台构建
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/citylife-linux-amd64 ./cmd/citylife
	GOOS=darwin GOARCH=amd64 go build -o bin/citylife-darwin-amd64 ./cmd/citylife
	GOOS=darwin GOARCH=arm64 go build -o bin/citylife-darwin-arm64 ./cmd/citylife
	GOOS=windows GOARCH=amd64 go build -o bin/citylife-windows-amd64.exe ./cmd/citylife
