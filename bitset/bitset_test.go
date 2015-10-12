package bitset

import "testing"

func TestSet(t *testing.T) {
	s := New(1)

	s.Set(0)

	if !s.Get(0) {
		t.Fatal("s.Get() should be", true)
	}
}

func TestSize(t *testing.T) {
	s := New(3)

	if s.Size() != 8 {
		t.Fatal("invalid size")
	}
}

func TestCount(t *testing.T) {
	s := New(3)

	s.Set(0)
	s.Set(1)

	if s.Count() != 2 {
		t.Fatal("s.Count() should be", 2)
	}
}

func TestAll(t *testing.T) {
	s1 := New(8)
	for i := 0; i < s1.Size(); i++ {
		s1.Set(i)
	}

	if !s1.All() {
		t.Fatal("s1.All() should return", true)
	}

	s2 := New(8)
	s2.Set(1)

	if s2.All() {
		t.Fatal("s1.All() should return", false)
	}
}

func TestAny(t *testing.T) {
	s1 := New(3)
	s1.Set(2)

	if !s1.Any() {
		t.Fatal("s1.Any() should return", true)
	}

	s2 := New(3)

	if s2.Any() {
		t.Fatal("s1.Any() should return", false)
	}
}

func TestNone(t *testing.T) {
	s1 := New(2)
	s1.Set(2)

	if s1.None() {
		t.Fatal("s1.None() should return", false)
	}

	s2 := New(3)

	if !s2.None() {
		t.Fatal("s1.None() should return", true)
	}
}

func TestString(t *testing.T) {
	s := New(4)
	s.Set(1)
	s.Set(2)

	if s.String() != "01100000" {
		t.Fatalf("s.String() should return %s, was %s", "01100000", s.String())
	}
}

func TestIndexTooBig(t *testing.T) {
	s := New(2)

	if err := s.Set(9); err == nil {
		t.Fatal(err)
	}
}
func TestIndexTooSmall(t *testing.T) {
	s := New(2)

	if err := s.Set(-1); err == nil {
		t.Fatal(err)
	}
}

func BenchmarkInsert(b *testing.B) {
	N := 10000000

	bs := New(N)

	for i := 0; i < bs.Size(); i++ {
		bs.Set(i)
	}
}
