package selection

import "sort"

func Sort(data sort.Interface) {
	n := data.Len()

	for i := 0; i < n; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if !data.Less(minIdx, j) {
				data.Swap(minIdx, j)
			}
		}

		data.Swap(i, minIdx)
	}
}
