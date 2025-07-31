package logger

import (
    "go.uber.org/zap"
)

type zapLogger struct {
	sugared *zap.SugaredLogger
}

func NewZapLogger() (*zapLogger, error) {
    core, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }
    return &zapLogger{sugared: core.Sugar()}, nil
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
    z.sugared.Infof(msg, keysAndValues...)
}

func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
    z.sugared.Errorf(msg, keysAndValues...)
}

func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
    z.sugared.Debugf(msg, keysAndValues...)
}

func (z *zapLogger) With(keysAndValues ...interface{}) *zapLogger {
    return &zapLogger{sugared: z.sugared.With(keysAndValues...)}
}

func (z *zapLogger) Sync() error {
    return z.sugared.Sync()
}
