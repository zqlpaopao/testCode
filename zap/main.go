package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func main(){
	//interface()类型的数据
	//Sugar()

	//请类型字段
	//Logger()

	//
	//ZdyMessageType()

	//写入文件
	LogWrite()
}


/**
* interface{}类型的转化
**/
func Sugar(){
	logger, _ := zap.NewProduction()
	// 默认 logger 不缓冲。
	// 但由于底层 api 允许缓冲，所以在进程退出之前调用 Sync 是一个好习惯。
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infof("message: %s", "https://github.com")

}

/**
* 对性能和类型安全要求严格的情况下，可以使用 Logger ，
* 它甚至比前者SugaredLogger更快，内存分配次数也更少，
* 但它仅支持强类型的结构化日志记录。
*/
func Logger(){
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		// 强类型字段
		zap.String("url", "https://github.com"),
		zap.Int("int", 3),
		zap.Duration("total", time.Second),
	)

	//二进制 zap.Binary
	//time zap.Time
}

func ZdyMessageType(){
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		atom,
	)
	logger := zap.New(core, zap.AddCaller(), zap.Development())
	defer logger.Sync()

	// 配置 zap 包的全局变量
	zap.ReplaceGlobals(logger)

	// 运行时安全地更改 logger 日记级别
	atom.SetLevel(zap.InfoLevel)
	sugar := logger.Sugar()
	// 问题 1: debug 级别的日志打印到控制台了吗?
	sugar.Debug("debug")
	sugar.Info("info")
	sugar.Warn("warn")
	sugar.DPanic("dPanic")
	// 问题 2: 最后的 error 会打印到控制台吗?
	sugar.Error("error")
}

func NewCustomEncoderConfig() zapcore.EncoderConfig{
		return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func LogWrite(){
	hook := lumberjack.Logger{
		Filename:   "./logs/spikeProxy1.log", // 日志文件路径
		MaxSize:    128,                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                       // 日志文件最多保存多少个备份
		MaxAge:     7,                        // 文件最多保存多少天
		Compress:   true,                     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	logger.Info("log 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
}