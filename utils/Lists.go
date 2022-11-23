package utils

// Map applies a function to all elements in a list and returns a mutated list.
func Map[T, V any](items []T, cb func(T) V) []V {
	result := make([]V, len(items))
	for i, item := range items {
		result[i] = cb(item)
	}
	return result
}

// Reverse returns a reversed version of a list.
func Reverse[T any](items []T) []T {
	result := make([]T, len(items))
	for i, item := range items {
		result[len(result)-i-1] = item
	}
	return result
}
