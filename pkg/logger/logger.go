package logger

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 自定义颜色（ANSI 转义码）
var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorGray   = "\033[37m"
)

// CustomFormatter 实现 logrus.Formatter 接口，自定义日志格式和颜色
type CustomFormatter struct{}

// Format 格式化日志条目，添加颜色和时间戳
func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	// 时间格式：1999.9.9 8:10
	timestamp := entry.Time.Format(" 2006.1.2 15:04:05")

	// 颜色根据级别变化
	var levelColor string
	switch entry.Level {
	case log.DebugLevel:
		levelColor = colorGray
	case log.InfoLevel:
		levelColor = colorGreen
	case log.WarnLevel:
		levelColor = colorYellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = colorRed
	default:
		levelColor = colorBlue
	}

	// [INFO]1999.9.9 8:10 message...
	level := strings.ToUpper(entry.Level.String())
	logLine := fmt.Sprintf("%s[%s]%s%s %s%s\n",
		levelColor, level, colorReset, timestamp, entry.Message, colorReset,
	)
	return []byte(logLine), nil
}

// InitLogger 初始化日志系统
func InitLogger(level string) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&CustomFormatter{})

	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
