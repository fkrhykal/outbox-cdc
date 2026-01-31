package utils

func Map[T any, P any](elements []T, mapper func(elements T) P) []P {
	result := make([]P, len(elements))
	for i, e := range elements {
		result[i] = mapper(e)
	}
	return result
}
