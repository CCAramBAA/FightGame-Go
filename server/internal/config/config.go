package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ==================== .env 加载 ====================

// findProjectRoot 向上查找包含 go.mod 的项目根目录
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		_, file, _, _ := runtime.Caller(0)
		dir = filepath.Dir(file)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return dir
}

// LoadEnv 从项目根目录 .env 文件加载环境变量（仅设置未设置的变量）
func LoadEnv() {
	projectRoot := findProjectRoot()
	// 从项目根目录的上级目录（即 FightGame 根目录）查找 .env
	root := filepath.Dir(projectRoot)
	envFile := filepath.Join(root, ".env")
	// 如果上级目录没有，则在当前项目根目录查找
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		envFile = filepath.Join(projectRoot, ".env")
	}

	file, err := os.Open(envFile)
	if err != nil {
		return // .env 文件不存在不是致命错误
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 去除行尾注释（值中可能包含 #）
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if key == "" || value == "" {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

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
	// 优先加载 .env 文件（Docker Compose 中环境变量已设置，不会被覆盖）
	LoadEnv()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
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
		RedisDB:        redisDB,
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
