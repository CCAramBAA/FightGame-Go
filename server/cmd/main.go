package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fightgame-server/internal/config"
	"fightgame-server/internal/cron"
	"fightgame-server/internal/handler"
	"fightgame-server/internal/logger"
	"fightgame-server/internal/middleware"
	"fightgame-server/internal/model"
	ws "fightgame-server/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// 加载配置
	cfg := config.Load()

	// 初始化 JWT
	middleware.InitJWT(cfg.JWTSecret)

	// 初始化数据库
	db, err := model.InitDB(cfg.DSN())
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database connected")

	// 自动迁移
	if err := model.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}
	logger.Info("Database migration completed")

	// 种子数据
	if err := model.SeedData(db); err != nil {
		logger.Warn("Seed data warning: %v", err)
	}

	// 初始化 Redis（非必需，不可用时不影响基本功能）
	var rdb *redis.Client
	rdb, err = model.InitRedis(cfg.RedisAddr(), cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Warn("Redis not available, some features may be limited: %v", err)
		rdb = nil
	}

	// 初始化 WebSocket Hub
	hub := ws.NewHub(cfg.JWTSecret)
	go hub.Run()
	logger.Info("WebSocket Hub started")

	// 创建 handler（传递 hub 用于在线状态等）
	userHandler := handler.NewUserHandler(db, rdb, cfg.JWTSecret)
	characterHandler := handler.NewCharacterHandler(db)
	roomHandler := handler.NewRoomHandler(hub)
	adminHandler := handler.NewAdminHandler(db, hub, cfg.JWTSecret)
	friendHandler := handler.NewFriendHandler(db, hub)
	shopHandler := handler.NewShopHandler(db)
	battleHandler := handler.NewBattleHandler(db)
	pveHandler := handler.NewPVEHandler(db)
	skinHandler := handler.NewSkinHandler(db)
	profileHandler := handler.NewProfileHandler(db)

	// 创建 Gin 引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// ===================== 公开路由 =====================
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
	})

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// ===================== 管理员登录（无需 AuthRequired） =====================
	api.POST("/admin/login", adminHandler.Login)

	// ===================== 需要登录的路由 =====================
	auth := api.Group("")
	auth.Use(middleware.AuthRequired())
	{
		// 用户
		auth.POST("/logout", userHandler.Logout)
		auth.POST("/refresh-token", userHandler.RefreshToken)

		// 角色 & 皮肤
		auth.GET("/characters", characterHandler.ListCharacters)
		auth.GET("/characters/:id", characterHandler.GetCharacter)
		auth.GET("/characters/:id/skills", characterHandler.ListSkills)
		auth.GET("/skins", skinHandler.ListSkins)
		auth.GET("/skins/my", skinHandler.GetUserSkins)

		// 玩家资产
		auth.GET("/profile", profileHandler.GetProfile)
		auth.GET("/profile/characters", profileHandler.GetMyCharacters)
		auth.GET("/profile/rank", profileHandler.GetRankScore)
		auth.GET("/profile/gold", shopHandler.GetGoldTransactions)

		// 商城
		auth.GET("/shop/items", shopHandler.ListShopItems)
		auth.POST("/shop/purchase", shopHandler.PurchaseItem)

		// 好友
		auth.GET("/friends", friendHandler.ListFriends)
		auth.POST("/friends/add", friendHandler.AddFriend)
		auth.DELETE("/friends/:id", friendHandler.RemoveFriend)
		auth.POST("/friends/invite", friendHandler.SendInvite)

		// PVE
		auth.GET("/pve/stages", pveHandler.ListStages)
		auth.GET("/pve/stages/:id", pveHandler.GetStage)
		auth.GET("/pve/progress", pveHandler.GetProgress)
		auth.POST("/pve/progress", pveHandler.SaveProgress)

		// 对局记录 & 回放
		auth.POST("/battle/record", battleHandler.RecordBattle)
		auth.GET("/battle/records", battleHandler.GetBattleRecords)
		auth.GET("/battle/replay/:id", battleHandler.GetReplay)

		// 房间
		auth.GET("/rooms", roomHandler.ListRooms)
		auth.POST("/rooms/create", roomHandler.CreateRoom)
		auth.POST("/rooms/join/:id", roomHandler.JoinRoom)
		auth.POST("/rooms/leave", roomHandler.LeaveRoom)

		// 密码修改
		auth.POST("/profile/password", profileHandler.ChangePassword)
	}

	// ===================== WebSocket =====================
	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWS(hub, c.Writer, c.Request)
	})

	// ===================== 管理后台路由（需管理员权限） =====================
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		// Dashboard
		admin.GET("/dashboard", adminHandler.Dashboard)

		// 用户管理
		admin.GET("/users", adminHandler.ListUsers)
		admin.PUT("/users/:id", adminHandler.UpdateUser)

		// 角色管理
		admin.GET("/characters", adminHandler.ListCharacters)
		admin.POST("/characters", adminHandler.CreateCharacter)
		admin.PUT("/characters/:id", adminHandler.UpdateCharacter)
		admin.DELETE("/characters/:id", adminHandler.DeleteCharacter)

		// 技能管理
		admin.GET("/skills", adminHandler.ListSkills)
		admin.POST("/skills", adminHandler.CreateSkill)
		admin.PUT("/skills/:id", adminHandler.UpdateSkill)
		admin.DELETE("/skills/:id", adminHandler.DeleteSkill)

		// 皮肤管理
		admin.GET("/skins", adminHandler.ListSkins)
		admin.POST("/skins", adminHandler.CreateSkin)
		admin.PUT("/skins/:id", adminHandler.UpdateSkin)
		admin.DELETE("/skins/:id", adminHandler.DeleteSkin)

		// 英雄特殊交互规则
		admin.GET("/special-rules", adminHandler.ListSpecialRules)
		admin.POST("/special-rules", adminHandler.CreateSpecialRule)
		admin.PUT("/special-rules/:id", adminHandler.UpdateSpecialRule)
		admin.DELETE("/special-rules/:id", adminHandler.DeleteSpecialRule)

		// 商城管理
		admin.GET("/shop/items", adminHandler.ListShopItems)
		admin.POST("/shop/items", adminHandler.CreateShopItem)
		admin.PUT("/shop/items/:id", adminHandler.UpdateShopItem)
		admin.DELETE("/shop/items/:id", adminHandler.DeleteShopItem)

		// PVE 关卡管理
		admin.GET("/pve/stages", adminHandler.ListStages)
		admin.POST("/pve/stages", adminHandler.CreateStage)
		admin.PUT("/pve/stages/:id", adminHandler.UpdateStage)
		admin.DELETE("/pve/stages/:id", adminHandler.DeleteStage)

		// 房间管理
		admin.GET("/rooms", adminHandler.ListRooms)
		admin.DELETE("/rooms/:id", adminHandler.ForceCloseRoom)

		// 日志
		admin.GET("/logs", adminHandler.ViewLogs)

		// 金币流水
		admin.GET("/gold-transactions", adminHandler.ListGoldTransactions)

		// 对局记录
		admin.GET("/battle-records", adminHandler.ListBattleRecords)

		// 缓存刷新
		admin.POST("/cache/refresh", adminHandler.RefreshCache)
	}

	// ===================== 启动定时任务 =====================
	cronMgr := cron.NewTaskManager(db)
	cronMgr.Start(hub.RoomMgr)
	logger.Info("Cron tasks started")

	// ===================== 启动 HTTP 服务 =====================
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// ===================== 优雅退出 =====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	cronMgr.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited gracefully")
}
