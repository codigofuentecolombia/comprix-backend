package utils

func GetLastElement[T any](items []T) T {
	return items[len(items)-1]
}
