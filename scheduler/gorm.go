package scheduler

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type GormScheduler struct {
	DB *gorm.DB
}

func NewGormScheduler(db *gorm.DB) Scheduler {
	return &GormScheduler{db}
}

func (s *GormScheduler) Do(c context.Context, t *Task, f func() error) {
	s.DB.WithContext(c).Create(t)
	defer func() {
		t.Duration = time.Since(t.CreatedAt).Truncate(time.Second).String()
		s.DB.WithContext(c).Updates(t)
	}()
	err := f()
	if err != nil {
		t.Error = err.Error()
	}
}
