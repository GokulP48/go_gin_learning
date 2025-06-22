package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	zap.SugaredLogger
}

func NewZapLogger(out io.Writer, level string) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "C"

	syslogCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(out),
		func(level string) zapcore.Level {
			switch level {
			case "info":
				return zap.InfoLevel
			case "error":
				return zap.ErrorLevel
			case "debug":
				return zap.DebugLevel
			default:
				return zapcore.InfoLevel
			}
		}(level),
	)

	log := zap.New(syslogCore, zap.AddCallerSkip(1), zap.AddCaller())

	return &ZapLogger{SugaredLogger: *log.Sugar()}
}
