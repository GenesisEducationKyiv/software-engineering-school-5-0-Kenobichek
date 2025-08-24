package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logInterval          = time.Second
	initialLogBurst      = 100
	sampleRateAfterBurst = 100
)

type zapLogger struct {
	sugar *zap.SugaredLogger
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

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *zapLogger) Infof(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

func (z *zapLogger) With(keysAndValues ...interface{}) *zapLogger {
	return &zapLogger{sugar: z.sugar.With(keysAndValues...)}
}

func (z *zapLogger) Sync() error {
	return z.sugar.Sync()
}
