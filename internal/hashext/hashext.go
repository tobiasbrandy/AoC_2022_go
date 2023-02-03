package hashext

import (
	"encoding/binary"
	"fmt"
	"hash/maphash"
	"io"

	"golang.org/x/exp/constraints"
)

type Hasher interface {
	Hash(io.Writer)
}

type num interface {
	constraints.Integer | constraints.Float | bool
}

func HashNum[T num](h io.Writer, n T) {
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

func HashNumPtr[T num](h io.Writer, n *T) {
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

func HashNumArr[T num](h io.Writer, arr []T) {
	switch vs := any(arr).(type) {
	case []int:
		l := len(arr)
		cast := make([]uint64, l, l)
		for i, v := range vs {
			cast[i] = uint64(v)
		}
		_ = binary.Write(h, binary.BigEndian, cast)
	case []uint:
		l := len(arr)
		cast := make([]uint64, l, l)
		for i, v := range vs {
			cast[i] = uint64(v)
		}
		_ = binary.Write(h, binary.BigEndian, cast)
	case []uintptr:
		l := len(arr)
		cast := make([]uint64, l, l)
		for i, v := range vs {
			cast[i] = uint64(v)
		}
		_ = binary.Write(h, binary.BigEndian, cast)
	default:
		_ = binary.Write(h, binary.BigEndian, arr)
	}
}

func HashString(h io.Writer, s string) {
	_, _ = h.Write([]byte(s))
}

func HashAny(h io.Writer, x any) {
	switch v := x.(type) {
	case Hasher:
		v.Hash(h)
	case string:
		HashString(h, v)
	case int:
		HashNum(h, v)
	case uint:
		HashNum(h, v)
	case uintptr:
		HashNum(h, v)
	case *int:
		HashNumPtr(h, v)
	case *uint:
		HashNumPtr(h, v)
	case *uintptr:
		HashNumPtr(h, v)
	case []int:
		HashNumArr(h, v)
	case []uint:
		HashNumArr(h, v)
	case []uintptr:
		HashNumArr(h, v)
	default:
		// Assume fixed number
		if err := binary.Write(h, binary.BigEndian, v); err != nil {
			return
		}

		panic(fmt.Sprint("value ", x, " cannot be hashed"))
	}
}

type HashMap[K Hasher, V any] struct {
	hash *maphash.Hash
	Keys map[uint64]K
	Vals map[uint64]V
}

func NewHashMap[K Hasher, V any]() HashMap[K, V] {
	return HashMap[K, V]{
		hash: &maphash.Hash{},
		Keys: make(map[uint64]K),
		Vals: make(map[uint64]V),
	}
}

func (m HashMap[K, V]) Hash(key K) uint64 {
	key.Hash(m.hash)
	ret := m.hash.Sum64()
	m.hash.Reset()
	return ret
}

func (m HashMap[K, V]) GetKey(key K) (V, bool) {
	ret, ok := m.Vals[m.Hash(key)]
	return ret, ok
}

func (m HashMap[K, V]) GetVal(key K) (V, bool) {
	ret, ok := m.Vals[m.Hash(key)]
	return ret, ok
}

func (m HashMap[K, V]) Set(key K, val V) {
	m.SetH(m.Hash(key), key, val)
}

func (m HashMap[K, V]) SetH(h uint64, key K, val V) {
	m.Keys[h] = key
	m.Vals[h] = val
}
