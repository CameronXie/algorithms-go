package merge

func Sort(s []int) []int {
	t := make([]int, len(s))
	copy(t, s)

	return sort(t)
}

func sort(s []int) []int {
	l := len(s)
	if l <= 1 {
		return s
	}

	m := l / 2
	return merge(sort(s[:m]), sort(s[m:]))
}

func merge(left, right []int) []int {
	res := make([]int, 0)
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] > right[j] {
			res = append(res, right[j])
			j++
			continue
		}

		res = append(res, left[i])
		i++
	}

	return append(res, append(left[i:], right[j:]...)...)
}
