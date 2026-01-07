package utils

func UniqueArray[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	var result []T

	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}
