# ==================== FightGame 项目 Makefile ====================

.PHONY: help dev install build clean docker-up docker-down docker-build

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==================== 开发 ====================

install: ## 安装所有依赖
	@echo ">>> 安装前端依赖..."
	cd frontend && npm install --registry=https://registry.npmmirror.com
	@echo ">>> 安装管理后台依赖..."
	cd admin && npm install --registry=https://registry.npmmirror.com
	@echo ">>> 安装后端依赖..."
	cd server && go mod download
	@echo ">>> 所有依赖安装完成!"

dev: ## 启动开发环境
	@echo ">>> 启动开发环境..."
	@echo ">>> 后端: http://localhost:8080"
	@echo ">>> 前端: http://localhost:3000"
	@echo ">>> 管理后台: http://localhost:3001"

dev-frontend: ## 启动前端开发服务器
	cd frontend && npm run dev

dev-admin: ## 启动管理后台开发服务器
	cd admin && npm run dev

dev-server: ## 启动后端开发服务器
	cd server && go run ./cmd/main.go

# ==================== 构建 ====================

build-frontend: ## 构建前端
	cd frontend && npm run build

build-admin: ## 构建管理后台
	cd admin && npm run build

build-server: ## 构建后端
	cd server && go build -o bin/server ./cmd/main.go

build: build-frontend build-admin build-server ## 构建所有

# ==================== Docker ====================

docker-up: ## 启动 Docker 环境
	docker-compose up -d

docker-down: ## 停止 Docker 环境
	docker-compose down

docker-build: ## 构建 Docker 镜像
	docker-compose build

docker-logs: ## 查看 Docker 日志
	docker-compose logs -f

docker-restart: docker-down docker-build docker-up ## 重启 Docker 环境

# ==================== 工具 ====================

ssl-generate: ## 生成开发用自签名 SSL 证书
	@echo ">>> 生成自签名 SSL 证书..."
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout docker/nginx/ssl/server.key \
		-out docker/nginx/ssl/server.crt \
		-subj "/C=CN/ST=Beijing/L=Beijing/O=FightGame/CN=localhost"
	@echo ">>> SSL 证书生成完成!"

clean: ## 清理构建产物
	rm -rf frontend/dist admin/dist server/bin server/logs
	@echo ">>> 清理完成!"

lint: ## 代码检查
	cd server && go vet ./...
	cd frontend && npx vue-tsc --noEmit
	cd admin && npx vue-tsc --noEmit

.DEFAULT_GOAL := help
