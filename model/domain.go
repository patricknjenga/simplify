package model

import (
	"context"
	"net/http"
)

type Query struct {
	Direction string
	Distinct  string
	Limit     int
	Offset    int
	OrderBy   string
	Search    map[string]any
}

type Index[T any] struct {
	Data  []T
	Query Query
}

type IRepository[T any] interface {
	Create(c context.Context, t *T) error
	Destroy(c context.Context, id int) error
	Index(c context.Context, q Query) (Index[T], error)
	Show(c context.Context, id int) (*T, error)
	Update(c context.Context, id int, t *T) error
}

type IHandler[T any] interface {
	RegisterRoutes()
	Create(w http.ResponseWriter, r *http.Request)
	Destroy(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type IService[T any] interface {
	Create(c context.Context, t *T) error
	Destroy(c context.Context, id int) error
	Index(c context.Context, q Query) (Index[T], error)
	Show(c context.Context, id int) (*T, error)
	Update(c context.Context, id int, t *T) error
}
