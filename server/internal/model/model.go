package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ===================== 账号与资产 =====================

// User 用户账号表
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	Gold      int64     `gorm:"default:0" json:"gold"`
	RankScore int       `gorm:"default:1000" json:"rank_score"`
	WinCount  int       `gorm:"default:0" json:"win_count"`
	LoseCount int       `gorm:"default:0" json:"lose_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GoldTransaction 金币流水表
type GoldTransaction struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Amount       int64     `gorm:"not null" json:"amount"`        // 正数为收入，负数为支出
	SourceType   string    `gorm:"size:30;not null" json:"source_type"` // pvp_win / pve_star / shop_buy
	SourceID     string    `gorm:"size:50" json:"source_id"`      // 来源标识（对局ID/关卡ID/商品ID）
	BalanceAfter int64     `gorm:"not null" json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

// ===================== 英雄与技能配置 =====================

// Character 角色表
type Character struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Story       string    `gorm:"type:text" json:"story"`
	UnlockType  string    `gorm:"size:20;default:'gold'" json:"unlock_type"` // gold / pve_stage
	UnlockPrice int64     `gorm:"default:0" json:"unlock_price"`             // 金币购买价格
	UnlockStageID uint    `gorm:"default:0" json:"unlock_stage_id"`          // PVE关卡解锁条件ID
	HP          int       `gorm:"default:1000" json:"hp"`
	Energy      int       `gorm:"default:100" json:"energy"`
	EnergyRegen int       `gorm:"default:5" json:"energy_regen"`    // 每秒回复
	Speed       int       `gorm:"default:200" json:"speed"`         // 移速
	Attack      int       `gorm:"default:100" json:"attack"`
	Defense     int       `gorm:"default:50" json:"defense"`
	AvatarPath  string    `gorm:"size:255" json:"avatar_path"`      // 立绘路径
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Skill 技能表
type Skill struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CharacterID   uint      `gorm:"index;not null" json:"character_id"`
	Name          string    `gorm:"size:50;not null" json:"name"`
	SkillType     string    `gorm:"size:30;default:'active'" json:"skill_type"` // active / passive / transform
	EnergyCost    int       `gorm:"default:10" json:"energy_cost"`
	CoolDown      int       `gorm:"default:5" json:"cooldown"`                 // 冷却秒数
	Damage        int       `gorm:"default:0" json:"damage"`
	Range         float64   `gorm:"default:150" json:"range"`                  // 作用范围（像素）
	PriorityLevel int       `gorm:"default:1" json:"priority_level"`           // 技能优先级层级（越大越高）
	Tags          string    `gorm:"size:255;default:''" json:"tags"`           // 逗号分隔标签: unstoppable,shield,steal,lifesteal,silence,armor_break,stun,knockup,slow,transform
	Description   string    `gorm:"type:text" json:"description"`
	EffectData    string    `gorm:"type:text" json:"effect_data"`              // JSON: 效果参数（护盾值、吸血比例等）
	SoundPath     string    `gorm:"size:255" json:"sound_path"`
	VFXPath       string    `gorm:"size:255" json:"vfx_path"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// HeroSpecialRule 英雄特殊交互规则表
type HeroSpecialRule struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	CharacterID   uint   `gorm:"index;not null" json:"character_id"`
	Name          string `gorm:"size:100;not null" json:"name"`
	Description   string `gorm:"type:text" json:"description"`
	PriorityOrder int    `gorm:"default:0;not null" json:"priority_order"` // 行内执行顺序（从小到大）
	EffectType    string `gorm:"size:50;not null" json:"effect_type"`      // steal_skill / pierce_shield / immune_damage / priority_override
	Condition     string `gorm:"type:text" json:"condition"`               // JSON: 触发条件
	EffectData    string `gorm:"type:text" json:"effect_data"`             // JSON: 效果参数
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Skin 皮肤表
type Skin struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CharacterID  uint      `gorm:"index;not null" json:"character_id"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	Price        int64     `gorm:"default:0" json:"price"` // 金币价格
	ResourcePath string    `gorm:"size:255" json:"resource_path"`  // 贴图路径
	PreviewPath  string    `gorm:"size:255" json:"preview_path"`   // 预览图路径
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ===================== 用户资产 =====================

