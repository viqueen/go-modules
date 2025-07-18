// Package collections provides generic utility functions for working with slices and other collections.
package collections

// Map transforms a slice of items of type T into a slice of items of type R using the provided mapper function.
// If the input slice is nil, it returns nil.
func Map[T any, R any](items []T, mapper func(T) R) []R {
	if items == nil {
		return nil
	}

	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}

	return result
}
