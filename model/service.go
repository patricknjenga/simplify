package model

import (
	"context"
)

type IService[T any] interface {
	Create(c context.Context, t T) error
	CreateBatch(c context.Context, t []T, b int) error
	Delete(c context.Context, id int) error
	DeleteAll(c context.Context) error
	DeleteBatch(c context.Context, ids []int) error
	Get(c context.Context, id int) (T, error)
	Query(c context.Context, q Query) (QueryResult[T], error)
	Update(c context.Context, id int, t T) error
}

type Service[T any] struct {
	Repository IRepository[T]
}

func NewModelService[T any](r IRepository[T]) IService[T] {
	return &Service[T]{
		Repository: r,
	}
}

func (s Service[T]) Create(c context.Context, t T) error {
	return s.Repository.Create(c, t)
}

func (s Service[T]) CreateBatch(c context.Context, t []T, b int) error {
	return s.Repository.CreateBatch(c, t, b)
}

func (s Service[T]) Delete(c context.Context, id int) error {
	return s.Repository.Delete(c, id)
}

func (s Service[T]) DeleteAll(c context.Context) error {
	return s.Repository.DeleteAll(c)
}

func (s Service[T]) DeleteBatch(c context.Context, ids []int) error {
	return s.Repository.DeleteBatch(c, ids)
}

func (s Service[T]) Get(c context.Context, id int) (T, error) {
	return s.Repository.Get(c, id)
}

func (s Service[T]) Query(c context.Context, q Query) (QueryResult[T], error) {
	return s.Repository.Query(c, q)
}

func (s Service[T]) Update(c context.Context, id int, t T) error {
	return s.Repository.Update(c, id, t)
}
