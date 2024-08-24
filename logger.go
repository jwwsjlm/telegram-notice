package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"telegram-notice/global"
)

// InitLogger 初始化日志记录器
func InitLogger() {
	// 获取日志写入器
	writeSyncer := getLogWriter()
	// 获取日志编码器
	encoder := getEncoder()
	// 创建日志核心处理器，设置日志级别为 Debug
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 创建一个结构化的日志记录器，添加调用者信息
	logger := zap.New(core, zap.AddCaller())
	// 包装日志记录器为一个更易用的 SugaredLogger
	global.LogZap = logger.Sugar()
	// 创建一个生产环境的日志记录器
	global.Log, _ = zap.NewProduction()
}

// getEncoder 获取日志编码器，用于定义日志的格式
func getEncoder() zapcore.Encoder {
	// 创建一个新的日志编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间的编码方式为 ISO8601 格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置日志级别的编码方式为大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 创建并返回 JSON 格式的编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 获取日志写入器，并同时写入到文件和控制台
func getLogWriter() zapcore.WriteSyncer {
	// 打开或创建日志文件，如果文件已存在则追加内容
	file, err := os.OpenFile("config/telebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 如果打开文件失败，记录错误并退出程序
		// 这里不能使用 global.LogZap，因为它可能尚未初始化
		zap.L().Fatal("Failed to open log file", zap.Error(err))
	}
	// 创建文件的写入同步器
	fileWriteSyncer := zapcore.AddSync(file)
	// 创建控制台的写入同步器，确保线程安全
	consoleWriteSyncer := zapcore.Lock(os.Stderr)
	// 使用 MultiWriteSyncer 同时写入到文件和控制台
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
	return multiWriteSyncer
}
