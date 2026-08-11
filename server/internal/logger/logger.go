package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
	recentLogs  []string
	mu          sync.RWMutex
	maxRecent   = 500
)

func Init() error {
	logsDir := filepath.Join(".", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("create logs dir: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(logsDir, "server.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	multiInfo := io.MultiWriter(os.Stdout, f)
	multiWarn := io.MultiWriter(os.Stdout, f)
	multiError := io.MultiWriter(os.Stdout, f)

	flag := log.Ldate | log.Ltime

	infoLogger = log.New(multiInfo, "[INFO]  ", flag)
	warnLogger = log.New(multiWarn, "[WARN]  ", flag)
	errorLogger = log.New(multiError, "[ERROR] ", flag)
	debugLogger = log.New(multiInfo, "[DEBUG] ", flag)
	logFile = f

	recentLogs = make([]string, 0, maxRecent)

	Info("Logger initialized")
	return nil
}

func recordRecent(level, msg string) {
	mu.Lock()
	defer mu.Unlock()
	entry := fmt.Sprintf("[%s] %s %s", level, time.Now().Format("2006-01-02 15:04:05"), msg)
	recentLogs = append(recentLogs, entry)
	if len(recentLogs) > maxRecent {
		recentLogs = recentLogs[len(recentLogs)-maxRecent:]
	}
}

func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if infoLogger != nil {
		infoLogger.Output(2, msg)
	}
	recordRecent("INFO", msg)
}

func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if warnLogger != nil {
		warnLogger.Output(2, msg)
	}
	recordRecent("WARN", msg)
}

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if errorLogger != nil {
		errorLogger.Output(2, msg)
	}
	recordRecent("ERROR", msg)
}

func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if debugLogger != nil {
		debugLogger.Output(2, msg)
	}
	recordRecent("DEBUG", msg)
}

func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if errorLogger != nil {
		errorLogger.Output(2, "FATAL: "+msg)
	}
	recordRecent("FATAL", msg)
	if logFile != nil {
		logFile.Sync()
	}
	os.Exit(1)
}

func Sync() {
	if logFile != nil {
		logFile.Sync()
	}
}

func ReadRecentLogs(max int) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	if max <= 0 || max > len(recentLogs) {
		max = len(recentLogs)
	}
	result := make([]string, max)
	copy(result, recentLogs[len(recentLogs)-max:])
	return result, nil
}
