package logx

import (
	"monitor-data/common/utils"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *zap.SugaredLogger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	logFilename := path.Join(utils.CurrentProjectRootPath(), "logs/logx_zap.log")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
