package redblacktree

import (
	"fmt"
	"github.com/CameronXie/algorithms-go/tree"
	"golang.org/x/exp/constraints"
	"io"
	"sync"
)

const (
	colourRed     = true
	colourBlack   = false
	positionLeft  = true
	positionRight = false
)

type Node[K constraints.Ordered, V any] struct {
	key   K
	value V

	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	colour bool
}

func (n *Node[K, V]) Traversal() []*Node[K, V] {
	l := []*Node[K, V]{n}

	for i := 0; i < len(l); i++ {
		if l[i].left != nil {
			l = append(l, l[i].left)
		}

		if l[i].right != nil {
			l = append(l, l[i].right)
		}
	}

	return l
}

func (n *Node[K, V]) String() string {
	colour := "BLACK"
	if n.colour == colourRed {
		colour = "RED"
	}

	return fmt.Sprintf("%v-%v(%v)", n.key, n.value, colour)
}

func (n *Node[K, V]) Left() tree.Node {
	return n.left
}

func (n *Node[K, V]) Right() tree.Node {
	return n.right
}

func (n *Node[K, V]) search(i K) (*Node[K, V], error) {
	if i > n.key {
		if n.right == nil {
			return nil, valueNotExistsError(i)
		}

		return n.right.search(i)
	}

	if i < n.key {
		if n.left == nil {
			return nil, valueNotExistsError(i)
		}

		return n.left.search(i)
	}

	return n, nil
}

func (n *Node[K, V]) insertNode(node *Node[K, V]) error {
	if node.key > n.key {
		if n.right == nil {
			n.addChildNode(node, positionRight)
			return nil
		}

		return n.right.insertNode(node)
	}

	if node.key < n.key {
		if n.left == nil {
			n.addChildNode(node, positionLeft)
			return nil
		}

		return n.left.insertNode(node)
	}

	return valueAlreadyExistsError(node.key)
}

func (n *Node[K, V]) getChildNodePosition(child *Node[K, V]) bool {
	if n.left != nil && n.left.key == child.key {
		return positionLeft
	}

	if n.right != nil && n.right.key == child.key {
		return positionRight
	}

	panic(invalidChildError(n.key, child.key))
}

func (n *Node[K, V]) addChildNode(newNode *Node[K, V], position bool) {
	if newNode != nil {
		newNode.parent = n
	}

	if position == positionLeft {
		n.left = newNode
		return
	}

	n.right = newNode
}

func (n *Node[K, V]) replaceChildNode(old, new *Node[K, V]) {
	p := n.getChildNodePosition(old)

	n.removeChildNode(old)
	n.addChildNode(new, p)
}

func (n *Node[K, V]) removeChildNode(childNode *Node[K, V]) {
	if n.left != nil && n.left.key == childNode.key {
		n.left = nil
		childNode.parent = nil
		return
	}

	if n.right != nil && n.right.key == childNode.key {
		n.right = nil
		childNode.parent = nil
		return
	}

	panic(invalidChildError(childNode.key, n.key))
}

func (n *Node[K, V]) getSibling() *Node[K, V] {
	if n.parent == nil {
		return nil
	}

	if n.parent.getChildNodePosition(n) == positionLeft {
		return n.parent.right
	}

	return n.parent.left
}

func (n *Node[K, V]) getUncle() *Node[K, V] {
	if n.parent == nil || n.parent.parent == nil {
		return nil
	}

	return n.parent.getSibling()
}

func (n *Node[K, V]) getMinimumNode() *Node[K, V] {
	if n.left == nil {
		return n
	}

	return n.left.getMinimumNode()
}

func isBlackNode[K constraints.Ordered, V any](n *Node[K, V]) bool {
	if n == nil || n.colour == colourBlack {
		return true
	}

	return false
}

type Tree[K constraints.Ordered, V any] struct {
	sync.RWMutex
	root *Node[K, V]
}

func (t *Tree[K, V]) ToList() []*Node[K, V] {
	t.RLock()
	defer t.RUnlock()

	if t.root == nil {
		return make([]*Node[K, V], 0)
	}

	return t.root.Traversal()
}

func (t *Tree[K, V]) Search(key K) (*Node[K, V], error) {
	t.RLock()
	defer t.RUnlock()

	if t.root == nil {
		return nil, valueNotExistsError(key)
	}

	return t.root.search(key)
}

func (t *Tree[K, V]) Insert(key K, value V) error {
	t.Lock()
	defer t.Unlock()

	if t.root == nil {
		t.root = &Node[K, V]{key: key, value: value, colour: colourBlack}
		return nil
	}

	newNode := &Node[K, V]{key: key, value: value, colour: colourRed}
	if err := t.root.insertNode(newNode); err != nil {
		return err
	}

	t.rebalanceAfterInsertion(newNode)
	return nil
}

func (t *Tree[K, V]) rebalanceAfterInsertion(n *Node[K, V]) {
	// it is root.
	if n.parent == nil {
		n.colour = colourBlack
		return
	}

	// parent is black.
	if n.parent.colour == colourBlack {
		return
	}

	// grandparent is not null.
	// parent and uncle are red.
	parent := n.parent
	grandparent := parent.parent
	uncle := n.getUncle()
	if uncle != nil && uncle.colour == colourRed {
		uncle.colour = colourBlack
		parent.colour = colourBlack
		grandparent.colour = colourRed
		t.rebalanceAfterInsertion(grandparent)
		return
	}

	// parent is red, and uncle is black.
	// parent is the inner child of grandparent.
	if grandparent.getChildNodePosition(parent) == positionLeft {
		// new node is outer grandchild.
		if parent.getChildNodePosition(n) == positionRight {
			t.rotateLeft(parent)
			parent = n
			grandparent = parent.parent
		}

		// new node is inner grandchild.
		t.rotateRight(grandparent)
		parent.colour = colourBlack
		parent.right.colour = colourRed
		return
	}

	// parent is the outer child of grandparent.
	// new node is inner grandchild.
	if parent.getChildNodePosition(n) == positionLeft {
		t.rotateRight(parent)
		parent = n
		grandparent = parent.parent
	}

	// new node is outer grandchild.
	t.rotateLeft(grandparent)
	parent.colour = colourBlack
	parent.left.colour = colourRed
}

