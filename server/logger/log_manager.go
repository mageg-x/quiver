// logger/logger.go
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LogDir      = "./logs"  // 日志根目录
	MaxSize     = 10        // 每个日志文件最大 10MB
	MaxBackups  = 7         // 保留 7 个备份
	MaxAge      = 30        // 30 天过期
	Compress    = true      // 压缩旧日志
	DefaultFile = "app.log" // 默认日志文件
)

var (
	loggers = make(map[string]*logrus.Logger)
	mu      sync.Mutex
)

// getLoggerKey 标准化 logger 名称（避免路径问题）
func getLoggerKey(name string) string {
	name = strings.TrimSuffix(name, ".log")
	name = filepath.Base(name)
	return fmt.Sprintf("%s.log", name)
}

// === 封装的 Logger 结构体 ===
type Logger struct {
	*logrus.Logger
}

// getCaller 获取调用 Infof/Debugf 等函数的源码位置
// skip = 2: 跳过 getCaller 和 外层日志函数（如 Infof）
func (l *Logger) getCaller() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???:0"
	}
	_, filename := filepath.Split(file)
	fn := runtime.FuncForPC(pc)
	funcName := fn.Name()
	parts := strings.Split(funcName, ".")
	return fmt.Sprintf("%s:%s:%d", filename, parts[len(parts)-1], line)
}

// === Debug 级别 ===
func (l *Logger) Debug(args ...any) {
	caller := l.getCaller()
	l.Logger.Debug(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Debugf(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Debugf("%s "+format, append([]any{caller}, args...)...)
}

// === Info 级别 ===
func (l *Logger) Info(args ...any) {
	caller := l.getCaller()
	l.Logger.Info(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Infof(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Infof("%s "+format, append([]any{caller}, args...)...)
}

// === Warn 级别 ===
func (l *Logger) Warn(args ...any) {
	caller := l.getCaller()
	l.Logger.Warn(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Warnf(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Warnf("%s "+format, append([]any{caller}, args...)...)
}

// === Error 级别 ===
func (l *Logger) Error(args ...any) {
	caller := l.getCaller()
	l.Logger.Error(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Errorf(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Errorf("%s "+format, append([]any{caller}, args...)...)
}

// === Fatal 级别（会 os.Exit(1)）===
func (l *Logger) Fatal(args ...any) {
	caller := l.getCaller()
	l.Logger.Fatal(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Fatalf(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Fatalf("%s "+format, append([]any{caller}, args...)...)
}

// === Panic 级别（会 panic）===
func (l *Logger) Panic(args ...any) {
	caller := l.getCaller()
	l.Logger.Panic(fmt.Sprintf("%s %v", caller, fmt.Sprint(args...)))
}

func (l *Logger) Panicf(format string, args ...any) {
	caller := l.getCaller()
	l.Logger.Panicf("%s "+format, append([]any{caller}, args...)...)
}

// === GetLogger：返回封装后的 Logger ===
func GetLogger(name string) *Logger {
	key := getLoggerKey(name)

	mu.Lock()
	defer mu.Unlock()

	// 如果已存在，直接返回封装的 Logger
	if logger, exists := loggers[key]; exists {
		return &Logger{logger}
	}

	// 创建日志目录
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		fmt.Printf("无法创建日志目录: %v\n", err)
		return &Logger{logrus.StandardLogger()} // fallback
	}

	// 日志文件路径
	logPath := filepath.Join(LogDir, key)

	// 创建新的 logrus.Logger
	logger := logrus.New()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   Compress,
	})
	logger.SetLevel(logrus.DebugLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 缓存原始 *logrus.Logger
	loggers[key] = logger

	// 返回封装的 Logger
	return &Logger{logger}
}
