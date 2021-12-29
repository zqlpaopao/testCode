/**
 * @Author: zhangsan
 * @Description:
 * @File:  zaplog
 * @Version: 1.0.0
 * @Date: 2021/3/2 上午10:05
 */

package src

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	logger *zap.Logger
	cfg  zap.Config
	Ti string
)
func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func FormatLog(args []interface{}) *zap.Logger {
	log := logger.With(ToJsonData(args))
	return log
}

func Debug(msg string, args ...interface{}) {
	resetLogger()
	FormatLog(args).Sugar().Debugf(msg)
}

func Warn(msg string, args ...interface{}) {
	FormatLog(args).Sugar().Warnf(msg)
}
func Info(msg string, args ...interface{}) {
	FormatLog(args).Sugar().Infof(msg)
}

func Error(msg string, args ...interface{}) {
	FormatLog(args).Sugar().Errorf(msg)
}

func ToJsonData(args []interface{}) zap.Field {
	det := make([]string, 0)
	if len(args) > 0 {
		for _, v := range args {
			det = append(det, fmt.Sprintf("%+v", v))
		}
	}
	z := zap.Any("params", det)
	return z
}

func resetLogger(){
	var ti  = time.Now().Format("2006-01-02 15:04")
	if Ti !=  ti{
		Ti = ti
		logger,_  = MakeCfg(ti).Build()
	}
}

func MakeCfg(ti string)zap.Config{
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "string-byte",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     formatEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"./zap.log"+ti},
		//ErrorOutputPaths: []string{"./zap.log"},
		InitialFields: map[string]interface{}{
			"app": "fire-timer",
		},
	}
}

func InitZapLog() {
	Ti = time.Now().Format("2006-01-02 15:04")
	var err error
	defer func() {
		err = logger.Sync()
	}()
	logger, err = MakeCfg(time.Now().Format("2006-01-02 15:04")).Build()
	if err != nil {
		panic("log init fail:" + err.Error())
	}
}



