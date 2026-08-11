# FightGame - 在线格斗游戏平台

## 项目概述

FightGame 是一个基于 Web 的实时对战格斗游戏平台，支持 PVP 房间对战、PVE AI 闯关、角色选择、技能系统和帧同步战斗。项目采用前后端分离架构，包含三个独立模块：游戏前端、管理后台和 Go 后端服务。

## 技术栈

| 层级     | 技术选型                                               |
| -------- | ------------------------------------------------------ |
| 游戏前端 | Vue 3 + TypeScript + Vite + Phaser 3 + Pinia + Howler  |
| 管理后台 | Vue 3 + TypeScript + Vite + Element Plus + Pinia       |
| 后端     | Go 1.21 + Gin + gorilla/websocket + GORM + golang-jwt  |
| 存储     | MySQL 8.0 + Redis 7                                    |
| 运维     | Docker + Docker Compose + Nginx + SSL                  |

## 项目架构

```
                    ┌──────────────────────────────────┐
                    │          Nginx (80/443)           │
                    │       反向代理 + SSL 终端          │
                    └──────────────┬───────────────────┘
           ┌──────────────────────┼──────────────────────┐
           ▼                      ▼                      ▼
   ┌───────────────┐     ┌───────────────┐     ┌───────────────┐
   │   前端 :3000   │     │  管理后台 :3001│     │  后端 :8080    │
   │  Vue3+Phaser3 │     │ Vue3+Element+ │     │   Go + Gin    │
   │  游戏客户端    │     │  运营管理平台   │     │  REST + WS    │
   └───────────────┘     └───────────────┘     └───────┬───────┘
                                                       │
                                               ┌───────┴───────┐
                                               ▼               ▼
                                         ┌─────────┐     ┌─────────┐
                                         │ MySQL 8 │     │ Redis 7 │
                                         │  :3306  │     │  :6379  │
                                         └─────────┘     └─────────┘
```

### 通信流程

- **HTTP/REST**: 前端 <-> 后端，用于登录注册、用户信息、房间列表等常规请求
- **WebSocket**: 前端 <-> 后端，用于实时对战通信（创建房间、帧同步、战斗结算）
- **JWT 认证**: 登录后获取 token，HTTP 通过 `Authorization: Bearer <token>` 头传递，WebSocket 通过 `?token=` 查询参数传递

---

## 当前开发进度

### 阶段一：项目基础架构 ✅

- [x] Go 后端脚手架（Gin 框架、GORM ORM、配置管理、日志系统、CORS 中间件）
- [x] Vue3 前端脚手架（Vite、Phaser3 游戏引擎、路由系统、Axios 封装）
- [x] Vue3 管理后台脚手架（Element Plus UI 库、路由守卫）
- [x] MySQL 15 张数据表设计（用户、角色、技能、皮肤、PVE/PVP、好友、商城）
- [x] Redis 集成
- [x] Docker Compose 一键部署（MySQL + Redis + Server + Nginx）
- [x] JWT 认证中间件（通用认证 + 管理员权限校验）

### 阶段二：用户系统 ✅

- [x] 用户注册/登录 API（bcrypt 密码加密、JWT 令牌签发）
- [x] 用户信息查询 API（资产、战绩、角色皮肤列表）
- [x] 用户资料修改 API（昵称、头像、密码）
- [x] 角色查询 API（角色列表、详情、技能信息）
- [x] 管理员 API（Dashboard 统计、用户管理、CRUD）
- [x] 前端登录/注册页面（合并标签切换，纯 CSS 深色主题）
- [x] 路由守卫（未登录拦截、已登录重定向）
- [x] 管理后台对接真实 API（Dashboard 实时统计、用户分页列表）

### 阶段三：战斗系统（核心） ✅

