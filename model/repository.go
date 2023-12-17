package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Query struct {
	Count     bool
	Direction string
	Distinct  string
	Limit     int
	Offset    int
	OrderBy   string
	Select    []string
	Where     map[string]interface{}
}

type IRepository[T interface{}] interface {
	Delete(c context.Context, t []T) error
	Get(c context.Context, q Query) (int64, []T, error)
	Post(c context.Context, t []T) error
	Put(c context.Context, t []T) error
}

type Repository[T interface{}] struct {
	Batch int
	DB    *gorm.DB
}

func NewGormRepository[T interface{}](b int, db *gorm.DB) IRepository[T] {
	return Repository[T]{b, db}
}

func (r Repository[T]) Delete(c context.Context, t []T) error {
	if len(t) == 0 {
		return nil
	}
	return r.DB.WithContext(c).Model(&t).Delete(&t).Error
}

func (r Repository[T]) Get(c context.Context, q Query) (int64, []T, error) {
	var (
		t       T
		builder = r.DB.WithContext(c).Model(&t)
		count   int64
		res     []T
	)

	if q.Where != nil {
		builder = builder.Where(q.Where)
	}
	if q.Count {
		builder.Count(&count)
	}
	if q.Limit != 0 {
		builder = builder.Limit(q.Limit)
	}
	if q.Offset != 0 {
		builder = builder.Offset(q.Offset)
	}
	if q.Select != nil {
		builder = builder.Select(q.Select)
	}
	if q.Distinct != "" {
		builder = builder.Distinct(q.Distinct)
	}
	if q.OrderBy != "" && q.Direction != "" {
		builder = builder.Order(fmt.Sprintf("%s %s", q.OrderBy, q.Direction))
	}
	return count, res, builder.Find(&res).Error
}

func (r Repository[T]) Post(c context.Context, t []T) error {
	if len(t) == 0 {
		return nil
	}
	return r.DB.WithContext(c).CreateInBatches(t, r.Batch).Error
}

func (r Repository[T]) Put(c context.Context, t []T) error {
	for _, v := range t {
		err := r.DB.WithContext(c).Updates(&v).Error
		if err != nil {
			return err
		}
	}
	return nil
}
