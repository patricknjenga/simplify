package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Resource[T any] struct{}

func New[T any](router *mux.Router, db *gorm.DB) *Resource[T] {
	return &Resource[T]{}
}

func (r Resource[T]) Init(db *gorm.DB, rt *mux.Router) error {
	var (
		t          T
		name       = reflect.TypeOf(t).Name()
		repository = NewGormRepository[T](db)
		service    = NewResourceService[T](repository)
		handler    = NewResourceHandler[T](rt.PathPrefix(fmt.Sprintf("/%s", name)).Subrouter(), service)
	)
	handler.Register()
	return db.AutoMigrate(&t)
}

func NewArr(r *mux.Router, db *gorm.DB, rs ...Resource[any]) error {
	for _, v := range rs {
		err := v.Init(db, r)
		if err != nil {
			return err
		}
	}
	r.HandleFunc("/Schema", func(w http.ResponseWriter, r *http.Request) {
		var res = map[string]map[string]string{}
		for _, v := range rs {
			var (
				resource = reflect.TypeOf(v)
				fields   = map[string]string{}
			)
			for i := 0; i < resource.NumField(); i++ {
				field := resource.Field(i)
				fields[field.Name] = field.Type.String()
			}
			res[resource.Name()] = fields
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return nil
}
