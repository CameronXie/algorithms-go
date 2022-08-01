package redblacktree

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTree_Insert(t *testing.T) {
	cases := map[string]struct {
		keys     []int
		expected map[int]bool
		err      error
	}{
		"grandparent node is not null, parent node and uncle node are red": {
			keys: []int{1, 2, 3, 4, 5, 6},
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
			keys: []int{5, 4, 3, 1, 2},
			expected: map[int]bool{
				1: true,
				2: false,
				3: true,
				4: false,
				5: false,
			},
		},
		"parent node is red, uncle node is black, and parent node is the outer child node of grandparent node": {
			keys: []int{1, 2, 3, 5, 4},
			expected: map[int]bool{
				1: false,
				2: false,
				3: true,
				4: false,
				5: true,
			},
		},
		"duplicated node error": {
			keys: []int{5, 5},
			expected: map[int]bool{
				5: false,
			},
			err: valueAlreadyExistsError(5),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree[int, string])
			var err error
			for v, k := range tc.keys {
				if err = tree.Insert(k, strconv.Itoa(v)); err != nil {
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
		keys   []int
		search int
		err    error
	}{
		"search exists key": {
			keys:   []int{1, 2, 3},
			search: 3,
		},
		"search not exists key": {
			keys:   []int{1, 2, 3},
			search: 5,
			err:    valueNotExistsError(5),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree[int, string])
			for i, k := range tc.keys {
				_ = tree.Insert(k, strconv.Itoa(i))
			}

			node, err := tree.Search(tc.search)

			a.Equal(tc.err, err)
			if tc.err == nil {
				a.EqualValues(tc.search, node.key)
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {
	cases := map[string]struct {
		keys     []int
		delete   int
		expected map[int]bool
		err      error
	}{
		"delete the only node": {
			keys:     []int{1},
			delete:   1,
			expected: map[int]bool{},
		},
		"delete no exists node": {
			keys:     []int{1},
			delete:   2,
			expected: map[int]bool{1: false},
			err:      valueNotExistsError(2),
		},
		"delete no exists node in empty tree": {
			keys:     []int{},
			delete:   1,
			expected: map[int]bool{},
			err:      valueNotExistsError(1),
		},
		"delete the inner child node": {
			keys:   []int{1, 2, 4, 3},
			delete: 4,
			expected: map[int]bool{
				1: false,
				2: false,
				3: false,
			},
		},
		"delete the outer child node": {
			keys:   []int{1, 2, 3, 4},
			delete: 3,
			expected: map[int]bool{
				1: false,
				2: false,
				4: false,
			},
		},
		"sibling note is red, delete node is the inner child of parent node": {
			keys:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			delete: 5,
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
			keys:   []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			delete: 6,
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
			keys:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			delete: 6,
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
			keys:   []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 9},
			delete: 7,
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
			keys:   []int{1, 2, 5, 6, 3, 4},
			delete: 6,
			expected: map[int]bool{
				1: false,
				2: false,
				3: false,
				4: true,
				5: false,
			},
		},
		"delete root, root has two child nodes and successor node has outer child node": {
			keys:   []int{1, 2, 5, 6, 3, 4},
			delete: 2,
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
			tree := new(Tree[int, string])
			for i, k := range tc.keys {
				_ = tree.Insert(k, strconv.Itoa(i))
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
		keys     []int
		expected string
		err      error
	}{
		"print a empty tree": {
			keys:     []int{},
			expected: "empty\n",
		},
		"print the tree": {
			keys:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: "4-3(BLACK)\n|---L: 2-1(BLACK)\n|   |---L: 1-0(BLACK)\n|   `---R: 3-2(BLACK)\n`---R: 6-5(BLACK)\n    |---L: 5-4(BLACK)\n    `---R: 8-7(RED)\n        |---L: 7-6(BLACK)\n        `---R: 9-8(BLACK)\n            `---R: 10-9(RED)\n",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree[int, string])
			for i, k := range tc.keys {
				_ = tree.Insert(k, strconv.Itoa(i))
			}

			var sb strings.Builder
			err := tree.Print(&sb)

			a.Equal(tc.err, err)
			a.Equal(tc.expected, sb.String())
		})
	}
}

func TestTree_Concurrent(t *testing.T) {
	cases := map[string]struct {
		insertKeys []int
		deleteKeys []int
		expected   []int
	}{
		"concurrent insert and delete": {
			insertKeys: []int{6, 7, 8},
			deleteKeys: []int{1, 2, 3},
			expected:   []int{4, 5, 6, 7, 8},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			tree := new(Tree[int, string])

			for i := range make([]struct{}, 5) {
				_ = tree.Insert(i+1, strconv.Itoa(i+1))
			}

			var wg sync.WaitGroup
			wg.Add(len(tc.insertKeys) + len(tc.deleteKeys))

			for _, i := range tc.insertKeys {
				go func(n int) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_ = tree.Insert(n, strconv.Itoa(n))
				}(i)
			}

			for _, i := range tc.deleteKeys {
				go func(n int) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_ = tree.Delete(n)
				}(i)
			}

			wg.Wait()
			a.ElementsMatch(tc.expected, toList(tree))
		})
	}
}

func toMap[K constraints.Ordered, V any](t *Tree[K, V]) map[K]bool {
	res := make(map[K]bool)
	for _, node := range t.ToList() {
		res[node.key] = node.colour
	}

	return res
}

func toList[K constraints.Ordered, V any](t *Tree[K, V]) []K {
	res := make([]K, 0)
	for _, node := range t.ToList() {
		res = append(res, node.key)
	}

	return res
}
