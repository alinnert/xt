package utils

// Min returns the smaller number. It's a generic version of math.Min
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
