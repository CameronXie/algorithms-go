package selection

func Sort(s []int) []int {
	t := make([]int, len(s))
	copy(t, s)

	l := len(t)
	for i := 0; i < l; i++ {
		minIndex := i
		for j := i + 1; j < l; j++ {
			if t[minIndex] > t[j] {
				minIndex = j
			}
		}

		t[i], t[minIndex] = t[minIndex], t[i]
	}

	return t
}
