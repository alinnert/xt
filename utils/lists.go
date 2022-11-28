package utils

// MapContext is the type of the parameter that's being passed to Map's callback function.
// It contains all relevant data for a map iteration.
type MapContext[T any] struct {
	Items               *[]T
	Item                *T
	Index               int
	FirstItem, LastItem bool
}

// Map applies a function to all elements in a list and returns a mutated list.
func Map[T, V any](items *[]T, cb func(*MapContext[T]) V) []V {
	result := make([]V, len(*items))
	for i, item := range *items {
		ctx := &MapContext[T]{
			Items:     items,
			Item:      &item,
			Index:     i,
			FirstItem: i == 0,
			LastItem:  i == len(*items)-1,
		}
		result[i] = cb(ctx)
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
