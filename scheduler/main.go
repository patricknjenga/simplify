package scheduler

import (
	"context"

	"gorm.io/gorm"
)

type Task struct {
	Duration string
	Error    string
	Name     string
	gorm.Model
}

type Scheduler interface {
	Do(c context.Context, t *Task, f func() error)
}
