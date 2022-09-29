package types

import "math"

type BitSet struct {
	v []int64
}

func NewBitSet(bits int) *BitSet {
	return &BitSet{
		v: make([]int64, int(math.Ceil(float64(bits)/64.0))),
	}
}

func (b *BitSet) Value(n int) bool {
	if n >= len(b.v)*64 {
		return false
	}

	return (b.v[n/64] & (1 << (n % 64))) != 0
}

func (b *BitSet) BitsSet() int {
	count := 0

	for i := 0; i < len(b.v)*64; i++ {
		if b.Value(i) {
			count++
		}
	}

	return count
}

func (b *BitSet) Set1(n int) {
	for len(b.v) <= n/64 {
		b.v = append(b.v, 0)
	}

	b.v[n/64] |= 1 << (n % 64)
}

func (b *BitSet) Set0(n int) {
	for len(b.v) <= n/64 {
		b.v = append(b.v, 0)
	}

	b.v[n/64] &= ^(1 << (n % 64))
}
