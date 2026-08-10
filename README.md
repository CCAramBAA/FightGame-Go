# ==================== FightGame 项目启动指南 ====================

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + Phaser 3 + TypeScript + Axios + Howler |
| 后端 | Golang + Gin + gorilla/websocket + GORM |
| 存储 | MySQL 8 + Redis |
| 运维 | Docker + Nginx + SSL |
| 额外 | 日志系统 + 管理后台 (Vue 3 + Go) |

## 项目结构

```
FightGame/
├── frontend/          # 游戏前端 (Vue3 + Phaser3)
│   ├── src/           # 源码
│   │   ├── game/      # Phaser3 游戏场景
│   │   ├── views/     # 页面视图
│   │   ├── api/       # API 请求
│   │   ├── store/     # Pinia 状态管理
│   │   ├── router/    # Vue Router
│   │   └── utils/     # 工具 (WebSocket/Audio)
│   └── vite.config.ts
├── admin/             # 管理后台 (Vue3 + Element Plus)
│   └── src/views/     # dashboard/user/game/log
├── server/            # Go 后端
│   ├── cmd/main.go    # 入口
│   └── internal/
│       ├── config/    # 配置管理
│       ├── model/     # 数据模型
│       ├── handler/   # API 处理器
│       ├── middleware/ # 中间件
│       ├── websocket/ # WebSocket
│       └── logger/    # 日志
├── docker/            # Docker 配置
│   ├── nginx/         # Nginx 配置 + SSL
│   ├── mysql/         # 数据库初始化脚本
│   ├── Dockerfile.server
│   ├── Dockerfile.frontend
│   └── Dockerfile.admin
├── docker-compose.yml
├── Makefile
└── .env
```

## 快速开始

### 前置条件
- Node.js >= 18
- Go >= 1.21
- Docker & Docker Compose (可选)
- MySQL 8 & Redis (不使用 Docker 时需要)

### 1. SSL 证书（开发环境）
```bash
# Windows
scripts\generate-ssl.bat

# Linux/Mac
bash scripts/generate-ssl.sh
```

### 2. 安装依赖
```bash
# 使用 Make
make install

# 或手动安装
cd frontend && npm install
cd ../admin && npm install
cd ../server && go mod download
```

### 3. 启动开发环境

**方式一：Docker Compose（推荐）**
```bash
docker-compose up -d
```
访问：
- 前端游戏: https://localhost
- 管理后台: https://localhost/admin
- API: https://localhost/api/health

**方式二：本地开发**
```bash
# 终端 1：启动后端
cd server && go run ./cmd/main.go

# 终端 2：启动前端
cd frontend && npm run dev

# 终端 3：启动管理后台
cd admin && npm run dev
```

访问：
- 前端: http://localhost:3000
- 管理后台: http://localhost:3001
- API: http://localhost:8080/api/health

### 4. 构建生产版本
```bash
make build
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/health | 健康检查 |
| POST | /api/login | 用户登录 |
| POST | /api/register | 用户注册 |
| GET | /api/user/info | 获取用户信息 |
| POST | /api/user/logout | 用户登出 |
| GET | /ws | WebSocket 连接 |

## 环境变量

参考 `.env` 文件，主要配置项：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_PORT | 后端端口 | 8080 |
| DB_HOST | 数据库地址 | 127.0.0.1 |
| DB_PORT | 数据库端口 | 3306 |
| DB_PASSWORD | 数据库密码 | root123 |
| REDIS_HOST | Redis 地址 | 127.0.0.1 |
| JWT_SECRET | JWT 密钥 | (需修改) |
