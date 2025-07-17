package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

// 初始化日志
func Init() error {
	// 创建log目录
	if err := os.MkdirAll("log", 0755); err != nil {
		return fmt.Errorf("创建log目录失败: %v", err)
	}

	// 获取当前日期
	today := time.Now().Format("2006-01-02")

	// 创建info日志文件
	infoFile, err := os.OpenFile(
		filepath.Join("log", fmt.Sprintf("%s-info.log", today)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("创建info日志文件失败: %v", err)
	}

	// 创建error日志文件
	errorFile, err := os.OpenFile(
		filepath.Join("log", fmt.Sprintf("%s-error.log", today)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("创建error日志文件失败: %v", err)
	}

	// 创建debug日志文件
	debugFile, err := os.OpenFile(
		filepath.Join("log", fmt.Sprintf("%s-debug.log", today)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("创建debug日志文件失败: %v", err)
	}

	// 初始化日志记录器
	infoLogger = log.New(infoFile, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(errorFile, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	debugLogger = log.New(debugFile, "[DEBUG] ", log.LstdFlags|log.Lshortfile)

	return nil
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}

// Fatal 记录致命错误并退出程序
func Fatal(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
	os.Exit(1)
}