- [x] **WebSocket 房间管理**：创建/加入/离开房间、准备状态、角色选择
- [x] **房间状态机**：`waiting → selecting → ready → playing → finished`
- [x] **房间列表实时广播**：创建/加入/离开/断线自动广播 `room_list_update`
- [x] **帧同步战斗**：客户端帧输入广播、服务端帧计数
- [x] **战斗引擎**（Phaser3 GameScene）：
  - HP 条 & 能量条 UI
  - WASD 四方向移动
  - J 普攻 / K 技能1（消耗 30 能量）/ L 技能2（消耗 50 能量）
  - S 防御键（减伤 70%）
  - 攻击距离判定（100px 范围）
  - 受击闪烁反馈
  - 能量自动回复（每秒 +2）
- [x] **战斗结算**：胜者判定、金币/分数结算
- [x] **PVE 战斗系统**：AI 引擎、PVE 关卡战斗、关卡进度保存
- [x] **战斗 HUD 层**：血条/能量/倒计时/技能栏悬停 UI，退出按钮

### 阶段四：前端页面体系（全部 14 页） ✅

#### 导航体系

| 页面 | 路由 | 说明 |
|------|------|------|
| **启动检测页** | `/` | 全屏海报 + LOGO + 连接检测 + 网络异常弹窗 |
| **登录注册页** | `/login` | 登录/注册标签切换，白色半透明卡片 |
| **新手教程** | `/tutorial` | 全屏遮罩 + 动画播放 + 分页/跳过 |
| **游戏大厅** | `/lobby` | 全局导航栏 + 5 大圆形入口（PVP/PVE/商城/教程/图鉴） |
| **PVP 房间列表** | `/pvp-rooms` | 房间卡片网格 + 创建/加入房间 |
| **PVP 英雄选人** | `/pvp-select` | 全英雄列表 + 立绘/皮肤预览 + 10 秒倒计时 |
| **PVE 关卡地图** | `/pve-stages` | 关卡点位（简单/普通/困难/BOSS） + 训练模式 |
| **PVE 英雄选人** | `/pve-select` | 已解锁英雄可选，未解锁灰色锁定 |
| **战斗场景** | `/game` | Phaser 画布 + HUD（血条/能量/计时/技能栏） + 结算弹窗 |
| **对局回放** | `/replay` | 回放 ID 输入 + 播放/暂停/倍速/帧控制 |
| **商城** | `/shop` | 英雄/皮肤双标签 + 购买/已拥有状态 |
| **英雄图鉴** | `/herodex` | 网格列表 + 右侧详情面板（立绘/技能/背景故事） |

#### 弹窗组件

| 组件 | 说明 |
|------|------|
| **设置** | 按键自定义 / 画质切换 / 音量调节 / 手柄检测 / 快捷跳转图鉴 |
| **好友列表** | 搜索添加 / 在线离线区分 / 对战邀请 / 删除好友 |
| **对局结算** | 胜负文字 + HP 对比 + PVP 积分 / PVE 金币星级 |
| **网络异常** | 连接失败全屏弹窗，退出程序 / 重新连接 |

#### 通用组件

| 组件 | 说明 |
|------|------|
| `ArtPlaceholder` | 缺失美术资源占位（灰色虚线方框 + 中文标签） |
| `GlobalNav` | 全局顶部导航栏（LOGO / 金币 / 段位 / 设置 / 好友） |
| `ModalOverlay` | 全屏半透明弹窗容器 |
| `LoadingSpinner` | 加载转圈动画 |
| `NetworkErrorDialog` | 网络异常弹窗 |

### 阶段五：待开发

- [ ] 好友系统完善（好友申请通知、在线状态推送）
- [ ] 商城后端对接测试
- [ ] 回放数据存储与回放引擎
- [ ] 设置项本地持久化（按键映射、画质、音量）
- [ ] 教程动画分段资源
- [ ] 英雄图鉴详情数据
- [ ] PVE 关卡配置数据完善
- [ ] 美术资源替换（背景/角色立绘/技能图标/UI 图标）

---

## 数据模型（15 张表）

