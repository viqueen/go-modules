package collections

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
