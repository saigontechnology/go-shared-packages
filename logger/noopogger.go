package logger

import "context"

type NoopLogger struct{}

func NewNoopLogger() Logger {
	return &NoopLogger{}
}

func (log *NoopLogger) Debug(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Info(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Warn(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Error(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) DPanic(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Panic(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Fatal(ctx context.Context, msg string, fields ...Field) {}

func (log *NoopLogger) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) DPanicw(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {}

func (log *NoopLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {}
