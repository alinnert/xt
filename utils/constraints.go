package utils

// Ordered is a constraint for generics that includes all types that support <, <=, >, and >=.
type Ordered interface {
	int | float32 | float64 | ~string
}
