package initizalize

import (
	"CrestedIbis/src/config/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// InitLogger 初始化日志配置
// 1. 设置日志文件
// 2. 格式化日志输出格式
// 3. 设置日志输出流（无file：标注输出流，有file：标注输出流和文件流）
// 4. 设置日志基本（默认Info）
func InitLogger(log *model.Log) *logrus.Logger {
	logger := logrus.New()

	// 格式化日志
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var output io.Writer
	// 设置日志输出流
	if log.File != "" {
		file, err := os.OpenFile(log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file: %s", log.File))
		}
		output = io.MultiWriter(file, os.Stdout)
	} else {
		output = os.Stdout
	}
	logger.Out = output

	// 设置日志级别
	if log.Level != "" {
		level, err := logrus.ParseLevel(log.Level)
		if err != nil {
			panic(fmt.Sprintf("Invalid log level: %s", log.Level))
		}
		logger.SetLevel(level)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
