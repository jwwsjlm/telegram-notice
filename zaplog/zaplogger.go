package zaplogger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

func ParseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
func InitLogger(logLevel zapcore.Level) {
	//获取编码器
	encoder := getEncoder()

	// 日志级别

	infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // logLevel及以上级别
		return lev >= logLevel && lev < zap.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	//info文件WriteSyncer
	infoFileWriteSyncer := getInfoWriterSyncer()
	//error文件WriteSyncer
	errorFileWriteSyncer := getErrorWriterSyncer()

	//生成core
	//multiWriteSyncer := zapcore.NewMultiWriteSyncer(writerSyncer, zapcore.AddSync(os.Stdout)) //AddSync将io.Writer转换成WriteSyncer的类型
	//同时输出到控制台 和 指定的日志文件中
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoLevel)
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errorLevel)

	//将infocore 和 errcore 加入core切片

	var coreArr []zapcore.Core
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	//生成Logger
	Logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCallerSkip(0), zap.AddStacktrace(zap.WarnLevel)) //zap.AddCaller() 显示文件名 和 行号
	SugaredLogger = Logger.Sugar()
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var level string
	switch l {
	case zapcore.DebugLevel:
		level = "[DEBUG]"
	case zapcore.InfoLevel:
		level = "[INFO]"
	case zapcore.WarnLevel:
		level = "[WARN]"
	case zapcore.ErrorLevel:
		level = "[ERROR]"
	case zapcore.DPanicLevel:
		level = "[DPANIC]"
	case zapcore.PanicLevel:
		level = "[PANIC]"
	case zapcore.FatalLevel:
		level = "[FATAL]"
	default:
		level = fmt.Sprintf("[LEVEL(%d)]", l)
	}
	enc.AppendString(level)
}

func shortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", caller.TrimmedPath()))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// TimeKey 是日志中时间字段的键名，默认是 "ts"。
		TimeKey: "time",
		// LevelKey 是日志级别字段的键名，默认是 "level"。
		LevelKey: "level",
		// NameKey 是日志名称字段的键名，默认是 "logger"。
		NameKey: "logger",
		// CallerKey 是调用者信息字段的键名，这里设置为 "caller"。
		CallerKey: "caller",
		// FunctionKey 指定函数名字段的键名，这里使用 zapcore.OmitKey 表示不记录函数名。
		FunctionKey: "func",
		// MessageKey 是日志消息字段的键名，默认是 "msg"。
		MessageKey: "msg",
		// StacktraceKey 是堆栈跟踪字段的键名，默认是 "stacktrace"。
		StacktraceKey: "stacktrace",
		// LineEnding 指定日志条目的行结束符，默认是 zapcore.DefaultLineEnding。
		LineEnding: zapcore.DefaultLineEnding,
		// EncodeLevel 指定日志级别如何被编码，默认是 zapcore.LowercaseLevelEncoder，表示使用小写字母。
		EncodeLevel: levelEncoder,
		// EncodeTime 指定时间如何被编码，默认是 zapcore.EpochTimeEncoder，表示使用 Unix 时间戳。
		EncodeTime: timeEncoder,
		// EncodeDuration 指定持续时间如何被编码，默认是 zapcore.SecondsDurationEncoder，表示使用秒。
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller 指定调用者信息如何被编码，默认是 zapcore.ShortCallerEncoder，表示使用简短的文件名和行号。
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// core 三个参数之  Encoder 获取编码器
func getEncoder() zapcore.Encoder {
	//自定义编码配置,下方NewJSONEncoder输出如下的日志格式
	//{"L":"[INFO]","T":"2022-09-16 14:24:59.552","C":"[prototest/main.go:113]","M":"name = xiaoli, age = 18"}
	return zapcore.NewJSONEncoder(NewEncoderConfig())

	//下方NewConsoleEncoder输出如下的日志格式
	//2022-09-16 14:26:02.933 [INFO]  [prototest/main.go:113] name = xiaoli, age = 18
	//return zapcore.NewConsoleEncoder(NewEncoderConfig())
}

// core 三个参数之  日志输出路径
func getInfoWriterSyncer() zapcore.WriteSyncer {
	//file, _ := os.Create("./server/zaplog/log.log")
	//或者将上面的NewMultiWriteSyncer放到这里来，进行返回
	//return zapcore.AddSync(file)

	//引入第三方库 Lumberjack 加入日志切割功能
	infoLumberIO := &lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    100, // megabytes
		MaxBackups: 100,
		MaxAge:     28,    // days
		Compress:   false, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(infoLumberIO)
}

func getErrorWriterSyncer() zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    100, // megabytes
		MaxBackups: 100,
		MaxAge:     28,    // days
		Compress:   false, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)
}

// Debugf 不再封装使用 - 在显示调用者文件名的时候，会全部显示调用者为logger/zaplogger.go - 所以如果要显示调用者文件名和行号，这里的封装就不合适了
// 直接使用 logger.Logger.Info(xxx)
// 或者   logger.SugaredLogger.Infof("xxx%s", name)
// logs.Debug(...) 再封装
func Debugf(format string, v ...interface{}) {
	Logger.Sugar().Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	Logger.Sugar().Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	Logger.Sugar().Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Logger.Sugar().Errorf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	Logger.Sugar().Panicf(format, v...)
}

// logs.Debug(...) 再封装
func Debug(format string, fileds ...zapcore.Field) {
	Logger.Debug(format, fileds...)
}

func Info(format string, fileds ...zapcore.Field) {
	Logger.Info(format, fileds...)
}

func Warn(format string, fileds ...zapcore.Field) {
	Logger.Warn(format, fileds...)
}

func Error(format string, fileds ...zapcore.Field) {
	Logger.Error(format, fileds...)
}

func Panic(format string, fileds ...zapcore.Field) {
	Logger.Panic(format, fileds...)
}
