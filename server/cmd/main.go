package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fightgame-server/internal/config"
	"fightgame-server/internal/handler"
	"fightgame-server/internal/logger"
	"fightgame-server/internal/middleware"
	"fightgame-server/internal/model"
	"fightgame-server/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.LogLevel, cfg.LogFile)
	defer logger.Sync()

	logger.Info("Starting FightGame Server...")

	// 初始化数据库
	db, err := model.InitDB(cfg.DSN())
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移
	if err := model.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis
	redisClient, err := model.InitRedis(cfg.RedisAddr(), cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 初始化 WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 设置 Gin
	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Unix()})
	})

	// API 路由组
	api := r.Group("/api")
	{
		userHandler := handler.NewUserHandler(db, redisClient)
		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)

		auth := api.Group("")
		auth.Use(middleware.Auth(redisClient))
		{
			auth.GET("/user/info", userHandler.GetUserInfo)
			auth.POST("/user/logout", userHandler.Logout)
		}
	}

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWS(wsHub, c.Writer, c.Request)
	})

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: r,
	}

	go func() {
		logger.Infof("Server listening on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