| 表名              | 说明                   | 核心字段                                |
| ----------------- | ---------------------- | --------------------------------------- |
| `users`           | 用户账号               | username, password, gold, rank_score    |
| `gold_transactions` | 金币流水             | source_type(pvp_win/shop_buy 等), amount |
| `characters`      | 角色配置               | unlock_type, hp, energy, speed, atk, def |
| `skills`          | 技能表                 | skill_type(主动/被动), cooldown, damage  |
| `hero_special_rules` | 英雄特殊交互规则     | steal_skill, pierce_shield, immune_damage |
| `skins`           | 皮肤                   | 关联 character_id                       |
| `user_characters` | 用户已解锁角色         | (user_id, character_id) 联合唯一索引    |
| `user_skins`      | 用户已拥有皮肤         | (user_id, skin_id) 联合唯一索引         |
| `pve_stages`      | PVE 关卡配置           | difficulty, gold_reward, unlock_character |
| `pve_progress`    | PVE 玩家存档           | (user_id, stage_id) 联合唯一索引        |
| `game_rooms`      | PVP 对战房间           | host/guest, status, host/guest_char_id  |
| `battle_records`  | 对战记录               | frame_data(JSON), winner_id             |
| `friend_relations` | 好友关系              | status(pending/accepted/blocked)        |
| `shop_items`      | 商城商品               | item_type(character/skin), price         |

## API 接口

### 公开接口（无需认证）

| 方法   | 路径                  | 说明         |
| ------ | --------------------- | ------------ |
| POST   | `/api/login`          | 用户登录     |
| POST   | `/api/register`       | 用户注册     |
| POST   | `/api/admin/login`    | 管理员登录   |
| GET    | `/api/characters`     | 角色列表     |
| GET    | `/api/characters/:id` | 角色详情     |
| GET    | `/api/health`         | 健康检查     |
| GET    | `/api/online`         | 在线人数统计 |

### 用户接口（需要 JWT 认证）

| 方法   | 路径                    | 说明           |
| ------ | ----------------------- | -------------- |
| GET    | `/api/user/info`        | 获取用户信息   |
| PUT    | `/api/user/profile`     | 修改资料       |
| PUT    | `/api/user/password`    | 修改密码       |
| GET    | `/api/profile/characters` | 我的角色列表 |
| GET    | `/api/skins/my`         | 我的皮肤列表   |
| GET    | `/api/rooms`            | 房间列表       |
| POST   | `/api/rooms/create`     | 创建房间       |
| POST   | `/api/rooms/join/:id`   | 加入房间       |
| GET    | `/api/shop/items`       | 商城商品列表   |
| POST   | `/api/shop/purchase`    | 购买商品       |
| GET    | `/api/pve/stages`       | PVE 关卡列表   |
| GET    | `/api/pve/stages/:id`   | 关卡详情       |
| POST   | `/api/pve/progress`     | 保存关卡进度   |
| GET    | `/api/friends`          | 好友列表       |
| POST   | `/api/friends/add`      | 添加好友       |
| DELETE | `/api/friends/:id`      | 删除好友       |
| GET    | `/api/battle/replay/:id`| 获取回放数据   |
| WS     | `/ws?token=<jwt>`       | WebSocket 连接 |

### 管理员接口（需要 JWT + AdminOnly 中间件）

| 方法   | 路径                    | 说明           |
| ------ | ----------------------- | -------------- |
| GET    | `/api/admin/dashboard`  | 数据概览       |
| GET    | `/api/admin/users`      | 用户列表       |
| GET    | `/api/admin/users/:id`  | 用户详情       |

## WebSocket 协议

### 连接

```
ws://localhost:8080/ws?token=<JWT_TOKEN>
```

### 客户端 → 服务端（发送消息类型）

| type              | 载荷                            | 说明             |
| ----------------- | ------------------------------- | ---------------- |
| `create_room`     | `{ room_id?, character_id? }`   | 创建房间         |
| `join_room`       | `{ room_id }`                   | 加入房间         |
| `leave_room`      | —                               | 离开房间         |
| `get_room_list`   | —                               | 请求房间列表     |
| `set_ready`       | `{ ready: bool }`               | 设置准备状态     |
| `select_character`| `{ character_id }`              | 选择角色         |
| `start_game`      | —                               | 开始游戏（房主） |
| `frame_input`     | `{ action, x, y, ... }`         | 帧同步输入       |
| `battle_over`     | `{ winner_id, result }`         | 战斗结算         |
| `invite_response` | `{ accepted, from_uid }`        | 好友邀请回复     |

