package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Logger 应用日志器
type Logger struct {
	level string
	mu    sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

// Init 初始化日志系统
func Init(level, logFile string) {
	once.Do(func() {
		instance = &Logger{level: level}

		// 创建日志目录
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Printf("Failed to create logs directory: %v", err)
		}

		// 设置日志输出
		if logFile != "" {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Printf("Failed to open log file: %v, using stdout", err)
			} else {
				log.SetOutput(file)
			}
		}

		log.SetFlags(log.LstdFlags | log.Lshortfile)
	})
}

// log 输出日志
func (l *Logger) log(level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	log.Printf("[%s] %s", level, msg)
}

// Sync 关闭日志资源
func Sync() {
	// 刷新日志缓冲区
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	if instance != nil {
		instance.log("INFO", format, args...)
	}
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	if instance != nil {
		instance.log("WARN", format, args...)
	}
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	if instance != nil {
		instance.log("ERROR", format, args...)
	}
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	if instance != nil {
		instance.log("DEBUG", format, args...)
	}
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Info(format, args...)
}

// Fatalf 致命错误日志并退出
func Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("[FATAL] %s", msg)
}
