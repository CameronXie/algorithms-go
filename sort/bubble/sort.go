package bubble

import (
	"sort"
)

func Sort(data sort.Interface) {
	n := data.Len()
	for i := n; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if !data.Less(j, j+1) {
				data.Swap(j, j+1)
			}
		}
	}
}
