package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	Logger = logrus.New()

	logLevel := os.Getenv("LOG_LEVEL")
	logFormat := os.Getenv("LOG_FORMAT")
	output := os.Getenv("LOG_OUTPUT")
	filePath := os.Getenv("LOG_FILE")

	if logLevel == "" {
		logLevel = "info"
	}
	if output == "" {
		output = "both"
	}
	if filePath == "" {
		filePath = "./logs/app.log"
	}

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		Logger.SetLevel(level)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	switch strings.ToLower(logFormat) {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	var writers []io.Writer

	if output == "stdout" || output == "both" {
		writers = append(writers, os.Stdout)
	}
	if output == "file" || output == "both" {
		dir := filepath.Dir(filePath)
		_ = os.MkdirAll(dir, 0o755)

		lumberjackLogger := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    1,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
			LocalTime:  true,
		}
		writers = append(writers, lumberjackLogger)
	}

	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	Logger.SetOutput(io.MultiWriter(writers...))
}
