package merge

type Interface interface {
	Less(i any) bool
}

func Sort[T Interface](data []T) []T {
	copied := make([]T, len(data))
	copy(copied, data)

	return divide(copied)
}

func divide[T Interface](data []T) []T {
	n := len(data)
	if n <= 1 {
		return data
	}

	pivot := n / 2
	return merge(divide(data[:pivot]), divide(data[pivot:]))
}

func merge[T Interface](left, right []T) []T {
	res := make([]T, 0)
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Less(right[j]) {
			res = append(res, left[i])
			i++
			continue
		}

		res = append(res, right[j])
		j++
	}

	return append(res, append(left[i:], right[j:]...)...)
}
