package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type IResource interface {
	Migrate(db *gorm.DB) error
	RegisterRoutes()
	Schema() (string, map[string]string)
}

type Resource[T any] struct {
	Handler    IHandler[T]
	Name       string
	Repository IRepository[T]
	Service    IService[T]
}

func New[T any](rt *mux.Router, db *gorm.DB) *Resource[T] {
	var (
		t T
		n = reflect.TypeOf(t).Name()
		r = NewGormRepository[T](db)
		s = NewResourceService[T](r)
		h = NewResourceHandler[T](rt.PathPrefix(fmt.Sprintf("/%s", n)).Subrouter(), s)
	)
	return &Resource[T]{h, n, r, s}
}

func (r *Resource[T]) RegisterRoutes() {
	r.Handler.RegisterRoutes()
}

func (r *Resource[T]) Migrate(db *gorm.DB) error {
	var t T
	return db.AutoMigrate(&t)
}

func (r *Resource[T]) Schema() (string, map[string]string) {
	var (
		fields = map[string]string{}
		t      T
	)
	for i := 0; i < reflect.TypeOf(t).NumField(); i++ {
		field := reflect.TypeOf(t).Field(i)
		fields[field.Name] = field.Type.String()
	}
	return reflect.TypeOf(t).Name(), fields
}

func NewArr(r *mux.Router, db *gorm.DB, rs ...IResource) error {
	var res = map[string]map[string]string{}
	for _, v := range rs {
		v.RegisterRoutes()
		err := v.Migrate(db)
		if err != nil {
			return err
		}
		name, fields := v.Schema()
		res[name] = fields
	}
	r.HandleFunc("/Schema", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return nil
}
