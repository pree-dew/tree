package ntree

import (
	"cmp"
	"fmt"
	"io"
	"strings"
)


// Tree is a generic n-ary tree
type Tree[K comparable, V any] struct {
	Root       *Node[K, V]
	Comparator func(x, y K) int
	size       int // total number of keys in the tree
	m          int // maximum number of keys in a node
}

type Node[K comparable, V any] struct {
	Parent   *Node[K, V]
	Children []*Node[K, V]
	Elements []*Element[K, V]
}

type Element[K comparable, V any] struct {
	Key   K
	Value V
}

// New returns a new n-ary tree
func New[K cmp.Ordered, V any](m int) *Tree[K, V] {
	return &Tree[K, V]{Comparator: cmp.Compare[K], m: m}
}

// Put inserts or updates a key-value pair into the tree
func (t *Tree[K, V]) Put(key K, value V) {
	ele := &Element[K, V]{Key: key, Value: value}
	if t.Root == nil {
		t.Root = &Node[K, V]{Elements: []*Element[K, V]{ele}}
		t.size++
		return
	}

	if t.insert(t.Root, ele) {
		t.size++
	}
}

// Get retrieves the value associated with the key from the tree
func (t *Tree[K, V]) Get(key K) (value V, found bool) {
	if t.Root == nil {
		return value, false
	}

	n, index, found := t.searchRecursive(t.Root, key)
	if found {
		return n.Elements[index].Value, true
	}

	return value, false
}

func (t *Tree[K, V]) GetNode(key K) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}

	n, _, found := t.searchRecursive(t.Root, key)
	return n, found
}

func (t *Tree[K, V]) Size() int {
	return t.size
}

func (n *Node[K, V]) Size() int {
	if n == nil {
		return 0
	}

	s := 1
	for _, c := range n.Children {
		s += c.Size()
	}

	return s
}

func (t *Tree[K, V]) Clear() {
	t.Root = nil
	t.size = 0
}

func (t *Tree[K, V]) Empty() bool {
	return t.size == 0
}

func (t *Tree[K, V]) Print(w io.Writer) {
	if t.Root == nil {
		return
	}

	t.print(w, t.Root, 0)
}

func (t *Tree[K, V]) print(w io.Writer, n *Node[K, V], level int) {
	if n == nil {
		return
	}

	for e := 0; e < len(n.Elements)+1; e++ {
		if e < len(n.Children) {
			t.print(w, n.Children[e], level+1)
		}

		if e < len(n.Elements) {
			w.Write([]byte(strings.Repeat("  ", level)))
			//fmt.Fprintf(w, "%v\n", len(n.Children))
			//w.Write([]byte(strings.Repeat("  ", level)))
			fmt.Fprintf(w, "%v\n", n.Elements[e].Key)
		}
	}
}

func (t *Tree[K, V]) Height() int {
	return t.height(t.Root)
}

func (t *Tree[K, V]) height(n *Node[K, V]) int {
	h := 0
	for ; n != nil; n = n.Children[0] {
		h++
		if len(n.Children) == 0 {
			break
		}
	}

	return h
}

func (t *Tree[K, V]) searchRecursive(n *Node[K, V], key K) (*Node[K, V], int, bool) {
	if t.Empty() {
		return nil, 0, false
	}

	for {
		ipos, found := t.search(n, key)
		if found {
			return n, ipos, true
		}

		if t.isLeaf(n) {
			return nil, -1, false
		}

		n = n.Children[ipos]
	}
}

func (t *Tree[K, V]) maxChildren() int {
	return t.m
}

func (t *Tree[K, V]) maxElements() int {
	return t.m - 1
}

func (t *Tree[K, V]) insert(n *Node[K, V], ele *Element[K, V]) bool {
	if t.isLeaf(n) {
		return t.insertIntoLeaf(n, ele)
	}

	return t.insertIntoChildren(n, ele)
}

func (t *Tree[K, V]) isLeaf(n *Node[K, V]) bool {
	return len(n.Children) == 0
}

