package registry

type Item[T any] struct {
	ID   string `json:"id"`
	Data T      `json:"data"`
}

type Filter[T any] func(item Item[T]) bool

type Registry[T any] interface {
	CreateOrUpdate(item Item[T]) (*Item[T], error)
	Read(id string) (*Item[T], error)
	Delete(id string) (*Item[T], error)
	ListItems(filter Filter[T]) ([]Item[T], error)
	ListIDs() ([]string, error)
}

func AllFilter[T any]() Filter[T] {
	return func(item Item[T]) bool {
		return true
	}
}
