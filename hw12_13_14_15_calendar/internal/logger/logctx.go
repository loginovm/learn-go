package logger

import "context"

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func Get(ctx context.Context) *Logger {
	if ctx == nil {
		panic("nil context passed to logger.Get()")
	}
	logger, _ := ctx.Value(loggerKey{}).(*Logger)
	if logger == nil {
		panic("context without logger passed to logger.Get()")
	}
	return logger
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	Get(ctx).Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	Get(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	Get(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...Field) {
	Get(ctx).Error(msg, fields...)
}