### 服务端 → 客户端（推送消息类型）

| type              | 说明               |
| ----------------- | ------------------ |
| `room_created`    | 房间创建成功       |
| `room_joined`     | 加入房间成功       |
| `player_joined`   | 有玩家加入         |
| `left_room`       | 已离开房间         |
| `room_update`     | 房间状态更新       |
| `room_closed`     | 房间关闭           |
| `room_list`       | 房间列表（请求响应）|
| `room_list_update`| 房间列表广播更新   |
| `game_countdown`  | 游戏倒计时 (3-2-1) |
| `game_start`      | 游戏开始           |
| `frame_input`     | 对手帧输入         |
| `battle_result`   | 对战结果           |
| `error`           | 错误消息           |

## 快速开始

### 前置条件

| 工具         | 最低版本 | 是否必须 | 说明                                               |
| ------------ | -------- | -------- | -------------------------------------------------- |
| **MySQL**    | 8.0      | **必须** | 本地运行或 Docker 启动，创建数据库 `fight_game`    |
| **Redis**    | 7.x      | 推荐     | 用于缓存和在线状态，没有也能运行                     |
| **Node.js**  | 18.x     | **必须** | 前端/管理后台运行时                                                     |
| **Go**       | 1.21     | **必须** | 后端运行时                                                   |

#### 安装 MySQL（如果没有）

**Windows 推荐方案：**

