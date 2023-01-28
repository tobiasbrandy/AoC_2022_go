package hashext

import (
	"encoding/binary"
	"io"

	"golang.org/x/exp/constraints"
)

type Hasher interface {
	Hash(io.Writer)
}

func HashNum[T constraints.Integer | constraints.Float | bool | []bool](h io.Writer, n T) {
	switch v := any(n).(type) {
	case int:
		_ = binary.Write(h, binary.BigEndian, uint64(v))
	case uint:
		_ = binary.Write(h, binary.BigEndian, uint64(v))
	case uintptr:
		_ = binary.Write(h, binary.BigEndian, uint64(v))
	default:
		_ = binary.Write(h, binary.BigEndian, n)
	}
}

func HashNumP[T constraints.Integer | constraints.Float | bool](h io.Writer, n *T) {
	switch v := any(n).(type) {
	case *int:
		_ = binary.Write(h, binary.BigEndian, uint64(*v))
	case *uint:
		_ = binary.Write(h, binary.BigEndian, uint64(*v))
	case *uintptr:
		_ = binary.Write(h, binary.BigEndian, uint64(*v))
	default:
		_ = binary.Write(h, binary.BigEndian, *n)
	}
}
