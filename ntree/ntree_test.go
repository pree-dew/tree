package ntree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func exampleTree() *Tree[int, string] {
	// Create a new n-ary tree with a maximum of 2 keys per node
	t := New[int, string](5)

	// Insert key-value pairs into the tree
	t.Put(1, "a")
	t.Put(2, "b")
	t.Put(3, "c")
	t.Put(4, "d")
	t.Put(5, "e")
	t.Put(6, "f")
	t.Put(7, "g")
	t.Put(8, "h")
	t.Put(9, "i")

	return t
}

func TestGetElement(t *testing.T) {
	tr := exampleTree()

	// Retrieve the value associated with the key
	value, found := tr.Get(3)
	assert.True(t, found, "key 3 should be found")
	assert.Equal(t, value, "c", "value should be c")
}

func TestGetElementNotFound(t *testing.T) {
	tr := exampleTree()

	// Retrieve the value associated with the key
	_, found := tr.Get(16)
	assert.False(t, found, "key 16 should not be found")
}

func TestSize(t *testing.T) {
	tr := exampleTree()
	assert.Equal(t, tr.Size(), 9, "size should be 9")
}

func TestHeight(t *testing.T) {
	tr := exampleTree()
	assert.Equal(t, tr.Height(), 2, "height should be 2")
}

func TestPrint(t *testing.T) {
	tr := exampleTree()
	tr.Print(os.Stdout)
}
