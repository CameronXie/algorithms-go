# Red-black Tree

A Golang implementation of Red-black tree.

## Features

* Tree operations `Insert`, `Search` and `Delete`
* Thread safe.
* Extensible - any `constraints.Ordered` can be used as node key.
* Supported print Tree.

## Prerequisite

* Require Golang version 1.18+

## Usage

```go
package main

import (
	"fmt"
	"github.com/CameronXie/algorithms-go/tree/redblacktree"
	"os"
)

func main() {
	// new red-black tree.
	tree := new(redblacktree.Tree[int, string])

	// insert 10 nodes.
	for i := range make([]int, 10) {
		// using integer 0 to 9 as key, using "Item_0" to "Item_9" as value.
		_ = tree.Insert(i, fmt.Sprintf("Item_%v", i))
	}

	// print tree.
	_ = tree.Print(os.Stdout)
	/*
        Output:

		3-Item_3(BLACK)
		|---L: 1-Item_1(BLACK)
		|   |---L: 0-Item_0(BLACK)
		|   `---R: 2-Item_2(BLACK)
		`---R: 5-Item_5(BLACK)
		    |---L: 4-Item_4(BLACK)
		    `---R: 7-Item_7(RED)
		        |---L: 6-Item_6(BLACK)
		        `---R: 8-Item_8(BLACK)
		            `---R: 9-Item_9(RED)
	*/

	// search node by key.
	n, _ := tree.Search(6)
	fmt.Println(n)
	// output: 6-Item_6(BLACK)

	// delete node (or root) by key.
	_ = tree.Delete(3)

	// print tree again.
	_ = tree.Print(os.Stdout)
	/* 
	    Output:
	
		4-Item_4(BLACK)
		|---L: 1-Item_1(BLACK)
		|   |---L: 0-Item_0(BLACK)
		|   `---R: 2-Item_2(BLACK)
		`---R: 5-Item_5(BLACK)
		    `---R: 7-Item_7(RED)
		        |---L: 6-Item_6(BLACK)
		        `---R: 8-Item_8(BLACK)
		            `---R: 9-Item_9(RED)
	*/
}
```
