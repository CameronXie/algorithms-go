# Treap

A Golang implementation of Treap.

## Features

* Treap operations `Search`, `Insert`, `Update`, `Pop`, and `Delete`.
* Thread safe.
* Supported print Treap.

## Prerequisite

* Require Golang version 1.18+

## Usage

```go
package main

import (
	"fmt"
	"github.com/CameronXie/algorithms-go/tree/treap"
	"os"
)

func main() {
	h := treap.New[string, int](func(i, j *treap.Node[string, int]) bool {
		if i.Priority() == j.Priority() {
			return i.Key() < j.Key()
		}

		return i.Priority() > j.Priority()
	})

	_ = h.Insert(treap.NewNode("A", 3))
	_ = h.Insert(treap.NewNode("B", 1))
	_ = h.Insert(treap.NewNode("C", 5))
	_ = h.Insert(treap.NewNode("D", 4))
	_ = h.Insert(treap.NewNode("E", 2))

	_ = h.Print(os.Stdout)
	/*
		output:

		C(5)
		|---L: A(3)
		|   `---R: B(1)
		`---R: D(4)
		    `---R: E(2)
	*/

	n, _ := h.Search("E")
	fmt.Println(n)
	// output E(2)

	_ = h.Update("E", 6)

	_ = h.Delete("B")

	_ = h.Print(os.Stdout)
	/*
		output:

		E(6)
		`---L: C(5)
		    |---L: A(3)
		    `---R: D(4)
	*/

	_ = h.Insert(treap.NewNode("F", 4))
	for i := h.Pop(); i != nil; i = h.Pop() {
		fmt.Printf("%v (%v)\n", i.Key(), i.Priority())
	}

	/*
		output:

		E (6)
		C (5)
		D (4)
		F (4)
		A (3)
	*/

	_ = h.Print(os.Stdout)
	// output: empty
}
```
