package util

func ArrayReduce[T any, L any](data []T, initialValue L, reduce func(a L, b T) L) L {
	result := initialValue

	for _, d := range data {
		result = reduce(result, d)
	}

	return result
}