- **[MySQL Community Server](https://dev.mysql.com/downloads/installer/)** — 下载 MSI 安装包，安装时设置 root 密码为 `root123`
- **[XAMPP](https://www.apachefriends.org/)** — 一键安装含 MySQL，修改 root 密码为 `root123`

**用 Docker 启动（最简单）：**

```bash
# 先安装 Docker Desktop: https://www.docker.com/products/docker-desktop/
# 然后只启动 MySQL + Redis：
docker-compose up -d mysql redis
```

安装 MySQL 后，创建数据库：
```sql
CREATE DATABASE IF NOT EXISTS fight_game
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
```

### 方式一：Docker Compose（一键部署，推荐）

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问地址：
- 游戏前端：https://localhost
- 管理后台：https://localhost/admin
- API：https://localhost/api/health

### 方式二：本地开发

```bash
# 1. 安装依赖
cd frontend && npm install
cd ../admin && npm install
cd ../server && go mod download

# 2. 启动 MySQL 和 Redis（需要本地已安装并运行）

# 3. 启动后端（终端 1）
cd server && go run ./cmd/main.go

# 4. 启动前端（终端 2）
cd frontend && npm run dev

# 5. 启动管理后台（终端 3）
cd admin && npm run dev
```

**Windows 用户可使用一键启动脚本：**

```batch
start-dev.bat
```

访问地址：
- 游戏前端：http://localhost:3000
- 管理后台：http://localhost:3001
- API 服务：http://localhost:8080

### 预置测试账号

首次启动后端时会自动创建测试数据（种子数据包含角色、技能和测试账号）：

| 角色      | 用户名      | 密码        | 说明                           |
| --------- | ----------- | ----------- | ------------------------------ |
| 普通玩家  | `player1`   | `123456`    | 已解锁全部角色，5000 金币 |
| 普通玩家  | `player2`   | `123456`    | 已解锁全部角色，5000 金币 |
| 管理员    | `admin`     | `admin123`  | 可登录管理后台                 |

**预置游戏角色（1 个）：**

| 角色     | HP   | 攻击 | 速度 | 类型           | 简介                           |
| -------- | ---- | ---- | ---- | -------------- | ------------------------------ |
| 烈焰战士 | 1200 | 120  | 180  | 近战爆发       | 攻高血厚，技能：烈焰斩/炎爆    |

> 测试之前先在其他浏览器或无痕窗口登录 `player1` 和 `player2`，就可以进行双人 PVP 对战。

### 测试流程（手把手）

双击 `start-dev.bat` 后，等待所有服务就绪浏览器自动打开，然后按以下步骤操作：

#### 第一步：打开两个浏览器窗口

| 窗口   | 方式                            | 用途     |
| ------ | ------------------------------- | -------- |
| 窗口 A | 正常浏览器标签页（已自动打开）  | 玩家1    |
| 窗口 B | 无痕窗口 或 另一个浏览器        | 玩家2    |

两个窗口都打开 `http://localhost:3000`。

#### 第二步：联网检测

前端打开后自动进入**启动检测页**：
- 自动检测服务器连接
- 连接成功 → 自动跳转登录页
- 连接失败 → 弹出网络异常弹窗（重连 / 退出）

#### 第三步：登录

- **窗口 A** → 输入 `player1` / `123456` → 登录 → 进入大厅
- **窗口 B** → 输入 `player2` / `123456` → 登录 → 进入大厅

#### 第四步：创建 & 加入房间

1. **窗口 A**（player1）→ 在大厅点击 **PVP 房间列表** → 进入房间列表页
2. 点击右下角 **创建房间** → 房间创建成功 → 跳转选人页
3. **窗口 B**（player2）→ 同样进入 PVP 房间列表，看到 player1 的房间
4. 点击 **加入房间** → 跳转选人页

#### 第五步：选择角色

1. 在选人页面左侧英雄列表中选择英雄
2. 右侧查看英雄立绘、皮肤、技能信息
3. 10 秒倒计时结束自动锁定，进入战斗

#### 第六步：战斗！

进入战斗场景后的操作键位：

| 按键 | 操作               | 说明                    |
| ---- | ------------------ | ----------------------- |
| W/A/S/D | 上下左右移动   | —                       |
| J    | 普攻               | 基础伤害，冷却 400ms    |
| K    | 技能1              | 消耗 30 能量，中等伤害  |
| L    | 技能2              | 消耗 50 能量，高伤害    |
| S    | 防御（按住）       | 减伤 70%                |

- 顶部居中显示 **99 秒倒计时**
- 左上/右上分别为双方 **血条 + 能量条**
- 底部四个 **技能图标**，显示冷却状态
- 击败对手后弹出 **对局结算弹窗**（胜负/HP 对比/积分/金币）

#### PVE 闯关测试

1. 从大厅点击 **PVE 闯关地图**
2. 选择关卡（简单/普通/困难/BOSS）→ 进入选人页
3. 选择已解锁英雄 → 点击 **开始战斗**
4. 对战 AI 敌人，战斗结束后结算金币和星级

#### 管理后台测试

浏览器打开 `http://localhost:3001` → 输入 `admin` / `admin123` 登录，查看 Dashboard 数据大屏和用户管理。

### 方式三：Makefile

```bash
# 查看所有命令
make

# 安装依赖
make install

# 构建生产版本
make build

# 开发模式
make dev-frontend   # 单独启动前端
make dev-admin      # 单独启动管理后台
make dev-server     # 单独启动后端

# Docker 操作
make docker-up      # 启动
make docker-down    # 停止
make docker-build   # 构建镜像
make docker-restart # 重启
```

## 环境变量

参见 `.env` 文件：

| 变量             | 说明            | 默认值                                               |
| ---------------- | --------------- | ---------------------------------------------------- |
| `SERVER_PORT`    | 后端端口        | `8080`                                               |
| `GIN_MODE`       | Gin 运行模式    | `debug`                                              |
| `DB_HOST`        | 数据库地址      | `127.0.0.1`                                          |
| `DB_PORT`        | 数据库端口      | `3306`                                               |
| `DB_USER`        | 数据库用户      | `root`                                               |
| `DB_PASSWORD`    | 数据库密码      | `root123`                                            |
| `DB_NAME`        | 数据库名        | `fight_game`                                         |
| `REDIS_HOST`     | Redis 地址      | `127.0.0.1`                                          |
| `REDIS_PORT`     | Redis 端口      | `6379`                                               |
| `REDIS_PASSWORD` | Redis 密码      | （空）                                               |
| `JWT_SECRET`     | JWT 签名密钥    | `change-this-to-a-random-secret-in-production`       |
| `LOG_LEVEL`      | 日志级别        | `info`                                               |
| `ALLOWED_ORIGINS`| CORS 允许域名   | `http://localhost:3000,http://localhost:3001`        |

## 构建部署

```bash
# 构建所有模块
make build

# 单独构建后端（输出到 server/bin/server）
make build-server

# 前端构建（输出到 frontend/dist）
cd frontend && npm run build

# 管理后台构建（输出到 admin/dist）
cd admin && npm run build
```

### 生产部署

使用 Docker Compose 进行生产部署前，请务必：

1. 修改 `.env` 中的 `JWT_SECRET` 为强随机字符串
2. 修改 `DB_PASSWORD` 为安全密码
3. 设置 `GIN_MODE=release`
4. 配置 SSL 证书到 `docker/nginx/ssl/`

```bash
docker-compose -f docker-compose.yml up -d --build
```

## 目录结构

```
FightGame/
├── frontend/                    # 游戏前端 (Vue3 + Phaser3)
│   ├── src/
│   │   ├── game/                # Phaser3 游戏引擎
│   │   │   ├── FightGame.ts     #   游戏入口配置
│   │   │   ├── BattleEngine.ts  #   战斗引擎核心逻辑
│   │   │   ├── AIEngine.ts      #   PVE AI 引擎
│   │   │   ├── enums.ts         #   枚举定义
│   │   │   └── scenes/          #   游戏场景
│   │   │       ├── BootScene.ts #     启动加载场景
│   │   │       ├── GameScene.ts #     主战斗场景 (HP条/能量/技能/帧同步)
│   │   │       └── PVEBattleScene.ts # PVE 战斗场景
│   │   ├── components/          # 通用组件
│   │   │   ├── ArtPlaceholder.vue     #   美术资源占位
│   │   │   ├── ModalOverlay.vue       #   弹窗容器
│   │   │   ├── GlobalNav.vue          #   全局导航栏
│   │   │   ├── LoadingSpinner.vue     #   加载动画
│   │   │   ├── NetworkErrorDialog.vue #   网络异常弹窗
│   │   │   ├── SettingsDialog.vue     #   系统设置弹窗
│   │   │   ├── FriendsDialog.vue      #   好友列表弹窗
│   │   │   └── BattleResultModal.vue  #   对局结算弹窗
│   │   ├── views/               # 页面视图 (14 页)
│   │   │   ├── ConnectionPage.vue     #   启动联网检测
│   │   │   ├── LoginRegisterPage.vue  #   登录/注册
│   │   │   ├── TutorialPage.vue       #   新手教程
│   │   │   ├── LobbyPage.vue          #   游戏大厅
│   │   │   ├── PVPRoomsPage.vue       #   PVP 房间列表
│   │   │   ├── PVPSelectPage.vue      #   PVP 英雄选人
│   │   │   ├── PVEStagePage.vue       #   PVE 关卡地图
│   │   │   ├── PVESelectPage.vue      #   PVE 英雄选人
│   │   │   ├── GamePage.vue           #   战斗场景 + HUD
│   │   │   ├── ReplayPage.vue         #   回放播放器
│   │   │   ├── ShopPage.vue           #   商城
│   │   │   └── HeroDexPage.vue        #   英雄图鉴
│   │   ├── api/                 # HTTP API 封装
│   │   ├── store/               # Pinia 状态管理
│   │   │   ├── user.ts          #   用户状态 (登录/信息/角色资产)
│   │   │   └── game.ts          #   游戏状态 (房间/战斗/帧同步)
│   │   ├── router/              # Vue Router 路由 + 路由守卫
│   │   └── utils/               # 工具函数
│   │       └── websocket.ts     #   WebSocket 客户端 (自动重连/消息分发)
│   └── vite.config.ts
│
├── admin/                       # 管理后台 (Vue3 + Element Plus)
│   ├── src/
│   │   ├── api/                 # API 封装
│   │   ├── router/              # 路由配置 (dashboard/user/game/log)
│   │   └── views/
│   │       ├── LoginPage.vue    #   管理员登录
│   │       ├── dashboard/
│   │       │   └── Dashboard.vue#   数据概览仪表盘
│   │       ├── user/
│   │       │   └── UserList.vue #   用户管理列表
│   │       ├── game/
│   │       │   └── GameManage.vue#  游戏管理
│   │       └── log/
│   │           └── LogView.vue  #   日志查看
│   └── vite.config.ts
│
├── server/                      # Go 后端
│   ├── cmd/
│   │   └── main.go              # 入口 (路由注册/优雅退出)
│   └── internal/
│       ├── config/config.go     # 环境变量配置
│       ├── model/model.go       # 15张数据表 + DB/Redis 初始化
│       ├── handler/
│       │   ├── user.go          #   用户登录/注册/信息
│       │   ├── admin.go         #   管理员登录/Dashboard/用户管理
│       │   ├── character.go     #   角色列表/详情/技能
│       │   ├── profile.go       #   个人资料/密码/角色
│       │   ├── skin.go          #   皮肤查询/拥有
│       │   ├── room.go          #   房间列表/创建/加入
│       │   ├── shop.go          #   商城商品/购买
│       │   ├── friend.go        #   好友增删查
│       │   ├── battle.go        #   回放数据
│       │   └── pve.go           #   PVE 关卡/进度
│       ├── middleware/
│       │   └── middleware.go    #   JWT 认证/管理员权限/CORS/日志
│       ├── websocket/
│       │   ├── websocket.go     #   Hub/Client/消息处理(12种消息类型)
│       │   └── room.go          #   GameRoom 状态机/RoomManager
│       └── logger/logger.go     #   Zap 日志系统
│
├── docker/                      # Docker 配置
│   ├── Dockerfile.server        # 后端镜像 (golang:1.21-alpine)
│   ├── Dockerfile.frontend      # 前端镜像 (node:20-alpine → nginx)
│   ├── Dockerfile.admin         # 管理后台镜像 (node:20-alpine → nginx)
│   ├── mysql/init.sql           # 数据库初始化脚本
│   └── nginx/
│       ├── nginx.conf           # Nginx 反向代理配置
│       └── ssl/                 # SSL 证书
│
├── scripts/                     # 工具脚本
│   ├── generate-ssl.bat/.sh     # SSL 证书生成
│   └── *.js                     # 内部检查脚本
│
├── docker-compose.yml           # 服务编排 (MySQL+Redis+Server+Nginx)
├── Makefile                     # 构建/开发/部署命令
├── start-dev.bat                # Windows 一键启动脚本
├── .env                         # 环境变量
└── README.md
```

## 核心设计

### 房间状态机

```
waiting ──→ selecting ──→ ready ──→ playing ──→ finished
   │            │           │           │
   └── 创建/加入   选择角色    双方就绪    战斗结束
```

### 战斗系统

- **移动**: WASD 四方向，带速度系数
- **普攻 (J)**: 基础伤害，冷却 400ms
- **技能1 (K)**: 消耗 30 能量，较高伤害
- **技能2 (L)**: 消耗 50 能量，最高伤害
- **防御 (S)**: 按住减伤 70%，松开恢复
- **攻击判定**: 100px 范围内命中判定
- **能量系统**: 初始 100，每秒回复 2 点
- **HP 系统**: 根据角色属性初始化，归零判负
- **对局时间**: 99 秒倒计时，超时按 HP 比例判胜负
- **HUD 层**: Vue 组件实现的战斗 UI（血条/能量条/倒计时/技能图标），与 Phaser 画布叠加

### 帧同步

客户端每帧（约 60fps）通过 WebSocket 发送操作指令，服务端标记帧号后广播给对手，确保双方所见一致。

### 前端路由守卫

- 需要登录的页面（`meta.requiresAuth`）未登录自动跳转 `/login`
- 已登录访问 `/login` 自动跳转 `/lobby`
- 启动页 `/` 自动检测 token 有效性，决定跳转登录或大厅

## 开发团队

本项目使用 DeepSeek AI 辅助开发。

## 许可证

MIT License
