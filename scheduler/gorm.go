package scheduler

import (
	"context"
	"time"

	"github.com/patricknjenga/simplify/logger"
	"gorm.io/gorm"
)

type GormScheduler struct {
	DB     *gorm.DB
	Logger logger.Logger
}

func NewGormScheduler(db *gorm.DB, logger logger.Logger) Scheduler {
	return &GormScheduler{db, logger}
}

func (s *GormScheduler) Do(c context.Context, t Task, f func() error) {
	t.StartedAt = time.Now()
	err := s.DB.WithContext(c).Create(&t).Error
	if err != nil {
		s.Logger.Error(err)
		return
	}
	defer func() {
		t.StoppedAt = time.Now()
		err = s.DB.WithContext(c).Updates(&t).Error
		if err != nil {
			s.Logger.Error(err)
		}
	}()
	err = f()
	if err != nil {
		s.Logger.Error(err)
	}
}
