# D-ary heap

A Golang implementation of D-ary heap.

## Features

* Heap operations `Len`, `Peek`, `Find`, `Pop`, `Push`, `Update` and `Remove`.
* Thread safe.
* Extensible - Implement `heap.Node` interface.

## Usage

```go
package main

import (
	"fmt"
	"github.com/CameronXie/algorithms-go/tree/heap"
)

type Item struct {
	Priority int
	Value    string
}

func (i Item) GetUniqueID() string {
	return i.Value
}

func (i Item) Less(data heap.Node) bool {
	return i.Priority > data.(Item).Priority
}

func main() {
	items := &[]Item{
		{Priority: 3, Value: "A"},
		{Priority: 1, Value: "B"},
		{Priority: 5, Value: "C"},
		{Priority: 4, Value: "D"},
	}

	dh := heap.New(3, items)

	_ = dh.Push(Item{Priority: 2, Value: "E"})

	n, _ := dh.Find("E")
	fmt.Println(n)
	// output: {2, E}

	dn, _ := dh.Remove("B")
	fmt.Println(dn)
	// output: {1 B}

	_ = dh.Update("E", func(old Item) Item {
		return Item{Priority: old.Priority + 4, Value: "E"}
	})

	for dh.Len() > 0 {
		fmt.Println(dh.Pop())
	}
	/*
		    output:

			{6 E}
			{5 C}
			{4 D}
			{3 A}
	*/
}
```
