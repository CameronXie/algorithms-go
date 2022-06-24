package selection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	cases := map[string]struct {
		input    testData
		expected testData
	}{
		"sort descending list": {
			input:    testData{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expected: testData{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		"sort random list": {
			input:    testData{6, 2, 7, 1, 9, 10, 8, 3, 5, 4},
			expected: testData{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		"sort empty list": {
			input:    testData{},
			expected: testData{},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)

			Sort(tc.input)
			a.Equal(tc.expected, tc.input)
		})
	}
}

type testData []int

func (l testData) Len() int {
	return len(l)
}

func (l testData) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l testData) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
