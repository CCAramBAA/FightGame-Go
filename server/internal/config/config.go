package config

import (
	"fmt"
	"os"
	"strings"
)

// Config 应用配置
type Config struct {
	ServerPort     string
	GinMode        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	RedisDB        int
	LogLevel       string
	LogFile        string
	JWTSecret      string
	AllowedOrigins []string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "fight_game"),
		RedisHost:      getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        0,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFile:        getEnv("LOG_FILE", "logs/server.log"),
		JWTSecret:      getEnv("JWT_SECRET", "fight-game-secret-key-change-in-production"),
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}
}

// DSN 生成 MySQL 连接字符串
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// RedisAddr 返回 Redis 地址
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
