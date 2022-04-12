package selection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	cases := map[string]struct {
		input    []int
		expected []int
	}{
		"sort descending list": {
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		"sort random list": {
			input:    []int{6, 2, 7, 1, 9, 10, 8, 3, 5, 4},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		"sort empty list": {
			input:    []int{},
			expected: []int{},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			original := make([]int, len(tc.input))
			copy(original, tc.input)

			sorted := Sort(tc.input)
			a.Equal(tc.expected, sorted)
			a.Equal(original, tc.input)
		})
	}
}
