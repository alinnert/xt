package xsd

// GetAncestorPathSegments returns multiple paths, one for every ancestor of the given element path.
func GetAncestorPathSegments(segments []string) [][]string {
	var result [][]string
	for i := range segments {
		result = append(result, segments[:(i+1)])
	}
	return result[1:]
}

// GetParent returns the parent element path for a given element path.
func GetParent(segments []string) []string {
	return segments[:len(segments)-1]
}
