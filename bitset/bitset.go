package bitset

import (
	"errors"
	"strings"
)

var errOutOfRange = errors.New("index out of range")

// Set ...
type Set struct {
	data []byte
}

// Set sets a bit to 1.
func (s *Set) Set(i int) error {
	if i < 0 || i >= s.Size() {
		return errOutOfRange
	}

	idx := i / 8
	s.data[idx] = s.data[idx] | (1 << uint(i%8))

	return nil
}

// Get returns whether a bit is set or not.
func (s *Set) Get(i int) bool {
	return (s.data[i/8] & (1 << uint(i%8))) != 0
}

// Size returns the number of bits, both ones and zeroes.
func (s *Set) Size() int {
	return len(s.data) * 8
}

// Count returns the number of bits set to one.
func (s *Set) Count() int {
	var count int
	for _, b := range s.data {
		for b > 0 {
			count = count + int(b&1)
			b >>= 1
		}
	}
	return count
}

// All tests whether all bits are set.
func (s *Set) All() bool {
	for _, b := range s.data {
		if (b ^ 0xff) > 0 {
			return false
		}
	}

	return true
}

// Any tests whether any bit is set.
func (s *Set) Any() bool {
	for _, b := range s.data {
		if (b | 0x00) > 0 {
			return true
		}
	}
	return false
}

// None tests if no bits are set.
func (s *Set) None() bool {
	for _, b := range s.data {
		if (b | 0x00) > 0 {
			return false
		}
	}

	return true
}

func (s *Set) String() string {
	slice := make([]string, s.Size())
	for _, b := range s.data {
		for b > 0 {
			if (b & 0x08) > 0 {
				slice = append(slice, "1")
			} else {
				slice = append(slice, "0")
			}
			b <<= 1
		}
		slice = append(slice, "0")
	}
	return strings.Join(slice, "")
}

// New returns a new bitset with a given size.
func New(n int) *Set {
	size := n / 8
	if r := n % 8; r > 0 {
		size = size + 1
	}

	return &Set{
		data: make([]byte, size),
	}
}
