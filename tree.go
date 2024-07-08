package tree

import "io"

type Tree[V any] interface {
	Empty() bool
	Size() int
	Height() int
	Print(io.Writer)
}
