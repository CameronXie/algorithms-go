package quick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	cases := map[string]struct {
		input    testData
		expected testData
	}{
		"sort three items": {
			input:    testData{2, 1, 3},
			expected: testData{1, 2, 3},
		},
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

func (d testData) Len() int {
	return len(d)
}

func (d testData) Less(i, j int) bool {
	return d[i] < d[j]
}

func (d testData) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
