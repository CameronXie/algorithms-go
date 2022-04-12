package insection

func Sort(input []int) []int {
	t := make([]int, len(input))
	copy(t, input)

	l := len(t)
	for i := 1; i < l; i++ {
		index, value := i, t[i]
		for j := i - 1; j >= 0; j-- {
			if t[j] < value {
				break
			}

			t[j+1] = t[j]
			index = j
		}

		t[index] = value
	}

	return t
}
