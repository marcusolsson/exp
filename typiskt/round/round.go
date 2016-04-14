package round

import (
	"errors"
	"log"
	"strings"
	"time"
)

// ErrEndOfWords is returned when there are no more words.
var ErrEndOfWords = errors.New("end of words")

// Word represents a word token.
type Word struct {
	Text string
	Done bool
	OK   bool
	Curr bool
}

// Event represents an event during the round.
type Event interface{}

// TypedEvent represents a typed keystroke.
type TypedEvent struct {
	Ch      rune
	Correct bool
	TypedAt time.Time
}

// CorrectedEvent represents a corrected keystroke.
type CorrectedEvent struct {
	Ch          rune
	CorrectedAt time.Time
}

// State represents a round state.
type State int

// Defines all possible round states.
const (
	RoundWaiting State = iota
	RoundStarted
	RoundStopped
)

// Round represents a round.
type Round struct {
	text     string
	Words    []Word
	Events   []Event
	State    State
	progress int
	Scratch  string

	StartedAt time.Time
	Took      time.Duration

	logger *log.Logger
}

// New returns a new instance of a Round.
func New(text string, logger *log.Logger) *Round {
	r := &Round{
		text:   text,
		logger: logger,
	}

	r.Reset()

	return r
}

func (r *Round) expected() rune {
	curr := r.Words[r.progress].Text
	if len(r.Scratch) < len(curr) {
		return rune(curr[len(r.Scratch)])
	}
	return ' '
}

// Advance add a new keystroke.
func (r *Round) Advance(rn rune) {

	if r.State == RoundStarted {
		r.logger.Printf("method=%s keystroke=%q expected=%q",
			"advance", rn, r.expected())

		r.Events = append(r.Events, TypedEvent{
			Ch:      rn,
			TypedAt: time.Now(),
			Correct: r.expected() == rn,
		})
	}

	r.Scratch += string(rn)
}

// Undo removes the last keystroke.
func (r *Round) Undo() {
	if len(r.Scratch) > 0 {
		rn := r.Scratch[len(r.Scratch)-1]
		r.Scratch = r.Scratch[:len(r.Scratch)-1]

		if r.State == RoundStarted {
			r.logger.Printf("method=%s keystroke=%q",
				"undo", rn)

			r.Events = append(r.Events, CorrectedEvent{
				Ch:          rune(rn),
				CorrectedAt: time.Now(),
			})
		}
	}
}

// Next moves to the next word.
func (r *Round) Next() error {
	if r.State == RoundStopped {
		r.Scratch = ""

		return nil
	}

	r.Words[r.progress].Curr = false
	r.Words[r.progress].Done = true

	if r.Scratch == r.Words[r.progress].Text {
		r.Words[r.progress].OK = true
	}

	r.logger.Printf("method=%s prev=%s correct=%v",
		"next_word", r.Words[r.progress].Text, r.Words[r.progress].OK)

	r.Events = append(r.Events, TypedEvent{
		Ch:      ' ',
		TypedAt: time.Now(),
		Correct: r.expected() == ' ',
	})

	r.Scratch = ""
	r.progress++

	if r.progress >= len(r.Words) {
		r.Stop()

		return ErrEndOfWords
	}

	r.Words[r.progress].Curr = true

	return nil
}

// Reset resets the round to its initial state.
func (r *Round) Reset() {
	r.State = RoundWaiting
	r.Scratch = ""
	r.Events = []Event{}

	r.Words = buildWords(r.text)
	r.progress = 0
}

// Start starts the round.
func (r *Round) Start() {
	r.State = RoundStarted
	r.StartedAt = time.Now()
}

// Stop stops the round.
func (r *Round) Stop() {
	r.Took = time.Since(r.StartedAt)
	r.State = RoundStopped
}

func buildWords(str string) []Word {
	var result []Word
	for _, s := range strings.Fields(str) {
		result = append(result, Word{Text: s})
	}
	return result
}

// Accuracy returns the percentage of correct keystrokes.
func (r *Round) Accuracy() float64 {
	all := float64(len(r.Typed()))
	bad := float64(len(r.Mistyped()))
	return 100.0 * (all - bad) / all
}

// Typed returns the events where a keystroke was made.
func (r *Round) Typed() []Event {
	var result []Event
	for _, e := range r.Events {
		if _, ok := e.(TypedEvent); ok {
			result = append(result, e)
		}
	}
	return result
}

// Mistyped returns the events where an incorrect keystroke was made.
func (r *Round) Mistyped() []Event {
	var result []Event
	for _, e := range r.Events {
		if te, ok := e.(TypedEvent); ok && !te.Correct {
			result = append(result, e)
		}
	}
	return result
}

// Corrections returns the events where a correction was made (e.g.
// Backspace).
func (r *Round) Corrections() []Event {
	var result []Event
	for _, e := range r.Events {
		if _, ok := e.(CorrectedEvent); ok {
			result = append(result, e)
		}
	}
	return result
}
