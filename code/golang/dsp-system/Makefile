.PHONY: run build test clean install help

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME=dsp-system
BUILD_DIR=./build
MAIN_FILE=main.go

# 颜色输出
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

## help: 显示帮助信息
help:
	@echo ''
	@echo '使用方法:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo '可用目标:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "  ${YELLOW}%-15s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## install: 安装依赖
install:
	@echo "安装依赖..."
	go mod download
	go mod tidy
	@echo "依赖安装完成!"

## run: 运行应用
run:
	@echo "启动 DSP 服务..."
	go run $(MAIN_FILE)

## build: 编译应用
build:
	@echo "编译应用..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

## build-linux: 编译 Linux 版本
build-linux:
	@echo "编译 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)-linux"

## test: 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告: coverage.html"

## bench: 运行基准测试
bench:
	@echo "运行基准测试..."
	go test -bench=. -benchmem ./...

## lint: 代码检查
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint 未安装，请运行: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

## clean: 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成!"

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t $(APP_NAME):latest .

## docker-run: 运行 Docker 容器
docker-run:
	@echo "运行 Docker 容器..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

## dev: 开发模式（热重载）
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air 未安装，请运行: go install github.com/cosmtrek/air@latest"; \
		echo "或直接运行: make run"; \
	fi

## test-request: 发送测试请求
test-request:
	@echo "发送测试竞价请求..."
	@curl -X POST http://localhost:8080/bid \
		-H "Content-Type: application/json" \
		-d @test_request.json

## health: 健康检查
health:
	@curl -s http://localhost:8080/health | jq .

## stats: 查看统计信息
stats:
	@curl -s http://localhost:8080/stats | jq .

