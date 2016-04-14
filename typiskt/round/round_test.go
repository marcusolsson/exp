package round

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestOneWord_Correct(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('e')
	r.AdvanceChar('s')
	r.AdvanceChar('t')
	r.AdvanceWord()

	r.Stop()

	if got := len(r.Typed()); got != 5 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 4)
	}
	if got := len(r.Mistyped()); got != 0 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 0)
	}
}

func TestOneWord_Mistyped(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('a')
	r.AdvanceChar('s')
	r.AdvanceChar('t')
	r.AdvanceWord()

	r.Stop()

	if got := len(r.Typed()); got != 5 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 5)
	}
	if got := len(r.Mistyped()); got != 1 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 1)
	}
}

func TestOneWord_EarlyNext(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('e')
	r.AdvanceChar('s')
	r.AdvanceWord()

	r.Stop()

	if got := len(r.Typed()); got != 4 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 4)
	}
	if got := len(r.Mistyped()); got != 1 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 1)
	}
}

func TestTwoWords_Correct(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test word", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('e')
	r.AdvanceChar('s')
	r.AdvanceChar('t')
	r.AdvanceWord()
	r.AdvanceChar('w')
	r.AdvanceChar('o')
	r.AdvanceChar('r')
	r.AdvanceChar('d')
	r.AdvanceWord()

	r.Stop()

	if got := len(r.Typed()); got != 10 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 10)
	}
	if got := len(r.Mistyped()); got != 0 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 0)
	}
}

func TestTwoWords_Mistyped(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test word", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('e')
	r.AdvanceChar('t')
	r.AdvanceWord()
	r.AdvanceChar('w')
	r.AdvanceChar('o')
	r.AdvanceChar('r')
	r.AdvanceChar('d')
	r.AdvanceWord()

	r.Stop()

	if got := len(r.Typed()); got != 9 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 9)
	}
	if got := len(r.Mistyped()); got != 2 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 2)
	}
}

func TestTwoWords_Stopped(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	r := New("test word", logger)

	r.Start()

	r.AdvanceChar('t')
	r.AdvanceChar('e')
	r.AdvanceChar('s')
	r.AdvanceChar('t')
	r.AdvanceWord()
	r.AdvanceChar('w')
	r.AdvanceChar('o')
	r.Stop()
	r.AdvanceChar('r')
	r.AdvanceChar('d')
	r.AdvanceWord()

	if got := len(r.Typed()); got != 7 {
		t.Errorf("len(r.Typed) = %d; want = %d", got, 7)
	}
	if got := len(r.Mistyped()); got != 0 {
		t.Errorf("len(r.Mistyped) = %d; want = %d", got, 0)
	}
}
