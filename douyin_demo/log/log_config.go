package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 使用方法
//
//	log.SetLog().Info("正常运行Info")
//	log.SetLog().Warn("这是一条警告")
//	log.SetLog().Error("出现重大bug")
var log *zap.Logger

func Init() {
	var coreArr []zapcore.Core

	// 获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 控制台输出
	consoleCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.Lock(os.Stdout)), lowPriority)

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.Lock(os.Stdout)), highPriority)

	coreArr = append(coreArr, consoleCore) // 将控制台输出的core添加到数组
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	log = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())

}

func SetLog() *zap.Logger {
	return log
}
