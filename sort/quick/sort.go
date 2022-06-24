package quick

import (
	"sort"
)

func Sort(data sort.Interface) {
	quickSort(data, 0, data.Len()-1)
}

func quickSort(data sort.Interface, min, max int) {
	if min >= max {
		return
	}

	pivot := max
	for left, right := min, max-1; left <= right; {
		for ; data.Less(left, max) && left < right; left++ {
		}

		for ; data.Less(max, right) && left < right; right-- {
		}

		if left != right {
			data.Swap(left, right)
			continue
		}

		if data.Less(max, right) {
			data.Swap(max, right)
			pivot = right
		}

		break
	}

	quickSort(data, 0, pivot-1)
	quickSort(data, pivot+1, max)
}
