// Package registry provides a generic registry interface for storing and retrieving items.
package registry

// Item represents a generic item with an ID and data.
type Item[T any] struct {
	ID   string `json:"id"`
	Data T      `json:"data"`
}

// Filter defines a function type for filtering items based on custom criteria.
type Filter[T any] func(item Item[T]) bool

// Registry defines the interface for a generic registry that can store and retrieve items.
type Registry[T any] interface {
	CreateOrUpdate(item Item[T]) (*Item[T], error)
	Read(id string) (*Item[T], error)
	Delete(id string) (*Item[T], error)
	ListItems(filter Filter[T]) ([]Item[T], error)
	ListIDs() ([]string, error)
}

// AllFilter returns a filter that accepts all items.
func AllFilter[T any]() Filter[T] {
	return func(_ Item[T]) bool {
		return true
	}
}
