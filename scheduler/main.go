package scheduler

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	Error     error
	Name      string
	StoppedAt *time.Time
	gorm.Model
}

type Scheduler interface {
	Do(c context.Context, t *Task, f func() error)
}
