package insertion

type Interface interface {
	Len() int
	Less(idx int, value any) bool
	Get(idx int) any
	Set(idx int, value any)
}

func Sort(data Interface) {
	n := data.Len()
	for i := 1; i < n; i++ {
		idx, pivot := i, data.Get(i)
		for j := i - 1; j >= 0; j-- {
			if data.Less(j, pivot) {
				break
			}

			data.Set(j+1, data.Get(j))
			idx = j
		}

		data.Set(idx, pivot)
	}
}
