package treap

import (
	"fmt"
	"github.com/CameronXie/algorithms-go/tree"
	"golang.org/x/exp/constraints"
	"io"
	"sync"
)

const (
	positionLeft  = true
	positionRight = false
)

type Node[K ~string, P constraints.Integer] struct {
	key      K
	priority P

	parent *Node[K, P]
	left   *Node[K, P]
	right  *Node[K, P]
}

func (n *Node[K, P]) String() string {
	return fmt.Sprintf("%v(%v)", n.key, n.priority)
}

func (n *Node[K, P]) Left() tree.Node {
	return n.left
}

func (n *Node[K, P]) Right() tree.Node {
	return n.right
}

func (n *Node[K, P]) Key() K {
	return n.key
}

func (n *Node[K, P]) Priority() P {
	return n.priority
}

func (n *Node[K, P]) traversal() []*Node[K, P] {
	l := []*Node[K, P]{n}

	for i := 0; i < len(l); i++ {
		current := l[i]
		if current.left != nil {
			l = append(l, current.left)
		}

		if current.right != nil {
			l = append(l, current.right)
		}
	}

	return l
}

func (n *Node[K, P]) search(key K) (*Node[K, P], error) {
	if key > n.key {
		if n.right == nil {
			return nil, valueNotExistsError(string(key))
		}

		return n.right.search(key)
	}

	if key < n.key {
		if n.left == nil {
			return nil, valueNotExistsError(string(key))
		}

		return n.left.search(key)
	}

	return n, nil
}

func (n *Node[K, P]) insert(node *Node[K, P]) error {
	if node.key > n.key {
		if n.right == nil {
			n.addChildNode(node, positionRight)
			return nil
		}

		return n.right.insert(node)
	}

	if node.key < n.key {
		if n.left == nil {
			n.addChildNode(node, positionLeft)
			return nil
		}

		return n.left.insert(node)
	}

	return valueAlreadyExistsError(string(node.key))
}

func (n *Node[K, P]) addChildNode(newNode *Node[K, P], position bool) {
	if newNode != nil {
		newNode.parent = n
	}

	if position == positionLeft {
		n.left = newNode
		return
	}

	n.right = newNode
}

func (n *Node[K, P]) replaceChildNode(old, new *Node[K, P]) {
	p := n.getChildNodePosition(old)

	n.removeChildNode(old)
	n.addChildNode(new, p)
}

func (n *Node[K, P]) removeChildNode(childNode *Node[K, P]) {
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

	panic(invalidChildError(string(childNode.key), string(n.key)))
}

func (n *Node[K, P]) getChildNodePosition(child *Node[K, P]) bool {
	if n.left != nil && n.left.key == child.key {
		return positionLeft
	}

	if n.right != nil && n.right.key == child.key {
		return positionRight
	}

	panic(invalidChildError(string(n.key), string(child.key)))
}

func NewNode[K ~string, P constraints.Integer](key K, priority P) *Node[K, P] {
	return &Node[K, P]{key: key, priority: priority}
}

type Treap[K ~string, P constraints.Integer] struct {
	root *Node[K, P]
	less func(i, j *Node[K, P]) bool
	mu   sync.RWMutex
}

func (t *Treap[K, P]) Search(key K) (*Node[K, P], error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.root == nil {
		return nil, valueNotExistsError(string(key))
	}

	return t.root.search(key)
}

func (t *Treap[K, P]) Print(w io.StringWriter) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.root == nil {
		_, err := w.WriteString("empty\n")
		return err
	}

	return tree.Print(t.root, w)
}

func (t *Treap[K, P]) Insert(n *Node[K, P]) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.root == nil {
		t.root = n
		return nil
	}

	if err := t.root.insert(n); err != nil {
		return err
	}

	t.up(n)
	return nil
}