// UserCharacter 用户角色解锁表
type UserCharacter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_char;not null" json:"user_id"`
	CharacterID uint      `gorm:"uniqueIndex:idx_user_char;not null" json:"character_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserSkin 用户皮肤持有表
type UserSkin struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_skin;not null" json:"user_id"`
	SkinID      uint      `gorm:"uniqueIndex:idx_user_skin;not null" json:"skin_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ===================== PVE 系统 =====================

// PVEStage PVE 关卡配置表
type PVEStage struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"size:100;not null" json:"name"`
	Difficulty      string `gorm:"size:20;default:'normal'" json:"difficulty"` // easy/normal/hard/boss
	StageOrder      int    `gorm:"default:0" json:"stage_order"`               // 关卡顺序
	UnlockCondition string `gorm:"type:text" json:"unlock_condition"`          // JSON: 解锁条件
	BossConfig      string `gorm:"type:text" json:"boss_config"`               // JSON: BOSS多阶段配置
	RewardGold      int64  `gorm:"default:0" json:"reward_gold"`               // 满星金币
	RewardCharID    uint   `gorm:"default:0" json:"reward_char_id"`            // 通关解锁角色
	StarGold1       int64  `gorm:"default:0" json:"star_gold_1"`               // 1星金币
	StarGold2       int64  `gorm:"default:0" json:"star_gold_2"`
	StarGold3       int64  `gorm:"default:0" json:"star_gold_3"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PVEProgress PVE 玩家存档表
type PVEProgress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_stage;not null" json:"user_id"`
	StageID     uint      `gorm:"uniqueIndex:idx_user_stage;not null" json:"stage_id"`
	Stars       int       `gorm:"default:0" json:"stars"`                    // 当前最高星级
	CompletedAt *time.Time `json:"completed_at"`                              // 通关时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ===================== PVP 房间与对战 =====================

// GameRoom 游戏房间
type GameRoom struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	HostID     uint      `gorm:"index" json:"host_id"`
	GuestID    *uint     `json:"guest_id"`
	HostCharID *uint     `json:"host_char_id"`     // 房主选的角色
	GuestCharID *uint    `json:"guest_char_id"`     // 客方选的角色
	Status     string    `gorm:"size:20;default:'waiting'" json:"status"` // waiting/selecting/playing/finished
	SelectDeadline *time.Time `json:"select_deadline"` // 选人截止时间
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BattleRecord 对战记录（帧数据存储）
type BattleRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomID     string    `gorm:"index;size:36" json:"room_id"`
	Player1ID  uint      `json:"player1_id"`
	Player2ID  uint      `json:"player2_id"`
	Char1ID    uint      `json:"char1_id"`
	Char2ID    uint      `json:"char2_id"`
	WinnerID   *uint     `json:"winner_id"`       // NULL=平局
	LoserID    *uint     `json:"loser_id"`
	IsDraw     bool      `gorm:"default:false" json:"is_draw"`
	Duration   int       `json:"duration"`         // 对局秒数
	ScoreDelta int       `gorm:"default:0" json:"score_delta"` // 段位分变动
	FrameData  string    `gorm:"type:longtext" json:"frame_data"` // 完整帧数据JSON
	CreatedAt  time.Time `json:"created_at"`
}

// ===================== 好友系统 =====================

// FriendRelation 好友关系表
type FriendRelation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	FriendID  uint           `gorm:"index;not null" json:"friend_id"`
	Status    string         `gorm:"size:20;default:'accepted'" json:"status"` // pending/accepted/blocked
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Friend    User           `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

// ===================== 商城 =====================

// ShopItem 商城商品配置表
type ShopItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ItemType    string    `gorm:"size:20;not null" json:"item_type"` // character / skin
	ItemID      uint      `gorm:"not null" json:"item_id"`           // character_id 或 skin_id
	Price       int64     `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ===================== 数据库初始化 =====================

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

