package scheduler

import (
	"context"
	"time"
)

type Task struct {
	ID        uint
	Name      string
	StartedAt time.Time
	StoppedAt time.Time
}

type Scheduler interface {
	Do(c context.Context, t *Task, f func() error)
}
