package heap

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDHeap_Len(t *testing.T) {
	cases := map[string]struct {
		expected int
	}{
		"item length": {
			expected: 5,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			a.Equal(tc.expected, dh.Len())
		})
	}
}

func TestDHeap_Peek(t *testing.T) {
	cases := map[string]struct {
		expected testItem
	}{
		"peek top item": {
			expected: testItem{Priority: 5, Value: "B"},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			length := len(*dh.nodes)
			item := dh.Peek()

			a.Equal(tc.expected, item)
			a.Equal(length, len(*dh.nodes))
		})
	}
}

func TestDHeap_Find(t *testing.T) {
	cases := map[string]struct {
		id       string
		expected *testItem
		err      error
	}{
		"find exists item": {
			id:       "C",
			expected: &testItem{Priority: 3, Value: "C"},
		},
		"find not exists item": {
			id:       "F",
			expected: nil,
			err:      errors.New("id F not found"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			item, err := setupTestData().Find(tc.id)

			if tc.err == nil {
				a.Equal(*tc.expected, item.(testItem))
				a.Nil(err)
				return
			}

			a.Nil(item)
			a.Equal(tc.err, err)
		})
	}
}

func TestDHeap_Pop(t *testing.T) {
	cases := map[string]struct {
		dh       *DHeap[testItem]
		expected *testItem
	}{
		"pop top item": {
			dh:       setupTestData(),
			expected: &testItem{Priority: 5, Value: "B"},
		},
		"pop from empty heap": {
			dh:       New(3, &[]testItem{}),
			expected: nil,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			length := len(*tc.dh.nodes)
			item := tc.dh.Pop()

			if tc.expected != nil {
				a.Equal(length-1, len(*tc.dh.nodes))
				a.Equal(*tc.expected, item.(testItem))
				return
			}

			a.Nil(item)
			a.Equal(0, len(*tc.dh.nodes))
		})
	}
}

func TestDHeap_Push(t *testing.T) {
	cases := map[string]struct {
		item     testItem
		expected []testItem
		err      error
	}{
		"push item": {
			item: testItem{Priority: 6, Value: "F"},
			expected: []testItem{
				{Priority: 6, Value: "F"},
				{Priority: 5, Value: "B"},
				{Priority: 4, Value: "E"},
				{Priority: 3, Value: "C"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
		},
		"push exists item": {
			item: testItem{Priority: 6, Value: "A"},
			expected: []testItem{
				{Priority: 5, Value: "B"},
				{Priority: 4, Value: "E"},
				{Priority: 3, Value: "C"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
			err: errors.New("id A already exists"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			err := dh.Push(tc.item)

			items := getItems(dh)
			a.Equal(tc.expected, items)
			a.Equal(tc.err, err)
		})
	}
}

func TestDHeap_Update(t *testing.T) {
	cases := map[string]struct {
		updatedID   string
		updatedItem testItem
		expected    []testItem
		err         error
	}{
		"update item": {
			updatedID:   "C",
			updatedItem: testItem{Priority: 6, Value: "C_Updated"},
			expected: []testItem{
				{Priority: 9, Value: "C_Updated"},
				{Priority: 5, Value: "B"},
				{Priority: 4, Value: "E"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
		},
		"update item not found": {
			updatedID: "F",
			expected: []testItem{
				{Priority: 5, Value: "B"},
				{Priority: 4, Value: "E"},
				{Priority: 3, Value: "C"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
			err: errors.New("id F not found"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			err := dh.Update(tc.updatedID, func(old testItem) testItem {
				return testItem{
					Priority: tc.updatedItem.Priority + old.Priority,
					Value:    tc.updatedItem.Value,
				}
			})

			if tc.err != nil {
				a.Equal(tc.err, err)
			}

			items := getItems(dh)
			a.Equal(tc.expected, items)
		})
	}
}

func TestDHeap_Remote(t *testing.T) {
	cases := map[string]struct {
		removeID    string
		removedItem testItem
		expected    []testItem
		err         error
	}{
		"remove item by id": {
			removeID:    "C",
			removedItem: testItem{Priority: 3, Value: "C"},
			expected: []testItem{
				{Priority: 5, Value: "B"},
				{Priority: 4, Value: "E"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
		},
		"remove item not found": {
			removeID: "F",
			err:      errors.New("id F not found"),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			node, err := dh.Remove(tc.removeID)

			if tc.err != nil {
				a.Equal(tc.err, err)
				a.Nil(node)
				return
			}

			items := getItems(dh)
			a.Equal(tc.expected, items)
			a.Equal(tc.removedItem, node.(testItem))
			a.Nil(err)
		})
	}
}

func TestDHeap_ConcurrentWrite(t *testing.T) {
	cases := map[string]struct {
		updates  map[string]testItem
		delete   []string
		expected []testItem
	}{
		"concurrent update and delete": {
			updates: map[string]testItem{
				"C": {Priority: 6, Value: "C_Updated"},
				"E": {Priority: 3, Value: "E_Updated"},
			},
			delete: []string{"B", "A"},
			expected: []testItem{
				{Priority: 6, Value: "C_Updated"},
				{Priority: 3, Value: "E_Updated"},
				{Priority: 1, Value: "D"},
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()

			var wg sync.WaitGroup
			wg.Add(len(tc.updates) + len(tc.delete))
			for id, value := range tc.updates {
				go func(id string, value testItem) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_ = dh.Update(id, func(_ testItem) testItem {
						return value
					})
				}(id, value)
			}

			for _, id := range tc.delete {
				go func(id string) {
					defer wg.Done()
					time.Sleep(time.Millisecond)
					_, _ = dh.Remove(id)
				}(id)
			}

			wg.Wait()
			a.Equal(tc.expected, getItems(dh))
		})
	}
}

func TestDHeap_ConcurrentRead(t *testing.T) {
	cases := map[string]struct {
		updates  map[string]testItem
		expected []testItem
	}{
		"concurrent read while update": {
			updates: map[string]testItem{
				"C": {Priority: 6, Value: "C_Updated"},
				"E": {Priority: 3, Value: "E_Updated"},
			},
			expected: []testItem{
				{Priority: 6, Value: "C_Updated"},
				{Priority: 5, Value: "B"},
				{Priority: 3, Value: "E_Updated"},
				{Priority: 2, Value: "A"},
				{Priority: 1, Value: "D"},
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			a := assert.New(t)
			dh := setupTestData()
			length := dh.len()

			var wg sync.WaitGroup
			wg.Add(len(tc.updates) * 2)
			for id, value := range tc.updates {
				go func(id string, value testItem) {
					defer wg.Done()
					_ = dh.Update(id, func(_ testItem) testItem {
						return value
					})
				}(id, value)

				go func() {
					defer wg.Done()
					a.Equal(length, dh.Len())
				}()
			}

			wg.Wait()
			a.Equal(tc.expected, getItems(dh))
		})
	}
}

func setupTestData() *DHeap[testItem] {
	return New(3, &[]testItem{
		{Priority: 2, Value: "A"},
		{Priority: 5, Value: "B"},
		{Priority: 3, Value: "C"},
		{Priority: 1, Value: "D"},
		{Priority: 4, Value: "E"},
	})
}

func getItems(dh *DHeap[testItem]) []testItem {
	items := make([]testItem, 0)
	for dh.len() > 0 {
		items = append(items, dh.Pop().(testItem))
	}

	return items
}

type testItem struct {
	Priority int
	Value    string
}

func (i testItem) GetUniqueID() string {
	return i.Value
}

func (i testItem) Less(data Node) bool {
	return i.Priority > data.(testItem).Priority
}
