package treap

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTreap_Insert(t *testing.T) {
	cases := map[string]struct {
		nodes    []*Node[string, int]
		expected map[string]int
		err      error
	}{
		"insert new node": {
			nodes: []*Node[string, int]{
				{
					key:      "A",
					priority: 2,
				},
				{
					key:      "C",
					priority: 1,
				},
				{
					key:      "B",
					priority: 3,
				},
			},
			expected: map[string]int{
				"B": 3,
				"A": 2,
				"C": 1,
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			for _, i := range tc.nodes {
				err := treap.Insert(i)
				a.Equal(tc.err, err)

				if tc.err != nil {
					return
				}
			}

			res := make(map[string]int)
			for _, i := range treap.root.traversal() {
				res[i.key] = i.priority
			}

			a.Equal(tc.expected, res)
		})
	}
}

func TestTreap_Search(t *testing.T) {
	cases := map[string]struct {
		key      string
		expected *Node[string, int]
		err      error
	}{
		"node exist": {
			key: "A",
			expected: &Node[string, int]{
				key:      "A",
				priority: 3,
			},
		},
		"node not exist": {
			key: "B",
			err: valueNotExistsError("B"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			nodes := []*Node[string, int]{
				{
					key:      "A",
					priority: 3,
				},
			}

			for _, node := range nodes {
				_ = treap.Insert(node)
			}

			node, err := treap.Search(tc.key)

			a.Equal(tc.expected, node)
			a.Equal(tc.err, err)
		})
	}
}

func TestTreap_Update(t *testing.T) {
	cases := map[string]struct {
		key      string
		priority int
		expected *Node[string, int]
		err      error
	}{
		"node exist": {
			key:      "A",
			priority: 10,
			expected: &Node[string, int]{
				key:      "A",
				priority: 10,
			},
		},
		"node not exist": {
			key: "B",
			err: valueNotExistsError("B"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			nodes := []*Node[string, int]{
				{
					key:      "A",
					priority: 3,
				},
			}

			for _, node := range nodes {
				_ = treap.Insert(node)
			}

			err := treap.Update(tc.key, tc.priority)

			a.Equal(tc.err, err)

			if tc.err != nil {
				return
			}

			node, rErr := treap.Search(tc.key)
			a.Equal(tc.expected, node)
			a.Nil(rErr)
		})
	}
}

func TestTreap_Delete(t *testing.T) {
	cases := map[string]struct {
		key      string
		expected map[string]int
		err      error
	}{
		"node exist": {
			key: "A",
			expected: map[string]int{
				"B": 3,
				"C": 1,
			},
		},
		"node not exist": {
			key: "D",
			err: valueNotExistsError("D"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			nodes := []*Node[string, int]{
				{
					key:      "A",
					priority: 2,
				},
				{
					key:      "C",
					priority: 1,
				},
				{
					key:      "B",
					priority: 3,
				},
			}

			for _, node := range nodes {
				_ = treap.Insert(node)
			}

			err := treap.Delete(tc.key)
			a.Equal(tc.err, err)

			if tc.err != nil {
				return
			}

			res := make(map[string]int)
			for _, i := range treap.root.traversal() {
				res[i.key] = i.priority
			}

			a.Equal(tc.expected, res)
		})
	}
}

func TestTreap_Print(t *testing.T) {
	cases := map[string]struct {
		nodes    []*Node[string, int]
		expected string
		err      error
	}{
		"insert new node": {
			nodes: []*Node[string, int]{
				{
					key:      "A",
					priority: 2,
				},
				{
					key:      "C",
					priority: 1,
				},
				{
					key:      "B",
					priority: 3,
				},
			},
			expected: "B(3)\n|---L: A(2)\n`---R: C(1)\n",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			for _, i := range tc.nodes {
				_ = treap.Insert(i)
			}

			var sb strings.Builder
			err := treap.Print(&sb)

			a.Equal(tc.err, err)
			a.Equal(tc.expected, sb.String())
		})
	}
}

func TestTreap_Concurrent(t *testing.T) {
	cases := map[string]struct {
		updates  map[string]int
		delete   []string
		expected map[string]int
	}{
		"concurrent update and delete": {
			updates: map[string]int{
				"B": 1,
				"C": 5,
			},
			delete: []string{"A"},
			expected: map[string]int{
				"B": 1,
				"C": 5,
				"D": 3,
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			treap := New[string, int](func(i, j int) bool {
				return i > j
			})

			nodes := []*Node[string, int]{
				{
					key:      "A",
					priority: 2,
				},
				{
					key:      "C",
					priority: 1,
				},
				{
					key:      "D",
					priority: 3,
				},
				{
					key:      "B",
					priority: 4,
				},
			}

			for _, node := range nodes {
				_ = treap.Insert(node)
			}

			var wg sync.WaitGroup
			wg.Add(len(tc.updates) + len(tc.delete))
			for k, v := range tc.updates {
				go func(key string, value int) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_ = treap.Update(key, value)
				}(k, v)
			}

			for _, k := range tc.delete {
				go func(key string) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_ = treap.Delete(key)
				}(k)
			}
			wg.Wait()

			res := make(map[string]int)
			for _, i := range treap.root.traversal() {
				res[i.key] = i.priority
			}

			a.Equal(tc.expected, res)
		})
	}
}