func (t *Treap[K, P]) Update(key K, priority P) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.root == nil {
		return valueNotExistsError(string(key))
	}

	n, err := t.root.search(key)
	if err != nil {
		return err
	}

	oldPriority := n.priority
	n.priority = priority
	if t.less(n, &Node[K, P]{key: key, priority: oldPriority}) {
		t.up(n)
		return nil
	}

	t.down(n)
	return nil
}

func (t *Treap[K, P]) Delete(key K) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.root == nil {
		return valueNotExistsError(string(key))
	}

	n, err := t.root.search(key)
	if err != nil {
		return err
	}

	t.bottom(n)
	t.delete(n)
	return nil
}

func (t *Treap[K, P]) Pop() *Node[K, P] {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.root == nil {
		return nil
	}

	root := t.root
	t.bottom(root)
	t.delete(root)
	return root
}

func (t *Treap[K, P]) up(node *Node[K, P]) {
	for parent := node.parent; parent != nil; parent = node.parent {
		if t.less(parent, node) {
			break
		}

		p := parent.getChildNodePosition(node)
		if p == positionLeft {
			t.rotateRight(parent)
			continue
		}

		t.rotateLeft(parent)
	}
}

func (t *Treap[K, P]) down(node *Node[K, P]) {
	for node.left != nil || node.right != nil {
		left, right := node.left, node.right

		if left == nil {
			if !t.less(node, right) {
				break
			}

			t.rotateLeft(node)
			continue
		}

		if right == nil {
			if !t.less(node, left) {
				break
			}

			t.rotateRight(node)
			continue
		}

		higherPriority := positionRight
		if t.less(left, right) {
			higherPriority = positionLeft
		}

		if higherPriority == positionRight && t.less(node, right) {
			t.rotateLeft(node)
			continue
		}

		if higherPriority == positionLeft && t.less(node, left) {
			t.rotateRight(node)
			continue
		}

		break
	}
}

func (t *Treap[K, P]) bottom(node *Node[K, P]) {
	for {
		if node.left == nil && node.right == nil {
			break
		}

		if node.left == nil {
			t.rotateLeft(node)
			continue
		}

		if node.right == nil {
			t.rotateRight(node)
			continue
		}

		if t.less(node.left, node.right) {
			t.rotateRight(node)
			continue
		}

		t.rotateLeft(node)
	}
}

func (t *Treap[K, P]) delete(node *Node[K, P]) {
	if node.parent == nil {
		t.root = nil
		return
	}

	node.parent.removeChildNode(node)
}

func (t *Treap[K, P]) rotateLeft(n *Node[K, P]) {
	rightChild := n.right

	n.removeChildNode(rightChild)
	if rightChild.left != nil {
		n.addChildNode(rightChild.left, positionRight)
	}

	t.replaceChildNote(n, rightChild)
	rightChild.addChildNode(n, positionLeft)
}

func (t *Treap[K, P]) rotateRight(n *Node[K, P]) {
	leftChild := n.left

	n.removeChildNode(leftChild)
	if leftChild.right != nil {
		n.addChildNode(leftChild.right, positionLeft)
	}

	t.replaceChildNote(n, leftChild)
	leftChild.addChildNode(n, positionRight)
}

func (t *Treap[K, P]) replaceChildNote(oldNote *Node[K, P], newNote *Node[K, P]) {
	parent := oldNote.parent
	if parent != nil {
		parent.replaceChildNode(oldNote, newNote)
		return
	}

	t.root = newNote
}

func New[K ~string, P constraints.Integer](less func(i, j *Node[K, P]) bool) *Treap[K, P] {
	return &Treap[K, P]{less: less}
}

func valueAlreadyExistsError(i string) error {
	return fmt.Errorf(`value %v already exists`, i)
}

func valueNotExistsError(i string) error {
	return fmt.Errorf(`value %v not exists`, i)
}

func invalidChildError(p, c string) error {
	return fmt.Errorf(`%v is not a child node of %v node`, c, p)
}
