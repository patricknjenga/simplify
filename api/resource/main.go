package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func New[T any](router *mux.Router, db *gorm.DB) {
	var (
		t          T
		name       = reflect.TypeOf(t).Name()
		repository = NewGormRepository[T](db)
		service    = NewResourceService[T](repository)
		handler    = NewResourceHandler[T](router.PathPrefix(fmt.Sprintf("/%s", name)).Subrouter(), service)
	)
	handler.Register()
}

func Schema(router *mux.Router, resources ...any) {
	router.HandleFunc("/Schema", func(w http.ResponseWriter, r *http.Request) {
		var res = map[string]map[string]string{}
		for _, v := range resources {
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
}
