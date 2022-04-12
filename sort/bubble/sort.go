package bubble

func Sort(s []int) []int {
	t := make([]int, len(s))
	copy(t, s)

	for i := len(t); i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if t[j] > t[j+1] {
				t[j], t[j+1] = t[j+1], t[j]
			}
		}
	}

	return t
}
