package model

import (
	"reflect"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Model[T interface{}] struct {
	Handler    IHandler
	Name       string
	Repository IRepository[T]
	Service    IService[T]
}

func New[T interface{}](rt *mux.Router, db *gorm.DB) *Model[T] {
	var (
		t T
		n = reflect.TypeOf(t).Name()
		r = NewGormRepository[T](1000, db)
		s = NewModelService[T](r)
		h = NewModelHandler[T](n, rt, s)
	)
	return &Model[T]{h, n, r, s}
}
