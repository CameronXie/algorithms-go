package heap

import (
	"fmt"
	"sync"
)

type Node interface {
	GetUniqueID() string
	Less(data Node) bool
}

type DHeap[T Node] struct {
	nodes *[]T
	d     int
	m     map[string]int
	mu    sync.RWMutex
}

func (dh *DHeap[T]) Len() int {
	dh.mu.RLock()
	defer dh.mu.RUnlock()

	return dh.len()
}

func (dh *DHeap[T]) Peek() Node {
	dh.mu.RLock()
	defer dh.mu.RUnlock()

	return (*dh.nodes)[0]
}

func (dh *DHeap[T]) Find(id string) (Node, error) {
	dh.mu.RLock()
	defer dh.mu.RUnlock()

	idx, ok := dh.m[id]
	if !ok {
		return nil, itemNotExistsError(id)
	}

	return (*dh.nodes)[idx], nil
}

func (dh *DHeap[T]) Pop() Node {
	dh.mu.Lock()
	defer dh.mu.Unlock()

	n := dh.len() - 1
	if n < 0 {
		return nil
	}

	dh.swap(0, n)
	dh.down(0, n)
	return dh.pop()
}

func (dh *DHeap[T]) Push(node T) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()

	id := node.GetUniqueID()
	if _, ok := dh.m[id]; ok {
		return itemAlreadyExistsError(id)
	}

	*dh.nodes = append(*dh.nodes, node)
	n := dh.len() - 1
	dh.m[id] = n
	dh.up(n)

	return nil
}

func (dh *DHeap[T]) Update(id string, updates func(old T) T) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()

	idx, ok := dh.m[id]
	if !ok {
		return itemNotExistsError(id)
	}

	old := (*dh.nodes)[idx]
	data := updates(old)
	delete(dh.m, id)
	(*dh.nodes)[idx] = data
	dh.m[data.GetUniqueID()] = idx
	dh.fix(idx, dh.len())

	return nil
}

func (dh *DHeap[T]) Remove(id string) (Node, error) {
	dh.mu.Lock()
	defer dh.mu.Unlock()

	idx, ok := dh.m[id]
	if !ok {
		return nil, itemNotExistsError(id)
	}

	n := dh.len() - 1
	if n != idx {
		dh.swap(n, idx)
		dh.fix(idx, n)
	}

	return dh.pop(), nil
}

func (dh *DHeap[T]) init() {
	n := dh.len()
	for i := 0; i < n; i++ {
		node := (*dh.nodes)[i]
		dh.m[node.GetUniqueID()] = i
	}

	for i := (n - 1) / dh.d; i >= 0; i-- {
		dh.down(i, n)
	}
}

func (dh *DHeap[T]) len() int {
	return len(*dh.nodes)
}

func (dh *DHeap[T]) swap(i, j int) {
	nodes := *dh.nodes
	nodes[i], nodes[j] = nodes[j], nodes[i]
	dh.m[nodes[i].GetUniqueID()] = i
	dh.m[nodes[j].GetUniqueID()] = j
}

func (dh *DHeap[T]) pop() Node {
	n := dh.len() - 1
	nodes := *dh.nodes
	last := nodes[n]

	delete(dh.m, last.GetUniqueID())
	*dh.nodes = nodes[:n]

	return last
}

func (dh *DHeap[T]) fix(idx, n int) {
	if !dh.down(idx, n) {
		dh.up(idx)
	}
}

func (dh *DHeap[T]) up(idx int) {
	nodes := *dh.nodes
	for {
		parent := (idx - 1) / dh.d
		if parent == idx || !nodes[idx].Less(nodes[parent]) {
			break
		}

		dh.swap(idx, parent)
		idx = parent
	}
}

func (dh *DHeap[T]) down(idx, n int) bool {
	current := idx
	nodes := *dh.nodes
	for {
		swapID := current
		for i := 1; i <= dh.d; i++ {
			childIdx := current*dh.d + i

			if childIdx >= n || childIdx < 0 {
				break
			}

			if nodes[childIdx].Less(nodes[swapID]) {
				swapID = childIdx
			}
		}

		if swapID == current {
			break
		}

		dh.swap(current, swapID)
		current = swapID
	}

	return current > idx
}

func itemAlreadyExistsError(id string) error {
	return fmt.Errorf(`id %v already exists`, id)
}

func itemNotExistsError(id string) error {
	return fmt.Errorf(`id %v not found`, id)
}

func New[T Node](d int, items *[]T) *DHeap[T] {
	dh := DHeap[T]{d: d, nodes: items, m: make(map[string]int)}
	dh.init()

	return &dh
}
