package model

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type IService[T interface{}] interface {
	Delete(c context.Context, t []T) error
	Get(c context.Context, q Query) (int64, []T, error)
	Post(c context.Context, t []T) error
	Put(c context.Context, t []T) error
}

type Service[T interface{}] struct {
	Repository IRepository[T]
	Validate   *validator.Validate
}

func NewModelService[T interface{}](r IRepository[T]) IService[T] {
	return &Service[T]{r, validator.New()}
}

func (s Service[T]) Delete(c context.Context, t []T) error {
	for _, v := range t {
		err := s.Validate.Struct(v)
		if err != nil {
			return err
		}
	}
	return s.Repository.Delete(c, t)
}

func (s Service[T]) Get(c context.Context, q Query) (int64, []T, error) {
	return s.Repository.Get(c, q)
}

func (s Service[T]) Post(c context.Context, t []T) error {
	for _, v := range t {
		err := s.Validate.Struct(v)
		if err != nil {
			return err
		}
	}
	return s.Repository.Post(c, t)
}

func (s Service[T]) Put(c context.Context, t []T) error {
	for _, v := range t {
		err := s.Validate.Struct(v)
		if err != nil {
			return err
		}
	}
	return s.Repository.Put(c, t)
}
