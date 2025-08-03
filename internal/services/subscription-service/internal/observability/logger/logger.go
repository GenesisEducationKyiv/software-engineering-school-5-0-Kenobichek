package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	logInterval 			= time.Second
	initialLogBurst			= 100
	sampleRateAfterBurst	= 100	
)

type sugarManager interface {
	Infof(msg string, keysAndValues ...interface{})
	Errorf(msg string, keysAndValues ...interface{})
	Debugf(msg string, keysAndValues ...interface{})
	Sync() error
}

type zapLogger struct {
	sugar sugarManager
}

func NewZapLogger() (*zapLogger, error) {

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	writer := zapcore.AddSync(os.Stdout)

	baseCore := zapcore.NewCore(
		encoder,
		writer,
		zap.DebugLevel,
	)

	sampledCore := zapcore.NewSamplerWithOptions(
		baseCore,
		logInterval,
		initialLogBurst,
		sampleRateAfterBurst,
	)

	logger := zap.New(sampledCore, zap.AddCaller())
	sugar := logger.Sugar()

	return &zapLogger{sugar: sugar}, nil
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.sugar.Infof(msg, keysAndValues...)
}

func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	z.sugar.Errorf(msg, keysAndValues...)
}

func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.sugar.Debugf(msg, keysAndValues...)
}

func (z *zapLogger) Sync() error {
	return z.sugar.Sync()
}
