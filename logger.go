package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"telegram-notice/global"
)

// sugarLogger 是一个全局变量，用于存储配置好的SugaredLogger实例。

// InitLogger 初始化日志记录器
func InitLogger() {
	// 获取日志写入器
	writeSyncer := getLogWriter()
	// 获取日志编码器
	encoder := getEncoder()
	// 创建日志核心处理器
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 创建一个结构化的日志记录器
	logger := zap.New(core, zap.AddCaller())
	// 包装日志记录器为一个更易用的SugaredLogger
	global.LogZap = logger.Sugar()
	global.Log, _ = zap.NewProduction()
}

// getEncoder 获取日志编码器，用于定义日志的格式
func getEncoder() zapcore.Encoder {
	// 创建一个新的日志编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间的编码方式为ISO8601格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置日志级别的编码方式为大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 创建并返回JSON格式的编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 获取日志写入器，并同时写入到文件和控制台
func getLogWriter() zapcore.WriteSyncer {
	// 打开或创建日志文件，如果文件已存在则追加内容
	file, err := os.OpenFile("config/telebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 如果打开文件失败，记录错误并退出程序
		global.LogZap.Fatal("Failed to open log file: ", err)
	}
	// 创建文件的写入同步器
	fileWriteSyncer := zapcore.AddSync(file)
	// 创建控制台的写入同步器
	consoleWriteSyncer := zapcore.Lock(os.Stderr) // 确保线程安全
	// 使用MultiWriteSyncer同时写入到文件和控制台
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
	return multiWriteSyncer
}

// simpleHttpGet 发起一个简单的HTTP GET请求，并记录相关日志
func simpleHttpGet(url string) {
	// 记录调试信息，尝试发起GET请求
	global.LogZap.Debugf("Trying to hit GET request for %s", url)
	// 发起GET请求
	resp, err := http.Get(url)
	if err != nil {
		// 如果请求失败，记录错误信息
		global.LogZap.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		// 如果请求成功，记录状态码和URL
		global.LogZap.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		// 关闭响应体
		resp.Body.Close()
	}
	global.LogZap.Infoln("我发送完毕啦哈哈哈", "你好呀")
}
