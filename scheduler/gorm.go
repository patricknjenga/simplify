package scheduler

import (
	"context"
	"time"

	"github.com/patricknjenga/simplify/errors"
	"gorm.io/gorm"
)

type GormScheduler struct {
	DB *gorm.DB
}

func NewGormScheduler(db *gorm.DB) Scheduler {
	return &GormScheduler{db}
}

func (s *GormScheduler) Do(c context.Context, t *Task, f func() error) {
	t.StartedAt = time.Now()
	err := s.DB.WithContext(c).Create(t).Error
	if err != nil {
		errors.Log(err)
		return
	}
	defer func() {
		t.StoppedAt = time.Now()
		err = s.DB.WithContext(c).Updates(t).Error
		if err != nil {
			errors.Log(err)
		}
	}()
	err = f()
	if err != nil {
		errors.Log(err)
	}
}
