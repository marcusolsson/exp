package main

import (
	"errors"
	"math/rand"
)

// EventType is the event type.
type EventType int

// Available event types
const (
	Joined EventType = iota
	Failed
)

// Update represents a change to the member list.
type Update struct {
	Member Member
	Type   EventType
	Round  int8
}

// Member is a node in a cluster.
type Member struct {
	Name    string
	Address string
}

// List contains the members of a cluster.
type List struct {
	Members map[string]Member
	Failed  map[string]Member
	Updates []Update
	Rounds  int
}

// NewList returns a new instance of a member list.
func NewList(rounds int) *List {
	return &List{
		Members: make(map[string]Member),
		Failed:  make(map[string]Member),
		Updates: make([]Update, 0),
		Rounds:  rounds,
	}
}

// Add adds a member to the member list.
func (l *List) Add(m Member) {
	if _, ok := l.Members[m.Address]; !ok {
		l.Members[m.Address] = m
		l.Updates = append(l.Updates, Update{Member: m, Type: Joined})
		delete(l.Failed, m.Address)
	}
}

// Remove removes a member from the member list.
func (l *List) Remove(m Member) {
	if _, ok := l.Members[m.Address]; ok {
		l.Failed[m.Address] = m
		l.Updates = append(l.Updates, Update{Member: m, Type: Failed})
		delete(l.Members, m.Address)
	}
}

// Merge updates the member list with updates.
func (l *List) Merge(updates []Update) {
	for _, u := range updates {
		switch u.Type {
		case Joined:
			l.Add(u.Member)
		case Failed:
			l.Remove(u.Member)
		}
	}
}

// IncrementRound increases round of all updates.
func (l *List) IncrementRound() {
	var updates []Update
	for _, v := range l.Updates {
		if v.Round < int8(l.Rounds) {
			clone := v
			clone.Round = clone.Round + 1
			updates = append(updates, clone)
		}
	}
	l.Updates = updates
}

// Random picks k random members from the member list.
func (l List) Random(k int, exclude ...Member) ([]Member, error) {
	var otherMembers []Member
	for _, m := range l.Members {
		include := true

		for _, e := range exclude {
			if m.Address == e.Address {
				include = false
			}
		}

		if include {
			otherMembers = append(otherMembers, m)
		}
	}

	if len(otherMembers) == 0 {
		return nil, errors.New("empty member list")
	}

	if len(otherMembers) < k {
		k = len(otherMembers)
	}

	var result []Member
	for i := 0; i < k; i++ {
		r := rand.Intn(len(otherMembers))
		result = append(result, otherMembers[r])
		otherMembers = append(otherMembers[:r], otherMembers[r+1:]...)
	}

	return result, nil
}
