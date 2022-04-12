package redblacktree

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTree_Insert(t *testing.T) {
	cases := map[string]struct {
		values   []int
		expected map[int]bool
		err      error
	}{
		"grandparent node is not null, parent node and uncle node are red": {
			values: []int{1, 2, 3, 4, 5, 6},
			expected: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: true,
				5: false,
				6: true,
			},
		},
		"parent node is red, uncle node is black, and parent node is the inner child node of grandparent node": {
			values: []int{5, 4, 3, 1, 2},
			expected: map[int]bool{
				1: true,
				2: false,
				3: true,
				4: false,
				5: false,
			},
		},
		"parent node is red, uncle node is black, and parent node is the outer child node of grandparent node": {
			values: []int{1, 2, 3, 5, 4},
			expected: map[int]bool{
				1: false,
				2: false,
				3: true,
				4: false,
				5: true,
			},
		},
		"duplicated node error": {
			values: []int{5, 5},
			expected: map[int]bool{
				5: false,
			},
			err: valueAlreadyExistsError(5),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree)
			var err error
			for _, v := range tc.values {
				if err = tree.Insert(v); err != nil {
					break
				}
			}

			a.Equal(tc.err, err)
			a.EqualValues(tc.expected, toMap(tree))
		})
	}
}

func TestTree_Search(t *testing.T) {
	cases := map[string]struct {
		insertValue []int
		search      int
		err         error
	}{
		"search exists value": {
			insertValue: []int{1, 2, 3},
			search:      3,
		},
		"search not exists value": {
			insertValue: []int{1, 2, 3},
			search:      5,
			err:         valueNotExistsError(5),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree)
			for _, v := range tc.insertValue {
				_ = tree.Insert(v)
			}

			node, err := tree.Search(tc.search)

			a.Equal(tc.err, err)
			if tc.err == nil {
				a.EqualValues(tc.search, node.value)
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {
	cases := map[string]struct {
		insertValues []int
		delete       int
		expected     map[int]bool
		err          error
	}{
		"delete the only node": {
			insertValues: []int{1},
			delete:       1,
			expected:     map[int]bool{},
		},
		"delete no exists node": {
			insertValues: []int{1},
			delete:       2,
			expected:     map[int]bool{1: false},
			err:          valueNotExistsError(2),
		},
		"delete no exists node in empty tree": {
			insertValues: []int{},
			delete:       1,
			expected:     map[int]bool{},
			err:          valueNotExistsError(1),
		},
		"delete the inner child node": {
			insertValues: []int{1, 2, 4, 3},
			delete:       4,
			expected: map[int]bool{
				1: false,
				2: false,
				3: false,
			},
		},
		"delete the outer child node": {
			insertValues: []int{1, 2, 3, 4},
			delete:       3,
			expected: map[int]bool{
				1: false,
				2: false,
				4: false,
			},
		},
		"sibling note is red, delete node is the inner child of parent node": {
			insertValues: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			delete:       5,
			expected: map[int]bool{
				1:  false,
				2:  false,
				3:  false,
				4:  false,
				6:  false,
				7:  true,
				8:  false,
				9:  false,
				10: true,
			},
		},
		"sibling node is red, deleted node is the outer child of parent node": {
			insertValues: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			delete:       6,
			expected: map[int]bool{
				1:  true,
				2:  false,
				3:  false,
				4:  true,
				5:  false,
				7:  false,
				8:  false,
				9:  false,
				10: false,
			},
		},
		"sibling node is black, and both siblings child nodes are black": {
			insertValues: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			delete:       6,
			expected: map[int]bool{
				1:  false,
				2:  false,
				3:  false,
				4:  false,
				5:  false,
				7:  false,
				8:  true,
				9:  false,
				10: true,
			},
		},
		"delete node is inner child, sibling node is black, and siblings node inner child node is red": {
			insertValues: []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 9},
			delete:       7,
			expected: map[int]bool{
				1:  false,
				2:  false,
				3:  false,
				4:  false,
				5:  false,
				6:  false,
				8:  false,
				9:  true,
				10: false,
			},
		},
		"delete node is outer child, sibling node is black, and siblings node outer child node is red": {
			insertValues: []int{1, 2, 5, 6, 3, 4},
			delete:       6,
			expected: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: true,
				5: false,
			},
		},
		"delete root, root has two child nodes and successor node has outer child node": {
			insertValues: []int{1, 2, 5, 6, 3, 4},
			delete:       2,
			expected: map[int]bool{
				1: false,
				3: false,
				4: false,
				5: true,
				6: false,
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree)
			for _, v := range tc.insertValues {
				_ = tree.Insert(v)
			}

			err := tree.Delete(tc.delete)

			res := toMap(tree)
			a.Equal(tc.err, err)
			a.EqualValues(tc.expected, res)
		})
	}
}

func TestTree_Print(t *testing.T) {
	cases := map[string]struct {
		insertValues []int
		expected     string
		err          error
	}{
		"print a empty tree": {
			insertValues: []int{},
			expected:     "empty\n",
		},
		"print the tree": {
			insertValues: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected:     "4(BLACK)\n|---L: 2(BLACK)\n|   |---L: 1(BLACK)\n|   `---R: 3(BLACK)\n`---R: 6(BLACK)\n    |---L: 5(BLACK)\n    `---R: 8(RED)\n        |---L: 7(BLACK)\n        `---R: 9(BLACK)\n            `---R: 10(RED)\n",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree)
			for _, v := range tc.insertValues {
				_ = tree.Insert(v)
			}

			var sb strings.Builder
			err := tree.Print(&sb)

			a.Equal(tc.err, err)
			a.Equal(tc.expected, sb.String())
		})
	}
}

func toMap(t *Tree) map[int]bool {
	res := make(map[int]bool)
	for _, node := range t.ToList() {
		res[node.value] = node.colour
	}

	return res
}
