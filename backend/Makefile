.PHONY: all build run clean test tidy lint

# ========== DevOpsStar 后端 Makefile ==========
# 国内 Go 代理加速（七牛云 goproxy.cn）
export GOPROXY := https://goproxy.cn,direct
export GOSUMDB := off

APP_NAME := devops-star
BUILD_DIR := bin
MAIN_FILE := main.go

all: build

# 下载依赖
deps:
	go mod download
	go mod tidy

# 构建二进制文件
build: deps
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(APP_NAME)"

# 运行服务
run: build
	@echo "🚀 启动 DevOpsStar 后端服务..."
	@$(BUILD_DIR)/$(APP_NAME)

# 热重载运行（需要 gin-hotreload 或 air）
dev:
	@which air > /dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

# 清理构建文件
clean:
	@rm -rf $(BUILD_DIR)
	@echo "🧹 清理完成"

# 运行测试
test:
	go test -v ./...

# 代码检查
lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

# 格式化代码
fmt:
	go fmt ./...

# 检查编译（不生成二进制）
check:
	go vet ./...
	go build -o /dev/null $(MAIN_FILE)

# 生成 Swagger 文档（需要 swag）
swagger:
	@which swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(MAIN_FILE) -o docs

# Docker 构建
docker-build:
	docker build -t devops-star:latest .

# Docker 运行
docker-run:
	docker run -p 8080:8080 --env-file .env devops-star:latest

# 数据库迁移
migrate:
	@echo "🗄️  运行数据库迁移..."
	go run scripts/migrate.go

help:
	@echo "DevOpsStar 后端构建命令"
	@echo ""
	@echo "  make build      - 构建二进制文件"
	@echo "  make run        - 构建并运行"
	@echo "  make dev        - 热重载开发模式"
	@echo "  make test       - 运行测试"
	@echo "  make lint       - 代码检查"
	@echo "  make clean      - 清理构建文件"
	@echo "  make swagger    - 生成 API 文档"
	@echo "  make docker-build - Docker 构建"
	@echo "  make help        - 显示帮助"
