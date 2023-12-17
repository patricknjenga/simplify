package model

import (
	"fmt"
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

func (m Model[T]) New(rt *mux.Router, db *gorm.DB) *Model[T] {
	var (
		t T
		n = reflect.TypeOf(t).Name()
		r = NewGormRepository[T](1000, db)
		s = NewModelService[T](r)
		h = NewModelHandler[T](rt.PathPrefix(fmt.Sprintf("/%s", n)).Subrouter(), s)
	)
	return &Model[T]{h, n, r, s}
}
