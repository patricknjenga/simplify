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
	init(db *gorm.DB, rt *mux.Router) error
	schema() (string, map[string]string)
}

type Resource[T any] struct{}

func New[T any]() *Resource[T] {
	return &Resource[T]{}
}

func (r *Resource[T]) init(db *gorm.DB, rt *mux.Router) error {
	var (
		t          T
		repository = NewGormRepository[T](db)
		service    = NewResourceService[T](repository)
		handler    = NewResourceHandler[T](rt.PathPrefix(fmt.Sprintf("/%s", reflect.TypeOf(t).Name())).Subrouter(), service)
	)

	handler.Register()
	return db.AutoMigrate(&t)
}

func (r *Resource[T]) schema() (string, map[string]string) {
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
		err := v.init(db, r)
		if err != nil {
			return err
		}
		name, fields := v.schema()
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
