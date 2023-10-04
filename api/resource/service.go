package resource

import "context"

type Service[T any] struct {
	Repository IRepository[T]
}

func NewResourceService[T any](r IRepository[T]) IService[T] {
	return &Service[T]{
		Repository: r,
	}
}

func (s Service[T]) Create(c context.Context, t *T) error {
	return s.Repository.Create(c, t)
}

func (s Service[T]) Destroy(c context.Context, id int) error {
	return s.Repository.Destroy(c, id)
}

func (s Service[T]) Index(c context.Context, q Query) (ResourceIndex[T], error) {
	return s.Repository.Index(c, q)
}

func (s Service[T]) Show(c context.Context, id int) (*T, error) {
	return s.Repository.Show(c, id)
}

func (s Service[T]) Update(c context.Context, id int, t *T) error {
	return s.Repository.Update(c, id, t)
}
