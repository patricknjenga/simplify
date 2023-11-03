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
		t.StoppedAt.Time = time.Now()
		s.DB.WithContext(c).Updates(t)
	}()
	t.Error = f().Error()
}
