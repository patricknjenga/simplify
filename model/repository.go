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
	Search    map[string]any
}

type QueryResult[T any] struct {
	Count int64
	Data  []T
	Query Query
}

type IRepository[T any] interface {
	Create(c context.Context, t T) error
	CreateBatch(c context.Context, t []T, b int) error
	Delete(c context.Context, id int) error
	DeleteAll(c context.Context) error
	DeleteBatch(c context.Context, ids []int) error
	Get(c context.Context, id int) (T, error)
	GetAll(c context.Context) ([]T, error)
	Query(c context.Context, q Query) (QueryResult[T], error)
	Update(c context.Context, id int, t T) error
}

type Repository[T any] struct {
	DB *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) IRepository[T] {
	return Repository[T]{db}
}

func (r Repository[T]) Create(c context.Context, t T) error {
	return r.DB.WithContext(c).Create(&t).Error
}

func (r Repository[T]) CreateBatch(c context.Context, t []T, b int) error {
	if len(t) == 0 {
		return nil
	}
	return r.DB.WithContext(c).CreateInBatches(&t, b).Error
}

func (r Repository[T]) Delete(c context.Context, id int) error {
	var t T
	return r.DB.WithContext(c).Delete(&t, id).Error
}

func (r Repository[T]) DeleteAll(c context.Context) error {
	var t T
	return r.DB.WithContext(c).Where("1 = 1").Delete(&t).Error
}

func (r Repository[T]) DeleteBatch(c context.Context, ids []int) error {
	var t T
	if len(ids) == 0 {
		return nil
	}
	return r.DB.WithContext(c).Delete(&t, ids).Error
}

func (r Repository[T]) Get(c context.Context, id int) (T, error) {
	var t T
	return t, r.DB.WithContext(c).First(&t, id).Error
}

func (r Repository[T]) GetAll(c context.Context) ([]T, error) {
	var t []T
	return t, r.DB.WithContext(c).Find(&t).Error
}

func (r Repository[T]) Query(c context.Context, q Query) (QueryResult[T], error) {
	var (
		t   T
		res = QueryResult[T]{
			Count: 0,
			Data:  []T{},
			Query: q,
		}
		builder = r.DB.Model(&t)
	)
	if q.Distinct != "" {
		builder = builder.Distinct(q.Distinct)
	}
	if q.OrderBy != "" && q.Direction != "" {
		builder = builder.Order(fmt.Sprintf("%s %s", q.OrderBy, q.Direction))
	}
	if q.Search != nil {
		builder = builder.Where(q.Search)
	}
	if q.Count {
		builder.Count(&res.Count)
	}
	return res, builder.Limit(q.Limit).Offset(q.Offset).Find(&res.Data).Error
}

func (r Repository[T]) Update(c context.Context, id int, t T) error {
	return r.DB.WithContext(c).Updates(&t).Error
}
