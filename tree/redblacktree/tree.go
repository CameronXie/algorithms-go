package redblacktree

import (
	"fmt"
	"io"
)

const (
	colourRed     = true
	colourBlack   = false
	positionLeft  = true
	positionRight = false
)

type Tree struct {
	root *Node
}

func (t *Tree) ToList() []*Node {
	if t.root == nil {
		return make([]*Node, 0)
	}

	return t.root.Traversal()
}

func (t *Tree) Search(i int) (*Node, error) {
	return t.root.search(i)
}

func (t *Tree) Insert(i int) error {
	if t.root == nil {
		t.root = &Node{value: i, colour: colourBlack}
		return nil
	}

	newNode := &Node{value: i, colour: colourRed}
	if err := t.root.insertNode(newNode); err != nil {
		return err
	}

	t.rebalanceAfterInsertion(newNode)
	return nil
}

func (t *Tree) rebalanceAfterInsertion(n *Node) {
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

func (t *Tree) Delete(i int) error {
	if t.root == nil {
		return valueNotExistsError(i)
	}

	deleteNode, err := t.Search(i)
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
	deleteNode.value = successor.value
	successor.parent.replaceChildNode(successor, successor.right)

	if deleteNode.colour == colourBlack && successor.right != nil {
		t.rebalanceAfterDeletion(successor.right)
	}

	return nil
}

func (t *Tree) rebalanceAfterDeletion(n *Node) {
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

func (t *Tree) rotateLeft(n *Node) {
	rightChild := n.right

	n.removeChildNode(rightChild)
	if rightChild.left != nil {
		n.addChildNode(rightChild.left, positionRight)
	}

	t.replaceChildNote(n, rightChild)
	rightChild.addChildNode(n, positionLeft)
}

func (t *Tree) rotateRight(n *Node) {
	leftChild := n.left

	n.removeChildNode(leftChild)
	if leftChild.right != nil {
		n.addChildNode(leftChild.right, positionLeft)
	}

	t.replaceChildNote(n, leftChild)
	leftChild.addChildNode(n, positionRight)
}

func (t *Tree) replaceChildNote(oldNote *Node, newNote *Node) {
	parent := oldNote.parent
	if parent != nil {
		parent.replaceChildNode(oldNote, newNote)
		return
	}

	t.root = newNote
}

func (t *Tree) Print(w io.StringWriter) error {
	if t.root == nil {
		_, err := w.WriteString("empty\n")
		return err
	}

	return t.root.Print(w)
}

type Node struct {
	value  int
	parent *Node
	left   *Node
	right  *Node
	colour bool
}

func (n *Node) Traversal() []*Node {
	l := []*Node{n}

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

func (n *Node) Print(w io.StringWriter) error {
	return n.print(w, "", "", false, false)
}

func (n *Node) print(w io.StringWriter, indent string, position string, isOpen bool, isLast bool) error {
	colour := "BLACK"
	if n.colour == colourRed {
		colour = "RED"
	}

	cornerSymbol := "|"
	if isLast {
		cornerSymbol = "`"
	}

	if position != "" {
		position = fmt.Sprintf("%v---%v: ", cornerSymbol, position)
	}

	if _, err := w.WriteString(fmt.Sprintf("%v%v%v(%v)\n", indent, position, n.value, colour)); err != nil {
		return err
	}

	if n.parent != nil {
		if isOpen {
			indent += "|   "
		} else {
			indent += "    "
		}
	}

	if n.left != nil {
		if err := n.left.print(w, indent, "L", true, n.right == nil); err != nil {
			return err
		}
	}

	if n.right != nil {
		if err := n.right.print(w, indent, "R", false, true); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) search(i int) (*Node, error) {
	if i > n.value {
		if n.right == nil {
			return nil, valueNotExistsError(i)
		}

		return n.right.search(i)
	}

	if i < n.value {
		if n.left == nil {
			return nil, valueNotExistsError(i)
		}

		return n.left.search(i)
	}

	return n, nil
}

func (n *Node) insertNode(node *Node) error {
	if node.value > n.value {
		if n.right == nil {
			n.addChildNode(node, positionRight)
			return nil
		}

		return n.right.insertNode(node)
	}

	if node.value < n.value {
		if n.left == nil {
			n.addChildNode(node, positionLeft)
			return nil
		}

		return n.left.insertNode(node)
	}

	return valueAlreadyExistsError(node.value)
}

func (n *Node) getChildNodePosition(child *Node) bool {
	if n.left != nil && n.left.value == child.value {
		return positionLeft
	}

	if n.right != nil && n.right.value == child.value {
		return positionRight
	}

	panic(invalidChildError(n.value, child.value))
}

func (n *Node) addChildNode(newNode *Node, position bool) {
	if newNode != nil {
		newNode.parent = n
	}

	if position == positionLeft {
		n.left = newNode
		return
	}

	n.right = newNode
}

func (n *Node) replaceChildNode(old, new *Node) {
	p := n.getChildNodePosition(old)

	n.removeChildNode(old)
	n.addChildNode(new, p)
}

func (n *Node) removeChildNode(childNode *Node) {
	if n.left != nil && n.left.value == childNode.value {
		n.left = nil
		childNode.parent = nil
		return
	}

	if n.right != nil && n.right.value == childNode.value {
		n.right = nil
		childNode.parent = nil
		return
	}

	panic(invalidChildError(childNode.value, n.value))
}

func (n *Node) getSibling() *Node {
	if n.parent == nil {
		return nil
	}

	if n.parent.getChildNodePosition(n) == positionLeft {
		return n.parent.right
	}

	return n.parent.left
}

func (n *Node) getUncle() *Node {
	if n.parent == nil || n.parent.parent == nil {
		return nil
	}

	return n.parent.getSibling()
}

func (n *Node) getMinimumNode() *Node {
	if n.left == nil {
		return n
	}

	return n.left.getMinimumNode()
}

func isBlackNode(n *Node) bool {
	if n == nil || n.colour == colourBlack {
		return true
	}

	return false
}

func valueAlreadyExistsError(i int) error {
	return fmt.Errorf(`value %v already exists`, i)
}

func valueNotExistsError(i int) error {
	return fmt.Errorf(`value %v not exists`, i)
}

func invalidChildError(p, c int) error {
	return fmt.Errorf(`%v is not a child node of %v node`, c, p)
}