func (t *Tree[K, V]) Delete(i K) error {
	t.Lock()
	defer t.Unlock()

	if t.root == nil {
		return valueNotExistsError(i)
	}

	deleteNode, err := t.root.search(i)
	if err != nil {
		return err
	}

	// node has no children.
	if deleteNode.left == nil && deleteNode.right == nil {
		t.rebalanceAfterDeletion(deleteNode)
		t.replaceChildNote(deleteNode, nil)
		return nil
	}

	// node has left child.
	if deleteNode.right == nil {
		leftChild := deleteNode.left
		colour := deleteNode.colour

		t.replaceChildNote(deleteNode, leftChild)
		if colour == colourBlack {
			t.rebalanceAfterDeletion(leftChild)
		}

		return nil
	}

	// node has right child.
	if deleteNode.left == nil {
		rightChild := deleteNode.right
		colour := deleteNode.colour

		t.replaceChildNote(deleteNode, rightChild)
		if colour == colourBlack {
			t.rebalanceAfterDeletion(rightChild)
		}

		return nil
	}

	// node has two children.
	successor := deleteNode.right.getMinimumNode()
	deleteNode.key = successor.key
	deleteNode.value = successor.value
	successor.parent.replaceChildNode(successor, successor.right)

	if deleteNode.colour == colourBlack && successor.right != nil {
		t.rebalanceAfterDeletion(successor.right)
	}

	return nil
}

func (t *Tree[K, V]) rebalanceAfterDeletion(n *Node[K, V]) {
	// node is root or is red.
	if n.colour == colourRed || n.parent == nil {
		n.colour = colourBlack
		return
	}

	isInnerChild := n.parent.getChildNodePosition(n) == positionLeft
	sibling := n.getSibling()

	if sibling == nil {
		return
	}

	// sibling node is red.
	if !isBlackNode(sibling) {
		sibling.colour = colourBlack
		sibling.parent.colour = colourRed

		if isInnerChild {
			t.rotateLeft(n.parent)
		} else {
			t.rotateRight(n.parent)
		}

		sibling = n.getSibling()
	}

	if sibling == nil {
		return
	}

	// sibling node is black and both sibling's children node are black.
	if isBlackNode(sibling.left) && isBlackNode(sibling.right) {
		sibling.colour = colourRed

		if n.parent.colour == colourRed {
			n.parent.colour = colourBlack
			return
		}

		t.rebalanceAfterDeletion(n.parent)
		return
	}

	// node is inner child
	if isInnerChild {
		// sibling node is black and sibling's inner child is red.
		if !isBlackNode(sibling.left) {
			sibling.left.colour = colourBlack
			sibling.colour = colourRed

			t.rotateRight(sibling)
			sibling = n.getSibling()
		}

		// sibling node is black and sibling's outer child is red.
		sibling.colour = n.parent.colour
		n.parent.colour = colourBlack
		sibling.right.colour = colourBlack
		t.rotateLeft(n.parent)

		return
	}

	// node is outer child, sibling is black and sibling's inner child is red.
	if !isBlackNode(sibling.right) {
		sibling.right.colour = colourBlack
		sibling.colour = colourRed
		t.rotateLeft(sibling)
		sibling = n.getSibling()
	}

	// node is outer child, sibling is black and sibling's outer child is red.
	sibling.colour = n.parent.colour
	n.parent.colour = colourBlack
	sibling.left.colour = colourBlack
	t.rotateRight(n.parent)
}

func (t *Tree[K, V]) rotateLeft(n *Node[K, V]) {
	rightChild := n.right

	n.removeChildNode(rightChild)
	if rightChild.left != nil {
		n.addChildNode(rightChild.left, positionRight)
	}

	t.replaceChildNote(n, rightChild)
	rightChild.addChildNode(n, positionLeft)
}

func (t *Tree[K, V]) rotateRight(n *Node[K, V]) {
	leftChild := n.left

	n.removeChildNode(leftChild)
	if leftChild.right != nil {
		n.addChildNode(leftChild.right, positionLeft)
	}

	t.replaceChildNote(n, leftChild)
	leftChild.addChildNode(n, positionRight)
}

func (t *Tree[K, V]) replaceChildNote(oldNote *Node[K, V], newNote *Node[K, V]) {
	parent := oldNote.parent
	if parent != nil {
		parent.replaceChildNode(oldNote, newNote)
		return
	}

	t.root = newNote
}

func (t *Tree[K, V]) Print(w io.StringWriter) error {
	t.RLock()
	defer t.RUnlock()

	if t.root == nil {
		_, err := w.WriteString("empty\n")
		return err
	}

	return tree.Print(t.root, w)
}

func valueAlreadyExistsError(i any) error {
	return fmt.Errorf(`key %v already exists`, i)
}

func valueNotExistsError(i any) error {
	return fmt.Errorf(`key %v not exists`, i)
}

func invalidChildError(p, c any) error {
	return fmt.Errorf(`%v is not a child node of %v node`, c, p)
}
