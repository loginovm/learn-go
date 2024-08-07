package logger

import (
	"fmt"
	"io"
	"time"
)

type Logger struct {
	h Handler
}

func New(level Level, w io.Writer) *Logger {
	return &Logger{
		h: NewTextHandler(level, w),
	}
}

func (l Logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l Logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

func (l Logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

func (l Logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

func (l Logger) log(level Level, msg string, fields ...Field) {
	r := Record{
		Lvl:    level,
		Msg:    msg,
		Fields: fields,
	}
	if r.Time.IsZero() {
		r.Time = time.Now()
	}
	if err := l.h.Handle(r); err != nil {
		fmt.Println(err)
	}
}
