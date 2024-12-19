package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFieldAppRole   = "app_role"
	logFieldMessage   = "message"
	logFieldTimestamp = "timestamp"
	logCallerSkip     = 1
)

type ZapLogger struct {
	cfg *loggerConfig
	zl  *zap.Logger
}

func NewZapLogger(cfg *loggerConfig) (Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = logFieldMessage
	encoderConfig.TimeKey = logFieldTimestamp
	if cfg.LogWithTimestamp {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encoderConfig.EncodeTime = func(_ time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("<disabled>")
		}
	}

	var samplingConfig *zap.SamplingConfig
	samplingConfig = nil
	if cfg.EnableSampling {
		samplingConfig = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}
	development := false
	if cfg.IsDevEnv() {
		development = true
	}
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(cfg.Level)),
		Development:      development,
		Sampling:         samplingConfig,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			logFieldAppRole: cfg.AppRole,
		},
	}

	zl, err := config.Build(zap.AddCallerSkip(logCallerSkip))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{cfg: cfg, zl: zl}, nil
}

func (log *ZapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Debug(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Debug(msg, fields...)
	}
}

func (log *ZapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Info(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Info(msg, fields...)
	}
}

func (log *ZapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Warn(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Warn(msg, fields...)
	}
}

func (log *ZapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Error(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Error(msg, fields...)
	}
}

func (log *ZapLogger) DPanic(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.DPanic(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.DPanic(msg, fields...)
	}
}

func (log *ZapLogger) Panic(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Panic(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Panic(msg, fields...)
	}
}

func (log *ZapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Fatal(msg, append(tracing.ToFieldSlice(), fields...)...)
		}
	} else {
		log.zl.Fatal(msg, fields...)
	}
}

func (log *ZapLogger) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Debugw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Debugw(msg, keysAndValues...)
	}
}

func (log *ZapLogger) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Infow(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Infow(msg, keysAndValues...)
	}
}

func (log *ZapLogger) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Warnw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Warnw(msg, keysAndValues...)
	}
}

func (log *ZapLogger) Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Errorw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Errorw(msg, keysAndValues...)
	}
}

func (log *ZapLogger) DPanicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().DPanicw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().DPanicw(msg, keysAndValues...)
	}
}

func (log *ZapLogger) Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Panicw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Panicw(msg, keysAndValues...)
	}
}

func (log *ZapLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if log.cfg.EnableTracing {
		tracing := extractTracingDataFromContext(ctx)
		if tracing != nil {
			log.zl.Sugar().Fatalw(msg, append(tracing.ToInterfaceSlice(), keysAndValues...)...)
		}
	} else {
		log.zl.Sugar().Fatalw(msg, keysAndValues...)
	}
}
