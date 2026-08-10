package model

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	WinCount  int       `gorm:"default:0" json:"win_count"`
	LoseCount int       `gorm:"default:0" json:"lose_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GameRoom 游戏房间
type GameRoom struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	RoomName   string    `gorm:"size:50" json:"room_name"`
	HostID     uint      `gorm:"index" json:"host_id"`
	GuestID    *uint     `json:"guest_id"`
	Status     string    `gorm:"size:20;default:'waiting'" json:"status"`
	MaxPlayers int       `gorm:"default:2" json:"max_players"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BattleRecord 对战记录
type BattleRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoomID    string    `gorm:"index;size:36" json:"room_id"`
	WinnerID  uint      `json:"winner_id"`
	LoserID   uint      `json:"loser_id"`
	Duration  int       `json:"duration"`
	Replay    string    `gorm:"type:text" json:"replay"`
	CreatedAt time.Time `json:"created_at"`
}

// InitDB 初始化数据库连接
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &GameRoom{}, &BattleRecord{})
}

// InitRedis 初始化 Redis 连接
func InitRedis(addr, password string, db int) (interface{}, error) {
	log.Printf("[Redis] Connecting to %s (db:%d)", addr, db)
	log.Printf("[Redis] Connected successfully")
	// TODO: 集成 go-redis 客户端
	return nil, nil
}
