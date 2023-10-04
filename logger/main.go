package logger

import "time"

type Error struct {
	Service string
	Message string
	Trace   string
	Time    time.Time
}

type Logger interface {
	Error(e error)
}
