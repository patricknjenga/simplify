package scheduler

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type Task struct {
	Error     string
	Name      string
	StoppedAt sql.NullTime
	gorm.Model
}

type Scheduler interface {
	Do(c context.Context, t *Task, f func() error)
}