// insertIntoLeaf inserts the element into the leaf node after
// finding the correct position for it in the elements slice.
func (t *Tree[K, V]) insertIntoLeaf(n *Node[K, V], ele *Element[K, V]) bool {
	ipos, found := t.search(n, ele.Key)
	if found {
		n.Elements[ipos] = ele
		return false
	}

	n.Elements = append(n.Elements, nil)
	copy(n.Elements[ipos+1:], n.Elements[ipos:])
	n.Elements[ipos] = ele
	t.split(n)
	return true
}

// insertIntoChildren finds the correct child node to insert the element
// into and recursively calls insert on that child node.
func (t *Tree[K, V]) insertIntoChildren(n *Node[K, V], ele *Element[K, V]) bool {
	ipos, found := t.search(n, ele.Key)
	if found {
		n.Elements[ipos] = ele
		return false
	}

	return t.insert(n.Children[ipos], ele)
}

// search finds the correct position for the key in the elements slice
// by using binary search. It returns the index and a boolean indicating
// whether the key was found.
func (t *Tree[K, V]) search(n *Node[K, V], key K) (int, bool) {
	if n == nil {
		return 0, false
	}

	lo, hi := 0, len(n.Elements)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		comp := t.Comparator(key, n.Elements[mid].Key)
		switch {
		case comp == 0:
			return mid, true
		case comp > 0:
			lo = mid + 1
		case comp < 0:
			hi = mid - 1
		}
	}

	return lo, false
}

func (t *Tree[K, V]) split(n *Node[K, V]) {
	if !t.shouldSplit(n) {
		return
	}

	if t.isRoot(n) {
		t.splitRoot()
		return
	}

	t.splitNonRoot(n)
}

func (t *Tree[K, V]) shouldSplit(n *Node[K, V]) bool {
	return len(n.Elements) > t.maxElements()
}

func (t *Tree[K, V]) isRoot(n *Node[K, V]) bool {
	return n == t.Root
}

func (t *Tree[K, V]) splitRoot() {
	mid := (t.m - 1) / 2
	left := &Node[K, V]{Elements: t.Root.Elements[:mid]}
	right := &Node[K, V]{Elements: t.Root.Elements[mid+1:]}

	// what if the root has children?
	if !t.isLeaf(t.Root) {
		left.Children = append([]*Node[K, V](nil), t.Root.Children[:mid+1]...)
		right.Children = append([]*Node[K, V](nil), t.Root.Children[mid+1:]...)
		for _, c := range left.Children {
			c.Parent = left
		}
		for _, c := range right.Children {
			c.Parent = right
		}
	}

	newRoot := &Node[K, V]{
		Elements: []*Element[K, V]{t.Root.Elements[mid]},
		Children: []*Node[K, V]{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	t.Root = newRoot
}

func (t *Tree[K, V]) splitNonRoot(n *Node[K, V]) {
	mid := (t.m - 1) / 2
	parent := n.Parent

	left := &Node[K, V]{Elements: n.Elements[:mid], Parent: parent}
	right := &Node[K, V]{Elements: n.Elements[mid+1:], Parent: parent}

	if !t.isLeaf(n) {
		left.Children = append([]*Node[K, V](nil), n.Children[:mid+1]...)
		right.Children = append([]*Node[K, V](nil), n.Children[mid+1:]...)
		for _, c := range left.Children {
			c.Parent = left
		}
		for _, c := range right.Children {
			c.Parent = right
		}
	}

	ipos, _ := t.search(parent, n.Elements[mid].Key)
	parent.Elements = append(parent.Elements, nil)
	copy(parent.Elements[ipos+1:], parent.Elements[ipos:])
	parent.Elements[ipos] = n.Elements[mid]

	parent.Children[ipos] = left

	parent.Children = append(parent.Children, nil)
	copy(parent.Children[ipos+2:], parent.Children[ipos+1:])
	parent.Children[ipos+1] = right

	t.split(parent)
}
