package logger

import "time"

type Record struct {
	Lvl    Level
	Time   time.Time
	Msg    string
	Fields []Field
}

type Field struct {
	Key   string
	Value string
}
