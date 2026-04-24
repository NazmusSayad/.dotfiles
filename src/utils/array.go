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

func SortArray(slice []string) []string {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[j] < slice[i] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}

	return slice
}
