package quick

func Sort(s []int) []int {
	t := make([]int, len(s))
	copy(t, s)

	sort(t)
	return t
}

func sort(s []int) {
	if len(s) <= 1 {
		return
	}

	index := partition(s)

	sort(s[:index])
	sort(s[index:])
}

func partition(s []int) int {
	index := len(s) - 1
	pivot := s[index]

	for left, right := 0, index-1; left < right; {
		for s[left] < pivot && left < right {
			left++
		}

		for s[right] >= pivot && right > left {
			right--
		}

		if left != right {
			s[left], s[right] = s[right], s[left]
			continue
		}

		if s[right] > pivot {
			index = right
		}
	}

	s[index], s[len(s)-1] = s[len(s)-1], s[index]
	return index
}