// AutoMigrate 自动迁移所有表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Character{},
		&Skill{},
		&HeroSpecialRule{},
		&Skin{},
		&UserCharacter{},
		&UserSkin{},
		&GoldTransaction{},
		&PVEStage{},
		&PVEProgress{},
		&GameRoom{},
		&BattleRecord{},
		&FriendRelation{},
		&ShopItem{},
	)
}

// InitRedis 初始化 Redis 连接
func InitRedis(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Printf("[Redis] Connected successfully to %s (db:%d)", addr, db)
	return rdb, nil
}

// ===================== 种子数据 =====================

// SeedData 初始化测试数据（幂等：已存在则跳过）
func SeedData(db *gorm.DB) error {
	log.Println("[Seed] Checking seed data...")

	// 创建测试角色（仅保留一个用于开发调试）
	characters := []Character{
		{Name: "烈焰战士", Description: "擅长近距离爆发输出，攻高血厚", HP: 1200, Energy: 100, EnergyRegen: 3, Speed: 180, Attack: 120, Defense: 60},
	}
	for i := range characters {
		if err := db.Where("name = ?", characters[i].Name).FirstOrCreate(&characters[i]).Error; err != nil {
			return fmt.Errorf("seed character %s: %w", characters[i].Name, err)
		}
	}

	// 为烈焰战士创建技能
	if db.Where("character_id = ?", characters[0].ID).First(&Skill{}).Error != nil {
		skills1 := []Skill{
			{CharacterID: characters[0].ID, Name: "烈焰斩", SkillType: "active", EnergyCost: 30, CoolDown: 3, Damage: 80, Range: 120, PriorityLevel: 1, Description: "向前挥出火焰剑气，造成80点伤害"},
			{CharacterID: characters[0].ID, Name: "炎爆", SkillType: "active", EnergyCost: 50, CoolDown: 8, Damage: 150, Range: 100, PriorityLevel: 2, Description: "蓄力释放火焰爆炸，造成150点伤害"},
			{CharacterID: characters[0].ID, Name: "战意", SkillType: "passive", EnergyCost: 0, CoolDown: 0, Damage: 0, Range: 0, PriorityLevel: 0, Description: "被动：生命低于30%时攻击力+20%"},
		}
		for _, s := range skills1 {
			db.Create(&s)
		}
	}

	// 创建测试用户
	seedUsers := []struct {
		Username string
		Password string
		Nickname string
		IsAdmin  bool
		Gold     int64
	}{
		{"player1", "123456", "玩家一号", false, 5000},
		{"player2", "123456", "玩家二号", false, 5000},
		{"admin", "admin123", "管理员", true, 0},
	}

	for _, su := range seedUsers {
		var existing User
		if err := db.Where("username = ?", su.Username).First(&existing).Error; err == nil {
			continue // 已存在则跳过
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(su.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seed user %s hash password: %w", su.Username, err)
		}

		user := User{
			Username: su.Username,
			Password: string(hashedPwd),
			Nickname: su.Nickname,
			IsAdmin:  su.IsAdmin,
			Gold:     su.Gold,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("seed user %s: %w", su.Username, err)
		}

		// 为非管理员玩家解锁所有角色
		if !su.IsAdmin {
			for _, ch := range characters {
				uc := UserCharacter{UserID: user.ID, CharacterID: ch.ID}
				db.Where("user_id = ? AND character_id = ?", user.ID, ch.ID).FirstOrCreate(&uc)
			}
		}

		log.Printf("[Seed] Created user: %s (id=%d)", su.Username, user.ID)
	}

	log.Println("[Seed] Seed data ready.")
	return nil
}
