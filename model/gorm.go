package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) IRepository[T] {
	return Repository[T]{db}
}

func (r Repository[T]) Create(c context.Context, t *T) error {
	return r.DB.WithContext(c).Create(&t).Error
}

func (r Repository[T]) Destroy(c context.Context, id int) error {
	var t T
	return r.DB.WithContext(c).Delete(&t, id).Error
}

func (r Repository[T]) Index(c context.Context, q Query) (Index[T], error) {
	var (
		t       T
		data    []T
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
	err := builder.Limit(q.Limit).Offset(q.Offset).Find(&data).Error
	return Index[T]{
		Data:  data,
		Query: q,
	}, err
}

func (r Repository[T]) Show(c context.Context, id int) (*T, error) {
	var t T
	return &t, r.DB.WithContext(c).First(&t, id).Error
}

func (r Repository[T]) Update(c context.Context, id int, t *T) error {
	return r.DB.WithContext(c).Updates(&t).Error
}
