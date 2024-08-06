package logger

import "fmt"

type Level int

const (
	LevelDebug Level = iota + 1
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return fmt.Sprintf("%d", int(l))
	}
}
